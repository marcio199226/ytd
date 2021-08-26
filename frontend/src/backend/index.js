// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

const backend = {
  "main": {
    "AppState": {
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
       * ForceQuit
       * @returns {Promise<void>} 
       */
      "ForceQuit": () => {
        return window.backend.main.AppState.ForceQuit();
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
       * OpenUrl
       * @param {string} arg1 - Go Type: string
       * @returns {Promise<Error>}  - Go Type: error
       */
      "OpenUrl": (arg1) => {
        return window.backend.main.AppState.OpenUrl(arg1);
      },
      /**
       * PreWailsInit
       * @returns {Promise<void>} 
       */
      "PreWailsInit": () => {
        return window.backend.main.AppState.PreWailsInit();
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
    }
  }

};
export default backend;
