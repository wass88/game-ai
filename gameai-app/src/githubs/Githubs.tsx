import React, { useEffect } from "react";
import { useState } from "react";
import { useParams, Link, useHistory } from "react-router-dom";
import API from "../api";
import * as APIType from "../api-types";
import "./Githubs.css";
import { Button } from "../components/Button";
import { Input } from "../components/Input";

export function GithubsPage() {
  const { id } = useParams();
  const [githubs] = API.useAPI(API.ai_githubs, [id]);
  return Githubs(githubs);
}
export function Githubs(ai_githubs: APIType.AIGithub[] | null) {
  const header = (
    <div>
      <h1>AIs</h1>
      <Button onClick={() => setShowPopup(true)}>
        <p>New AI Config</p>
      </Button>
    </div>
  );
  const [showPopup, setShowPopup] = useState(false);
  const popup = PopupGithub(showPopup, setShowPopup);
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

export function PopupGithub(
  show: boolean,
  set: React.Dispatch<React.SetStateAction<boolean>>
) {
  const [github, setGithub] = useState("");
  const [branch, setBranch] = useState("master");
  const [post, setPost] = useState(false);
  const [sending, setSending] = useState(false);
  const history = useHistory();
  useEffect(() => {
    (async () => {
      if (!post) return;
      setSending(true);
      const resp = await API.post_ai_github(1, github, branch);
      setSending(false);
      history.push(`/ai/${resp.ai_github_id}`);
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [post]);

  function formatBranch(branch: string) {
    return branch.replace(/\s/, "");
  }
  function formatGithub(github: string) {
    return github.replace(/\s/, "");
  }
  function onChangeGithub(e: any) {
    setGithub(formatGithub(e.target.value));
  }
  function onChangeBranch(e: any) {
    setBranch(formatBranch(e.target.value));
  }
  function validateGithub(): string | undefined {
    if (github.indexOf("/") < 0) {
      return "Need '/' in user/reponame";
    }
  }
  function validateBranch(): string | undefined {
    if (branch === "") {
      return "branch is Empty";
    }
  }
  function validate() {
    return validateGithub() || validateBranch();
  }
  function createAI() {
    setPost(true);
  }
  function url() {
    return `https://github.com/${github}/tree/${branch}`;
  }
  return (
    <div
      className="popup"
      onClick={() => set(false)}
      style={{ display: show ? "flex" : "none" }}
    >
      <div onClick={(e) => e.stopPropagation()}>
        <h2>New AI Config</h2>
        <Input
          value={github}
          errMsg={validateGithub()}
          onChange={onChangeGithub}
        >
          Github Repository: username/reponame (eg: wass80/reversi-random)
        </Input>
        <Input
          value={branch}
          errMsg={validateBranch()}
          onChange={onChangeBranch}
        >
          Branch Name: (eg: master)
        </Input>
        <p>
          Github Page:{" "}
          <a target="_blank" rel="noopener noreferrer" href={url()}>
            {url()}
          </a>
        </p>
        <Button
          onClick={createAI}
          disabled={validate() !== undefined || sending}
        >
          Create AI
        </Button>
      </div>
    </div>
  );
}
