package model

type ItemDetail struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Desc  string `json:"desc"`
	Count int    `json:"count,omitempty"`
}

type ManufactureItem struct {
	Item   []ItemDetail `json:"item"`
	Action string       `json:"action"`
}

type HouseItem struct {
	Type string       `json:"type"`
	Item []ItemDetail `json:"item"`
}

type ItemInfo struct {
	ItemKey    string            `json:"itemKey"`
	ItemName   string            `json:"itemName"`
	ItemIcon   string            `json:"itemIcon"`
	ItemDesc   string            `json:"itemDesc"`
	Fishing    string            `json:"fishing,omitempty"`
	Node       []string          `json:"node,omitempty"`
	Shop       []string          `json:"shop,omitempty"`
	House      []HouseItem       `json:"house,omitempty"`
	Gathering  []string          `json:"gathering,omitempty"`
	Processing []ManufactureItem `json:"processing,omitempty"`
	Cooking    [][]ItemDetail    `json:"cooking,omitempty"`
	Alchemy    [][]ItemDetail    `json:"alchemy,omitempty"`
	MakeList   []ItemDetail      `json:"makelist,omitempty"`
}

type ItemString struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
type ItemMaking struct {
	Node       []ItemString `json:"node"`
	Fishing    []ItemString `json:"fishin"`
	Alchemy    []ItemString `json:"alchemy"`
	Cooking    []ItemString `json:"cooking"`
	Processing []ItemString `json:"processing"`
	House      []ItemString `json:"house"`
	Gathring   []ItemString `json:"gathring"`
}

type ParserReturn struct {
	ItemInfo    *ItemInfo          `json:"itemInfo"`
	ItemMaking  *ItemMaking        `json:"itemMaking"`
	ProductNote *map[string]string `json:"productNote"`
}

type SearchIndexItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type SearchIndex struct {
	Locale string
	Data   map[string]ItemRaw
}

type DynamicStrings struct {
	Approach    map[string]string `json:"approach"`
	Manufacture map[string]string `json:"manufacture"`
	Workshop    []string          `json:"workshop"`
}
