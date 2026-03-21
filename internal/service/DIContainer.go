package service

import (
	model "bdoPF/internal/model"
	"context"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DIContainer struct {
	AppCtx         *context.Context
	Addr           string
	Env            string
	locale         string
	ResourcePath   model.ResourcePath
	independencies map[string]interface{}
}

func NewDiContainer() *DIContainer {
	return &DIContainer{independencies: make(map[string]interface{})}
}

func (di *DIContainer) Register(name string, instance interface{}) {
	di.independencies[name] = instance
}

func Resolve[T any](di *DIContainer, name string) *T {
	instance, ok := di.independencies[name]

	if !ok {
		var zero T
		return &zero
	}
	return instance.(*T)
}

func (di *DIContainer) SetLocale(locale string) {
	di.locale = locale
}

func (di *DIContainer) GetLocale() string {
	return di.locale
}

func (di *DIContainer) SetAddr(addr string) {
	di.Addr = addr
}

func (di *DIContainer) GetAddr() string {
	return di.Addr
}

func (di *DIContainer) GetImgPath() map[string]string {
	fh := di.GetFileHandler()

	if fh != nil {
	}

	if di.Env == "dev" {
		return map[string]string{
			"env":  di.Env,
			"icon": di.ResourcePath.Icon,
			"png":  di.ResourcePath.Png,
		}
	}
	return map[string]string{
		"icon": fh.PathJoin(di.Addr, di.ResourcePath.AssetsPath, di.ResourcePath.Icon),
		"png":  fh.PathJoin(di.Addr, di.ResourcePath.AssetsPath, di.ResourcePath.Png),
	}
}

func (di *DIContainer) ListIndependencies() map[string]interface{} {
	temp := make(map[string]interface{})

	for k, v := range di.independencies {
		temp[k] = v
	}
	return temp
}

func (di *DIContainer) SetAppCtx(appCtx *context.Context) {
	di.AppCtx = appCtx
}

func (di *DIContainer) GetAppCtx() *context.Context {
	return di.AppCtx
}

func (di *DIContainer) GetResourcePath() *model.ResourcePath {
	return &di.ResourcePath
}

func (di *DIContainer) SetAssetsPath() {
	assetsPath := map[string]string{"dev": "frontend/public", "production": "public"}

	resourcePath := model.ResourcePath{}

	env := runtime.Environment(*di.GetAppCtx()).BuildType
	di.Env = env

	exePath, _ := os.Executable()

	if env == "dev" {
		exePath = filepath.Dir(exePath)
		exePath = filepath.Dir(exePath)
	}

	rootPath := filepath.Dir(exePath)
	os.MkdirAll(assetsPath[env], os.ModePerm)

	resourcePath.RootPath = rootPath
	resourcePath.AssetsPath = assetsPath[env]
	resourcePath.File = "gamecommondata"
	resourcePath.Icon = "icons"
	resourcePath.Locale = "locales"
	resourcePath.Png = "product_icon_png"
	di.ResourcePath = resourcePath
}

func (di *DIContainer) GetFileHandler() *FileHandler {
	return Resolve[FileHandler](di, "fileHandler")
}

func (di *DIContainer) GetHttpServer() *HttpServer {
	return Resolve[HttpServer](di, "httpServer")
}
