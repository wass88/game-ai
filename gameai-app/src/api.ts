import * as APIType from "./api-types";

import { useState, useEffect } from "react";

const Url = "http://localhost:8000";

const API = {
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
    console.log(api, args);
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
