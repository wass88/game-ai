import React from "react";
import { useParams } from "react-router-dom";
import API from "../api";

export default function Matches() {
  const { id } = useParams();
  const [matches] = API.useAPI(API.matches, [id]);
  if (matches === null) {
    return <p> Loading ... </p>;
  }
  const matchList = matches.map((match) => (
    <p key={match.id}> {match.state} </p>
  ));
  return <div> {matchList} </div>;
}
