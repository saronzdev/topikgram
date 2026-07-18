package validator

import "testing"

func TestEmail(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"user@example.com", true},
		{"test@test.co.uk", true},
		{"", false},
		{"not-an-email", false},
		{"@domain.com", false},
		{"user@", false},
		{"a@b.c", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Email(tt.input); got != tt.want {
				t.Errorf("Email(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestUsername(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"john_doe", true},
		{"abc123", true},
		{"a", false},
		{"ab", false},
		{"a_very_long_username_1", false},
		{"UPPERCASE", false},
		{"with-hyphen", false},
		{"valid_name_123", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Username(tt.input); got != tt.want {
				t.Errorf("Username(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestPassword(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"short", false},
		{"1234567", false},
		{"12345678", true},
		{"a valid password here", true},
		{"exactly24characterslong!", true},
		{"this password is way too long for the limit", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Password(tt.input); got != tt.want {
				t.Errorf("Password(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		value string
		max   int
		want  bool
	}{
		{"hello", 10, true},
		{"hello world", 5, false},
		{"  spaced  ", 10, true},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := MaxLength(tt.value, tt.max); got != tt.want {
				t.Errorf("MaxLength(%q, %d) = %v, want %v", tt.value, tt.max, got, tt.want)
			}
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		value string
		min   int
		want  bool
	}{
		{"hello", 3, true},
		{"hi", 3, false},
		{"  a  ", 1, true},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := MinLength(tt.value, tt.min); got != tt.want {
				t.Errorf("MinLength(%q, %d) = %v, want %v", tt.value, tt.min, got, tt.want)
			}
		})
	}
}
