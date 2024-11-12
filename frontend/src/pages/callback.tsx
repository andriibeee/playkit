import { useLocation, useSearch } from "wouter";
import {  useAuthState } from "../state/auth";
import { useEffect } from "react";
import axios from "axios";

interface ICallbackResponse {
  token: string;
}

export const CallbackPage = () => {
  const searchString = useSearch();
  const auth = useAuthState();
  const [, setLocation] = useLocation();
  useEffect(() => {
    if (searchString !== "") {
      axios
        .get<ICallbackResponse>("/api/auth/callback", {
          params: Object.fromEntries(
            Array.from(new URLSearchParams(searchString).entries()),
          ),
        })
        .then((r) => {
          if (r.status === 200) {
            auth.authenticate(r.data.token);
            localStorage.setItem("token", r.data.token);
            setLocation("/");
          } else {
            throw new Error(r.statusText);
          }
        });
    }
  }, [searchString]);
  return null;
};
