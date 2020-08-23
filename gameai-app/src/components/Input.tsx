import React, { useEffect } from "react";
import { useState } from "react";
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
