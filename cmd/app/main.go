package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func log(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

var cfg Config

func dimGen(item *[]dimItem, level int, cfgDim *ConfigDim, size int) int {
	for i := range *item {
		di := &(*item)[i]
		di.title = titleGen(cfgDim.Length)
		size++
		if level < len(cfgDim.Hier) {
			di.children = make([]dimItem, cfgDim.Hier[level].Count)
			level++
			size = dimGen(&di.children, level, cfgDim, size)
		}
	}
	return size
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

func printCol(e *excelize.File, sn string, d []dimItem, t, size, zeroPos, pos, level, width int) {

	for _, di := range d {
		cell, _ := excelize.CoordinatesToCellName(zeroPos+t*size*width+pos*width, level+1)
		e.SetCellValue(sn, cell, di.title)
		cellm, _ := excelize.CoordinatesToCellName(zeroPos+t*size*width+pos*width+width-1, level+1)
		e.MergeCell(sn, cell, cellm)
		pos++
		printCol(e, sn, di.children, t, size, zeroPos, pos, level, width)
	}

}

func printRow(e *excelize.File, sn string, d []dimItem, t, size, zeroPos, pos, level, height int) {

	for _, di := range d {
		cell, _ := excelize.CoordinatesToCellName(level+1, zeroPos+t*size*height+pos*height)
		e.SetCellValue(sn, cell, di.title)
		cellm, _ := excelize.CoordinatesToCellName(level+1, zeroPos+t*size*height+pos*height+height-1)
		e.MergeCell(sn, cell, cellm)
		// if level == 1 {
		// 	log(zeroPos, t, size, pos, height)
		// 	log(cell, cellm)
		// }
		pos++
		printRow(e, sn, di.children, t, size, zeroPos, pos, level, height)
	}

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
		s.name = cfg.Sheet[i].Name
		if cfg.Sheet[i].Ref != "" {
			for si, se := range cfg.Sheet {
				if cfg.Sheet[i].Ref == se.Name {
					s.cfgSheet = &cfg.Sheet[si]
					s.dims = b[si].dims
					break
				}
			}
		} else {
			s.cfgSheet = &cfg.Sheet[i]
			s.dims = &sheetDim{}
			s.dims.row = make([]dim, len(s.cfgSheet.Row))
			for ii := range s.dims.row {
				r := &s.dims.row[ii]
				r.cfgDim = &s.cfgSheet.Row[ii]
				r.root = make([]dimItem, r.cfgDim.Root)
				r.size = dimGen(&r.root, 0, r.cfgDim, 0)
				if r.cfgDim.Index {
					r.size = 1
				}
			}
			s.dims.col = make([]dim, len(s.cfgSheet.Col))
			for ii := range s.dims.col {
				c := &s.dims.col[ii]
				c.cfgDim = &s.cfgSheet.Col[ii]
				c.root = make([]dimItem, c.cfgDim.Root)
				c.size = dimGen(&c.root, 0, c.cfgDim, 0)
			}
		}
	}
	log("Dims generated")

	e := excelize.NewFile()
	for i := range b {
		s := &b[i]
		if len(e.GetSheetList()) == i {
			e.NewSheet(s.name)
		} else {
			e.SetSheetName(e.GetSheetName(i), s.name)
		}

		width := 1
		for _, c := range s.dims.col {
			width *= c.size
		}
		height := 1
		for _, r := range s.dims.row {
			height *= r.size
		}
		// log(width, height)

		if len(s.cfgSheet.Row)+width > 16384 {
			fmt.Println("Column count is more than limit of 16,384:", len(s.cfgSheet.Row)+width)
			os.Exit(2)
		}
		if len(s.cfgSheet.Col)+1+height > 1048576 {
			fmt.Println("Row count is more than limit of 1,048,576:", len(s.cfgSheet.Col)+1+height)
			os.Exit(2)
		}

		streamWriter, err := e.NewStreamWriter(s.name)
		if err != nil {
			fmt.Println(err)
		}
		for rowID := len(s.cfgSheet.Col) + 2; rowID < height+len(s.cfgSheet.Col)+2; rowID++ {
			row := make([]interface{}, width)
			for colID := 0; colID < width; colID++ {
				row[colID] = rand.Float64() * math.Pow10(rand.Intn(6))
			}
			cell, _ := excelize.CoordinatesToCellName(len(s.cfgSheet.Row)+1, rowID)
			if err := streamWriter.SetRow(cell, row); err != nil {
				fmt.Println(err)
			}
		}
		if err := streamWriter.Flush(); err != nil {
			fmt.Println(err)
		}

	}
	log("Streamed ok")

	if err := e.SaveAs(cfg.File); err != nil {
		fmt.Println(err)
	}
	log("Save streamed ok")

	e, err = excelize.OpenFile(cfg.File)

	for i := range b {
		s := &b[i]
		if len(e.GetSheetList()) == i {
			e.NewSheet(s.name)
		} else {
			e.SetSheetName(e.GetSheetName(i), s.name)
		}

		for ii, r := range s.cfgSheet.Row {
			cell1, _ := excelize.CoordinatesToCellName(ii+1, 1)
			e.SetCellValue(s.name, cell1, r.Title)
			cell1m, _ := excelize.CoordinatesToCellName(ii+1, len(s.cfgSheet.Col))
			e.MergeCell(s.name, cell1, cell1m)
			cell2, _ := excelize.CoordinatesToCellName(ii+1, len(s.cfgSheet.Col)+1)
			c, _, _ := excelize.SplitCellName(cell2)
			e.SetCellValue(s.name, cell2, c)
		}

		for ii, c := range s.dims.col {

			times := 1
			// if ii > 0 {
			for iii := 0; iii < ii; iii++ {
				times *= s.dims.col[iii].size
			}
			// }
			width := 1
			// if ii < len(s.dims.col)-1 {
			for iii := len(s.dims.col) - 1; iii > ii; iii-- {
				width *= s.dims.col[iii].size
			}
			// }

			// if len(s.cfgSheet.Row)+times*s.dims.col[ii].size*width > 16384 {
			// 	fmt.Println("Column count is more than limit of 16,384")
			// 	os.Exit(2)
			// }

			for t := 0; t < times; t++ {
				printCol(e, s.name, c.root, t, c.size, len(s.cfgSheet.Row)+1, 0, ii, width)
			}

			if width == 1 {
				for iii := 0; iii < times*s.dims.col[len(s.dims.col)-1].size; iii++ {
					cell, _ := excelize.CoordinatesToCellName(len(s.cfgSheet.Row)+1+iii, ii+2)
					e.SetCellValue(s.name, cell, strconv.Itoa(iii+1))
				}
				style, _ := e.NewStyle(`{
					"alignment": {
						"horizontal": "center",
						"vertical": "center"
					}
				}`)
				fromCell, _ := excelize.CoordinatesToCellName(1, 1)
				toCell, _ := excelize.CoordinatesToCellName(len(s.cfgSheet.Row)+1+times*s.dims.col[len(s.dims.col)-1].size, ii+2)
				e.SetCellStyle(s.name, fromCell, toCell, style)
			}

		}

		for ii, r := range s.dims.row {

			times := 1
			// if ii > 0 {
			for iii := 0; iii < ii; iii++ {
				times *= s.dims.row[iii].size
			}
			// }
			height := 1
			// if ii < len(s.dims.row)-1 {
			for iii := len(s.dims.row) - 1; iii > ii; iii-- {
				height *= s.dims.row[iii].size
			}
			// }

			// if len(s.cfgSheet.Col)+1+times*s.dims.row[ii].size*height > 1048576 {
			// 	fmt.Println("Row count is more than limit of 1,048,576")
			// 	os.Exit(2)
			// }

			fromCellCol := 1
			if r.cfgDim.Index {
				for index := 0; index < height; index++ {
					cell, _ := excelize.CoordinatesToCellName(1, len(s.cfgSheet.Col)+2+index)
					e.SetCellValue(s.name, cell, strconv.Itoa(index+1))
				}
				fromCellCol = 2
			} else {
				for t := 0; t < times; t++ {
					printRow(e, s.name, r.root, t, r.size, len(s.cfgSheet.Col)+2, 0, ii, height)
				}
				if height == 1 {
					style, _ := e.NewStyle(`{
						"alignment": {
							"vertical": "top"
						}
					}`)
					fromCell, _ := excelize.CoordinatesToCellName(fromCellCol, len(s.cfgSheet.Col)+2)
					toCell, _ := excelize.CoordinatesToCellName(len(s.cfgSheet.Row), len(s.cfgSheet.Col)+1+times*s.dims.row[len(s.dims.row)-1].size*height)
					e.SetCellStyle(s.name, fromCell, toCell, style)
				}
			}

		}

	}
	log("Dims ok")

	if err := e.SaveAs(cfg.File); err != nil {
		fmt.Println(err)
	}
	log("Save all ok")

}
