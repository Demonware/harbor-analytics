// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analyst
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

package outputgen

import (
	"fmt"
	"log"
	"time"

	"github.com/jung-kurt/gofpdf"
)

const (
	pdfOutFile     = OutDir + "/report.pdf"
	pdfOrientation = "Portrait"
	pdfUnit        = "mm" //millimeters
	pdfSize        = "A4" //DIN
	pdfFontDir     = ""
)

/*PDFSection describes a
 *standardized section of the
 *analytics report PDF to be generated.
 *Each section consitst of a title, a description
 *and a number of charts.
 *In the document each section will be presented
 *in this order, (i.e. first the title, then the description then the charts),
 *whereas the charts (within one section)
 *will be printed below each other with nothing inbetween.
 */
type PDFSection struct {
	Title       string
	Description string
	ChartFiles  []string
}

/*BuildPDF creates a PDF file
 *with the given PDFSections.
 *The resulting PDF will be written to the
 *outDir directory.
 */
func BuildPDF(sections ...PDFSection) error {

	//Create new PDF and ass one page
	pdf := gofpdf.New(pdfOrientation, pdfUnit, pdfSize, pdfFontDir)
	pdf.AddPage()
	setHeader(pdf)

	for _, section := range sections {
		addSection(pdf, section)
	}

	//Write file
	err := pdf.OutputFileAndClose(pdfOutFile)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func addSection(pdf *gofpdf.Fpdf, section PDFSection) {

	// Setting the hight to 0 means the
	// hight will be calculated so that the ratio
	// is maintained for the given 200mm wifth.
	hight := 0.0
	width := 200.0

	//Since flow is set to true, the y position will be adjusted
	//sutomatically so that a chart is placed
	//underneath the existing content
	flow := true
	yPosition := 0.0
	xPosition := 0.0

	link := 0
	linkString := ""

	for _, chartFile := range section.ChartFiles {
		pdf.ImageOptions(
			chartFile, xPosition, yPosition, width, hight, flow, gofpdf.ImageOptions{}, link, linkString)
	}
}

func setHeader(pdf *gofpdf.Fpdf) {

	titleFontName := "Arial"
	titleFontWeigth := "B"
	titleFontSize := 16.0
	pdf.SetFont(titleFontName, titleFontWeigth, titleFontSize)
	pdf.Cell(40, 10, "Harbor Analytics Report")

	//Linebreak
	pdf.Ln(0)

	subtitleFontName := "Arial"
	subtitleFontWeigth := "B"
	subtitleFontSize := 12.0

	pdf.SetFont(subtitleFontName, subtitleFontWeigth, subtitleFontSize)
	pdf.Cell(40, 20, fmt.Sprintf("Date of Document Creation: %s", time.Now().Format("2006-01-02 15:04")))
}
