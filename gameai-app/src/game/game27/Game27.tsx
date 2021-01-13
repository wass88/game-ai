import React from "react";
import "./Game27.css";
import {GamePlayer} from "../GamePlayer"

export function Game27(props: {
  first: string;
  second: string;
  record: string;
}) {
  return <>
    <GamePlayer first={props.first} second={props.second} record={props.record} GameView={Game27View}></GamePlayer>
  </>
}

function Game27View(props: {
  first: string;
  second: string;
  record: string[];
  nth_played: number
}) {
  const {first, second, record, nth_played} = props;

  const FIRST = 1;
  const SECOND = -1;
  const SIZE = 9;
  let board : [number, string][][] = Array.from(new Array(SIZE), () => []);
  board[0] = Array.from(new Array(SIZE), (_, i) => [FIRST, `f${i}`]);
  board[SIZE-1] = Array.from(new Array(SIZE), (_, i) => [SECOND, `s${i}`]);

  if (record[0] == null || record[0] === "") {
    return <></>
  }
  record.some((r, i) : boolean => {
    if (i >= nth_played) {
      return true;
    }

    const current_player = i % 2 === 0 ? FIRST : SECOND;
    if (r === "pass") {
      return false;
    }

    const [move, nc, nk] = r.split(" ");
    if (move !== "move") {
      throw `Invalid act ${move}`
    }
    const c = parseInt(nc, 10);
    const k = parseInt(nk, 10);

    const towers = board.filter((c) => c.length > 0 && c[0][0] === current_player).length;

    const d = c + towers * current_player;
    board[d] = board[c].slice(0, k).concat(board[d]);
    board[c] = board[c].slice(k);

    return false;
  })
  let first_score = board[SIZE-1].length;
  let second_score = board[0].length;
  return <>
  <h3>
    <div className="chip board board-8"></div>{first}:{first_score}
    vs
    <div className="chip board board-0"></div>{second}:{second_score}
    </h3>
  <div className="game27-board">
    {
      Array.from(new Array(SIZE)).map((_, i) => {
        return <div key={`board-${i}`} className={`chip board board-${i}`}
          style={{ gridArea: `-2 / ${i+1} / -1 / ${i+2}`}}
        ></div>
      })
    }
    {
      board.map((c, i) => 
        c.reverse().map((t, j) => {
          const color = t[0] === FIRST ? "red" : "blue";
          const pos = j;
          return <div key={t[1]} className={`chip ${color}`}
            style={{ gridArea: `${-pos-2} / ${i+1} / ${-pos-3} / ${i+2}`}}
          ></div>
        })
      ).flat()
    }
  </div>
  </>

}