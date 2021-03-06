import React from "react";
import "./Reversi.css";
import {GamePlayer} from "../GamePlayer"

function ReversiView(props: {
  first: string;
  second: string;
  record: string[];
  nth_played: number
}) {
  const {first, second, record, nth_played} = props;
  const EMPTY = 0;
  const FIRST = 1;
  const SECOND = 2;
  const SIZE = 8;
  const board = Array.from(new Array(SIZE), () =>
    Array.from(new Array(SIZE), () => EMPTY)
  );
  board[4][3] = board[3][4] = FIRST;
  board[3][3] = board[4][4] = SECOND;
  const D8Y = [1, 0, -1, -1, -1, 0, 1, 1];
  const D8X = [1, 1, 1, 0, -1, -1, -1, 0];

  if (record[0] != null || record[0] !== "") {
    record.some((r, i): boolean => {
      if (i >= nth_played) {
        return true;
      }
      const current_player = i % 2 === 0 ? FIRST : SECOND;
      if (r === "pass") {
        return false;
      }
      const inst = r.split(" ");
      if (inst[0] === "") {
        console.warn("empty inst", inst);
        return false;
      }
      if (inst[0] !== "put") {
        console.error("Unknown inst", inst);
        return false;
      }
      const y = Number.parseInt(inst[1], 10);
      const x = Number.parseInt(inst[2], 10);
      console.log(board)
      board[y][x] = current_player;
      for (let d = 0; d < 8; d++) {
        for (let l = 1; l <= SIZE; l++) {
          const ny = y + D8Y[d] * l;
          const nx = x + D8X[d] * l;
          if (ny < 0 || nx < 0 || ny >= SIZE || nx >= SIZE) {
            break;
          }
          const pos = board[ny][nx];
          if (pos === EMPTY) {
            break;
          }
          if (pos === current_player) {
            for (let v = 1; v < l; v++) {
              const vy = y + D8Y[d] * v;
              const vx = x + D8X[d] * v;
              board[vy][vx] = current_player;
            }
            break;
          }
        }
      }
      return false;
    });
  }
  let [first_pieces, second_pieces] = board.reduce(
    (fs, row) =>
      row.reduce(
        (fs, cell) => [
          fs[0] + (cell === FIRST ? 1 : 0),
          fs[1] + (cell === SECOND ? 1 : 0),
        ],
        fs
      ),
    [0, 0]
  );
  if (first_pieces === 0) {
    second_pieces = SIZE * SIZE;
  }
  if (second_pieces === 0) {
    first_pieces = SIZE * SIZE;
  }
  return (
    <>
      <h3>
        <div className="piece black inline-block"></div>
        {first} : {first_pieces} vs
        <div className="piece white inline-block"></div>
        {second} : {second_pieces}
      </h3>
      <div className="pane">
        <table className="grow0">
          <tbody>
            {board.map((row, i) => {
              return (
                <tr key={i}>
                  {row.map((cell, j) => {
                    return (
                      <td className="cell" key={j}>
                        <div>
                          <div
                            className={
                              "piece " +
                              (cell === FIRST
                                ? "black"
                                : cell === SECOND
                                ? "white"
                                : "empty")
                            }
                          ></div>
                        </div>
                      </td>
                    );
                  })}
                </tr>
              );
            })}
          </tbody>
        </table>
     </div>
    </>
  );
}
export function Reversi(props: {
  first: string;
  second: string;
  record: string;
}) {
  return <>
    <GamePlayer first={props.first} second={props.second} record={props.record} GameView={ReversiView}></GamePlayer>
  </>
}
