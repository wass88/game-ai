import React from "react";
import API from "../api";

export default function Home() {
  const [you] = API.useAPI(API.you, []);
  return (
    <>
      <p>ゲームAI対戦ベース</p>
      {(() => {
        if (!you) return <p></p>;
        if (you.login) {
          return <p>Welcome {you.name}</p>;
        }
        return <a href="/github/login">Login</a>;
      })()}
    </>
  );
}
