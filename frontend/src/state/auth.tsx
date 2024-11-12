import { createContext, useContext, useState, useEffect } from "react";

interface IAuthState {
  authenticate(tok: string): void;
  logOut(): void;
  isAuthenticated: boolean;
  getToken(): string | null;
}

const AuthContext = createContext<IAuthState>({
  authenticate: (tok: string) => {
    throw new Error("Unimplemented");
  },
  logOut: () => {
    throw new Error("Unimplemented");
  },
  getToken: () => null,
  isAuthenticated: false,
});

export const useAuthState = () => useContext(AuthContext);

export const AuthState = ({
  isAuthenticated,
  children,
}: {
  isAuthenticated: boolean;
  children: JSX.Element;
}) => {
  const [authenticated, setAuthenticated] = useState<boolean>(isAuthenticated);
  const [loading, setLoading] = useState<boolean>(true);

  const authenticate = (tok: string) => {
    localStorage.setItem("token", tok);
    setAuthenticated(true);
  };

  const logOut = () => {
    console.log("logOut");
    localStorage.removeItem("token");
    setTimeout(() => void setAuthenticated(false), 0);
  };

  const getToken = () => localStorage.getItem("token");

  useEffect(() => {
    if (!authenticated) {
      const token = getToken();
      if (typeof token === "string") authenticate(token);
    }
    setLoading(false);
  }, []);

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated: authenticated,
        authenticate,
        logOut,
        getToken,
      }}
    >
      {!loading && children}
    </AuthContext.Provider>
  );
};
