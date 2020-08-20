import React from "react";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  useParams,
} from "react-router-dom";
import Home from "./home/Home";
import { MatchesPage } from "./matches/Matches";

import "./App.css";

function Hey() {
  const { id } = useParams();
  console.log("Hey");
  return <p>Hey {id}</p>;
}
export default function App() {
  return (
    <Router>
      <div id="outer">
        <header>
          <h1>Game AI</h1>
          <Link to="/game/:id/matches">Matches</Link>
        </header>
        <main>
          <Switch>
            <Route path="/" exact>
              <Home />
            </Route>
            <Route path="/hey/:id" exact>
              <Hey />
            </Route>
            <Route path="/game/:id/matches" exact>
              <MatchesPage />
            </Route>
          </Switch>
        </main>
      </div>
    </Router>
  );
}
