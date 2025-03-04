import { AppBar, Toolbar, Typography, Button } from "@mui/material";
import { loginWithGoogle, logout, getToken } from "../api/auth";

function Navbar() {
  const token = getToken();

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" sx={{ flexGrow: 1 }}>
          YouTube Recommender
        </Typography>
        {token ? (
          <Button color="inherit" onClick={logout}>
            Log out
          </Button>
        ) : (
          <Button color="inherit" onClick={loginWithGoogle}>
            Log in with Google
          </Button>
        )}
      </Toolbar>
    </AppBar>
  );
}

export default Navbar;
