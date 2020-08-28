import React from "react";
import { useState } from "react";
import { useParams, Link, useHistory } from "react-router-dom";
import API from "../api";
import * as APIType from "../api-types";
import "./Githubs.css";
import { Button, Input, Popup, useStateValidate, Format } from "../components";

export function GithubsPage() {
  const { id } = useParams();
  const [githubs] = API.useAPI(API.ai_githubs, [id]);
  return Githubs(githubs);
}
export function Githubs(ai_githubs: APIType.AIGithub[] | null) {
  const [show, setShow] = useState(false);
  const popup = (
    <Popup show={show} setShow={setShow}>
      <FormGithub setShow={setShow}></FormGithub>
    </Popup>
  );
  const header = (
    <div>
      <h1>AIs</h1>
      <Button onClick={() => setShow(true)}>
        <p>New AI Config</p>
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
  return (
    <div key={ai_github.id}>
      <a
        href={
          "https://github.com/" +
          ai_github?.github +
          "/tree/" +
          ai_github?.branch
        }
      >
        <h2>
          {ai_github?.github} ({ai_github?.branch}){" "}
          {ai_github?.latest_ai?.state}
        </h2>
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

  function validate() {
    return github.err || branch.err;
  }
  function url() {
    return `https://github.com/${github}/tree/${branch}`;
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
        Github Page:{" "}
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
