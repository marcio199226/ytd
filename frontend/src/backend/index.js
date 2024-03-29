// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

const backend = {
  "main": {
    "AppState": {
      /**
       * AddToConvertQueue
       * @param {any} arg1 - Go Type: models.GenericEntry
       * @returns {Promise<Error>}  - Go Type: error
       */
      "AddToConvertQueue": (arg1) => {
        return window.backend.main.AppState.AddToConvertQueue(arg1);
      },
      /**
       * AddToDownload
       * @param {string} arg1 - Go Type: string
       * @param {boolean} arg2 - Go Type: bool
       * @returns {Promise<Error>}  - Go Type: error
       */
      "AddToDownload": (arg1, arg2) => {
        return window.backend.main.AppState.AddToDownload(arg1, arg2);
      },
      /**
       * Deadline
       * @returns {Promise<any|boolean>}  - Go Type: time.Time
       */
      "Deadline": () => {
        return window.backend.main.AppState.Deadline();
      },
      /**
       * Done
       * @returns {Promise<any>}  - Go Type: <-chan struct {}
       */
      "Done": () => {
        return window.backend.main.AppState.Done();
      },
      /**
       * Err
       * @returns {Promise<Error>}  - Go Type: error
       */
      "Err": () => {
        return window.backend.main.AppState.Err();
      },
      /**
       * ForceQuit
       * @returns {Promise<void>} 
       */
      "ForceQuit": () => {
        return window.backend.main.AppState.ForceQuit();
      },
      /**
       * GetAll
       * @returns {Promise<any>}  - Go Type: *main.AppState
       */
      "GetAll": () => {
        return window.backend.main.AppState.GetAll();
      },
      /**
       * GetAppConfig
       * @returns {Promise<any>}  - Go Type: *models.AppConfig
       */
      "GetAppConfig": () => {
        return window.backend.main.AppState.GetAppConfig();
      },
      /**
       * GetEntryById
       * @param {any} arg1 - Go Type: models.GenericEntry
       * @returns {Promise<any>}  - Go Type: *models.GenericEntry
       */
      "GetEntryById": (arg1) => {
        return window.backend.main.AppState.GetEntryById(arg1);
      },
      /**
       * InitializeListeners
       * @returns {Promise<void>} 
       */
      "InitializeListeners": () => {
        return window.backend.main.AppState.InitializeListeners();
      },
      /**
       * IsFFmpegInstalled
       * @returns {Promise<string|Error>}  - Go Type: string
       */
      "IsFFmpegInstalled": () => {
        return window.backend.main.AppState.IsFFmpegInstalled();
      },
      /**
       * IsSupportedUrl
       * @param {string} arg1 - Go Type: string
       * @returns {Promise<boolean>}  - Go Type: bool
       */
      "IsSupportedUrl": (arg1) => {
        return window.backend.main.AppState.IsSupportedUrl(arg1);
      },
      /**
       * Lock
       * @returns {Promise<void>} 
       */
      "Lock": () => {
        return window.backend.main.AppState.Lock();
      },
      /**
       * OpenUrl
       * @param {string} arg1 - Go Type: string
       * @returns {Promise<Error>}  - Go Type: error
       */
      "OpenUrl": (arg1) => {
        return window.backend.main.AppState.OpenUrl(arg1);
      },
      /**
       * PreWailsInit
       * @param {any} arg1 - Go Type: context.Context
       * @returns {Promise<void>} 
       */
      "PreWailsInit": (arg1) => {
        return window.backend.main.AppState.PreWailsInit(arg1);
      },
      /**
       * ReadSettingBoolValue
       * @param {string} arg1 - Go Type: string
       * @returns {Promise<boolean|Error>}  - Go Type: bool
       */
      "ReadSettingBoolValue": (arg1) => {
        return window.backend.main.AppState.ReadSettingBoolValue(arg1);
      },
      /**
       * ReadSettingValue
       * @param {string} arg1 - Go Type: string
       * @returns {Promise<string|Error>}  - Go Type: string
       */
      "ReadSettingValue": (arg1) => {
        return window.backend.main.AppState.ReadSettingValue(arg1);
      },
      /**
       * ReloadNewLanguage
       * @returns {Promise<void>} 
       */
      "ReloadNewLanguage": () => {
        return window.backend.main.AppState.ReloadNewLanguage();
      },
      /**
       * RemoveEntry
       * @param {any} arg1 - Go Type: map[string]interface {}
       * @returns {Promise<Error>}  - Go Type: error
       */
      "RemoveEntry": (arg1) => {
        return window.backend.main.AppState.RemoveEntry(arg1);
      },
      /**
       * SaveSettingBoolValue
       * @param {string} arg1 - Go Type: string
       * @param {boolean} arg2 - Go Type: bool
       * @returns {Promise<Error>}  - Go Type: error
       */
      "SaveSettingBoolValue": (arg1, arg2) => {
        return window.backend.main.AppState.SaveSettingBoolValue(arg1, arg2);
      },
      /**
       * SaveSettingValue
       * @param {string} arg1 - Go Type: string
       * @param {string} arg2 - Go Type: string
       * @returns {Promise<Error>}  - Go Type: error
       */
      "SaveSettingValue": (arg1, arg2) => {
        return window.backend.main.AppState.SaveSettingValue(arg1, arg2);
      },
      /**
       * SelectDirectory
       * @returns {Promise<string|Error>}  - Go Type: string
       */
      "SelectDirectory": () => {
        return window.backend.main.AppState.SelectDirectory();
      },
      /**
       * ShowWindow
       * @returns {Promise<void>} 
       */
      "ShowWindow": () => {
        return window.backend.main.AppState.ShowWindow();
      },
      /**
       * StartDownload
       * @param {any} arg1 - Go Type: map[string]interface {}
       * @returns {Promise<Error>}  - Go Type: error
       */
      "StartDownload": (arg1) => {
        return window.backend.main.AppState.StartDownload(arg1);
      },
      /**
       * Unlock
       * @returns {Promise<void>} 
       */
      "Unlock": () => {
        return window.backend.main.AppState.Unlock();
      },
      /**
       * Update
       * @param {boolean} arg1 - Go Type: bool
       * @returns {Promise<void>} 
       */
      "Update": (arg1) => {
        return window.backend.main.AppState.Update(arg1);
      },
      /**
       * Value
       * @param {number} arg1 - Go Type: interface {}
       * @returns {Promise<number>}  - Go Type: interface {}
       */
      "Value": (arg1) => {
        return window.backend.main.AppState.Value(arg1);
      },
    }
    "NgrokService": {
      /**
       * Deadline
       * @returns {Promise<any|boolean>}  - Go Type: time.Time
       */
      "Deadline": () => {
        return window.backend.main.NgrokService.Deadline();
      },
      /**
       * Done
       * @returns {Promise<any>}  - Go Type: <-chan struct {}
       */
      "Done": () => {
        return window.backend.main.NgrokService.Done();
      },
      /**
       * Err
       * @returns {Promise<Error>}  - Go Type: error
       */
      "Err": () => {
        return window.backend.main.NgrokService.Err();
      },
      /**
       * ExitCode
       * @returns {Promise<number>}  - Go Type: int
       */
      "ExitCode": () => {
        return window.backend.main.NgrokService.ExitCode();
      },
      /**
       * Exited
       * @returns {Promise<boolean>}  - Go Type: bool
       */
      "Exited": () => {
        return window.backend.main.NgrokService.Exited();
      },
      /**
       * GetPublicUrl
       * @param {any} arg1 - Go Type: chan main.NgrokTunnelInfo
       * @returns {Promise<void>} 
       */
      "GetPublicUrl": (arg1) => {
        return window.backend.main.NgrokService.GetPublicUrl(arg1);
      },
      /**
       * IsNgrokInstalled
       * @returns {Promise<string|Error>}  - Go Type: string
       */
      "IsNgrokInstalled": () => {
        return window.backend.main.NgrokService.IsNgrokInstalled();
      },
      /**
       * KillProcess
       * @returns {Promise<Error>}  - Go Type: error
       */
      "KillProcess": () => {
        return window.backend.main.NgrokService.KillProcess();
      },
      /**
       * MonitorNgrokProcess
       * @returns {Promise<void>} 
       */
      "MonitorNgrokProcess": () => {
        return window.backend.main.NgrokService.MonitorNgrokProcess();
      },
      /**
       * Pid
       * @returns {Promise<number>}  - Go Type: int
       */
      "Pid": () => {
        return window.backend.main.NgrokService.Pid();
      },
      /**
       * SetAuthToken
       * @returns {Promise<Error>}  - Go Type: error
       */
      "SetAuthToken": () => {
        return window.backend.main.NgrokService.SetAuthToken();
      },
      /**
       * StartProcess
       * @param {boolean} arg1 - Go Type: bool
       * @returns {Promise<any>}  - Go Type: main.NgrokProcessResult
       */
      "StartProcess": (arg1) => {
        return window.backend.main.NgrokService.StartProcess(arg1);
      },
      /**
       * String
       * @returns {Promise<string>}  - Go Type: string
       */
      "String": () => {
        return window.backend.main.NgrokService.String();
      },
      /**
       * Success
       * @returns {Promise<boolean>}  - Go Type: bool
       */
      "Success": () => {
        return window.backend.main.NgrokService.Success();
      },
      /**
       * Sys
       * @returns {Promise<number>}  - Go Type: interface {}
       */
      "Sys": () => {
        return window.backend.main.NgrokService.Sys();
      },
      /**
       * SysUsage
       * @returns {Promise<number>}  - Go Type: interface {}
       */
      "SysUsage": () => {
        return window.backend.main.NgrokService.SysUsage();
      },
      /**
       * SystemTime
       * @returns {Promise<any>}  - Go Type: time.Duration
       */
      "SystemTime": () => {
        return window.backend.main.NgrokService.SystemTime();
      },
      /**
       * UserTime
       * @returns {Promise<any>}  - Go Type: time.Duration
       */
      "UserTime": () => {
        return window.backend.main.NgrokService.UserTime();
      },
      /**
       * Value
       * @param {number} arg1 - Go Type: interface {}
       * @returns {Promise<number>}  - Go Type: interface {}
       */
      "Value": (arg1) => {
        return window.backend.main.NgrokService.Value(arg1);
      },
    }
  }

  "offline": {
    "OfflinePlaylistService": {
      /**
       * AddTrackToPlaylist
       * @param {Array.<any>} arg1 - Go Type: []map[string]interface {}
       * @returns {Promise<boolean|Error>}  - Go Type: bool
       */
      "AddTrackToPlaylist": (arg1) => {
        return window.backend.offline.OfflinePlaylistService.AddTrackToPlaylist(arg1);
      },
      /**
       * CreateNewPlaylist
       * @param {string} arg1 - Go Type: string
       * @returns {Promise<any|Error>}  - Go Type: models.OfflinePlaylist
       */
      "CreateNewPlaylist": (arg1) => {
        return window.backend.offline.OfflinePlaylistService.CreateNewPlaylist(arg1);
      },
      /**
       * CreateNewPlaylistWithTracks
       * @param {string} arg1 - Go Type: string
       * @param {Array.<string>} arg2 - Go Type: []string
       * @returns {Promise<any|Error>}  - Go Type: models.OfflinePlaylist
       */
      "CreateNewPlaylistWithTracks": (arg1, arg2) => {
        return window.backend.offline.OfflinePlaylistService.CreateNewPlaylistWithTracks(arg1, arg2);
      },
      /**
       * ExportPlaylist
       * @param {string} arg1 - Go Type: string
       * @param {string} arg2 - Go Type: string
       * @returns {Promise<boolean|Error>}  - Go Type: bool
       */
      "ExportPlaylist": (arg1, arg2) => {
        return window.backend.offline.OfflinePlaylistService.ExportPlaylist(arg1, arg2);
      },
      /**
       * GetPlaylistByUUID
       * @param {string} arg1 - Go Type: string
       * @returns {Promise<any>}  - Go Type: models.OfflinePlaylist
       */
      "GetPlaylistByUUID": (arg1) => {
        return window.backend.offline.OfflinePlaylistService.GetPlaylistByUUID(arg1);
      },
      /**
       * GetPlaylists
       * @param {boolean} arg1 - Go Type: bool
       * @returns {Promise<Array.<any>|Error>}  - Go Type: []models.OfflinePlaylist
       */
      "GetPlaylists": (arg1) => {
        return window.backend.offline.OfflinePlaylistService.GetPlaylists(arg1);
      },
      /**
       * RemovePlaylist
       * @param {string} arg1 - Go Type: string
       * @returns {Promise<boolean|Error>}  - Go Type: bool
       */
      "RemovePlaylist": (arg1) => {
        return window.backend.offline.OfflinePlaylistService.RemovePlaylist(arg1);
      },
      /**
       * RemoveTrackFromPlaylist
       * @param {string} arg1 - Go Type: string
       * @param {any} arg2 - Go Type: models.OfflinePlaylist
       * @returns {Promise<any|Error>}  - Go Type: models.OfflinePlaylist
       */
      "RemoveTrackFromPlaylist": (arg1, arg2) => {
        return window.backend.offline.OfflinePlaylistService.RemoveTrackFromPlaylist(arg1, arg2);
      },
      /**
       * SetConfig
       * @param {any} arg1 - Go Type: models.AppConfig
       * @returns {Promise<void>} 
       */
      "SetConfig": (arg1) => {
        return window.backend.offline.OfflinePlaylistService.SetConfig(arg1);
      },
    }
  }

};
export default backend;
