package pdf

import (
	"github.com/jung-kurt/gofpdf"
	"time"
)

type MyMoveForm struct {
	pdf  *gofpdf.Fpdf
	data PPMData
}

type PPMData struct {
	Name    string
	Phone   string
	DODID   string
	Branch  string
	Rank    string
	Email   string
	Address string
}

type formField struct {
	label string
	value string
}

const horizontalMargin = 15.0
const topMargin = 10.0
const totalWidth = 210.0
const bodyWidth = totalWidth - (horizontalMargin * 2)
const fieldHeight = 5.0
const fontFace = "Helvetica"

func NewMyMoveForm(data PPMData) *MyMoveForm {
	pdf := gofpdf.New("P", "mm", "A4", "")
	return &MyMoveForm{pdf: pdf, data: data}
}

func (m *MyMoveForm) GeneratePPM() error {
	m.pdf.SetMargins(horizontalMargin, topMargin, horizontalMargin)

	m.pdf.SetHeaderFunc(func() {
		m.pdf.SetFont(fontFace, "B", 17)
		m.pdf.Cell(bodyWidth*0.75, fieldHeight*2, "SHIPMENT SUMMARY WORKSHEET - PPM")
		m.setFieldLabelFont()
		m.pdf.Cell(bodyWidth*0.25, fieldHeight, "Date Prepared (YYYYMMDD)")
		m.setFieldValueFont()
		m.pdf.SetXY(horizontalMargin+bodyWidth*0.75, m.pdf.GetY()+5)
		m.pdf.Cell(bodyWidth*0.25, fieldHeight, time.Now().Format("2006-01-02"))
		m.pdf.SetLineWidth(1.0)
		m.pdf.Line(0, 20, totalWidth, 20)
		m.pdf.Ln(-1)
	})

	m.pdf.AddPage()
	m.addSectionHeader("MEMBER OR EMPLOYEE INFORMATION")

	// Name + phone
	nameFieldWidth := bodyWidth * 0.75
	nameLabelWidth := 40.0
	m.setFieldLabelFont()
	m.pdf.Cell(nameLabelWidth, fieldHeight, "Name")
	m.pdf.Ln(-1)
	m.pdf.SetFont(fontFace, "", 8)
	m.pdf.Cell(nameLabelWidth, fieldHeight, "(Last, First, Middle Initial)")

	m.pdf.SetXY(m.pdf.GetX(), m.pdf.GetY()-fieldHeight)
	m.pdf.SetFont(fontFace, "", 18)
	m.pdf.Cell(nameFieldWidth-nameLabelWidth, fieldHeight*2, m.data.Name)

	m.setFieldLabelFont()
	m.pdf.Cell(bodyWidth-nameFieldWidth, fieldHeight, "Preferred Phone Number")
	m.pdf.SetXY(nameFieldWidth+horizontalMargin, m.pdf.GetY()+fieldHeight)
	m.setFieldValueFont()
	m.pdf.Cell(bodyWidth-nameFieldWidth, fieldHeight, m.data.Phone)
	m.pdf.Ln(-1)

	m.drawGrayLineFull(2)

	// More stuff
	row := []formField{
		formField{label: "DoD ID", value: m.data.DODID},
		formField{label: "Service Branch/Agency", value: m.data.Branch},
		formField{label: "Rank/Grade", value: m.data.Rank},
		formField{label: "Preferred Email", value: m.data.Email},
	}
	m.addFormRow(row, bodyWidth)
	m.drawGrayLineFull(2)

	// Address
	m.setFieldLabelFont()
	m.pdf.Cell(bodyWidth*0.3, fieldHeight, "Preferred W2 Mailing Address")
	m.setFieldValueFont()
	m.pdf.Cell(bodyWidth*0.7, fieldHeight, m.data.Address)
	m.pdf.Ln(-1)

	// Not the right data
	m.addSectionHeader("ORDERS/ACCOUNTING INFORMATION")
	m.addFormRow(row, bodyWidth)

	m.addSectionHeader("ENTITLEMENTS/MOVE SUMMARY")
	y := m.pdf.GetY()
	entitlements := []formField{
		formField{label: "Entitlement", value: "12321 lbs"},
		formField{label: "Pro-Gear", value: "12321 lbs"},
		formField{label: "Spouse Pro-Gear", value: "12321 lbs"},
		formField{label: "Total Weight", value: "12321 lbs"},
	}
	m.addTable("Maximum Weight Entitlement", entitlements, bodyWidth*0.46, fieldHeight)

	middleX := totalWidth * 0.5
	m.pdf.SetXY(middleX, y)
	row = []formField{
		formField{label: "Authorized Origin", value: "Ft. Bragg"},
		formField{label: "Authorized Destination", value: "Pentagon"},
	}
	m.addFormRow(row, bodyWidth*0.5)
	m.drawGrayLine(2, middleX, totalWidth-horizontalMargin)
	m.pdf.SetX(middleX)
	m.addFormRow(row, bodyWidth*0.5)
	m.drawGrayLine(2, middleX, totalWidth-horizontalMargin)

	m.addSectionHeader("FINANCE/PAYMENT")

	return m.pdf.OutputFileAndClose("ppm.pdf")
}

func (m *MyMoveForm) addSectionHeader(title string) {
	m.pdf.Ln(2)
	m.pdf.SetFont(fontFace, "B", 10)
	m.pdf.SetFillColor(221, 231, 240)
	m.pdf.CellFormat(0, 7, title, "", 1, "L", true, 0, "")
	m.pdf.Ln(1)
}

func (m *MyMoveForm) setFieldLabelFont() {
	m.pdf.SetFont(fontFace, "B", 10)
}

func (m *MyMoveForm) setFieldValueFont() {
	m.pdf.SetFont(fontFace, "B", 11)
}

func (m *MyMoveForm) drawGrayLineFull(margin float64) {
	m.drawGrayLine(margin, horizontalMargin, totalWidth-horizontalMargin)
}

func (m *MyMoveForm) drawGrayLine(margin, x1, x2 float64) {
	m.pdf.SetDrawColor(221, 231, 240)
	m.pdf.SetLineWidth(0.2)
	m.pdf.Ln(margin)
	m.pdf.Line(x1, m.pdf.GetY(), x2, m.pdf.GetY())
	m.pdf.Ln(margin)
}

func (m *MyMoveForm) addFormRow(fields []formField, width float64) {
	x := m.pdf.GetX()
	fieldWidth := width / float64(len(fields))

	// Add labels
	m.setFieldLabelFont()
	for _, field := range fields {
		m.pdf.Cell(fieldWidth, fieldHeight, field.label)
	}
	m.pdf.Ln(-1)
	m.pdf.SetX(x)

	// Add values
	m.setFieldValueFont()
	for _, field := range fields {
		m.pdf.Cell(fieldWidth, fieldHeight, field.value)
	}
	m.pdf.Ln(-1)
	m.pdf.SetX(x)
}

func (m *MyMoveForm) addTable(header string, fields []formField, width, cellHeight float64) {
	m.pdf.SetFont(fontFace, "B", 10)
	m.pdf.CellFormat(width, cellHeight, header, "1", 1, "", false, 0, "")

	m.pdf.SetFont(fontFace, "", 10)
	for _, field := range fields {
		m.pdf.CellFormat(width/2, cellHeight, field.label, "LTB", 0, "", false, 0, "")
		m.pdf.CellFormat(width/2, cellHeight, field.value, "TRB", 1, "", false, 0, "")
	}
}
