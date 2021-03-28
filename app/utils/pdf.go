package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf-examples/project/booking/app/models"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

const licenseKey = `
-----BEGIN UNIDOC LICENSE KEY-----
Free trial license keys are available at: https://unidoc.io/
-----END UNIDOC LICENSE KEY-----
`

const companyName = "Company Name"

func init() {
	err := license.SetLicenseKey(licenseKey, companyName)
	if err != nil {
		panic(err)
	}
}

func GeneratePDF(booking models.Booking) (string, error) {
	c := creator.New()
	c.SetPageMargins(50, 50, 50, 50)

	// Create report fonts.
	// UniPDF supports a number of font-families, which can be accessed using model.
	font, err := model.NewStandard14Font(model.HelveticaName)
	if err != nil {
		return "", err
	}

	fontBold, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		return "", err
	}

	title := c.NewStyledParagraph()
	title.SetTextAlignment(creator.TextAlignmentCenter)
	tc := title.SetText("Hotel Booking Information")
	tc.Style.Font = fontBold

	if err := c.Draw(title); err != nil {
		return "", err
	}

	table := c.NewTable(2)
	table.SetMargins(0, 0, 20, 0)

	if err := table.SetColumnWidths(0.2, 0.8); err != nil {
		return "", err
	}

	drawCell := func(text string, font *model.PdfFont, align creator.CellHorizontalAlignment) {
		if text == "" {
			text = "-"
		}

		p := c.NewStyledParagraph()
		p.Append(text).Style.Font = font
		p.SetMargins(0, 0, 0, 10)

		cell := table.NewCell()
		cell.SetHorizontalAlignment(align)
		_ = cell.SetContent(p)
	}

	address := fmt.Sprintf("%s\n\n%s, %s", booking.Hotel.Address, booking.Hotel.City, booking.Hotel.State)
	smoking := "No"

	if booking.Smoking {
		smoking = "Yes"
	}

	drawCell("Name", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell(booking.User.Name, font, creator.CellHorizontalAlignmentLeft)

	drawCell("Hotel", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell(booking.Hotel.Name, font, creator.CellHorizontalAlignmentLeft)

	drawCell("Address", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell(address, font, creator.CellHorizontalAlignmentLeft)

	drawCell("Beds", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell(fmt.Sprintf("%d bed", booking.Beds), font, creator.CellHorizontalAlignmentLeft)

	drawCell("Nights", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell(fmt.Sprintf("%d nights", booking.Nights()), font, creator.CellHorizontalAlignmentLeft)

	drawCell("Smoking Area", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell(smoking, font, creator.CellHorizontalAlignmentLeft)

	drawCell("Check In", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell(booking.CheckInStr, font, creator.CellHorizontalAlignmentLeft)

	drawCell("Check Out", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell(booking.CheckOutStr, font, creator.CellHorizontalAlignmentLeft)

	c.DrawHeader(func(block *creator.Block, args creator.HeaderFunctionArgs) {
		desc := c.NewStyledParagraph()
		desc.SetText("Hotel Booking Information")
		desc.SetMargins(50, 0, 0, 0)

		_ = block.Draw(desc)
	})

	if err := c.Draw(table); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s_%s.pdf", strings.ReplaceAll(booking.Hotel.Name, " ", "_"),
		strings.ReplaceAll(booking.User.Name, " ", "_"))
	outputPath := filepath.Join(os.TempDir(), fileName)

	if err := c.WriteToFile(outputPath); err != nil {
		return "", err
	}

	return outputPath, nil
}
