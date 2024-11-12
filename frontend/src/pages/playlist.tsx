import axios from "axios";
import { useEffect, useRef } from "react";
import { twMerge } from "tailwind-merge";
import { IVideo, usePlaylist } from "../state/playlist";
import YouTubePlayer from "yt-player";
import { useAuthState } from "../state/auth";
import { PlayListVideo } from "../components/playlist";

function decodeISODurationToTime(isoDuration: string) {
  const [, hours = 0, minutes = 0] =
    isoDuration.match(/T(?:(\d+)H)?(?:(\d+)M)?/) || [];
  return `${String(hours).padStart(2, "0")}:${String(minutes).padStart(2, "0")}`;
}

export const PlaylistView = () => {
  const { updatePlaylist, videos } = usePlaylist();
  const ws = useRef<WebSocket | null>();
  const auth = useAuthState();
  const playerRef = useRef<HTMLDivElement | null>(null);
  const ytRef = useRef<YouTubePlayer | null>(null);
  const currentVideo = useRef<string>("");

  const ensurePlayer = (cb: (pl: YouTubePlayer) => void) => {
    if (!ytRef.current) ytRef.current = new YouTubePlayer(playerRef.current!);
    cb(ytRef.current);
  };

  useEffect(() => {
    ws.current =
      ws.current ||
      new WebSocket(
        "ws://localhost:3000/api/playlist/stream?token=" + auth.getToken(),
      );
    ws.current.onerror = (e) => {
      console.error(e);
    };

    ws.current.onmessage = (e) => {
      if (ytRef.current?.getState() === "playing") {
        ytRef.current?.once("ended", () => {
          updatePlaylist(JSON.parse(e.data) as IVideo[]);
        });
        return;
      }
      updatePlaylist(JSON.parse(e.data) as IVideo[]);
    };

    if (ytRef.current) {
      ytRef.current.on("ended", () => {
        const headers = {
          Authorization: "Bearer " + auth.getToken(),
        };
        axios({
          method: "DELETE",
          url: "/api/playlist/" + currentVideo.current,
          headers,
        })
          .then(() => {})
          .catch((e) => console.error(e));
      });
    }
  }, []);

  useEffect(() => {
    ensurePlayer(() => {
      if (videos.length) {
        ytRef.current!.load(videos[0].id, true);
        currentVideo.current = videos[0].uuid;
      }
    });
  }, [videos]);

  return (
    <div
      className={twMerge(
        "flex  justify-center ",
        "w-full h-screen",
        "bg-[#EEEEEE]",
      )}
    >
      <div ref={playerRef} className="w-1/2 self-center" />
      <div
        className={twMerge("flex flex-col w-1/2", "gap-3", "overflow-scroll")}
      >
        {videos.map((video) => (
          <div key={video.uuid} className="w-full h-fit-content">
            <PlayListVideo
              thumbnail={video.thumbnail}
              title={video.title}
              duration={decodeISODurationToTime(video.duration)}
              author={video.author}
              onDelete={() => {
                axios({
                  method: "DELETE",
                  url: "/api/playlist/" + video.uuid,
                  headers: {
                    Authorization: "Bearer " + auth.getToken(),
                  },
                });
              }}
            />
          </div>
        ))}
      </div>
    </div>
  );
};
