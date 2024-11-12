import { twMerge } from "tailwind-merge";

export interface IPlaylistVideoProps {
  thumbnail: string;
  duration: string;

  title: string;
  author: string;

  onDelete: () => void;
}

const DeleteIcon = () => (
    <svg
    width="16"
    height="16"
    viewBox="0 0 16 16"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <g clipPath="url(#clip0_2_113)">
      <path
        d="M15 8C15.0026 6.31557 14.3952 4.68718 13.29 3.416L3.416 13.291C4.43021 14.1698 5.67567 14.7386 7.004 14.9295C8.33234 15.1205 9.68757 14.9255 10.9083 14.368C12.1289 13.8104 13.1636 12.9138 13.8891 11.7847C14.6145 10.6557 15.0001 9.34199 15 8ZM2.71 12.584L12.584 2.709C11.2462 1.55067 9.51928 0.942016 7.75082 1.00552C5.98235 1.06902 4.30356 1.79996 3.05226 3.05126C1.80096 4.30256 1.07002 5.98135 1.00652 7.74982C0.943016 9.51828 1.55167 11.2462 2.71 12.584ZM16 8C16 10.1217 15.1571 12.1566 13.6569 13.6569C12.1566 15.1571 10.1217 16 8 16C5.87827 16 3.84344 15.1571 2.34315 13.6569C0.842855 12.1566 0 10.1217 0 8C0 5.87827 0.842855 3.84344 2.34315 2.34315C3.84344 0.842855 5.87827 0 8 0C10.1217 0 12.1566 0.842855 13.6569 2.34315C15.1571 3.84344 16 5.87827 16 8Z"
        fill="#DC2626"
      />
    </g>
    <defs>
      <clipPath id="clip0_2_113">
        <rect width="16" height="16" fill="white" />
      </clipPath>
    </defs>
  </svg>
)

export const PlayListVideo = ({
  thumbnail,
  duration,
  title,
  author,
  //  orderedBy,
  onDelete,
}: IPlaylistVideoProps) => (
  <div
    className={twMerge(
      "w-full h-full min-h-20",
      "font-semibold text-xs",
      "bg-darkGray/[.06]",
      "relative",
      "flex gap-2 items-center",
      "p-3",
      "rounded-md",
    )}
  >
    <div className={"w-1/5 relative"}>
      <img
        src={thumbnail}
        alt={"Thumbnail for the video " + title}
        className="rounded-md"
      />
      <div
        className={twMerge(
          "absolute bottom-1 right-1",
          "bg-darkGray/40",
          "p-1",
          "rounded-md",
        )}
      >
        <span className={twMerge("text-white", "font-semibold", "font-base")}>
          {duration}
        </span>
      </div>
    </div>
    <div className={"w-full flex flex-col"}>
      <span className={"font-semibold text-darkGray text-xl"}>{title}</span>
      <span className={"font-semibold text-darkGray text-lg"}>{author}</span>
    </div>
    <div
      className={twMerge(
        "absolute right-3 bottom-1",
        "flex gap-1.5 place-items-center",
      )}
    >
      {/*<span className={
          'text-darkGray text-s'
         }>Ordered by {orderedBy}</span>*/}
      <button
        className={twMerge(
          "w-6 h-6",
          "grid place-items-center",
          "bg-[#EEEEEE]",
          "rounded-md",
        )}
        onClick={() => onDelete()}
      >
        <DeleteIcon />
      </button>
    </div>
  </div>
);
