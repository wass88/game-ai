import React, {FunctionComponent, useState} from "react";
import { Button } from "../components/Button";
import "./GamePlayer.css"

export function GamePlayer(props: {
  first: string;
  second: string;
  record: string;
  GameView: FunctionComponent<{first: string, second: string, record: string[], nth_played: number}>;
}) {
  const first = props.first;
  const second = props.second;
  const record = props.record.split("\n");
  const GameView = props.GameView;
  const [nth_played, set_nth_played] = useState(record.length);
  return (
    <div className="game-player">
        <div className="game-view">
            {GameView({first, second, record, nth_played})}
        </div>
        <div className="game-button-1">
          <Button onClick={() => set_nth_played(nth_played > 0 ? nth_played - 1 : nth_played)}>Prev</Button>
        </div>
        <div className="game-button-2">
          <Button onClick={() => set_nth_played(nth_played < record.length ? nth_played + 1 : nth_played)}>Next</Button>
        </div>
        <div className="game-records">
          {record.map((e, i) => {
            return (
              <span
                onClick={() => set_nth_played(i + 1)}
                key={i}
                className={nth_played === i + 1 ? "highlight" : ""}
              >
                {i % 2 === 1 ? "(" : ""}
                {e}
                {i % 2 === 1 ? ")" : ""}
              </span>
            );
          })}
        </div>
    </div>
  )
}