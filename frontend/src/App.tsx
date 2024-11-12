import { Route, Switch } from "wouter";
import { PlaylistState } from "./state/playlist";
import { CallbackPage } from "./pages/callback.tsx";
import { PlaylistView } from "./pages/playlist.tsx";
import { useAuthState } from "./state/auth.tsx";
import { AuthPage } from "./pages/auth.tsx";

const IndexPage = () => {
  const auth = useAuthState();
  if (auth.isAuthenticated)
    return (
      <PlaylistState>
        <PlaylistView />
      </PlaylistState>
    );
  return <AuthPage />;
};

export const App = () => (
  <Switch>
    <Route path="/" component={IndexPage} />
    <Route path="/auth/callback" component={CallbackPage} />

    <Route>404: No such page!</Route>
  </Switch>
);
