import { Typography } from "@mui/material";
import { useEffect, useState } from "react";
import "./App.css";

// import reactLogo from "./assets/react.svg"; //TODO: 消す

import { useTheme } from "@mui/material/styles";

import Stack from "@mui/material/Stack";

import Grid from "@mui/material/Grid";

import FavoriteIcon from "@mui/icons-material/Favorite";
import { Box, Button, TextField } from "@mui/material";
import Avatar from "@mui/material/Avatar";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import CardHeader from "@mui/material/CardHeader";
import IconButton from "@mui/material/IconButton";

import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";

import { useCookies } from "react-cookie";
import { useNavigate } from "react-router-dom";

function App() {
  const theme = useTheme();
  const appBarHeight = theme.mixins.toolbar.minHeight;

  type User = {
    name: string;
    icon_img: string;
  };

  type Tweet = {
    id: string;
    user_name: string;
    created_at: string;
    text: string;
    icon_img: string;
  };

  const [user, setUser] = useState<User>({
    name: "--",
    icon_img: "",
  });

  const [tweets, setTweets] = useState<Tweet[]>([]);

  const [inputTweetText, setInputTweetText] = useState("");

  const [cookies] = useCookies(["user_id"]);

  // tweetを投稿
  const handleSubmit: React.FormEventHandler<HTMLFormElement> = (event) => {
    event.preventDefault();

    // 空投稿はNG
    if (inputTweetText == "") {
      return;
    }

    fetch("http://localhost:8006/add_tweets", {
      method: "POST",
      body: JSON.stringify({
        user_id: cookies.user_id,
        text: inputTweetText,
      }),
    })
      .then((response) => response.json())
      .then((data) => setTweets(data))
      .catch((error) => console.error(error));
  };

  const navigate = useNavigate();
  useEffect(() => {
    if (!cookies.user_id) {
      navigate("/");
    }

    // ユーザー情報を取得
    fetch("http://localhost:8006/user", {
      method: "POST",
      // TODO: http://localhost:8006/user/user:6のようにGETにする？(脆弱性？)
      body: JSON.stringify({
        user_id: cookies.user_id,
      }),
    })
      .then((response) => response.json())
      .then((data) => setUser(data))
      .catch((error) => console.error(error));

    // tweet一覧を取得
    fetch("http://localhost:8006/tweets")
      .then((response) => response.json())
      .then((data) => setTweets(data))
      .catch((error) => console.error(error));
  }, []);

  return (
    <div className="App">
      {/* ヘッダー */}
      <AppBar position="fixed">
        <Toolbar>
          <Typography variant="h6" color="inherit" noWrap>
            Twitterもどき
          </Typography>
        </Toolbar>
      </AppBar>
      {/* ヘッダーの分の余白を挿入 */}
      <div style={{ marginTop: appBarHeight, paddingTop: 20 }} />
      <Grid container>
        {/* 画面左側 */}
        <Grid xs={8}>
          {tweets.map((tweet) => (
            <div key={tweet.id}>
              <Card sx={{ margin: 2 }}>
                <CardHeader
                  avatar={
                    <Avatar
                      src={"data:image/png;base64," + tweet.icon_img}
                      sx={{ width: 50, height: 50 }}
                    />
                  }
                  title={tweet.user_name}
                  subheader={tweet.created_at}
                />
                <CardContent>
                  <Typography variant="body2" color="text.secondary">
                    {tweet.text}
                  </Typography>
                </CardContent>
                <CardActions disableSpacing sx={{ justifyContent: "flex-end" }}>
                  <IconButton aria-label="add to favorites">
                    <FavoriteIcon />
                  </IconButton>
                </CardActions>
              </Card>
            </div>
          ))}
        </Grid>
        {/* 画面右側(fixedでスクロールしても動かないようにする) */}
        <Grid xs={4} position="fixed" sx={{ right: 7 }}>
          <Stack direction="row" alignItems="center" spacing={2}>
            <Avatar
              src={"data:image/png;base64," + user.icon_img}
              sx={{ width: 50, height: 50 }}
            />
            <div>{user.name} さん</div>
          </Stack>
          <form onSubmit={handleSubmit}>
            <Box display="flex" flexDirection="column">
              <TextField
                label="呟きたいことを入力"
                value={inputTweetText}
                onChange={(event) => setInputTweetText(event.target.value)}
                margin="normal"
                variant="outlined"
                multiline
                rows={4}
              />
              <Button type="submit" variant="contained" color="primary">
                投稿
              </Button>
            </Box>
          </form>
        </Grid>
      </Grid>
    </div>
  );
}

export default App;
