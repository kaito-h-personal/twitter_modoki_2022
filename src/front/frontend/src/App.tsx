import { useEffect, useState } from "react";
import "./App.css";
// import reactLogo from "./assets/react.svg";

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
      <div>👤&lt;{posts.hello}</div>
      <div>👤&lt;起きた〜</div>
      <div>👤&lt;お腹すいた〜</div>
    </div>
  );
}

export default App;
