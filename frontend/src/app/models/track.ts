export interface Track {
  id: string;
  author: string;
  downloadProgress: number;
  duration: number;
  filesize: number;
  name: string;
  playlistId: string;
  status: "queued" | "placeholder" | "pending" | "processing" | "downloading" | "downloaded" | "failed";
  statusError: string;
  thumbnails: string[];
  isConvertedToMp3: boolean;
  converting: {
    status: "converting" | "converted" | "queued" | "failed";
    error: string;
    attempts: number;
  }
  url: string;
}
