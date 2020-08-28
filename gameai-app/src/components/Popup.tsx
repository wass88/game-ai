import React from "react";
import "./Popup.css";

export function Popup(props: React.ComponentProps<any>) {
  if (!props.show) {
    return <></>;
  }
  return (
    <div
      className="popup"
      onClick={() => props.setShow(false)}
      style={{ display: props.show ? "flex" : "none" }}
    >
      <div onClick={(e) => e.stopPropagation()}>{props.children}</div>
    </div>
  );
}
