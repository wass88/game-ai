import React, { useEffect } from "react";
import { useState } from "react";
import "./Button.css";

export function Button(props: React.ComponentProps<any>) {
  return (
    <button onClick={props.onClick} disabled={props.disabled}>
      {props.children}
    </button>
  );
}
