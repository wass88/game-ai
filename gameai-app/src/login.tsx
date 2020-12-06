import React from "react";

export const LoginUserContext = React.createContext<any>(null)

export const isVisitor = function(you : any) {
    return you == null || you.id == null || you.authority === "visitor"
}

