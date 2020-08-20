import React from "react";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import Home from "./home/Home";
import { MatchesPage } from "./matches/Matches";
import { MatchPage } from "./match/Match";

import "./App.css";

export default function App() {
  return (
    <Router>
      <div id="outer">
        <header>
          <h1>Game AI</h1>
          <Link to="/game/1/matches">Matches</Link>
        </header>
        <main>
          <Switch>
            <Route path="/" exact>
              <Home />
            </Route>
            <Route path="/game/:id/matches" exact>
              <MatchesPage />
            </Route>
            <Route path="/matches/:id" exact>
              <MatchPage />
            </Route>
          </Switch>
        </main>
      </div>
    </Router>
  );
}
