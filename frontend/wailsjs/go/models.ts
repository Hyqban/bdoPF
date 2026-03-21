export namespace model {
	
	export class ItemDetail {
	    id: string;
	    name: string;
	    icon: string;
	    desc: string;
	    count?: number;
	
	    static createFrom(source: any = {}) {
	        return new ItemDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	        this.desc = source["desc"];
	        this.count = source["count"];
	    }
	}
	export class HouseItem {
	    type: string;
	    item: ItemDetail[];
	
	    static createFrom(source: any = {}) {
	        return new HouseItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.item = this.convertValues(source["item"], ItemDetail);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ManufactureItem {
	    item: ItemDetail[];
	    action: string;
	
	    static createFrom(source: any = {}) {
	        return new ManufactureItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.item = this.convertValues(source["item"], ItemDetail);
	        this.action = source["action"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ItemInfo {
	    itemKey: string;
	    itemName: string;
	    itemIcon: string;
	    itemDesc: string;
	    fishing?: string;
	    node?: string[];
	    shop?: string[];
	    house?: HouseItem[];
	    gathering?: string[];
	    processing?: ManufactureItem[];
	    cooking?: ItemDetail[][];
	    alchemy?: ItemDetail[][];
	    makelist?: ItemDetail[];
	
	    static createFrom(source: any = {}) {
	        return new ItemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.itemKey = source["itemKey"];
	        this.itemName = source["itemName"];
	        this.itemIcon = source["itemIcon"];
	        this.itemDesc = source["itemDesc"];
	        this.fishing = source["fishing"];
	        this.node = source["node"];
	        this.shop = source["shop"];
	        this.house = this.convertValues(source["house"], HouseItem);
	        this.gathering = source["gathering"];
	        this.processing = this.convertValues(source["processing"], ManufactureItem);
	        this.cooking = this.convertValues(source["cooking"], ItemDetail);
	        this.alchemy = this.convertValues(source["alchemy"], ItemDetail);
	        this.makelist = this.convertValues(source["makelist"], ItemDetail);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ItemRaw {
	    id: string;
	    name: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new ItemRaw(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	    }
	}
	export class LatestApp {
	    version: string;
	    download: boolean;
	    downloadUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new LatestApp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.download = source["download"];
	        this.downloadUrl = source["downloadUrl"];
	    }
	}
	
	export class ResourcePath {
	    RootPath: string;
	    AssetsPath: string;
	    File: string;
	    Icon: string;
	    Locale: string;
	    Png: string;
	
	    static createFrom(source: any = {}) {
	        return new ResourcePath(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RootPath = source["RootPath"];
	        this.AssetsPath = source["AssetsPath"];
	        this.File = source["File"];
	        this.Icon = source["Icon"];
	        this.Locale = source["Locale"];
	        this.Png = source["Png"];
	    }
	}
	export class ResponseMsg {
	    code: string;
	    msg: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new ResponseMsg(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}

}

export namespace service {
	
	export class Window {
	    onTop: boolean;
	    width: number;
	    height: number;
	    maxWidth: number;
	    maxHeight: number;
	    minWidth: number;
	    minHeight: number;
	    isFullScreen: boolean;
	    isWidgetMode: boolean;
	    defaultWidgetWidth: number;
	    defaultWidgetHeight: number;
	    widgetWidth: number;
	    widgetHeight: number;
	
	    static createFrom(source: any = {}) {
	        return new Window(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.onTop = source["onTop"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.maxWidth = source["maxWidth"];
	        this.maxHeight = source["maxHeight"];
	        this.minWidth = source["minWidth"];
	        this.minHeight = source["minHeight"];
	        this.isFullScreen = source["isFullScreen"];
	        this.isWidgetMode = source["isWidgetMode"];
	        this.defaultWidgetWidth = source["defaultWidgetWidth"];
	        this.defaultWidgetHeight = source["defaultWidgetHeight"];
	        this.widgetWidth = source["widgetWidth"];
	        this.widgetHeight = source["widgetHeight"];
	    }
	}
	export class Config {
	    appName: string;
	    version: string;
	    window: Window;
	    theme: string;
	    locale: string;
	    newVersion: model.LatestApp;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appName = source["appName"];
	        this.version = source["version"];
	        this.window = this.convertValues(source["window"], Window);
	        this.theme = source["theme"];
	        this.locale = source["locale"];
	        this.newVersion = this.convertValues(source["newVersion"], model.LatestApp);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DIContainer {
	    AppCtx?: any;
	    Addr: string;
	    Env: string;
	    ResourcePath: model.ResourcePath;
	
	    static createFrom(source: any = {}) {
	        return new DIContainer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AppCtx = source["AppCtx"];
	        this.Addr = source["Addr"];
	        this.Env = source["Env"];
	        this.ResourcePath = this.convertValues(source["ResourcePath"], model.ResourcePath);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FileHandler {
	    DI?: DIContainer;
	
	    static createFrom(source: any = {}) {
	        return new FileHandler(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DI = this.convertValues(source["DI"], DIContainer);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class HttpServer {
	    DI?: DIContainer;
	    Addr: string;
	
	    static createFrom(source: any = {}) {
	        return new HttpServer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DI = this.convertValues(source["DI"], DIContainer);
	        this.Addr = source["Addr"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

