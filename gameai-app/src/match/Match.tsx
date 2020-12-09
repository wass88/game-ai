import React from "react";
import { useParams, Link } from "react-router-dom";
import API from "../api";
import * as APIType from "../api-types";
import { Reversi } from "../game/reversi/Reversi";
import "./Match.css";

export function MatchPage() {
  const { id } = useParams();
  const [match] = API.useAPI(API.match, [id]);
  return Match(match);
}
export function Match(match: APIType.Match | null) {
  return (
    <>
      <Link to={"/games/" + match?.game?.id + "/matches"}>
        <h1>Match of {match?.game?.name}</h1>
      </Link>
      {(() => {
        if (match)
          return (
            <>
              <Reversi
                first={match?.results[0].ai?.ai_github.github || "first"}
                second={match?.results[0].ai?.ai_github.github || "second"}
                record={match?.record || ""}
              />
              {MatchDesc(match)}
            </>
          );
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
    </>
  );
}

export function MatchDesc(match: APIType.Match) {
  return (
    <React.Fragment key={match.id}>
      <Link className="no-line" to={"/matches/" + match.id}>
        <div className="match">
          <p className="head">
            [{match.state}] #{match.id}
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
                  {result.ai?.ai_github?.github} ({result.ai?.ai_github?.branch}
                  )
                </p>
              </div>
            ))}
          </div>
        </div>
      </Link>
    </React.Fragment>
  );
}
