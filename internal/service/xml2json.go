package service

import (
	model "bdoPF/internal/model"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	rt "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/beevik/etree"
)

var maxWrokers = runtime.NumCPU() * 2

type GameData struct {
	Di          *DIContainer
	defaultLang string
	savePath    string
}

func NewGameData(di *DIContainer) *GameData {
	gd := GameData{}

	gd.Di = di

	return &gd
}

func (gd *GameData) OpenFolderDialog() map[string]string {
	selectedDir, _ := rt.OpenDirectoryDialog(*gd.Di.GetAppCtx(), rt.OpenDialogOptions{
		Title: "Please select a folder",
	})
	return map[string]string{
		"folderPath": selectedDir,
	}
}

func (gd *GameData) XmlToJson(xmlFolderPath string, defaultLang string) map[string]string {
	start := time.Now()
	gd.defaultLang = defaultLang
	// frontend/public or public
	gd.savePath = gd.Di.GetResourcePath().AssetsPath

	fh := gd.Di.GetFileHandler()

	response := make(map[string]string)
	response["msg"] = "Got a wrong during xml to json."

	if fh == nil {
		return response
	}

	exist := fh.pathExists(xmlFolderPath)

	if !exist {
		response["msg"] = "Selected path is invalid."
		return response
	}

	configMap := make(map[string]interface{})

	files := gd.ListDifferentLangXml(xmlFolderPath, gd.defaultLang)

	for k, v := range files {
		path := fh.PathJoin(gd.Di.GetResourcePath().RootPath, gd.savePath, gd.Di.GetResourcePath().File, k)
		isExist := fh.pathExists(path)

		if !isExist {
			err := os.MkdirAll(path, 0755)

			if err != nil {
				response["msg"] = "Failed to create json file."
				return response
			}
		}

		configMap["lang"] = k
		configMap["xmlFiles"] = v
		configMap["savePath"] = path
		configMap["xmlPath"] = xmlFolderPath

		gd.run(configMap)
	}

	elapsed := time.Since(start).Truncate(time.Second)
	fmt.Printf("\n============================================\n")
	fmt.Printf("Finished, cost: %s\n", elapsed)
	fmt.Printf("============================================\n")
	response["msg"] = fmt.Sprintf("Finished, cost %s", elapsed)
	return response
}

func (gd *GameData) ListDifferentLangXml(path string, defaultLang string) map[string][]string {
	xmlSlice := make(map[string][]string)

	files, _ := os.ReadDir(path)

	for _, f := range files {
		if !f.IsDir() {
			fileNameSplit := strings.Split(f.Name(), "_")

			if len(fileNameSplit) == 3 {
				xmlSlice[fileNameSplit[1]] = append(xmlSlice[fileNameSplit[1]], f.Name())
			} else if len(fileNameSplit) == 1 {
				xmlSlice[defaultLang] = append(xmlSlice[defaultLang], f.Name())
			}
		}
	}
	return xmlSlice
}

func (gd *GameData) run(configMap map[string]interface{}) {
	searchIndex := make(map[string]model.SearchIndexItem)
	var searchIndexMutex sync.RWMutex

	fh := gd.Di.GetFileHandler()

	if fh == nil {
		return
	}

	var wg sync.WaitGroup
	limiter := make(chan struct{}, maxWrokers)

	// Clean map
	// searchIndex = make(map[string]model.SearchIndexItem)

	for _, el := range configMap["xmlFiles"].([]string) {

		wg.Add(1)
		limiter <- struct{}{}

		go func(el string) {
			defer wg.Done()
			defer func() { <-limiter }()

			var fname string

			fileNameSplit := strings.Split(el, "_")
			if len(fileNameSplit) == 3 {
				fname = fileNameSplit[2]
			} else if len(fileNameSplit) == 1 {
				fname = fileNameSplit[0]
			}
			fname = strings.Split(fname, ".")[0] + ".json"

			xmlFilePath := fh.PathJoin(configMap["xmlPath"].(string), el)
			jsonFilePath := fh.PathJoin(configMap["savePath"].(string), fname)

			file, err := os.Open(xmlFilePath)
			if err != nil {
				return
			}

			xml := etree.NewDocument()
			defer file.Close()

			if _, err := xml.ReadFrom(file); err != nil {
				return
			}

			root := xml.Root()
			parserReturn := gd.chooseParser(root)

			switch {
			case parserReturn.ItemInfo != nil:
				gd.createJsonFile(parserReturn, "itemInfo", jsonFilePath)

				searchIndexMutex.Lock()
				searchIndex[parserReturn.ItemInfo.ItemKey] = model.SearchIndexItem{
					Id:   parserReturn.ItemInfo.ItemKey,
					Name: parserReturn.ItemInfo.ItemName,
					Icon: parserReturn.ItemInfo.ItemIcon,
				}
				searchIndexMutex.Unlock()
			case parserReturn.ItemMaking != nil:
				gd.createJsonFile(parserReturn, "itemmaking", jsonFilePath)
			case parserReturn.ProductNote != nil:
				gd.createJsonFile(parserReturn, "productNote", jsonFilePath)
			}
		}(el)
	}
	wg.Wait()
	gd.GenerateSearchIndex(&searchIndex, fh.PathJoin(configMap["savePath"].(string), "search_index.json"))
}

func (gd *GameData) createJsonFile(parserValue model.ParserReturn, key string, jsonFilePath string) {
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(jsonFilePath, flag, 0644)

	if err != nil {
		fmt.Printf("%s Open or Create is failed!!! \n", jsonFilePath)
	}
	defer f.Close()

	switch key {
	case "itemInfo":
		gd.createAndWriteJson(f, parserValue.ItemInfo)
	case "itemmaking":
		gd.createAndWriteJson(f, parserValue.ItemMaking)
	case "productNote":
		gd.createAndWriteJson(f, parserValue.ProductNote)
	}

}

func (gd *GameData) createAndWriteJson(f *os.File, data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		fmt.Println("JSON MarshalIndent is failed!!!")
		return
	}

	writer := bufio.NewWriter(f)
	n, err := f.Write(jsonBytes)
	if err := writer.Flush(); err != nil {
		fmt.Println("Flush JSON to file is failed!!!")
		return
	}

	if err != nil {
		fmt.Println("Write JSON to file is failed!!!")
		return
	}

	if n != len(jsonBytes) {
		fmt.Println("Write JSON to file is incomplete!!!")
		return
	}
}

func (gd *GameData) chooseParser(root *etree.Element) model.ParserReturn {
	parser := model.ParserReturn{}

	rootTag := root.Tag

	switch rootTag {
	case "itemInfo":
		parser.ItemInfo = gd.itemInfoParser(root)
		parser.ItemMaking = nil
		parser.ProductNote = nil
	case "itemmaking":
		parser.ItemInfo = nil
		parser.ItemMaking = gd.itemmakingParser(root)
		parser.ProductNote = nil
	case "productNote":
		parser.ItemInfo = nil
		parser.ItemMaking = nil
		parser.ProductNote = gd.productNotePaser(root)
	}
	return parser
}

func (gd *GameData) itemInfoParser(root *etree.Element) *model.ItemInfo {
	itemInfo := model.ItemInfo{}

	for _, child := range root.ChildElements() {
		switch child.Tag {
		case "itemKey":
			itemInfo.ItemKey = strings.TrimSpace(child.Text())
		case "itemName":
			itemInfo.ItemName = strings.TrimSpace(child.Text())
		case "itemIcon":
			itemInfo.ItemIcon = strings.TrimSpace(child.Text())
		case "itemDesc":
			itemInfo.ItemDesc = gd.NormalizeWhitespace(child.Text())
		case "fishing":
			itemInfo.Fishing = "1"
		case "node":
			itemInfo.Node = append(itemInfo.Node, gd.node(child))
		case "shop":
			itemInfo.Shop = append(itemInfo.Shop, gd.getCharacterName(child))
		case "house":
			itemInfo.House = append(itemInfo.House, gd.house(child))
		case "collect":
			itemInfo.Gathering = append(itemInfo.Gathering, gd.getCharacterName(child))
		case "manufacture":
			itemInfo.Processing = append(itemInfo.Processing, gd.manufacture(child))
		case "cook":
			itemInfo.Cooking = append(itemInfo.Cooking, gd.parseItem(child.ChildElements()))
		case "alchemy":
			itemInfo.Alchemy = append(itemInfo.Alchemy, gd.parseItem(child.ChildElements()))
		case "makelist":
			itemInfo.MakeList = append(itemInfo.MakeList, gd.parseItem(child.ChildElements())...)
		}
	}
	return &itemInfo
}

func (gd *GameData) NormalizeWhitespace(desc string) string {
	return strings.TrimSpace(strings.ReplaceAll(desc, "\\n", "\n"))
}

func (gd *GameData) StripDataPlaceholders(desc string) string {
	// \s*-\s*\{.*?\}    - {...}
	// \s*\(\{.*?\}\)    ({data})
	// \s*\{\{.*?\}\}    {{text}}
	// \s*\{.*?\}        {data}
	// processedText := strings.ReplaceAll(desc, "\\n", "\n")

	// patternA := `\s*\(\{.*?\}\)`
	// reA := regexp.MustCompile(patternA)
	// result := reA.ReplaceAllString(processedText, "")

	// patternB := `\s*-\s*\{.*?\}|\s*\(\{.*?\}\)|\s*\{\{.*?\}\}|\s*\{.*?\}`
	// reB := regexp.MustCompile(patternB)

	// result = reB.ReplaceAllString(result, "")
	// fmt.Println(result)
	// return result
	pattern := `\s*-\s*\{.*?\}|\s*\(\{.*?\}\)|\s*\{\{.*?\}\}|\s*\{.*?\}`

	cleanedText := gd.NormalizeWhitespace(desc)
	re := regexp.MustCompile(pattern)
	result := re.ReplaceAllString(cleanedText, "")
	return result

}

func (gd *GameData) safeExtracText(parent *etree.Element, tag string) string {
	child := parent.FindElement(tag)

	if child == nil {
		return ""
	}
	return strings.TrimSpace(child.Text())
}

func (gd *GameData) safeExtractAttr(el *etree.Element, attr string) string {
	return el.SelectAttrValue(attr, "N/A")
}

func (gd *GameData) parseItem(item []*etree.Element) []model.ItemDetail {
	tempData := []model.ItemDetail{}

	for _, el := range item {

		item := model.ItemDetail{
			Id:   gd.safeExtracText(el, "id"),
			Name: gd.safeExtracText(el, "name"),
			Icon: gd.safeExtracText(el, "icon"),
			Desc: gd.StripDataPlaceholders(gd.safeExtracText(el, "desc")),
		}
		count, _ := strconv.Atoi(gd.safeExtracText(el, "count"))

		if count >= 1 {
			item.Count = count
		}

		tempData = append(tempData, item)
	}
	return tempData
}

func (gd *GameData) node(node *etree.Element) string {
	return strings.TrimSpace(gd.safeExtractAttr(node, "region"))
}

func (gd *GameData) getCharacterName(shop *etree.Element) string {
	return strings.TrimSpace(shop.FindElement("character").FindElement("name").Text())
}

func (gd *GameData) house(house *etree.Element) model.HouseItem {
	houseItem := model.HouseItem{}

	houseItem.Type = gd.safeExtractAttr(house, "type")
	houseItem.Item = append(houseItem.Item, gd.parseItem(house.ChildElements())...)
	return houseItem
}

// Do it at the frontend
func (gd *GameData) actionConvert(action string) string {
	if action == "N/A" {
		return action
	}
	return action
}

func (gd *GameData) manufacture(manufacture *etree.Element) model.ManufactureItem {
	processing := model.ManufactureItem{}

	processing.Action = gd.actionConvert(gd.safeExtractAttr(manufacture, "action"))
	processing.Item = append(processing.Item, gd.parseItem(manufacture.ChildElements())...)
	return processing
}

func (gd *GameData) GenerateSearchIndex(searchIndexData *map[string]model.SearchIndexItem, savePath string) {
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(savePath, flag, 0644)

	if err != nil {
		fmt.Printf("%s Open or Create is failed!!! \n", savePath)
	}
	defer f.Close()

	gd.createAndWriteJson(f, &searchIndexData)
}

func (gd *GameData) itemmakingParser(root *etree.Element) *model.ItemMaking {
	itemmaking := model.ItemMaking{}

	for _, child := range root.ChildElements() {
		switch child.Tag {
		case "nodeProduct":
			itemmaking.Node = append(itemmaking.Node, gd.extractItem(child.ChildElements())...)
		case "fishing":
			itemmaking.Fishing = append(itemmaking.Fishing, gd.extractItem(child.ChildElements())...)
		case "alchemy":
			itemmaking.Alchemy = append(itemmaking.Alchemy, gd.extractItem(child.ChildElements())...)
		case "cook":
			itemmaking.Cooking = append(itemmaking.Cooking, gd.extractItem(child.ChildElements())...)
		case "manufacture":
			itemmaking.Processing = append(itemmaking.Processing, gd.extractItem(child.ChildElements())...)
		case "housecraft":
			itemmaking.House = append(itemmaking.House, gd.extractItem(child.ChildElements())...)
		case "collect":
			itemmaking.Gathring = append(itemmaking.Gathring, gd.extractItem(child.ChildElements())...)
		}
	}
	return &itemmaking
}

func (gd *GameData) extractItem(items []*etree.Element) []model.ItemString {
	tempData := []model.ItemString{}

	for _, el := range items {
		tempData = append(tempData, model.ItemString{
			Id:   gd.safeExtractAttr(el, "key"),
			Name: gd.safeExtractAttr(el, "name"),
			Icon: gd.safeExtractAttr(el, "icon"),
		})
	}
	return tempData
}

func (gd *GameData) productNotePaser(root *etree.Element) *map[string]string {
	itemStrings := map[string]string{}

	for _, el := range root.ChildElements()[0].ChildElements() {
		itemStrings[gd.safeExtractAttr(el, "index")] = gd.safeExtractAttr(el, "name")

	}
	return &itemStrings
}
