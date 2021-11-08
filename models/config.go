package models

import (
	"encoding/json"
	"fmt"
	"os/user"
	"reflect"
	"strconv"
	"sync"
	_ "unsafe"

	"github.com/wailsapp/wails/v2"
)

//go:noescape
//go:linkname dbReadSettingBoolValue ytd/db.DbReadSettingBoolValue
func dbReadSettingBoolValue(name string) (bool, error)

//go:noescape
//go:linkname dbReadSetting ytd/db.DbReadSetting
func dbReadSetting(Name string) (string, error)

type AppStats struct {
	sync.Mutex
	DownloadingCount uint
}

func (stats *AppStats) IncDndCount() {
	stats.Lock()
	defer stats.Unlock()
	stats.DownloadingCount++
}

func (stats *AppStats) DecDndCount() {
	stats.Lock()
	defer stats.Unlock()
	stats.DownloadingCount--
}

type TelegramConfig struct {
	Share    bool
	Username string
}

type PublicServerConfig struct {
	Enabled      bool
	Ngrok        NgrokConfig
	VerifyAppKey bool
	AppKey       string
}

type NgrokConfig struct {
	Authtoken string
	Auth      NgrokAuthConfig
}

type NgrokAuthConfig struct {
	Enabled  bool
	Username string
	Password string
}

type AppConfig struct {
	runtime                     *wails.Runtime
	ClipboardWatch              bool
	DownloadOnCopy              bool
	ConcurrentDownloads         bool
	ConcurrentPlaylistDownloads bool
	MaxParrallelDownloads       uint
	BaseSaveDir                 string
	ConvertToMp3                bool
	CleanWebmFiles              bool
	RunInBackgroundOnClose      bool
	CheckForUpdates             bool
	StartAtLogin                bool
	Language                    string
	Telegram                    TelegramConfig
	PublicServer                PublicServerConfig
	Proxy                       interface{}
}

func NewAppConfig(watch bool, dldOnCopy bool, cDownloads bool, cPlaylistDownloads bool, mpDownloads uint, baseSaveDir string) AppConfig {
	return AppConfig{
		ClipboardWatch:              watch,
		DownloadOnCopy:              dldOnCopy,
		ConcurrentDownloads:         cDownloads,
		ConcurrentPlaylistDownloads: cPlaylistDownloads,
		MaxParrallelDownloads:       mpDownloads,
		ConvertToMp3:                false,
		CleanWebmFiles:              false,
		BaseSaveDir:                 baseSaveDir,
		RunInBackgroundOnClose:      false,
		CheckForUpdates:             false,
		StartAtLogin:                false,
		Language:                    "en",
		Telegram:                    TelegramConfig{Share: false, Username: ""},
		PublicServer:                PublicServerConfig{Enabled: false},
	}
}

func (cfg *AppConfig) Init() *AppConfig {
	usr, _ := user.Current()
	defaultAppCfg := NewAppConfig(true, true, true, true, 3, fmt.Sprintf("%v/songs", usr.HomeDir))

	cfg = new(AppConfig)
	cfg.ClipboardWatch = getConfigValue(defaultAppCfg, "ClipboardWatch").(bool)
	cfg.DownloadOnCopy = getConfigValue(defaultAppCfg, "DownloadOnCopy").(bool)
	cfg.ConcurrentDownloads = getConfigValue(defaultAppCfg, "ConcurrentDownloads").(bool)
	cfg.ConcurrentPlaylistDownloads = getConfigValue(defaultAppCfg, "ConcurrentPlaylistDownloads").(bool)
	cfg.ConvertToMp3 = getConfigValue(defaultAppCfg, "ConvertToMp3").(bool)
	cfg.CleanWebmFiles = getConfigValue(defaultAppCfg, "CleanWebmFiles").(bool)
	cfg.MaxParrallelDownloads = getConfigValue(defaultAppCfg, "MaxParrallelDownloads").(uint)
	cfg.BaseSaveDir = getConfigValue(defaultAppCfg, "BaseSaveDir").(string)
	cfg.RunInBackgroundOnClose = getConfigValue(defaultAppCfg, "RunInBackgroundOnClose").(bool)
	cfg.CheckForUpdates = getConfigValue(defaultAppCfg, "CheckForUpdates").(bool)
	cfg.StartAtLogin = getConfigValue(defaultAppCfg, "StartAtLogin").(bool)
	cfg.Language = getConfigValue(defaultAppCfg, "Language").(string)
	cfg.Telegram = getConfigValue(defaultAppCfg, "Telegram").(TelegramConfig)
	cfg.PublicServer = getConfigValue(defaultAppCfg, "PublicServer").(PublicServerConfig)

	return cfg
}

func (cfg *AppConfig) SetRuntime(runtime *wails.Runtime) {
	cfg.runtime = runtime
}

func (cfg *AppConfig) Set(name string, val interface{}) error {
	switch name {
	case "BaseSaveDir":
		cfg.SetBaseSaveDir(val)
	case "ClipboardWatch":
		cfg.SetClipboardWatch(val)
	case "DownloadOnCopy":
		cfg.SetDownloadOnCopy(val)
	case "ConcurrentDownloads":
		cfg.SetConcurrentDownloads(val)
	case "ConcurrentPlaylistDownloads":
		cfg.SetConcurrentPlaylistDownloads(val)
	case "MaxParrallelDownloads":
		cfg.SetMaxParrallelDownloads(val)
	case "ConvertToMp3":
		cfg.SetConvertToMp3(val)
	case "CleanWebmFiles":
		cfg.SetCleanWebmFiles(val)
	case "RunInBackgroundOnClose":
		cfg.SetRunInBackgroundOnClose(val)
	case "CheckForUpdates":
		cfg.SetCheckForUpdates(val)
	case "StartAtLogin":
		cfg.SetStartAtLogin(val)
	case "Language":
		cfg.SetLanguage(val)
	case "Telegram":
		cfg.SetTelegram(val)
	case "PublicServer":
		cfg.SetPublicServer(val)
	}
	return nil
}

func (cfg *AppConfig) SetBaseSaveDir(val interface{}) error {
	cfg.BaseSaveDir = val.(string)
	return nil
}

func (cfg *AppConfig) SetClipboardWatch(val interface{}) error {
	cfg.ClipboardWatch = val.(bool)
	return nil
}

func (cfg *AppConfig) SetDownloadOnCopy(val interface{}) error {
	cfg.DownloadOnCopy = val.(bool)
	return nil
}

func (cfg *AppConfig) SetConcurrentDownloads(val interface{}) error {
	cfg.ConcurrentDownloads = val.(bool)
	return nil
}

func (cfg *AppConfig) SetConcurrentPlaylistDownloads(val interface{}) error {
	cfg.ConcurrentPlaylistDownloads = val.(bool)
	return nil
}

func (cfg *AppConfig) SetMaxParrallelDownloads(val interface{}) error {
	v, _ := strconv.ParseUint(val.(string), 10, 8)
	cfg.MaxParrallelDownloads = uint(v)
	return nil
}

func (cfg *AppConfig) SetConvertToMp3(val interface{}) error {
	cfg.ConvertToMp3 = val.(bool)
	return nil
}

func (cfg *AppConfig) SetCleanWebmFiles(val interface{}) error {
	cfg.CleanWebmFiles = val.(bool)
	return nil
}

func (cfg *AppConfig) SetRunInBackgroundOnClose(val interface{}) error {
	cfg.RunInBackgroundOnClose = val.(bool)
	return nil
}

func (cfg *AppConfig) SetCheckForUpdates(val interface{}) error {
	cfg.CheckForUpdates = val.(bool)
	return nil
}

func (cfg *AppConfig) SetStartAtLogin(val interface{}) error {
	cfg.StartAtLogin = val.(bool)
	return nil
}

func (cfg *AppConfig) SetLanguage(val interface{}) error {
	cfg.Language = val.(string)
	return nil
}

func (cfg *AppConfig) SetTelegram(val interface{}) error {
	var telegram TelegramConfig
	json.Unmarshal([]byte(val.(string)), &telegram)
	cfg.Telegram = telegram
	return nil
}

func (cfg *AppConfig) SetPublicServer(val interface{}) error {
	var config PublicServerConfig
	json.Unmarshal([]byte(val.(string)), &config)
	cfg.PublicServer = config

	cfg.runtime.Events.Emit("ngrok:configured")
	return nil
}

func (cfg *AppConfig) IsNgrokApiKeyEnabled() bool {
	return cfg.PublicServer.VerifyAppKey
}

func (cfg *AppConfig) GetNgrokApiKkey() string {
	return cfg.PublicServer.AppKey
}

func getConfigValue(defaultAppCfg AppConfig, name string) interface{} {
	var data interface{}
	t := reflect.ValueOf(defaultAppCfg).FieldByName(name).Kind()
	switch t {
	case reflect.Bool:
		if cfgVal, err := dbReadSettingBoolValue(name); err != nil {
			data = reflect.ValueOf(defaultAppCfg).FieldByName(name).Interface()
		} else {
			data = cfgVal
		}
	case reflect.String:
		if cfgVal, err := dbReadSetting(name); err != nil {
			data = reflect.ValueOf(defaultAppCfg).FieldByName(name).Interface()
		} else {
			data = cfgVal
		}
	case reflect.Uint:
		if cfgVal, err := dbReadSetting(name); err != nil {
			data = reflect.ValueOf(defaultAppCfg).FieldByName(name).Interface()
		} else {
			val, _ := strconv.ParseUint(cfgVal, 10, 8)
			data = uint(val)
		}
	case reflect.Struct:
		if cfgVal, err := dbReadSetting(name); err != nil {
			data = reflect.ValueOf(defaultAppCfg).FieldByName(name).Interface()
		} else {
			// unmarshall from db to the correct struct type
			switch name {
			case "Telegram":
				var telegram TelegramConfig
				json.Unmarshal([]byte(cfgVal), &telegram)
				data = telegram

			case "PublicServer":
				var config PublicServerConfig
				json.Unmarshal([]byte(cfgVal), &config)
				data = config
			}
		}
	}
	return data
}
