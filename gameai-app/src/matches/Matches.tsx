import React, { useState, useContext } from "react";
import { useParams, Link } from "react-router-dom";
import API from "../api";
import * as APIType from "../api-types";
import "./Matches.css";
import { MatchDesc } from "../match/Match";
import { Button, Popup, Select } from "../components";
import { isVisitor, LoginUserContext } from "../login"
import { useQuery } from "../Query"

export function MatchesPage() {
  let { gameID } = useParams();
  gameID = parseInt(gameID, 10);
  let { page } = useQuery("page");
  let pageV = page ? parseInt(page, 10) : undefined;
  const [matches] = API.useAPI(API.matches, [gameID, pageV]);
  return Matches(gameID, matches, pageV);
}
export function Matches(gameID: number, matches: {pages: number, matches: APIType.Match[]} | null, page?: number) {
  const you = useContext(LoginUserContext)
  const [show, setShow] = useState(false);
  const [created, setCreated] = useState(false);
  const popup = (
    <Popup show={show} setShow={setShow}>
      <FormMatch gameID={gameID} setShow={setShow} setCreated={setCreated} />
    </Popup>
  );
  let head = <>
    <h1>Match Results</h1>
    <Link to={`/games/${gameID}/githubs`}>List of AIs</Link>
  </>
  if (!matches) {
    return (
      <>
        {head}
        <p> Loading ... </p>;
      </>
    );
  }
  const matchList = matches.matches.map((match) => MatchDesc(match));
  return (
    <>
      {head}
      {(() => {
        if (created) {
          return <p>Match Created</p>;
        }
        return <></>;
      })()}
      {(() => {
        if (isVisitor(you)) {
          return <p>Matchを作る権限がありません。</p>
        }
        return <Button onClick={() => setShow(true)}>
        <p>Create Match</p>
      </Button>
      })()}
      <div className="match-list">{matchList}</div>
      {
        (() => {
          let pageV = (page ? page : 0)
          if ( pageV < matches.pages - 1) {
            return <Link to={`?page=${pageV + 1}`}>Next Page...</Link>
          }
        })()
      }
      {popup}
    </>
  );
}

function FormMatch(props: any) {
  const [ais] = API.useAPI(API.latest_ai, [props.gameID]);
  const [playoutAI1, setPlayoutAI1] = useState<any>(undefined);
  const [playoutAI2, setPlayoutAI2] = useState<any>(undefined);
  const options = ais?.map((ai: any) => {
    return {
      id: ai.id,
      name: ai.github + ":" + ai.branch,
    };
  });
  function disabled() {
    return playoutAI1 === undefined || playoutAI2 === undefined;
  }
  const [createMatch, sending] = API.useCallAPI<any, any>(
    API.post_match,
    [props.gameID, [playoutAI1?.id, playoutAI2?.id]],
    (resp) => {
      props.setShow(false);
      props.setCreated(true);
    }
  );
  return (
    <>
      <h1>Create Match</h1>
      <Select options={options} onChange={setPlayoutAI1} />
      <Select options={options} onChange={setPlayoutAI2} />
      <Button onClick={() => createMatch()} disabled={disabled() || sending}>
        <p>Enqueue Match</p>
      </Button>
    </>
  );
}
