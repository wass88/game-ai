import React from "react";
import { useState , useContext } from "react";
import { useParams, useHistory, Link} from "react-router-dom";
import API from "../api";
import * as APIType from "../api-types";
import "./Githubs.css";
import { Button, Input, Popup, useStateValidate, Format } from "../components";
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import utc from 'dayjs/plugin/utc'
import { isVisitor, LoginUserContext } from "../login"

export function GithubsPage() {
  let { gameID } = useParams();
  gameID = parseInt(gameID, 10);
  const [githubs] = API.useAPI(API.ai_githubs, [gameID]);
  return Githubs(gameID, githubs);
}
export function Githubs(gameID: number, ai_githubs: APIType.AIGithub[] | null) {
  const you = useContext(LoginUserContext)
  const [show, setShow] = useState(false);
  const popup = (
    <Popup show={show} setShow={setShow}>
      <FormGithub gameID={gameID}></FormGithub>
    </Popup>
  );
  const header = (
    <div>
      <h1>AIs</h1>
      <Link to={`/games/${gameID}/matches`}>List of Matches</Link>
      {
        (() => {
          if (isVisitor(you)) {
            return <p> AIを作る権限がありません。{you?.github_name ?? ""} </p>
          }
          return <Button onClick={() => setShow(true)}>
            <p>Create New Config of AI</p>
          </Button>
        })()
      }
    </div>
  );
  if (ai_githubs === null) {
    return (
      <div>
        {header} <p>Loading...</p>
      </div>
    );
  }
  return (
    <div>
      {header}
      {ai_githubs.map((ai_github) => AIGithubDesc(ai_github))}
      {popup}
    </div>
  );
}

export function AIGithubDesc(ai_github: APIType.AIGithub) {
  let updated_at = ai_github?.latest_ai?.updated_at;
  dayjs.extend(relativeTime)
  dayjs.extend(utc);
  let before_updated_at = (updated_at) ? dayjs(updated_at).from(dayjs.utc()) : "none"
 
  return (
    <div key={ai_github.id}>
      <a className="no-decoration"
        href={
          "https://github.com/" +
          ai_github?.github +
          "/tree/" +
          ai_github?.branch
        }
      >
        <h2>
          {ai_github?.github} ({ai_github?.branch})
        </h2>
      </a>
      <a className="no-decoration"
        href={
          "https://github.com/" +
          ai_github?.github +
          "/tree/" +
          ai_github?.latest_ai?.commit
        }
      >
        <p>
          {" ["} {ai_github?.latest_ai?.state ?? "finding"} {"] "}
          { "Rate: " }{ai_github?.latest_ai?.rate ?? "(not rated)"} {" "}
          {" Commit ID: "} {ai_github?.latest_ai?.commit.substr(0, 6 ?? "none" )}
          {" Last Update: "} {before_updated_at}
        </p>
      </a>
    </div>
  );
}

function FormGithub(props: any) {
  const {gameID} = props;
  const github = useStateValidate("", Format.trimSpace, (s) => {
    if (s.split("/").length === 1) {
      return "Need '/' in user/reponame";
    }
  });
  const branch = useStateValidate("master", Format.trimSpace, (s) => {
    if (s === "") {
      return "branch is Empty";
    }
  });

  const history = useHistory();

  const [createAI, sending] = API.useCallAPI(
    API.post_ai_github,
    [gameID, github.value, branch.value],
    (resp) => {
      history.push(`/ai/${resp.ai_github_id}`);
    }
  );

  function url() {
    return `https://github.com/${github.value}/tree/${branch.value}`;
  } //validate() !== undefined || sending
  return (
    <>
      <h2>New AI Config</h2>
      <Input value={github.value} errMsg={github.err} onChange={github.change}>
        Github Repository: username/reponame (eg: wass80/reversi-random)
      </Input>
      <Input value={branch.value} errMsg={branch.err} onChange={branch.change}>
        Branch Name: (eg: master)
      </Input>
      <p>
        Please Check the Github Page:{" "}
        <a target="_blank" rel="noopener noreferrer" href={url()}>
          {url()}
        </a>
      </p>

      <Button onClick={createAI} disabled={sending}>
        Create AI
      </Button>
    </>
  );
}
