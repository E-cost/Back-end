package files

import (
	"Ecost/internal/order/model"
	"Ecost/internal/utils/types"
	"fmt"
	"os"
	"strconv"

	"github.com/jung-kurt/gofpdf/v2"
)

func GenerateRusPDF(
	pdf *gofpdf.Fpdf,
	imagePath string,
	newOrder *model.Order,
	contactInfo types.ContactInfo,
	shortId string,
	totalSum float64) {

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Println("Error: Image file does not exist") // TODO: to add a logger
		return
	}

	pdf.AddPage()
	pdf.SetMargins(10, 10, 10)

	opt := gofpdf.ImageOptions{
		ReadDpi:   true,
		ImageType: "PNG",
	}

	imageWidth := 45.0
	imageHeight := 30.0
	pdf.ImageOptions(imagePath, 10, 10, imageWidth, imageHeight, false, opt, 0, "")

	pdf.SetFont("DejaVu", "", 18)
	pdf.SetXY(0, imageHeight)
	pdf.CellFormat(210, 10, "Информация для заказа", "", 0, "C", false, 0, "")

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.SetXY(10, imageHeight+20)
	pdf.Cell(40, 10, "Номер заказа:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 10, shortId)
	pdf.Ln(5)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(40, 10, "Клиент:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 10, contactInfo.Name+" "+contactInfo.Surname)
	pdf.Ln(5)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(40, 10, "Почта:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 10, contactInfo.Email)
	pdf.Ln(5)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(40, 10, "Телефон:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 10, contactInfo.ContactPhone)
	pdf.Ln(5)

	pdf.SetXY(10, imageHeight+40)
	pdf.SetFont("DejaVuBold", "", 9)
	pdf.Ln(10)

	colWidths := []float64{30, 30, 45, 30, 20, 30}

	tableWidth := 0.0
	for _, width := range colWidths {
		tableWidth += width
	}
	startX := (210 - tableWidth) / 2
	pdf.SetX(startX)

	headersRUS := []string{"№", "Артикул", "Название", "Масса нетто", "Кол-во", "Цена"}
	for i, header := range headersRUS {
		pdf.CellFormat(colWidths[i], 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	counter := 0

	pdf.SetFont("DejaVu", "", 8)
	for _, item := range newOrder.Items {
		counter++
		pdf.SetX(startX)
		pdf.CellFormat(colWidths[0], 8, strconv.Itoa(counter), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[1], 8, item.Article, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[2], 8, item.Name, "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[3], 8, item.NetWeight, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[4], 8, strconv.Itoa(item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[5], 8, fmt.Sprintf("%.2f", item.PriceBYN), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(2)
	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(40, 10, "Сумма:")
	pdf.SetFont("DejaVu", "", 10)
	pdf.CellFormat(0, 10, fmt.Sprintf("%.2f BYN", totalSum), "", 1, "R", false, 0, "")

	pdf.Ln(10)
	pdf.SetFont("DejaVuBold", "", 12)
	pdf.Cell(0, 10, "Условия")
	pdf.Ln(10)

	pdf.SetFont("DejaVu", "", 10)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(5, 5, "1.")
	pdf.SetFont("DejaVu", "", 10)
	pdf.MultiCell(0, 5, " Оформление заказа производится по строгой предоплате", "", "L", false)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(5, 5, "2.")
	pdf.SetFont("DejaVu", "", 10)
	pdf.MultiCell(0, 5, " Ожидание заказа не должно превысить 1 месяц с дня оплаты", "", "L", false)

	pdf.SetFont("DejaVuBold", "", 10)
	pdf.Cell(5, 5, "3.")
	pdf.SetFont("DejaVu", "", 10)
	pdf.MultiCell(0, 5, " Компания берет всю ответственность за логистику и доставку товара", "", "L", false)
}
