import React from "react";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import Home from "./home/Home";
import { MatchesPage } from "./matches/Matches";
import { MatchPage } from "./match/Match";
import { GithubsPage } from "./githubs/Githubs";
import API from "./api";


import "./App.css";

export const LoginUserContext = React.createContext<any>(null)

export default function App() {
  const [you] = API.useAPI(API.you, []);
  return (
    <Router>
      <LoginUserContext.Provider value={you}>
      <div id="outer">
        <div id="inner">
          <header>
            <Link to="/"> <h1>Game AI</h1> </Link>
            {(()=>{
              if (you != null) {
                return <p className="userinfo"> you are {you.name}</p>
              }
            })() }
          </header>
          <main>
            <Switch>
              <Route path="/" exact>
                <Home />
              </Route>
              <Route path="/games/:gameID/matches" exact>
                <MatchesPage />
              </Route>
              <Route path="/matches/:id" exact>
                <MatchPage />
              </Route>
              <Route path="/games/:gameID/githubs" exact>
                <GithubsPage />
              </Route>
              <Route path="*" status={404}>
                <p>Not found</p>
              </Route>
            </Switch>
          </main>
        </div>
      </div>
      </LoginUserContext.Provider>
    </Router>
  );
}
