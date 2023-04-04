import { Typography } from "@mui/material";
import { useEffect, useState } from "react";
import "./App.css";

// import reactLogo from "./assets/react.svg"; //TODO: 消す

import FavoriteIcon from "@mui/icons-material/Favorite";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import ShareIcon from "@mui/icons-material/Share";
import { Button, TextField } from "@mui/material";
import Avatar from "@mui/material/Avatar";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import CardContent from "@mui/material/CardContent";
import CardHeader from "@mui/material/CardHeader";
import { red } from "@mui/material/colors";
import IconButton from "@mui/material/IconButton";

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
        name: "John Doe",
        email: "johndoe@example.com",
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
    auther: number;
    created_at: string;
    text: string;
  };

  return (
    <div className="App">
      <div>
        {tweets.map((tweet) => (
          <div key={tweet.id}>
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
                title={tweet.auther}
                subheader={tweet.created_at}
              />
              {/* <CardMedia
          component="img"
          height="194"
          image="/static/images/cards/paella.jpg"
          alt="Paella dish"
        /> */}
              <CardContent>
                <Typography variant="body2" color="text.secondary">
                  {tweet.text}
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
        ))}
      </div>
      <form onSubmit={handleSubmit}>
        <TextField
          label="Name"
          value={name}
          onChange={(event) => setName(event.target.value)}
          margin="normal"
          variant="outlined"
        />
        <Button type="submit" variant="contained" color="primary">
          Submit
        </Button>
      </form>
    </div>
  );
}

export default App;
