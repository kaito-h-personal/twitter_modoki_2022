import React from "react";
import { Link } from "react-router-dom";

const MyComponent: React.FC = () => {
  return (
    <div>
      <h1>Login Page</h1>
      <Link to="/App">
        <p>To App Page</p>
      </Link>
    </div>
  );
};

export default MyComponent;
