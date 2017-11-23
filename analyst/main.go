// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analyst
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

package main

import (
	"log"
	"os"

	"github.com/Demonware/harbor-analytics/analyst/configreader"
	"github.com/Demonware/harbor-analytics/analyst/outputgen"
	"github.com/Demonware/harbor-analytics/analyst/parser"
)

func main() {

	registry, err := parser.CSVsToRegistry()
	if err != nil {
		log.Fatal(err.Error())
	}

	//Create output directory if non-exist
	if _, err := os.Stat(outputgen.OutDir); os.IsNotExist(err) {
		os.Mkdir(outputgen.OutDir, outputgen.OutDirMode)
	}

	var chartsPaths []string
	chartStatsFunctions := configreader.GetStatsMethodsFromConfig(registry)
	for _, chartStatsFunction := range chartStatsFunctions {
		chartPath, err := outputgen.BuildBarChart(chartStatsFunction.Call())
		if err != nil {
			log.Fatalf("\nFailed to generate chart :: %s", err.Error())
		}
		chartsPaths = append(chartsPaths, chartPath)
	}

	outputgen.BuildPDF(outputgen.PDFSection{
		Title:       "",
		Description: "",
		ChartFiles:  chartsPaths,
	})

}
