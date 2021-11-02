import { Entry } from './entry';
import { OfflinePlaylist } from './offline-playlist';

export interface AppConfig {
  BaseSaveDir: string;
  ClipboardWatch: boolean;
  ConcurrentDownloads: boolean;
  ConcurrentPlaylistDownloads: boolean;
  DownloadOnCopy: boolean;
  MaxParrallelDownloads: number;
  ConvertToMp3: boolean;
  CleanWebmFiles: boolean;
	RunInBackgroundOnClose: boolean;
	CheckForUpdates: boolean;
	StartAtLogin: boolean;
  Language: string;
  Telegram: {
    Share: boolean;
    Username: string;
  },
  PublicServer: {
    Enabled: boolean;
    Ngrok: {
      Authtoken: string;
      Auth: {
        Enabled: boolean;
        Username: string;
        Password: string;
      },
    },
    VerifyAppKey: boolean;
    AppKey: string;
  }
}

export interface AppStats {
  DownloadingCount: number;
}

export interface AppState {
  entries: Entry[];
  offlinePlaylists: OfflinePlaylist[];
  config: AppConfig;
  appVersion: string;
  stats: AppStats;
}


export interface BackendCallbacks {
  main: {
    AppState: {
      GetAppConfig: () => Promise<any>;
      SelectDirectory: () => Promise<string>
      IsSupportedUrl: (url: string) => Promise<boolean>;
      AddToDownload: (url: string, isFromClipboard: boolean) => Promise<any>;
      StartDownload: (entry: Entry) => Promise<any>;
      ReadSettingBoolValue: (name: string) => Promise<any>;
      ReadSettingValue: (name: string) => Promise<any>;
      SaveSettingBoolValue: (name: string, val: boolean) => Promise<any>;
      SaveSettingValue: (name: string, val: string) => Promise<any>;
      RemoveEntry: (entry: Entry) => Promise<any>;
      IsFFmpegInstalled: () => Promise<boolean>;
      OpenUrl: (url: string) => Promise<any>;
      ReloadNewLanguage: () => Promise<void>;
      Update: (restart: boolean) => Promise<any>;
      ShowWindow: () => Promise<any>;
      ForceQuit: () => Promise<any>;
    },
    OfflinePlaylistService: {
      CreateNewPlaylist: (name: string) =>  Promise<OfflinePlaylist>;
      RemovePlaylist: (uuid: string) => Promise<boolean>;
      RemoveTrackFromPlaylist: (tid: string, playlist: OfflinePlaylist) => Promise<OfflinePlaylist>;
      AddTrackToPlaylist: (playlists: OfflinePlaylist[]) => Promise<boolean>;
      ExportPlaylist: (uuid: string, path: string) => Promise<string>;
    },
    NgrokService: {
      IsNgrokInstalled: () =>  Promise<boolean>;
    }
  }
}
