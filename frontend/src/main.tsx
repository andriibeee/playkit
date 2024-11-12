import { StrictMode } from "react";
import "./index.css";
import { createRoot } from "react-dom/client";
import { App } from "./App.tsx";
import { AuthState } from "./state/auth";

const token = localStorage.getItem("token");

createRoot(document.getElementById("root")!).render(
  <AuthState
    isAuthenticated={
      token !== null && typeof token === "string" && token === ""
    }
  >
    <StrictMode>
      <App />
    </StrictMode>
  </AuthState>,
);
