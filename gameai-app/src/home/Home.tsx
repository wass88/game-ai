import React, {useContext} from "react";
import { Link } from "react-router-dom";
import { LoginUserContext } from "../login";

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

              <p>
                Reversi: {" "}
                <Link to="/games/1/matches">Matches</Link> {" "}
                <Link to="/games/1/githubs">AIs</Link>
              </p>

              <p>
                戦い方
              </p>
              <ol>
                <li>AIを作る</li>
                <li>Dockerfileを書く</li>
                <li>Githubにpush</li>
                <li>AIを登録</li>
                <li>戦いが始まる</li>
              </ol>
              <p><a href="https://github.com/wass88/reversi-random">wass88/reversi-random</a>を参照のこと</p>
            </>
          );
        }
        return <a href="/github/login">Login</a>;
      })()}
    </>
  );
}
