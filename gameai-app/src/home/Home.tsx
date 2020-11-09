import React, {useContext} from "react";
import { Link } from "react-router-dom";
import {LoginUserContext} from "../App";

export default function Home() {
  const you = useContext(LoginUserContext)
  return (
    <>
      <p>ゲームAI対戦サイト</p>
      {(() => {
        if (!you) return <p>loading...</p>;
        if (you.login) {
          return (
            <>
              <p>
                Welcome {you.name} <a href="/github/logout"> Logout </a>
              </p>
              <p>Your account is {you.authority}.</p>

              Reversi: {" "}
              <Link to="/games/1/matches">Matches</Link> {" "}

              <Link to="/games/1/githubs">AIs</Link>
            </>
          );
        }
        return <a href="/github/login">Login</a>;
      })()}
    </>
  );
}
