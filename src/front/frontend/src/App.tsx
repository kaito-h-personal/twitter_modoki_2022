import {
  Box,
  Button,
  Container,
  Grid,
  Link,
  TextField,
  Typography,
} from "@mui/material";
import { useEffect, useState } from "react";
import "./App.css";

// import reactLogo from "./assets/react.svg"; //TODO: 消す

function App() {
  const [count, setCount] = useState(0);

  const [posts, setPosts] = useState([]);
  useEffect(() => {
    fetch("http://localhost:8006/", { method: "GET" })
      .then((res) => res.json())
      .then((data) => {
        setPosts(data);
      });
  }, []);

  return (
    <div className="App">
      <Container maxWidth="xs">
        <Box
          sx={{
            marginTop: 8,
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
          }}
        >
          <Typography component="h1" variant="h4">
            ログイン
          </Typography>

          <Box component="form" noValidate sx={{ mt: 1 }}>
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="メールアドレス"
              name="email"
              autoComplete="email"
              autoFocus
            />

            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="パスワード"
              type="password"
              id="password"
              autoComplete="current-password"
            />

            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
            >
              ログイン
            </Button>

            <Grid container>
              <Grid item xs>
                <Link href="#" variant="body2">
                  パスワードを忘れた
                </Link>
              </Grid>

              <Grid item>
                <Link href="#" variant="body2">
                  新規登録
                </Link>
              </Grid>
            </Grid>
          </Box>
        </Box>
      </Container>
      <div>👤&lt;{posts.hello}</div>
      <div>👤&lt;起きた〜</div>
      <div>👤&lt;お腹すいた〜</div>
    </div>
  );
}

export default App;
