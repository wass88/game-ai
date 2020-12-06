import React from "react";
import { useState } from "react";
import { useParams, useHistory, Link} from "react-router-dom";
import API from "../api";
import * as APIType from "../api-types";
import "./Githubs.css";
import { Button, Input, Popup, useStateValidate, Format } from "../components";
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

export function GithubsPage() {
  const { gameID } = useParams();
  const [githubs] = API.useAPI(API.ai_githubs, [gameID]);
  return Githubs(gameID, githubs);
}
export function Githubs(gameID: number, ai_githubs: APIType.AIGithub[] | null) {
  const [show, setShow] = useState(false);
  const popup = (
    <Popup show={show} setShow={setShow}>
      <FormGithub setShow={setShow}></FormGithub>
    </Popup>
  );
  const header = (
    <div>
      <h1>AIs</h1>
      <Link to={`/games/${gameID}/matches`}>List of Matches</Link>
      <Button onClick={() => setShow(true)}>
        <p>Create New Config of AI</p>
      </Button>
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
  let before_updated_at = (updated_at) ? dayjs(updated_at).fromNow() : "none"
 
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

function FormGithub(setShow: any) {
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
    [1, github.value, branch.value],
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
