import { Typography } from "@mui/material";
import { useEffect, useState } from "react";
import "./App.css";

// import reactLogo from "./assets/react.svg"; //TODO: 消す

import FavoriteIcon from "@mui/icons-material/Favorite";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import ShareIcon from "@mui/icons-material/Share";
import Avatar from "@mui/material/Avatar";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import CardHeader from "@mui/material/CardHeader";
import { red } from "@mui/material/colors";
import IconButton, { IconButtonProps } from "@mui/material/IconButton";
import { styled } from "@mui/material/styles";

interface ExpandMoreProps extends IconButtonProps {
  expand: boolean;
}

const ExpandMore = styled((props: ExpandMoreProps) => {
  const { expand, ...other } = props;
  return <IconButton {...other} />;
})(({ theme, expand }) => ({
  transform: !expand ? "rotate(0deg)" : "rotate(180deg)",
  marginLeft: "auto",
  transition: theme.transitions.create("transform", {
    duration: theme.transitions.duration.shortest,
  }),
}));

function App() {
  const [tweets, setTweets] = useState([]);

  useEffect(() => {
    fetch("http://localhost:8006/tweets")
      .then((response) => response.json())
      .then((data) => setTweets(data));
  }, []);

  return (
    <div className="App">
      <div>
        <h1>Tweets</h1>
        <ul>
          {tweets.map((tweet) => (
            <li key={tweet.id}>
              <p>{tweet.text}</p>
              <p>Auther: {tweet.auther}</p>
              <p>Created At: {tweet.created_at}</p>
            </li>
          ))}
        </ul>
      </div>
      {/* ログインフォームここから
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
      ログインフォームここまで */}

      <Card sx={{ maxWidth: 345 }}>
        <CardHeader
          avatar={
            <Avatar sx={{ bgcolor: red[500] }} aria-label="recipe">
              Ri
            </Avatar>
          }
          action={
            <IconButton aria-label="settings">
              <MoreVertIcon />
            </IconButton>
          }
          title="梨香子"
          subheader="2023/02/04"
        />
        {/* <CardMedia
          component="img"
          height="194"
          image="/static/images/cards/paella.jpg"
          alt="Paella dish"
        /> */}
        <CardContent>
          <Typography variant="body2" color="text.secondary">
            起きた〜
          </Typography>
        </CardContent>
        <CardActions disableSpacing>
          <IconButton aria-label="add to favorites">
            <FavoriteIcon />
          </IconButton>
          <IconButton aria-label="share">
            <ShareIcon />
          </IconButton>
        </CardActions>
      </Card>

      <Card sx={{ maxWidth: 345 }}>
        <CardHeader
          avatar={
            <Avatar sx={{ bgcolor: red[500] }} aria-label="recipe">
              K
            </Avatar>
          }
          action={
            <IconButton aria-label="settings">
              <MoreVertIcon />
            </IconButton>
          }
          title="Kouji"
          subheader="2023/02/02"
        />
        {/* <CardMedia
          component="img"
          height="194"
          image="/static/images/cards/paella.jpg"
          alt="Paella dish"
        /> */}
        <CardContent>
          <Typography variant="body2" color="text.secondary">
            お腹すいた〜
          </Typography>
        </CardContent>
        <CardActions disableSpacing>
          <IconButton aria-label="add to favorites">
            <FavoriteIcon />
          </IconButton>
          <IconButton aria-label="share">
            <ShareIcon />
          </IconButton>
        </CardActions>
      </Card>
    </div>
  );
}

export default App;
