package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/xuri/excelize/v2"
)

var cfg Config

func dimGen(item *[]dimItem, level int, cfgDim *ConfigDim) {
	for i := range *item {
		di := &(*item)[i]
		di.title = titleGen(cfgDim.Length)
		if level < len(cfgDim.Hier) {
			di.children = make([]dimItem, cfgDim.Hier[level].Count)
			level++
			dimGen(&di.children, level, cfgDim)
		}
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func titleGen(length []int) string {
	n := length[0] + rand.Intn(length[1]-length[0])
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
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
	for i := range b {
		s := &b[i]
		if cfg.Sheet[i].Ref != "" {
			for si, se := range cfg.Sheet {
				if cfg.Sheet[i].Ref == se.Name {
					s.cfgSheet = &cfg.Sheet[si]
					break
				}
			}
		} else {
			s.cfgSheet = &cfg.Sheet[i]
			s.row = make([]dim, len(s.cfgSheet.Row))
			for ii := range s.row {
				r := &s.row[ii]
				r.cfgDim = &s.cfgSheet.Row[ii]
				r.root = make([]dimItem, r.cfgDim.Root)
				dimGen(&r.root, 0, r.cfgDim)
			}
			s.col = make([]dim, len(s.cfgSheet.Col))
			for ii := range s.col {
				c := &s.col[ii]
				c.cfgDim = &s.cfgSheet.Col[ii]
				c.root = make([]dimItem, c.cfgDim.Root)
				dimGen(&c.root, 0, c.cfgDim)
			}
		}
	}

	// fmt.Println(b)

	e := excelize.NewFile()
	for i := range b {
		sName := cfg.Sheet[i].Name
		if len(e.GetSheetList()) == i {
			e.NewSheet(sName)
		} else {
			e.SetSheetName(e.GetSheetName(i), sName)
		}
		s := &b[i]

		for i, r := range s.cfgSheet.Row {
			cell1, _ := excelize.CoordinatesToCellName(i+1, 1)
			e.SetCellValue(sName, cell1, r.Title)
			cell1m, _ := excelize.CoordinatesToCellName(i+1, len(s.cfgSheet.Col))
			e.MergeCell(sName, cell1, cell1m)
			cell2, _ := excelize.CoordinatesToCellName(i+1, len(s.cfgSheet.Col)+1)
			c, _, _ := excelize.SplitCellName(cell2)
			e.SetCellValue(sName, cell2, c)
		}

		// count at every level of hier recoursive

		// + titles rows
		// + ABC under titles rows
		// dims cols
		// 123 under dims cols
		// dims rows
		// fill values

		// for i := range s.row {

		// 	r := &s.row[i]
		// 	if r.cfgDim.Index {
		// 		n := 1
		// 		for _, dr := range s.row {
		// 			n = n * dr.cfgDim.Root * dr.cfgDim.Count
		// 		}

		// 	}
		// 	// r.cfgDim.Count
		// }
	}

	if err := e.SaveAs(cfg.File); err != nil {
		fmt.Println(err)
	}

	// index := e.NewSheet("Sheet2")
	// fmt.Println(index)

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
