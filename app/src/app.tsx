import { Switch, Route } from 'wouter'
import { Layout } from './pages/Layout'
import { Auth } from './components/Auth'
import { UserProfile } from './components/UserProfile'
import { PostView } from './components/PostView'

export function App() {
  return (
    <Switch>
      <Route path="/" component={Layout} />
      <Route path="/login">{() => <Auth mode="login" />}</Route>
      <Route path="/register">{() => <Auth mode="register" />}</Route>
      <Route path="/p/:id" component={PostView} />
      <Route path="/u/:username" component={UserProfile} />
    </Switch>
  )
}
