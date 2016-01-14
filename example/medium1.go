package example

import (
	"fmt"
	gr "github.com/mikeshimura/goreport"
	//"io/ioutil"
	"strconv"
	//"strings"
)

func Medium1Report() {
	r := gr.CreateGoReport()
	//Page Total Function
	r.PageTotal = true
	r.SumWork["amountcum"] = 0.0
	r.SumWork["g1item"] = 0.0
	r.SumWork["g1cum"] = 0.0
	r.SumWork["g2cum"] = 0.0
	font1 := gr.FontMap{
		FontName: "IPAexゴシック",
		FileName: "ttf//ipaexg.ttf",
	}
	fonts := []*gr.FontMap{&font1}
	r.SetFonts(fonts)
	d := new(M1Detail)
	r.RegisterBand(gr.Band(*d), gr.Detail)
	h := new(M1Header)
	r.RegisterBand(gr.Band(*h), gr.PageHeader)
	s := new(M1Summary)
	r.RegisterBand(gr.Band(*s), gr.Summary)
	s1 := new(M1G1Summary)
	r.RegisterGroupBand(gr.Band(*s1), gr.GroupSummary, 1)
	s2 := new(M1G2Summary)
	r.RegisterGroupBand(gr.Band(*s2), gr.GroupSummary, 2)
	r.Records = ReadText()
	fmt.Printf("Records %v \n", r.Records)
	r.SetPage("A4", "mm", "L")
	r.FooterY = 190
	r.Execute("medium1.pdf")
}

type M1Detail struct {
}

func (h M1Detail) GetHeight(report gr.GoReport) float64 {
	return 10
}
func (h M1Detail) Execute(report gr.GoReport) {
	cols := report.Records[report.DataPos].([]string)
	report.Font("IPAexゴシック", 12, "")
	y := 2.0
	report.Cell(15, y, cols[0])
	report.Cell(30, y, cols[1])
	report.Cell(60, y, cols[2])
	report.Cell(90, y, cols[3])
	report.Cell(120, y, cols[4])
	report.CellRight(135, y, 25, cols[5])
	report.CellRight(160, y, 20, cols[6])
	amt := ParseFloatNoError(cols[5]) * ParseFloatNoError(cols[6])
	report.SumWork["g1item"] += 1.0
	report.SumWork["amountcum"] += amt
	report.SumWork["g1cum"] += amt
	report.SumWork["g2cum"] += amt
	report.CellRight(180, y, 30, strconv.FormatFloat(amt, 'f', 2, 64))
}
func (h M1Detail) BreakCheckBefore(report gr.GoReport) int {
	if report.DataPos == 0 {
		//max no
		return 2
	}
	curr := report.Records[report.DataPos].([]string)
	before := report.Records[report.DataPos-1].([]string)
	if curr[0] != before[0] {
		return 2
	}
	if curr[2] != before[2] {
		return 1
	}
	return 0
}
func (h M1Detail) BreakCheckAfter(report gr.GoReport) int {
	if report.DataPos == len(report.Records)-1 {
		//max no
		return 2
	}
	curr := report.Records[report.DataPos].([]string)
	after := report.Records[report.DataPos+1].([]string)
	if curr[0] != after[0] {
		return 2
	}
	if curr[2] != after[2] {
		return 1
	}
	return 0
}

type M1Header struct {
}

func (h M1Header) GetHeight(report gr.GoReport) float64 {
	return 30
}
func (h M1Header) Execute(report gr.GoReport) {
	report.Font("IPAexゴシック", 14, "")
	report.Cell(50, 15, "Sales Report")
	report.Font("IPAexゴシック", 12, "")
	report.Cell(225, 20, "page")
	report.CellRight(230, 20, 15, strconv.Itoa(report.Page))
	report.Cell(250, 20, "of")
	report.CellRight(255, 20, 15, "{#TotalPage#}")
	y := 23.0
	report.Cell(15, y, "D No")
	report.Cell(30, y, "Dept")
	report.Cell(60, y, "Order")
	report.Cell(90, y, "Stock")
	report.Cell(120, y, "Name")
	report.CellRight(135, y, 25, "Unit Price")
	report.CellRight(160, y, 20, "Qty")
	report.CellRight(190, y, 20, "Amount")
}

type M1G1Summary struct {
}

func (h M1G1Summary) GetHeight(report gr.GoReport) float64 {
	//Conditional print  if item==1 not print
	if report.SumWork["g1item"] == 1.0 {
		fmt.Println("return 0")
		return 0
	}
	fmt.Println("return 10")
	return 10
}
func (h M1G1Summary) Execute(report gr.GoReport) {
	//Conditional print  if item==1 not print
	if report.SumWork["g1item"] != 1.0 {
		report.Cell(80, 2, "Item")
		report.CellRight(100, 2, 10, strconv.FormatFloat(
			report.SumWork["g1item"], 'f', 0, 64))
		report.Cell(150, 2, "Order Total")
		report.CellRight(180, 2, 30, strconv.FormatFloat(
			report.SumWork["g1cum"], 'f', 2, 64))
	}
	report.SumWork["g1item"] = 0.0
	report.SumWork["g1cum"] = 0.0
}

type M1G2Summary struct {
}

func (h M1G2Summary) GetHeight(report gr.GoReport) float64 {
	return 10
}
func (h M1G2Summary) Execute(report gr.GoReport) {
	report.Cell(150, 2, "Dept Total")
	report.CellRight(180, 2, 30, strconv.FormatFloat(
		report.SumWork["g2cum"], 'f', 2, 64))
	report.SumWork["g2cum"] = 0.0
	//Force New Page
	fmt.Println("report.NewPage")
	report.NewPage(false)
}

type M1Summary struct {
}

func (h M1Summary) GetHeight(report gr.GoReport) float64 {
	return 10
}
func (h M1Summary) Execute(report gr.GoReport) {
	report.Cell(160, 2, "Total")
	report.CellRight(180, 2, 30, strconv.FormatFloat(
		report.SumWork["amountcum"], 'f', 2, 64))
}
