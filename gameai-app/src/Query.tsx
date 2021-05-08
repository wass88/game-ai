import {useLocation} from "react-router"

export function useQuery(...keys: string[]) : { [keys: string]: string } {
  const params = new URLSearchParams(useLocation().search);
  return keys.reduce((acc, key) => {return {...acc, [key]: params.get(key)}}, {});
}