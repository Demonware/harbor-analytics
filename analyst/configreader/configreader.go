// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analyst
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

package configreader

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/demonware/harbor-analytics/analyst/registry"

	yaml "gopkg.in/yaml.v2"
)

const configFile = "../analyst.yaml"

/*AnalystConfig is the representation of the analyst.yaml confi file*/
type AnalystConfig struct {
	Charts []map[string]interface{} `yaml:"charts"`
}

const (
	chartTitleTemplateConfigParameter = "titleTemplate"
	statsMethodNameConfigParameter    = "statsMethodName"
	timePeriodInDatsConfigParameter   = "timePeriodInDays"

	configPlaceholderStartDate = "startDate"
	timeFormat                 = "2006-01-02"
)

func parseConfigFile() *AnalystConfig {
	fullConfig := AnalystConfig{}
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to access config file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &fullConfig)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
	}

	return &fullConfig
}

/*getStartDateFromPeriod returns the date of
the current date minus the given time period given in days*/
func getStartDateFromPeriod(timePeriodInDate int) time.Time {
	//Go back 0 years, 0 months and <timePeriodInDate> days
	return time.Now().AddDate(0, 0, -1*timePeriodInDate)
}

/*formatTemplateString will turn the placeholder strings in a given
templateString into concrete values and return the new string.*/
func formatTemplateString(templateString string, startDate time.Time) string {

	genericTemplateFormat := "{{ %s }}"

	var concreteString string

	//Replace startDate placeholder
	concreteString = strings.Replace(
		templateString,
		fmt.Sprintf(genericTemplateFormat, configPlaceholderStartDate),
		startDate.Format(timeFormat), -1)

	return concreteString
}

/*GetStatsMethodsFromConfig is the entrypoint for converting
the analyst config file into references to executable statistical methods
of the registry.Registry type.
The returned method references can be executed by running their
Call() method.
See registryreflector.GetAllChartStatsMethods() for more info.
*/
func GetStatsMethodsFromConfig(registry registry.Registry) []ChartStatsMethod {
	return getAllChartStatsMethods(registry, parseConfigFile(), getStartDateFromPeriod)
}
