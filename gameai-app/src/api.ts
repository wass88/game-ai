import * as APIType from "./api-types";

import { useState, useEffect } from "react";

const Url = "http://localhost:3000";

const API = {
  async post_ai_github(
    game_id: number,
    user_id: number,
    github: string,
    branch: string
  ): Promise<{ ai_github_id: number }> {
    const url = Url + "/api/ai-githubs";
    const body = { game_id, user_id, github, branch };
    const conf = {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    };

    const resp = await fetch(url, conf);
    if (resp.status !== 200) {
      throw new Error("Failed call API (" + resp.status + ") " + url);
    }
    const data = await resp.json();
    return data;
  },
  async ai_githubs(game_id: number): Promise<APIType.AIGithub[]> {
    const url = Url + "/api/games/" + game_id + "/ai-githubs";
    const resp = await fetch(url);
    if (resp.status !== 200) {
      throw new Error("Failed call API (" + resp.status + ") " + url);
    }
    const data = await resp.json();
    return data;
  },
  async matche(match_id: number): Promise<APIType.Match> {
    const url = Url + "/api/matches/" + match_id;
    const resp = await fetch(url);
    if (resp.status !== 200) {
      throw new Error("Failed call API (" + resp.status + ") " + url);
    }
    const data = await resp.json();
    return data;
  },
  async matches(game_id: number): Promise<APIType.Match[]> {
    const url = Url + "/api/games/" + game_id + "/matches";
    const resp = await fetch(url);
    if (resp.status !== 200) {
      throw new Error("Failed call API (" + resp.status + ") " + url);
    }
    const data = await resp.json();
    return data;
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
