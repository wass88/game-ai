import React from "react";
import { useParams } from "react-router-dom";
import API from "../api";
import * as APIType from "../api-types";
import "./Matches.css";

export function MatchesPage() {
  const { id } = useParams();
  const [matches] = API.useAPI(API.matches, [id]);
  return Matches(matches);
}
export function Matches(matches: APIType.Match[] | null) {
  let head = <h1>Match Result</h1>;
  let body = <p> Loading ... </p>;
  if (matches) {
    const matchList = matches.map((match) => (
      <div className="match" key={match.id}>
        <p className="head">
          [{match.state}] {match.id}
        </p>
        {(() => {
          if (match.exception !== "")
            return <p className="exception">{match.exception}</p>;
          return;
        })()}
        <div className="results">
          {match.results?.map((result, i) => (
            <div className="tr" key={i}>
              <p className="score">{result.result ?? "??"}</p>
              <p className="ai">
                {result.ai?.ai_github?.github} ({result.ai?.ai_github?.branch})
              </p>
            </div>
          ))}
        </div>
      </div>
    ));
    body = <div className="match-list">{matchList}</div>;
  }
  return (
    <div>
      {" "}
      {head} {body}
    </div>
  );
}
