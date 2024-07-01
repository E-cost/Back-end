package files

import (
	"Ecost/internal/order/model"
	"Ecost/internal/utils/types"
	"fmt"
	"strconv"

	"github.com/jung-kurt/gofpdf/v2"
)

func GenerateEnPDF(
	pdf *gofpdf.Fpdf,
	imagePath string,
	newOrder *model.Order,
	contactInfo *types.ContactInfo,
	orderId string,
	totalSum float64) {

	pdf.AddPage()
	pdf.SetMargins(10, 10, 10)

	opt := gofpdf.ImageOptions{
		ReadDpi:   true,
		ImageType: "PNG",
	}

	imageWidth := 60.0
	imageHeight := 60.0
	pdf.ImageOptions(imagePath, 10, 10, imageWidth, imageHeight, false, opt, 0, "")

	pdf.SetFont("DejaVu", "", 18)
	pdf.SetXY(0, imageHeight)
	pdf.CellFormat(210, 10, "Ordering information", "", 0, "C", false, 0, "")

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.SetXY(10, imageHeight+20)
	pdf.Cell(40, 10, "Order ID:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 10, orderId)
	pdf.Ln(5)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(40, 10, "Name:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 10, contactInfo.Name)
	pdf.Ln(5)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(40, 10, "Surname:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 10, contactInfo.Surname)
	pdf.Ln(5)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(40, 10, "Email:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 10, contactInfo.Email)
	pdf.Ln(5)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(40, 10, "Phone number:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 10, contactInfo.ContactPhone)
	pdf.Ln(5)

	pdf.SetXY(10, imageHeight+40)
	pdf.SetFont("DejaVuBold", "", 9)
	pdf.Ln(12)

	colWidths := []float64{30, 30, 45, 30, 20, 30}

	tableWidth := 0.0
	for _, width := range colWidths {
		tableWidth += width
	}
	startX := (210 - tableWidth) / 2
	pdf.SetX(startX)

	pdf.SetX(startX)
	headersEN := []string{"Category", "Product", "Name", "Net Weight", "Qty", "Price"}
	for i, header := range headersEN {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("DejaVu", "", 8)
	for _, item := range newOrder.Items {
		pdf.SetX(startX)
		pdf.CellFormat(colWidths[0], 10, item.Category, "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[1], 10, item.Product, "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[2], 10, item.Name, "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[3], 10, item.NetWeight, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[4], 10, strconv.Itoa(item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[5], 10, fmt.Sprintf("%.2f", item.PriceBYN), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(2)
	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(40, 10, "Total:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.CellFormat(0, 10, fmt.Sprintf("%.2f BYN", totalSum), "", 1, "R", false, 0, "")

	pdf.Ln(10)
	pdf.SetFont("DejaVuBold", "", 12)
	pdf.Cell(0, 10, "Terms&Conditions")
	pdf.Ln(10)

	pdf.SetFont("DejaVu", "", 10)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(5, 5, "1.")
	pdf.SetFont("DejaVu", "", 10)
	pdf.MultiCell(0, 5, " Order placement is made by strict prepayment", "", "L", false)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(5, 5, "2.")
	pdf.SetFont("DejaVu", "", 10)
	pdf.MultiCell(0, 5, " The waiting time for the order should not exceed 1 month from the date of payment", "", "L", false)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(5, 5, "3.")
	pdf.SetFont("DejaVu", "", 10)
	pdf.MultiCell(0, 5, " The company takes full responsibility for logistics and delivery of goods", "", "L", false)

}
