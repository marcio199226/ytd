export interface Track {
  id: string;
  author: string;
  downloadProgress: number;
  duration: number;
  filesize: number;
  name: string;
  playlistId: string;
  status: "pending" | "processing" | "downloading" | "downloaded" | "failed";
  statusError: string;
  thumbnails: string[];
  isConvertedToMp3: boolean;
  url: string;
}
