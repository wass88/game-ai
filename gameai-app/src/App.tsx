import React from "react";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import Home from "./home/Home";
import { MatchesPage } from "./matches/Matches";
import { MatchPage } from "./match/Match";
import { GithubsPage } from "./githubs/Githubs";
import { AIPage } from "./ai/AI";
import API from "./api";
import { LoginUserContext } from "./login"

import "./App.css";

export default function App() {
  const [you] = API.useAPI(API.you, []);
  return (
    <Router>
      <LoginUserContext.Provider value={you}>
      <div id="outer">
        <div id="inner">
          <header>
            <Link to="/"> <h1 className="">
              <div className="head-icon"></div>
              <div className="title">Game AI</div>
            </h1> </Link>
            {(()=>{
              if (you != null && you.name != null) {
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
              <Route path="/ai/:id" exact>
                <AIPage />
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
