import { Typography } from "@mui/material";
import { useEffect, useState } from "react";
import "./App.css";

// import reactLogo from "./assets/react.svg"; //TODO: 消す

import { useTheme } from "@mui/material/styles";

import Stack from "@mui/material/Stack";

import Grid2 from "@mui/material/Unstable_Grid2";

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

function App() {
  const [tweets, setTweets] = useState<Tweet[]>([]);

  const [name, setName] = useState("");

  const handleSubmit = (event: React.MouseEvent) => {
    event.preventDefault();
    console.log(`Name: ${name}`);
    fetch("http://localhost:8006/add_tweets", {
      method: "POST",
      // headers: {
      //   "Content-Type": "application/json",
      // },
      body: JSON.stringify({
        user_id: "user:6", // 決めうち
        text: name,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        setTweets(data);
      })
      .catch((error) => console.error(error));
  };

  useEffect(() => {
    fetch("http://localhost:8006/tweets")
      .then((response) => response.json())
      .then((data) => setTweets(data));
  }, []);

  type Tweet = {
    id: string;
    user_name: string;
    created_at: string;
    text: string;
    icon_img: string;
  };

  const theme = useTheme();
  const appBarHeight = theme.mixins.toolbar.minHeight;

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
      <div style={{ marginTop: appBarHeight }} />
      <Grid2 container>
        {/* 画面左側 */}
        <Grid2 xs={8}>
          {tweets.map((tweet) => (
            <div key={tweet.id}>
              <Card sx={{ margin: 2 }}>
                <CardHeader
                  avatar={
                    <Avatar
                      alt="Remy Sharp"
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
        </Grid2>
        {/* 画面右側 */}
        <Grid2 xs={4}>
          <Stack direction="row" alignItems="center" spacing={2}>
            <Avatar alt="x" src="x" sx={{ width: 50, height: 50 }} />
            <div>xx さん</div>
          </Stack>
          <form onSubmit={handleSubmit}>
            <Box display="flex" flexDirection="column">
              <TextField
                label="呟きたいことを入力"
                value={name}
                onChange={(event) => setName(event.target.value)}
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
        </Grid2>
      </Grid2>
    </div>
  );
}

export default App;
