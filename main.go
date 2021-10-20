package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var cfg Config

func dimGen(item []dimItem) {
	for _, di := range item {
		di.title = titleGen()
		di.children = make([]dimItem, 0) // r.cfgDim.Hier[0].)
		dimGen(di.children)
	}
}

func titleGen() string {
	return ""
}

func main() {

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	d, _ := ioutil.ReadAll(f)
	json.Unmarshal(d, &cfg)

	b := make([]sheet, len(cfg.Sheet))
	for i, s := range b {
		s.cfgSheet = &cfg.Sheet[i]
		s.row = make([]dim, len(s.cfgSheet.Row))
		for ii, r := range s.row {
			r.cfgDim = &s.cfgSheet.Row[ii]
			r.root = make([]dimItem, r.cfgDim.Root)

			for l := 0; l < len(r.cfgDim.Hier); l++ {

			}
			// r.title
		}
		s.col = make([]dim, len(s.cfgSheet.Col))
		for ii, c := range s.row {
			c.cfgDim = &s.cfgSheet.Col[ii]
		}
	}

	// rdim := make([]dim, cfg.)

	// fmt.Println(cfg)

	// e := excelize.NewFile()
	// index := e.NewSheet("Sheet2")
	// // Set value of a cell.
	// e.SetCellValue("Sheet2", "A2", "Hello world.")
	// e.SetCellValue("Sheet1", "B2", 100)
	// // Set active sheet of the workbook.
	// e.SetActiveSheet(index)
	// // Save spreadsheet by the given path.
	// if err := e.SaveAs(cfg.File); err != nil {
	// 	fmt.Println(err)
	// }

	// e.NewStreamWriter()

	// 	file, _ := os.Open("conf.json")
	// defer file.Close()
	// decoder := json.NewDecoder(file)
	// configuration := Configuration{}
	// err := decoder.Decode(&configuration)
	// if err != nil {
	//   fmt.Println("error:", err)
	// }

}
