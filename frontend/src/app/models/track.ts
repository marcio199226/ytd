export interface Track {
  id: string;
  author: string;
  downloaded: boolean;
  downloadProgress: number;
  duration: number;
  filesize: number;
  name: string;
  playlistId: string;
  status: "pending" | "processing" | "downloading" | "downloaded" | "failed";
  statusError: string;
  thumbnails: string[];
  url: string;
}
