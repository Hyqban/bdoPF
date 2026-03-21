package service

import (
	"bdoPF/internal/model"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Updater struct {
	DI         *DIContainer
	system     string
	systemType string
	// appVersion    string
	latestVersion string
	downloadUrl   string
	changeLog     string
	CurrentExe    string
	newExe        string
}

func NewUpdater(di *DIContainer) *Updater {
	up := Updater{}
	up.DI = di
	up.checkSystemType()

	return &up
}

func (up *Updater) checkSystemType() {
	up.system = runtime.GOOS
	up.systemType = runtime.GOARCH
}

func (up *Updater) concatAppName(appName string) string {
	tempName := ""

	switch up.system {
	case "windows":
		tempName = appName + "_" + up.system + "_" + up.systemType + ".exe"
	case "linux":
		tempName = appName + "_" + up.system + "_" + up.systemType
	}
	return tempName
}

func (up *Updater) HasLatestVersion(newVersion, appVersion string) {
	newv := strings.Split(newVersion, ".")
	appv := strings.Split(appVersion, ".")

	for i := 0; i < 3; i++ {
		nv, _ := strconv.Atoi(newv[i])
		av, _ := strconv.Atoi(appv[i])

		if av > nv {
			return
		}

		if nv > av {
			up.latestVersion = newVersion
			return
		}
	}
}

func (up *Updater) parseResponse(body map[string]any, resp *model.ResponseMsg) {
	// v1.0.1
	tag_name, ok := body["tag_name"].(string)

	if !ok {
		return
	}

	// up.HasLatestVersion(tag_name[1:], up.appVersion)
	up.HasLatestVersion(tag_name[1:], DEFAULT_VERSION)

	if up.latestVersion == "" {
		resp.Code = "100"
		resp.Msg = "This version is up to date."
		return
	}

	assets, _ := body["assets"].([]any)

	for _, v := range assets {
		info := v.(map[string]any)
		// [bdoPF windows amd64.ex]
		exeSlice := strings.Split(info["name"].(string), "_")

		if exeSlice[1] == up.system {
			switch up.system {
			case "windows":
				// bdoPF_windows_amd64 or bdoPF_windows_arm64
				// [amd64 exe]
				if strings.Split(exeSlice[2], ".")[0] == up.systemType {
					up.downloadUrl = info["browser_download_url"].(string)
					break
				}
			case "linux":
				//: bdoPF_linux_amd64 or bdoPF_linux_arm64
				// amd64
				if exeSlice[2] == up.systemType {
					up.downloadUrl = info["browser_download_url"].(string)
					break
				}
			}
		}
	}
}

func (up *Updater) AppCheckForUpdates() model.ResponseMsg {
	// {
	// 	"message":"Not Found",
	// 	"documentation_url":"https://docs.github.com/rest/releases/releases#get-the-latest-release",
	// 	"status":"404"
	// }
	url := "https://api.github.com/repos/bahyqn/bdoPF/releases/latest"

	header := map[string]string{
		"X-GitHub-Api-Version": "2022-11-28",
		"Accept":               "application/vnd.github+json",
	}
	var responseMsg model.ResponseMsg
	responseBody, ok := NewRequest(&responseMsg, "GET", url, header)

	if !ok {
		return responseMsg
	}

	up.parseResponse(responseBody, &responseMsg)

	if up.downloadUrl != "" && up.latestVersion != "" {
		cf := Resolve[Config](up.DI, "config")
		cf.NewVersion.Version = up.latestVersion
		cf.NewVersion.DownloadUrl = up.downloadUrl

		_ = cf.SaveConfig()

		responseMsg.Code = "200"
		responseMsg.Msg = "New version available."
		responseMsg.Data = map[string]string{
			"downloadUrl": cf.NewVersion.DownloadUrl,
			"newVersion":  cf.NewVersion.Version,
		}
	}
	return responseMsg
}

func (up *Updater) DownloadUpdates() model.ResponseMsg {
	var responseMsg model.ResponseMsg

	cf := Resolve[Config](up.DI, "config")
	// up.CurrentExe = cf.AppName + "_" + up.systemType + ".exe"
	up.CurrentExe = up.concatAppName(cf.AppName)
	latestAppPath := filepath.Join("tmp", up.CurrentExe)

	header := map[string]string{
		"X-GitHub-Api-Version": "2022-11-28",
		"Accept-Encoding":      "gzip, deflate",
	}
	bodyBytes, ok := NewRequestForDownload(&responseMsg, "GET", cf.NewVersion.DownloadUrl, header)

	// fmt.Println("ok: ", ok)
	if !ok {
		return responseMsg
	}

	dir := filepath.Dir(latestAppPath)
	fmt.Println("latestAppPath: ", latestAppPath)
	fmt.Println("dir: ", dir)

	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			responseMsg.Code = "100"
			responseMsg.Msg = fmt.Sprintf("Failed to create directory for download: %v", err)
		}
	}

	out, err := os.Create(latestAppPath)
	if err != nil {
		responseMsg.Code = "100"
		responseMsg.Msg = fmt.Sprintf("Failed to create file for download: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, bytes.NewReader(bodyBytes)); err != nil {
		responseMsg.Code = "100"
		responseMsg.Msg = fmt.Sprintln("Download failed.")
	}

	cf.NewVersion.Download = true
	_ = cf.SaveConfig()

	responseMsg.Code = "200"
	responseMsg.Msg = "Download completed successfully."
	return responseMsg
}

func (up *Updater) StartUpdate() {
	StartUpdateApp(up)
}
