package helpers

import (
	"Ecost/internal/order/model"
	"Ecost/internal/order/service/helpers/files"
	"Ecost/internal/utils/types"
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf/v2"
)

func GetPDF(
	newOrder *model.Order,
	contactInfo types.ContactInfo,
	orderId string, shortId string,
	totalSum float64) ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")

	fontDir := "/usr/local/src/internal/order/service/helpers/fonts"
	pdf.AddUTF8Font("DejaVu", "", fontDir+"/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVuBold", "", fontDir+"/DejaVuSans-Bold.ttf")
	pdf.SetFont("DejaVu", "", 16)
	imagePath := "/usr/local/src/internal/order/service/helpers/logo/logo.png"

	files.GenerateRusPDF(pdf, imagePath, newOrder, contactInfo, shortId, totalSum)
	//files.GenerateEnPDF(pdf, imagePath, orderDto, contactInfo, orderId, totalSum)

	// pdfPath := fmt.Sprintf("%s.pdf", shortId)
	// err := pdf.OutputFileAndClose(pdfPath)
	// if err != nil {
	// 	return "", err
	// }

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF output: %w", err)
	}

	return buf.Bytes(), nil
}
