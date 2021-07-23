package models

import (
	"fmt"
	"os/user"
	"reflect"
	"strconv"
	_ "unsafe"
)

//go:noescape
//go:linkname dbReadSettingBoolValue ytd/db.DbReadSettingBoolValue
func dbReadSettingBoolValue(name string) (bool, error)

//go:noescape
//go:linkname dbReadSetting ytd/db.DbReadSetting
func dbReadSetting(Name string) (string, error)

type AppConfig struct {
	ClipboardWatch              bool   `json:"clipboardWatch"`
	DownloadOnCopy              bool   `json:"downloadOnCopy"`
	ConcurrentDownloads         bool   `json:"concurrentDownloads"`
	ConcurrentPlaylistDownloads bool   `json:"concurrentPlaylistDownloads"`
	MaxParrallelDownloads       uint   `json:"maxParrallelDownloads"`
	BaseSaveDir                 string `json:"baseSaveDir"`
	Proxy                       interface{}
}

func NewAppConfig(watch bool, dldOnCopy bool, cDownloads bool, cPlaylistDownloads bool, mpDownloads uint, baseSaveDir string) AppConfig {
	return AppConfig{
		ClipboardWatch:              watch,
		DownloadOnCopy:              dldOnCopy,
		ConcurrentDownloads:         cDownloads,
		ConcurrentPlaylistDownloads: cPlaylistDownloads,
		MaxParrallelDownloads:       mpDownloads,
		BaseSaveDir:                 baseSaveDir,
	}
}

func (cfg AppConfig) Init() AppConfig {
	usr, _ := user.Current()
	defaultAppCfg := NewAppConfig(true, true, true, true, 3, fmt.Sprintf("%v/songs", usr.HomeDir))

	/* 	if clipboardWatch, err := dbReadSettingBoolValue("ClipboardWatch"); err != nil {
	   		fmt.Printf("FALLITO ClipboardWatch %s\n", err)
	   		cfg.ClipboardWatch = defaultAppCfg.ClipboardWatch
	   	} else {
	   		cfg.ClipboardWatch = clipboardWatch
	   	}

	   	if downloadOnCopy, err := dbReadSettingBoolValue("DownloadOnCopy"); err != nil {
	   		fmt.Printf("FALLITO DownloadOnCopy %s\n", err)
	   		cfg.DownloadOnCopy = defaultAppCfg.DownloadOnCopy
	   	} else {
	   		cfg.DownloadOnCopy = downloadOnCopy
	   	}

	   	if concurrentDownloads, err := dbReadSettingBoolValue("ConcurrentDownloads"); err != nil {
	   		fmt.Printf("FALLITO ConcurrentDownloads %s\n", err)
	   		cfg.ConcurrentDownloads = defaultAppCfg.ConcurrentDownloads
	   	} else {
	   		cfg.ConcurrentDownloads = concurrentDownloads
	   	}

	   		if maxParrallelDownloads, err := dbReadSetting("MaxParrallelDownloads"); err != nil {
	   	   		fmt.Printf("FALLITO MaxParrallelDownloads %s \n", err)
	   	   		cfg.MaxParrallelDownloads = defaultAppCfg.MaxParrallelDownloads
	   	   	} else {
	   	   		val, _ := strconv.ParseUint(maxParrallelDownloads, 10, 8)
	   	   		cfg.MaxParrallelDownloads = uint(val)
	   	   	} */

	cfg.ClipboardWatch = getConfigValue(defaultAppCfg, "ClipboardWatch").(bool)
	cfg.DownloadOnCopy = getConfigValue(defaultAppCfg, "DownloadOnCopy").(bool)
	cfg.ConcurrentDownloads = getConfigValue(defaultAppCfg, "ConcurrentDownloads").(bool)
	cfg.MaxParrallelDownloads = getConfigValue(defaultAppCfg, "MaxParrallelDownloads").(uint)
	cfg.BaseSaveDir = getConfigValue(defaultAppCfg, "BaseSaveDir").(string)

	return cfg
}

func getConfigValue(defaultAppCfg AppConfig, name string) interface{} {
	var data interface{}
	fmt.Printf("%s is of type %s\n", name, reflect.ValueOf(defaultAppCfg).FieldByName(name).Type())
	t := reflect.ValueOf(defaultAppCfg).FieldByName(name).Kind()
	fmt.Println(t)
	switch t {
	case reflect.Bool:
		if cfgVal, err := dbReadSettingBoolValue(name); err != nil {
			fmt.Printf("FALLITO %s %s \n", name, err)
			data = reflect.ValueOf(defaultAppCfg).FieldByName(name).Interface()
		} else {
			fmt.Printf("TROVATO %s %t \n", name, cfgVal)
			data = cfgVal
		}
	case reflect.String:
		if cfgVal, err := dbReadSetting(name); err != nil {
			fmt.Printf("FALLITO %s %s \n", name, err)
			data = reflect.ValueOf(defaultAppCfg).FieldByName(name).Interface()
		} else {
			data = cfgVal
		}
	case reflect.Uint:
		if cfgVal, err := dbReadSetting(name); err != nil {
			fmt.Printf("FALLITO %s %s \n", name, err)
			data = reflect.ValueOf(defaultAppCfg).FieldByName(name).Interface()
		} else {
			fmt.Printf("TROVATO %s %s \n", name, cfgVal)
			val, _ := strconv.ParseUint(cfgVal, 10, 8)
			data = uint(val)
		}
	}
	return data
}
