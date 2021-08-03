import { Entry } from './entry';

export interface AppConfig {
  BaseSaveDir: string;
  ClipboardWatch: boolean;
  ConcurrentDownloads: boolean;
  ConcurrentPlaylistDownloads: boolean;
  DownloadOnCopy: boolean;
  MaxParrallelDownloads: number;
}

export interface AppState {
  entries: Entry[];
  config: AppConfig;
}


export interface BackendCallbacks {
  AppState: {
    GetAppConfig: () => Promise<any>;
    SelectDirectory: () => Promise<string>
  }
  isSupportedUrl: (url: string) => Promise<boolean>;
  addToDownload: (url: string, isFromClipboard: boolean) => Promise<any>;
  startDownload: (entry: Entry) => Promise<any>;
  readSettingBoolValue: (name: string) => Promise<any>;
  readSettingValue: (name: string) => Promise<any>;
  saveSettingBoolValue: (name: string, val: boolean) => Promise<any>;
  saveSettingValue: (name: string, val: string) => Promise<any>;
  removeEntry: (entry: Entry) => Promise<any>
}
