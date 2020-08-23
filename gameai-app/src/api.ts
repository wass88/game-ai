import * as APIType from "./api-types";

import { useState, useEffect } from "react";

const Url = "http://localhost:3000";

async function fetch_with_cookie(
  method: string,
  path: string,
  body: any
): Promise<any> {
  const url = Url + path;
  const conf: any = {
    method: method,
    headers: { "Content-Type": "application/json" },
    credentials: "same-origin" as "same-origin",
  };
  if (method === "POST") {
    conf.body = JSON.stringify(body);
  }

  const resp = await fetch(url, conf);
  if (resp.status !== 200) {
    throw new Error("Failed call API (" + resp.status + ") " + url);
  }
  const data = await resp.json();
  console.log("Data", data);
  return data;
}

const API = {
  async you(): Promise<any> {
    return fetch_with_cookie("GET", "/api/you", {});
  },
  async post_ai_github(
    game_id: number,
    user_id: number,
    github: string,
    branch: string
  ): Promise<{ ai_github_id: number }> {
    return fetch_with_cookie("POST", "/api/ai-githubs", {
      game_id,
      user_id,
      github,
      branch,
    });
  },
  async ai_githubs(game_id: number): Promise<APIType.AIGithub[]> {
    return fetch_with_cookie(
      "GET",
      "/api/games/" + game_id + "/ai-githubs",
      {}
    );
  },
  async matche(match_id: number): Promise<APIType.Match> {
    return fetch_with_cookie("GET", "/api/matches/" + match_id, {});
  },
  async matches(game_id: number): Promise<APIType.Match[]> {
    return fetch_with_cookie("GET", "/api/games/" + game_id + "/matches", {});
  },
  useAPI<P extends any[], T>(
    api: (...args: P) => Promise<T>,
    args: P
  ): [T | null, React.Dispatch<React.SetStateAction<T | null>>] {
    const [data, setData] = useState<T | null>(null);
    useEffect(() => {
      (async () => {
        const res = await api(...args);
        setData(res);
      })();
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [api, ...args]);
    return [data, setData];
  },
};

export default API;
