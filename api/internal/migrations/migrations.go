package migrations

import (
	"context"
	"embed"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"topikgram/api/internal/store"
)

//go:embed *.sql
var migrationFiles embed.FS

type Migrator struct {
	conn *store.Pool
}

func New(conn *store.Pool) *Migrator {
	return &Migrator{conn: conn}
}

func (m *Migrator) Up(ctx context.Context) error {
	// Crear tabla de tracking si no existe
	if _, err := m.conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INT PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	// Obtener migraciones ya aplicadas
	rows, err := m.conn.Query(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return err
		}
		applied[v] = true
	}

	// Leer archivos de migración
	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var migrations []struct {
		version int
		name    string
	}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		// Extraer número: "001_initial.sql" → 1
		parts := strings.SplitN(name, "_", 2)
		if len(parts) < 2 {
			continue
		}
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			continue // ignora archivos sin número
		}
		migrations = append(migrations, struct {
			version int
			name    string
		}{v, name})
	}

	// Ordenar por versión
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	// Ejecutar las pendientes
	for _, mig := range migrations {
		if applied[mig.version] {
			continue
		}

		content, err := migrationFiles.ReadFile(mig.name)
		if err != nil {
			return fmt.Errorf("read %s: %w", mig.name, err)
		}

		// Transacción por migración
		tx, err := m.conn.Begin(ctx)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, string(content)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("migrate %s: %w", mig.name, err)
		}

		if _, err := tx.Exec(ctx,
			"INSERT INTO schema_migrations (version) VALUES ($1)",
			mig.version,
		); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("track %s: %w", mig.name, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit %s: %w", mig.name, err)
		}

		log.Printf("applied migration %s", mig.name)
	}

	return nil
}
