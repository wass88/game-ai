import React from "react";
import API from "../api";

export default function Home() {
  const [you] = API.useAPI(API.you, []);
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
            </>
          );
        }
        return <a href="/github/login">Login</a>;
      })()}
    </>
  );
}
