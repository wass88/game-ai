import * as APIType from "./api-types";

import { useState, useEffect } from "react";

const Url = `${window.location.protocol}//${window.location.hostname}:${window.location.port}`;

function encodeQuery(query: {[key: string]: string}) {
  const ret = [];
  for (let d in query)
    ret.push(encodeURIComponent(d) + '=' + encodeURIComponent(query[d]));
  return '?' + ret.join('&');
}

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
  console.log("Fetch Data", data, url);
  return data;
}

const API = {
  async logout(): Promise<void> {
    return fetch_with_cookie("POST", "/github/logout", {});
  },
  async you(): Promise<any> {
    return fetch_with_cookie("GET", "/api/you", {});
  },
  async post_match(game_id: number, ai_id: number[]): Promise<any> {
    return fetch_with_cookie("POST", "/api/matches", {
      game_id,
      ai_id,
    });
  },
  async post_ai_github(
    game_id: number,
    github: string,
    branch: string
  ): Promise<{ ai_github_id: number }> {
    return fetch_with_cookie("POST", "/api/ai-githubs", {
      game_id,
      github,
      branch,
    });
  },
  async latest_ai(game_id: number): Promise<any> {
    return fetch_with_cookie("GET", "/api/games/" + game_id + "/latest-ai", {});
  },
  async ai_githubs(game_id: number): Promise<APIType.AIGithub[]> {
    return fetch_with_cookie(
      "GET",
      "/api/games/" + game_id + "/ai-githubs",
      {}
    );
  },
  async match(match_id: number): Promise<APIType.Match> {
    return fetch_with_cookie("GET", "/api/matches/" + match_id, {});
  },
  async matches(game_id: number, page?: number): Promise<{pages: number, matches: APIType.Match[]}> {
    console.log("Y", {game_id, page});
    let query : any = {};
    if (page) {
      query.page = page.toString()
    }
    let url = "/api/games/" + game_id + "/matches" + encodeQuery(query);
    return fetch_with_cookie("GET", url, {});
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
  useCallAPI<P extends any[], T>(
    api: (...args: P) => Promise<T>,
    args: P,
    after: (resp: T) => any
  ): [
    () => void,
    boolean,
    T | null,
    React.Dispatch<React.SetStateAction<T | null>>
  ] {
    const [call, setCall] = useState(false);
    const [sending, setSending] = useState(false);
    const [res, setRes] = useState<T | null>(null);
    useEffect(() => {
      (async () => {
        if (!call) return;
        setSending(true);
        const resp = await api(...args);
        setRes(resp);
        setSending(false);
        after(resp);
      })();
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [call]);
    return [() => setCall(true), sending, res, setRes];
  },
};

export default API;
