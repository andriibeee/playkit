import { createContext, useContext, useEffect, useRef, useState } from "react";
import axios, { AxiosRequestConfig, AxiosResponse } from "axios";
import { useAuthState } from "./auth";

export interface IVideo {
  uuid: string;
  id: string;
  title: string;
  author: string;
  thumbnail: string;
  duration: string;
}

interface IPlaylistState {
  videos: Array<IVideo>;
  updatePlaylist(a: Array<IVideo>): void;
}

const PlaylistContext = createContext<IPlaylistState>({
  videos: [],
  updatePlaylist(a: Array<IVideo>) {
    throw new Error("Not implemented");
  },
});

export const usePlaylist = () => useContext(PlaylistContext);

export const PlaylistState = ({ children }: { children: JSX.Element }) => {
  const auth = useAuthState();
  const [videos, setVideos] = useState<Array<IVideo>>([]);
  const updatePlaylist = (a: Array<IVideo>) => {
    setVideos(a);
  };
  const fetchPlaylist = async () => {
    if (!auth.isAuthenticated) return;
    let opts: AxiosRequestConfig = {
      headers: auth.isAuthenticated
        ? {
            Authorization: "Bearer " + auth.getToken(),
          }
        : {},
    };

    try {
      let req = await axios.get<IVideo[]>("/api/playlist", opts);
      updatePlaylist(req.data);
    } catch (e) {
      auth.logOut();
    }
  };
  useEffect(() => void fetchPlaylist(), []);
  return (
    <PlaylistContext.Provider
      value={{
        videos,
        updatePlaylist,
      }}
    >
      {children}
    </PlaylistContext.Provider>
  );
};
