package excel

import (
	"fmt"
	"wechatbot/pkg/errors"

	"github.com/xuri/excelize/v2"
)

type Excel struct {
	File *excelize.File
}

func New(sheetName string) *Excel {
	f := excelize.NewFile()
	if sheetName != "" {
		sheetList := f.GetSheetList()
		oldName := sheetList[0]
		f.SetSheetName(oldName, sheetName)
	}
	return &Excel{File: f}
}

var pow26tab = [...]int{1, 26, 676, 17576, 456976, 11881376}

func pow26(n int) int {
	if 0 <= n && n <= 5 {
		return pow26tab[n]
	}
	return pow26(n-1) * 26
}

func (e *Excel) ColumnToLetter(i int) string {
	if i < 0 {
		return ""
	}
	var r []rune
	for c := 0; ; c++ {
		mod := i%pow26(c+1) + 1
		r = append([]rune{rune(mod/pow26(c) + 64)}, r...)
		i -= mod
		if i <= 0 {
			break
		}
	}
	return string(r)
}

func (e *Excel) LetterToColumn(s string) (int, error) {
	if s == "" {
		return 0, errors.New("argument is empty string")
	}
	var r int
	for i, c := range s {
		if 'A' <= c && c <= 'Z' {
			r += (int(c) - 64) * pow26(len(s)-i-1)
		} else if 'a' <= c && c <= 'z' {
			r += (int(c) - 96) * pow26(len(s)-i-1)
		} else {
			return 0, errors.New("must not contain non-alphabetic characters")
		}
	}
	return r - 1, nil
}

func (e *Excel) NewExcel() *excelize.File {
	f := excelize.NewFile()
	return f
}

func (e *Excel) AppendSheet(sheetName string, rows [][]interface{}) (index int) {
	if sheetName != "Sheet1" {
		index = e.File.NewSheet(sheetName)
	}
	for rowIndex, row := range rows {
		for cellIndex, cell := range row {
			letter, _ := excelize.ColumnNumberToName(cellIndex + 1)
			e.File.SetCellValue(sheetName, fmt.Sprintf("%s%v", letter, rowIndex+1), cell)
		}
	}
	return
}

func (e *Excel) Save(p string) error {
	err := e.File.SaveAs(p)
	return err
}

func (e *Excel) FitToColumn(sheetName string) error {
	cols, err := e.File.GetCols(sheetName)
	if err != nil {
		return err
	}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			// cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			cellWidth := len(rowCell) + 2
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			return err
		}
		err = e.File.SetColWidth(sheetName, name, name, float64(largestWidth))
	}
	return err
}
