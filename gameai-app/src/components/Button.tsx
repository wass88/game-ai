import React from "react";
import "./Button.css";

export function Button(props: React.ComponentProps<any>) {
  return (
    <button onClick={props.onClick} disabled={props.disabled || false}>
      {props.children}
    </button>
  );
}
