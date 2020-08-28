import React, { useState } from "react";
import "./Input.css";

export function Input(props: React.ComponentProps<any>) {
  return (
    <label>
      <p>{props.children}</p>
      {(() => {
        if (props.errMsg) {
          return <p className="error">{props.errMsg}</p>;
        }
        return <></>;
      })()}
      <input
        onChange={props.onChange}
        value={props.value}
        type={props.type ?? "text"}
      />
    </label>
  );
}

export function useStateValidate(
  init: string,
  format: (text: string) => string,
  validate: (text: string) => string | undefined
) {
  const [res, setRes] = useState(init);
  return {
    value: res,
    err: validate(res),
    change: (e: any) => {
      setRes(format(e.target.value));
    },
  };
}

export const Format = {
  trimSpace(s: string) {
    return s.replace(/\s/, "");
  },
};
