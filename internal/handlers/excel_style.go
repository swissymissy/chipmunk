package handlers

import "github.com/xuri/excelize/v2"

type ExcelStyle struct {
	Title        int
	Header       int
	Qualified    int
	AtRisk       int
	NotQualified int
}

// excel styles
func NewExcelStyles(f *excelize.File) ExcelStyle {
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#EEEEEE"}, Pattern: 1},
	})

	qualified, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#C6EFCE"}, Pattern: 1},
	})
	atRisk, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFEB9C"}, Pattern: 1},
	})
	notQualified, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFC7CE"}, Pattern: 1},
	})

	return ExcelStyle{
		Title:        titleStyle,
		Header:       headerStyle,
		Qualified:    qualified,
		AtRisk:       atRisk,
		NotQualified: notQualified,
	}
}
