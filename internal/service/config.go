package service

import (
	"bdoPF/internal/model"
	"encoding/json"
	"fmt"
	"os"
)

var CONFIG_PATH = "config.json"
var DEFAULT_VERSION = "1.0.1"

// var DEFAULT_VERSION = "0.0.9"

type Config struct {
	DI         *DIContainer    `json:"-"`
	AppName    string          `json:"appName"`
	Version    string          `json:"version"`
	Window     Window          `json:"window"`
	Theme      string          `json:"theme"`
	Locale     string          `json:"locale"`
	NewVersion model.LatestApp `json:"newVersion"`
}

func NewConfig(di *DIContainer) *Config {
	cf := loadAndValidateConfig(di)

	if cf == nil {
		panic(fmt.Sprintln("Fatal: failed to load configuration."))
	}

	return cf
}

func (cf *Config) ReadConfig() Config {
	return *cf
}

func getDefaultConfig() *Config {
	return &Config{
		AppName: "bdoPF",
		Version: DEFAULT_VERSION,
		Theme:   "lightskyblue",
		Locale:  "en",
		NewVersion: model.LatestApp{
			Version:     "",
			Download:    false,
			DownloadUrl: "",
		},
		Window: Window{
			OnTop:               false,
			Width:               420,
			Height:              668,
			MaxWidth:            1920,
			MaxHeight:           1080,
			MinWidth:            420,
			MinHeight:           560,
			IsFullScreen:        false,
			IsWidgetMode:        false,
			DefaultWidgetWidth:  200,
			DefaultWidgetHeight: 100,
			WidgetWidth:         200,
			WidgetHeight:        100,
		},
	}
}

func (cf *Config) enforceSystemFields() {
	cf.Version = DEFAULT_VERSION
	if cf.Version == cf.NewVersion.Version {
		cf.NewVersion.DownloadUrl = ""
		cf.NewVersion.Version = ""
		cf.NewVersion.Download = false
	}

}

func writeConfigToFile(cfg *Config, filePath string) error {
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file to %s: %w", filePath, err)
	}
	return nil
}

func (cfg *Config) ReceiveConfigUpdate(raw map[string]any) {

	b, err := json.Marshal(raw)

	if err != nil {
		fmt.Println(err)
		return
	}

	var config Config
	err = json.Unmarshal(b, &config)

	if err != nil {
		fmt.Println(err)
		return
	}

	writeConfigToFile(&config, CONFIG_PATH)
}

func loadAndValidateConfig(di *DIContainer) *Config {
	cf := getDefaultConfig()
	cf.DI = di

	content, err := os.ReadFile(CONFIG_PATH)

	if os.IsNotExist(err) {
		fmt.Printf("Config file does not exist, creating default: %s\n", CONFIG_PATH)

		if err := cf.SaveConfig(); err != nil {
			return nil
		}
		return cf
	}

	if err != nil {
		return nil
	}

	err = json.Unmarshal(content, cf)
	if err != nil {
		return nil
	}

	cf.enforceSystemFields()
	err = cf.SaveConfig()
	if err != nil {
		fmt.Printf("Warning: failed to update config file with missing fields: %v\n", err)
	}

	return cf
}

func (cf *Config) SaveConfig() error {
	fh := Resolve[FileHandler](cf.DI, "fileHandler")

	content, err := json.MarshalIndent(cf, "", "	")

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	fpath := fh.PathJoin(CONFIG_PATH)
	if err := os.WriteFile(fpath, content, 0644); err != nil {
		return fmt.Errorf("failed to write config file to %s: %w", fpath, err)
	}
	return nil
}

func (cfg *Config) StartupPrepare(envPath string) {
	// locales
	// dynamic_string.json
	fh := cfg.DI.GetFileHandler()

	// dynamic_strings.json
	// tpath := fh.PathJoin(envPath, "dynamic_strings.json")
	somethingPrepare(fh, DynamicStringsMap, envPath, "dynamic_strings.json")

	// de.json, en.json, fr.json, sp.json
	os.MkdirAll(fh.PathJoin(cfg.DI.ResourcePath.AssetsPath, cfg.DI.ResourcePath.Locale), os.ModePerm)

	for _, el := range LocalesMap {
		temppath, _ := el["locale"].(string)
		somethingPrepare(fh, el, envPath, cfg.DI.ResourcePath.Locale, temppath+".json")
	}

}

func somethingPrepare(fh *FileHandler, tmap map[string]any, fPath ...string) {
	tpath := fh.PathJoin(fPath...)
	data, err := os.ReadFile(tpath)
	if os.IsNotExist(err) {
		// byt, _ := json.MarshalIndent(tmap, "", " ")
		// err = os.WriteFile(tpath, byt, 0644)
		f, err := os.Create(tpath)
		if err != nil {
			fmt.Println("write err:", err)
			return
		}
		defer f.Close()

		enc := json.NewEncoder(f)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "    ")

		if err := enc.Encode(tmap); err != nil {
			fmt.Println("encode err:", err)
		}
		return
	}

	var dsm map[string]any
	err = json.Unmarshal(data, &dsm)
	if err != nil {
		fmt.Println("line 197: ", err)
		return
	}

	changed := mergeDefault(dsm, tmap)
	if changed {
		fmt.Println("changed.")
		f, err := os.Create(tpath)
		if err != nil {
			fmt.Println("write err:", err)
			return
		}
		defer f.Close()

		enc := json.NewEncoder(f)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "    ")

		if err := enc.Encode(dsm); err != nil {
			fmt.Println("encode err:", err)
		}
	}
}

func mergeDefault(dst, def map[string]any) (changed bool) {
	for k, defVal := range def {
		dstVal, exists := dst[k]

		if !exists {
			dst[k] = defVal
			changed = true
			continue
		}

		dstMap, ok1 := dstVal.(map[string]any)
		defMap, ok2 := defVal.(map[string]any)

		if ok1 && ok2 {
			if mergeDefault(dstMap, defMap) {
				changed = true
			}
		}
	}
	return
}
