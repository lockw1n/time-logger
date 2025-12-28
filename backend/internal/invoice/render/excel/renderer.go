package excel

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func Render(invoice InvoiceView) ([]byte, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	// === Column widths (PDF-aligned) ===
	f.SetColWidth(sheet, "A", "A", 18) // Date
	f.SetColWidth(sheet, "B", "B", 42) // Notes
	f.SetColWidth(sheet, "C", "C", 12) // Hours

	// === Styles ===
	titleStyle,
		headerStyle,
		cellDateStyle,
		cellNotesStyle,
		cellHoursStyle,
		subtotalStyle,
		subtotalHoursStyle,
		totalStyle,
		totalHoursStyle := styles(f)

	row := 1
	var subtotalCells []string

	for _, activity := range invoice.Activities {
		// === Block title ===
		titleCell := fmt.Sprintf("A%d", row)
		f.SetCellValue(sheet, titleCell, activity.Title)
		f.MergeCell(sheet, titleCell, fmt.Sprintf("C%d", row))
		f.SetCellStyle(sheet, titleCell, titleCell, titleStyle)
		f.SetRowHeight(sheet, row, 24)
		row++

		// === Table header ===
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Date")
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Notes")
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "Hours")
		f.SetCellStyle(
			sheet,
			fmt.Sprintf("A%d", row),
			fmt.Sprintf("C%d", row),
			headerStyle,
		)
		row++

		startRow := row

		// === Data rows ===
		for _, entry := range activity.Entries {
			a := fmt.Sprintf("A%d", row)
			b := fmt.Sprintf("B%d", row)
			c := fmt.Sprintf("C%d", row)

			f.SetCellValue(sheet, a, entry.DateFormatted)
			f.SetCellValue(sheet, b, entry.TicketCode)
			f.SetCellValue(sheet, c, entry.Hours)

			f.SetCellStyle(sheet, a, a, cellDateStyle)
			f.SetCellStyle(sheet, b, b, cellNotesStyle)
			f.SetCellStyle(sheet, c, c, cellHoursStyle)

			row++
		}

		// === Subtotal ===
		subtitleCell := fmt.Sprintf("A%d", row)
		subtitleEndCell := fmt.Sprintf("B%d", row)
		valueCell := fmt.Sprintf("C%d", row)

		f.SetCellValue(sheet, subtitleCell, "Subtotal "+activity.Title)
		f.MergeCell(sheet, subtitleCell, subtitleEndCell)
		f.SetCellFormula(
			sheet,
			valueCell,
			fmt.Sprintf("SUM(C%d:C%d)", startRow, row-1),
		)

		// Apply styles correctly (NO overwriting)
		f.SetCellStyle(sheet, subtitleCell, subtitleEndCell, subtotalStyle)
		f.SetCellStyle(sheet, valueCell, valueCell, subtotalHoursStyle)

		subtotalCells = append(subtotalCells, valueCell)

		row += 2 // spacing between blocks
	}

	// === Grand total ===
	totalLabelCell := fmt.Sprintf("A%d", row)
	totalLabelEndCell := fmt.Sprintf("B%d", row)
	totalValueCell := fmt.Sprintf("C%d", row)

	f.SetCellValue(sheet, totalLabelCell, "Total over all projects and tasks")
	f.MergeCell(sheet, totalLabelCell, totalLabelEndCell)
	f.SetCellFormula(
		sheet,
		totalValueCell,
		fmt.Sprintf("SUM(%s)", strings.Join(subtotalCells, ",")),
	)

	// Apply styles correctly
	f.SetCellStyle(sheet, totalLabelCell, totalLabelEndCell, totalStyle)
	f.SetCellStyle(sheet, totalValueCell, totalValueCell, totalHoursStyle)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func styles(f *excelize.File) (
	title,
	header,
	cellDate,
	cellNotes,
	cellHours,
	subtotal,
	subtotalHours,
	total,
	totalHours int,
) {
	// === Group title (PDF label bar) ===
	title, _ = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#E08B4F"},
		},
		Border: []excelize.Border{
			{Type: "top", Style: 2},
			{Type: "bottom", Style: 2},
		},
	})

	// === Table header row ===
	header, _ = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Border: []excelize.Border{
			{Type: "bottom", Style: 1},
		},
	})

	// === Date cells ===
	cellDate, _ = f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			WrapText:   false,
		},
	})

	// === Notes cells ===
	cellNotes, _ = f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			WrapText:   false,
		},
	})

	// === Regular hours cells ===
	cellHours, _ = f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			WrapText:   false,
		},
		NumFmt: 2, // 0.00
	})

	// === Subtotal label (A–B) ===
	subtotal, _ = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "top", Style: 2},
		},
	})

	// === Subtotal hours (C) ===
	subtotalHours, _ = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			WrapText:   false,
		},
		NumFmt: 2,
		Border: []excelize.Border{
			{Type: "top", Style: 2},
		},
	})

	// === Grand total label (A–B) ===
	total, _ = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 13,
		},
		Border: []excelize.Border{
			{Type: "top", Style: 2},
		},
	})

	// === Grand total hours (C) ===
	totalHours, _ = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 13,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			WrapText:   false,
		},
		NumFmt: 2,
		Border: []excelize.Border{
			{Type: "top", Style: 2},
		},
	})

	return
}
