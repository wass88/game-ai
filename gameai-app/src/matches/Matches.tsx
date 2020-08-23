import React from "react";
import { useParams } from "react-router-dom";
import API from "../api";
import * as APIType from "../api-types";
import "./Matches.css";
import { MatchDesc } from "../match/Match";

export function MatchesPage() {
  const { id } = useParams();
  const [matches] = API.useAPI(API.matches, [id]);
  return Matches(matches);
}
export function Matches(matches: APIType.Match[] | null) {
  let head = <h1>Match Results</h1>;
  let body = <p> Loading ... </p>;
  if (matches) {
    const matchList = matches.map((match) => MatchDesc(match));
    body = <div className="match-list">{matchList}</div>;
  }
  return (
    <>
      {" "}
      {head} {body}
    </>
  );
}
