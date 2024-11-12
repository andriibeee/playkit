import { Logo } from "../components/logo";
import { twMerge } from "tailwind-merge";

import mod from "../assets/mod.png";
import flicker from "../assets/flicker.jpg";
import { PlayListVideo } from "../components/playlist";

const TwitchChatMessage = () => (
  <div
    className={twMerge(
      "w-full h-full",
      "bg-[#18181B]",
      "font-semibold text-xs",
      "p-3",
      "flex items-center gap-1",
      "rounded-md",
    )}
  >
    <span className="text-[#ADADB8]">04:20</span>
    <img src={mod} width={16} height={16} />
    <span>
      <span className="text-[#5F9EA0]">pu$$y3ater</span>
      <span className="text-white">
        : !play https://www.youtube.com/watch?v=D1sZ_vwqwcE
      </span>
    </span>
  </div>
);

export function AuthPage() {
  return (
    <div
      className={twMerge(
        "flex flex-col justify-center items-center",
        "w-full h-screen",
        "bg-[#EEEEEE]",
      )}
    >
      <Logo />
      <span className={"font-semibold text-xl text-[#666666]"}>
        Allow your viewers to populate your on-stream playlist!
      </span>

      <div
        className={twMerge(
          "w-3/5 max-h-screen",
          "bg-[#D9D9D9]",
          "rounded-md",
          "mt-8 p-3",
          "flex items-center flex-col",
        )}
      >
        <span className={twMerge("font-semibold text-s text-[#666666]", "p-3")}>
          Your twitch chat:
        </span>
        <div className="min-w-lg max-h-8">
          <TwitchChatMessage />
        </div>
        <span
          className={twMerge(
            "font-semibold",
            "text-s",
            "p-3",
            "text-[#666666]",
          )}
        >
          One moment later your playlist:
        </span>
        <PlayListVideo
          thumbnail={flicker}
          duration={"4:20"}
          title={"Flicker"}
          author={"Porter Robinson"}
          onDelete={() => {}}
        />
      </div>
      <form className="p-3" method="GET" action="/api/auth/login">
        <button
          type="submit"
          className={
            "text-white bg-gray-800 hover:bg-gray-900 focus:outline-none focus:ring-4 focus:ring-gray-300 font-medium rounded-full text-sm px-5 py-2.5 me-2 mb-2 dark:bg-gray-800 dark:hover:bg-gray-700 dark:focus:ring-gray-700 dark:border-gray-700"
          }
        >
          Log in via Twitch
        </button>
      </form>
    </div>
  );
}
