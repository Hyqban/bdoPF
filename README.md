# 🚀 bdoPF

FIRST TIME USING: [GUIDE](docs/guide.md)

# ✨ Key Features
- 🌳 Crafting Notes: Conveniently view recipes and quantities for crafting items.
- ⏰ Boss Schedule: Quickly and clearly view a table of all world boss spawn times. Never miss a fight again.
- ✍️ Quick Notes: Jot down important information, tips, or reminders while you're playing. Your personal BDO diary!
- 🌐 Useful Links: Quick access to all the most useful BDO-related websites, saving you from endless searching.
- 🛠 Useful Tools: A collection of handy utilities planned for future implementation. 

# 🌎 Language & Usability
- Supported Languages(NA): The app currently includes built-in support for **English**, **French**, **Germany** and **Spanish**.
- Cross-Server Usability: While this app was developed using data from the NA (North American) server, it is **fully usable for players on other servers**. Please follow the setup guide for a seamless experience.

## ⚠️ Regarding Windows Security Warning
You might see a security warning from Windows, because this app isn't officially certified by Windows or its partners. Rest assured, the app is safe and works offline.

## Issues
1. Always no top is not working on linux.


## Required
### Both windows and linux
 - https://go.dev/
- https://wails.io/docs/gettingstarted/installation
- https://nodejs.org/en/download
- https://angular.dev/installation
```
Go
Wails
Node.js
Angular cli

go install github.com/wailsapp/wails/v2/cmd/wails@latest

wails doctor
```
### Linux

```
// Debian 13
sudo apt update && sudo apt install -y \
  build-essential \
  pkg-config \
  libgtk-3-dev \
  libwebkit2gtk-4.1-dev \
  libnss3 \
  ca-certificates


// wails
echo 'export GOPATH=$HOME/go' >> ~/.bashrc

source ~/.bashrc
```

## Run
### Windows
```
wails dev
```

### linux
```
wails dev -tags webkit2_41
```


## build
### Windows
```
// Automatic compilation，ignore the following if you use this.
wails build

wails build -platform windows/amd64 -o bdoPF_windows_amd64.exe

wails build -platform windows/arm64 -o bdoPF_windows_arm64.exe
```

### Linux
```
// Automatic compilation，ignore the following if you use this.
// windows/amd64
// windows/arm64
// linux/amd64
make all

wails build -tags webkit2_41 -platform windows/amd64 -o bdoPF_windows_amd64.exe

wails build -tags webkit2_41 -platform windows/arm64 -o bdoPF_windows_arm64.exe

wails build -tags webkit2_41 -platform linux/amd64 -o bdoPF_linux_amd64

wails build -tags webkit2_41 -platform linux/arm64 -o bdoPF_linux_arm64
```

## Debug after build
### Windows（Perhaps, I have forgotten）
https://wails.io/docs/reference/cli/
```
wails build -debug
```

### Linux（Lauch terminal）
```
./bdoPF_xxx_xxx
...
output
...
```