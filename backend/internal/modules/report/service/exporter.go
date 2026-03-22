package service

import (
	"fmt"
	"io"
	"strings"

	"github.com/koperasi-gresik/backend/internal/modules/report/dto"
	"github.com/xuri/excelize/v2"
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) ExportBalanceSheetExcel(res *dto.BalanceSheetResponse, writer io.Writer) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Neraca"
	f.SetSheetName("Sheet1", sheet)

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 15)
	f.SetColWidth(sheet, "B", "B", 40)
	f.SetColWidth(sheet, "C", "C", 20)

	// Headers
	f.SetCellValue(sheet, "A1", "LAPORAN NERACA (BALANCE SHEET)")
	f.SetCellValue(sheet, "A3", "Kode Akun")
	f.SetCellValue(sheet, "B3", "Nama Akun")
	f.SetCellValue(sheet, "C3", "Saldo")

	row := 4

	var writeTree func(nodes []*dto.AccountNode, level int)
	writeTree = func(nodes []*dto.AccountNode, level int) {
		for _, n := range nodes {
			indent := strings.Repeat("    ", level)
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), n.Code)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), indent+n.Name)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), n.EndingBalance)
			row++
			writeTree(n.Children, level+1)
		}
	}

	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "ASET")
	row++
	writeTree(res.Assets, 1)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "TOTAL ASET")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), res.TotalAssets)
	row += 2

	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "KEWAJIBAN")
	row++
	writeTree(res.Liabilities, 1)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "TOTAL KEWAJIBAN")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), res.TotalLia)
	row += 2

	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "MODAL")
	row++
	writeTree(res.Equity, 1)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "TOTAL MODAL")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), res.TotalEquity)

	return f.Write(writer)
}

func (e *Exporter) ExportProfitLossExcel(res *dto.ProfitLossResponse, writer io.Writer) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Laba Rugi"
	f.SetSheetName("Sheet1", sheet)

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 15)
	f.SetColWidth(sheet, "B", "B", 40)
	f.SetColWidth(sheet, "C", "C", 20)

	// Headers
	f.SetCellValue(sheet, "A1", "LAPORAN LABA RUGI (PROFIT & LOSS)")
	f.SetCellValue(sheet, "A3", "Kode Akun")
	f.SetCellValue(sheet, "B3", "Nama Akun")
	f.SetCellValue(sheet, "C3", "Saldo")

	row := 4

	var writeTree func(nodes []*dto.AccountNode, level int)
	writeTree = func(nodes []*dto.AccountNode, level int) {
		for _, n := range nodes {
			indent := strings.Repeat("    ", level)
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), n.Code)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), indent+n.Name)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), n.PeriodMovement)
			row++
			writeTree(n.Children, level+1)
		}
	}

	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "PENDAPATAN")
	row++
	writeTree(res.Revenues, 1)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "TOTAL PENDAPATAN")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), res.TotalRevenue)
	row += 2

	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "BEBAN")
	row++
	writeTree(res.Expenses, 1)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "TOTAL BEBAN")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), res.TotalExpenses)
	row += 2

	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "LABA (RUGI) BERSIH")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), res.NetProfit)

	return f.Write(writer)
}
