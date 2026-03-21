export interface WindowSize {
    w: number;
    h: number;
}
// Search
export interface SearchResult {
    score: number;
    matchString: string;
    item: { id: string; name: string; icon: string };
}

export interface SearchHistory {
    query: string;
    searchResult: Item[];
    currentItem: ItemInfo;
    breadCrumbs: Item[];
    breadCrumbsIndex: number;
    breadCrumbsLength: number;
    amount: number[];
}

export interface SearchResultItem {
    id: string;
    name: string;
    icon: string;
}

export interface BreadCrumbs {
    data: Item[];
    amount: number[];
    index: number;
    length: number;
}

export interface Item {
    id: string;
    name: string;
    icon: string;
    desc?: string;
    count?: number;
}

export interface HouseItem {
    type: string;
    item: Item[];
}

export interface processingItem {
    action: string;
    item: Item[];
}

export interface ItemInfo {
    itemKey: string;
    itemName: string;
    itemIcon: string;
    itemDesc: string;
    shop?: string[];
    node?: string[];
    house?: HouseItem[];
    gathering?: string[];
    cooking?: Item[][];
    alchemy?: Item[][];
    processing?: processingItem[];
    fishing?: string;
    makelist?: Item[];
}

export interface DynamicStrings {
    apporach: Record<string, string>;
    manufacture: Record<string, string>;
    workshop: Record<string, string>;
}

export interface AppWindow {
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
}

export interface Config {
    appName: string;
    version: string;
    theme: string;
    locale: string;
    window: AppWindow;
    newVersion: LatestApp;
}

export interface LatestApp {
    version: string;
    download: boolean;
    downloadUrl: string;
}

export interface WindowSizeChange {
    widthBeforeEnterWidget: number;
    heightBeforeEnterWidget: number;
    minWidthBeforeEnterWidget: number;
    minHeightBeforeEnterWidget: number;
}

export interface RecipeAmountInterface {
    open: boolean;
    items: Item[];
    amountItems: Item[];
    amount: number;
    averageYield: number;
}
