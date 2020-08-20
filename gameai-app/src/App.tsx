import React from "react";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import Home from "./home/Home";
import Matches from "./matches/Matches";

import "./App.css";

function hey() {
  return <p>Hey</p>;
}
export default function App() {
  return (
    <Router>
      <p>Menu</p>
      <Link to="/hey">Hey</Link>
      <Link to="/game/1/matches">game</Link>
      <Switch>
        <Route path="/" exact component={Home} />
        <Route path="/hey" component={hey} />
        <Route path="/game/:id/matches" component={Matches} />
      </Switch>
    </Router>
  );
}
