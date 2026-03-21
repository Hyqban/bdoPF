import { Injectable, signal, WritableSignal } from '@angular/core';
import { DomSanitizer, SafeUrl } from '@angular/platform-browser';
import { GetAddr, GetImgPath } from '../../../wailsjs/go/service/DIContainer';
import { SearchResultItem, ItemInfo, BreadCrumbs, Item } from '../shared/models/model';
import { ReadFileById, ReadDynamicStrings } from '../../../wailsjs/go/service/FileHandler';

const DEFAULT_ITEM_INFO: ItemInfo = {
    itemKey: '',
    itemName: '',
    itemDesc: '',
    itemIcon: '',
    shop: [],
    node: [],
    house: [],
    gathering: [],
    cooking: [],
    alchemy: [],
    processing: [],
    fishing: '',
    makelist: [],
};

@Injectable({
    providedIn: 'root',
})
export class SearchService {
    // searchHistory = new SearchHistory();
    env = '';
    addr = '';
    query = signal('');
    searchResults = signal<Item[]>([]);
    searchHistoryes: WritableSignal<Item[]> = signal<Item[]>([]);
    currentItem: WritableSignal<ItemInfo> = signal<ItemInfo>({
        itemKey: '',
        itemName: '',
        itemIcon: '',
        itemDesc: '',
        gathering: [],
    });
    breadCrumbs = signal<BreadCrumbs>({
        data: [],
        amount: [1],
        index: 0,
        length: 0,
    });
    imgPath: Record<string, string> = {};
    // dynamicStrings: DynamicStrings = {
    //     apporach: {},
    //     manufacture: {},
    //     workshop: {},
    // };

    constructor(private sanitizer: DomSanitizer) {}

    ngOnInit() {
        this.getAddr();
        this.getImgPath();
    }

    getAddr(): void {
        GetAddr()
            .then((addr: string) => {
                this.addr = addr;
            })
            .catch((err) => {
                console.log(err);
            });
    }

    getImgPath(): void {
        GetImgPath().then((res) => {
            this.env = res['env'];
            this.imgPath = {
                png: res['png'],
                icon: res['icon'],
            };
        });
    }

    imgPathJoin(folderPath: string, iconPath: string): string {
        if (folderPath && iconPath) {
            if (this.env === 'dev') {
                // dev:   http://127.0.0.1:51780\public\product_icon_png/00000874.png
                // build: http://127.0.0.1:51780/public/127.0.0.1:51780\public\product_icon_png/00000874.png
                return (
                    'http://' + this.addr + '/public/' + folderPath + '/' + iconPath.toLowerCase()
                );
            }
            return 'http://' + folderPath + '/' + iconPath.toLowerCase();
        }
        return '';
    }

    setCurrentItem(el: ItemInfo) {
        console.log(el);
    }

    addSearchResults(searchResults: SearchResultItem[]) {
        this.searchResults.set(searchResults);
    }

    cleanSearchResult() {
        this.searchResults.set([]);
    }

    addItemToSearchHistory(ele: Item) {
        for (let i = 0; i < this.searchHistoryes().length; i++) {
            if (this.searchHistoryes()[i].name === ele.name) {
                return;
            }
        }
        this.searchHistoryes.update((el) => [ele, ...el]);
    }

    cleanSearchHistory() {
        this.searchHistoryes.set([]);
    }

    async selectItem(ele: Item): Promise<Record<string, any>> {
        this.cleanBreadCrumbs();

        this.breadCrumbs.update((el) => {
            el.data.push(ele);
            el.length += 1;

            return { ...el };
        });

        let itemInfo: Record<string, any> = await ReadFileById(ele.id);
        this.addItemToSearchHistory(ele);
        return itemInfo;
    }

    nextQueryAndSetCurrentItem(id: string) {
        ReadFileById(id).then((res) => {
            if (res.itemKey) {
                this.currentItem.set(res as ItemInfo);
            }
        });
    }

    nextQuery(ele: Item) {
        ReadFileById(ele.id).then((res) => {
            if (res.itemKey) {
                this.currentItem.set(res as ItemInfo);

                this.breadCrumbs.update((el) => {
                    el.amount.push(Number(ele?.count || 1));
                    el.data.push(ele);
                    el.index += 1;
                    el.length += 1;
                    return { ...el };
                });
            }
        });
        // this.nextQueryAndSetCurrentItem(ele.id);
    }

    addBreadCrumb(bd: SearchResultItem, amount: string) {
        this.breadCrumbs.update((el) => {
            const newData = [...el.data, bd];
            const newAmout = [...el.amount, Number(amount)];
            el.index += 1;
            el.length += 1;

            return {
                ...el,
                data: newData,
                amount: newAmout,
            };
        });
    }

    selectBreadCrumb(index: number) {
        // last icon
        if (index + 1 === this.breadCrumbs().length) {
            return;
        }

        this.breadCrumbs.update((el) => {
            if (index === 0) {
                el.data.length = 1;
                el.amount.length = 1;
                el.index = 0;
                el.length = 1;

                this.nextQueryAndSetCurrentItem(
                    this.breadCrumbs().data[this.breadCrumbs().index].id,
                );
                return { ...el };
            }

            if (index + 1 < el.length) {
                const newData = el.data.slice(0, index + 1);
                const newAmout = el.amount.slice(0, index + 1);
                el.index = index;
                el.length = index + 1;

                this.nextQueryAndSetCurrentItem(
                    this.breadCrumbs().data[this.breadCrumbs().index].id,
                );
                return {
                    ...el,
                    data: newData,
                    amount: newAmout,
                };
            }
            return { ...el };
        });
    }

    calculateDeltaAmount(count: number): number {
        const len = this.breadCrumbs().length;
        let total: number = 1;

        if (len === 1) {
            return count;
        }
        if (len > 1) {
            total = this.breadCrumbs().amount[len - 1] * count;
        }
        return total;
    }

    totalAmout(count: number): number {
        // console.log('count: ', count);
        const len = this.breadCrumbs().length;
        let total: number = 1;

        if (len === 1) {
            return count;
        }
        if (len > 1) {
            this.breadCrumbs().amount.forEach((el) => {
                total *= el;
            });
            total *= count;
        }
        // console.log('total: ', total);
        return total;
    }

    cleanBreadCrumbs() {
        this.breadCrumbs.update((el) => {
            el.data.length = 0;
            el.amount.length = 1;
            el.index = 0;
            el.length = 0;

            return { ...el };
        });
    }

    isEmpty(obj: any): boolean {
        for (const prop in obj) {
            return false;
        }
        return true;
    }

    cleanAllCache() {
        this.query.set('');
        this.cleanSearchHistory();
        this.cleanSearchResult();
        this.currentItem.set(DEFAULT_ITEM_INFO);
        this.cleanBreadCrumbs();
    }
}
