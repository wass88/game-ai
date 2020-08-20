import React from "react";
import { useParams, Link } from "react-router-dom";
import API from "../api";
import * as APIType from "../api-types";
import "./Match.css";

export function MatchPage() {
  const { id } = useParams();
  const [match] = API.useAPI(API.matche, [id]);
  return Match(match);
}
export function Match(match: APIType.Match | null) {
  return (
    <div>
      <h1>Match of {match?.game?.name}</h1>
      <p>{match?.record ?? "<None>"}</p>
      {(() => {
        if (match) return MatchDesc(match);
        return <p>None</p>;
      })()}
      {match?.results.map((result, i) => (
        <div key={i}>
          <h2>{result.ai?.ai_github.github} Stderr</h2>
          <code>
            <pre>{result.stderr}</pre>
          </code>
        </div>
      ))}
    </div>
  );
}

export function MatchDesc(match: APIType.Match) {
  return (
    <Link to={"/matches/" + match.id}>
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
    </Link>
  );
}
