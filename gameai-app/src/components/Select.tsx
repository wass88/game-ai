import React from "react";
import "./Select.css";
const ReactSuperSelect = require("react-super-select");

export function Select(props: React.ComponentProps<any>) {
  return (
    <ReactSuperSelect
      dataSource={props.options}
      onChange={props.onChange || (() => null)}
      searchable={props.searchable || true}
    />
  );
}
