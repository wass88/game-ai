import React from "react";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import Home from "./home/Home";
import { MatchesPage } from "./matches/Matches";
import { MatchPage } from "./match/Match";
import { GithubsPage } from "./githubs/Githubs";

import "./App.css";

export default function App() {
  return (
    <Router>
      <div id="outer">
        <div id="inner">
          <header>
            <h1>Game AI</h1>

            <Link to="/games/1/matches">Matches</Link>

            <Link to="/games/1/githubs">AIs</Link>

            {/* TODO User info */}
            {/* TODO  Game List */}
            {/* TODO  Links of Current Game */}
          </header>
          <main>
            <Switch>
              <Route path="/" exact>
                <Home />
              </Route>
              <Route path="/games/:id/matches" exact>
                <MatchesPage />
              </Route>
              <Route path="/matches/:id" exact>
                <MatchPage />
              </Route>
              <Route path="/games/:id/githubs" exact>
                <GithubsPage />
              </Route>
              <Route path="*" status={404}>
                <p>Not found</p>
              </Route>
            </Switch>
          </main>
        </div>
      </div>
    </Router>
  );
}
