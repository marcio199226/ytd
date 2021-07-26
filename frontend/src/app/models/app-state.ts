import { Entry } from './entry';

interface AppConfig {
  baseSaveDir: string;
  clipboardWatch: boolean;
  concurrentDownloads: boolean;
  concurrentPlaylistDownloads: boolean;
  downloadOnCopy: boolean;
}

export interface AppState {
  entries: Entry[];
  config: AppConfig;
}
