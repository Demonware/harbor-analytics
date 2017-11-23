// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analyst
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

package configreader

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/demonware/harbor-analytics/analyst/outputgen"
	"github.com/demonware/harbor-analytics/analyst/registry"
)

/*ChartStatsMethod defines the structure of the container
in which references to statistical methods of the regitry.Registry type
and their parameters can be stored.
The struct contains the reference to the callable wrapper method of
a concrete statistical method and a parameter of interface type
registry.StatsMethodParameters.*/
type ChartStatsMethod struct {
	methodWrapper func(registry.StatsMethodParameters) outputgen.BarChartable
	parameters    registry.StatsMethodParameters
	title         string
}

/*Call implements the call of a method stored in a ChartStatsMethod struct.
When executed, it will call the wrapper method stored in the struct
and hand the given parameter wrapper to it.
It will also set the title of the chartable.*/
func (c ChartStatsMethod) Call() outputgen.BarChartable {
	barChartable := c.methodWrapper(c.parameters)
	barChartable.SetTitle(c.title)
	return barChartable
}

/*getChartStatsMethod converts one given statistical method configuration
into an actual method reference incl. parameters wrapped in a method wrapper
of type ChartStatsMethod.

This method will perform a number of reflections on the registry package
and specifically the methods implemented on the registry.Registry type.*/
func getChartStatsMethod(concreteRegistry registry.Registry, chartConfig map[string]interface{}) ChartStatsMethod {

	callableName := chartConfig[statsMethodNameConfigParameter].(string)
	callable := reflect.ValueOf(&concreteRegistry).MethodByName(callableName)
	if !callable.IsValid() {
		log.Fatalf("\nFailed to find Registry method with name \"%s\".", callableName)
	}

	callableWrapperName := fmt.Sprintf("%sWrapper", callableName)
	callableWrapper := reflect.ValueOf(&concreteRegistry).MethodByName(callableWrapperName)
	if !callableWrapper.IsValid() {
		log.Fatalf("\nFailed to find Registry method with name \"%s\".", callableWrapperName)
	}

	callableParametersStructType := callable.Type().In(0).Elem()
	//The type of the parameter handed to a stats function will be a pointer-type
	callableParametersStructPointer := reflect.New(callableParametersStructType)
	callableParametersStruct := callableParametersStructPointer.Elem()

	for configParameter, configValue := range chartConfig {

		if configParameter == statsMethodNameConfigParameter {
			continue
		}

		if configParameter == chartTitleTemplateConfigParameter {
			//This needs to be set last since it depends
			//on other variables which might not have been
			//parsed yet.
			//Will be set after loop over config parameters
			continue
		}

		if configParameter == timePeriodInDatsConfigParameter {
			if periodInDays, isInt := configValue.(int); isInt {
				log.Printf("\nSet StartDate on parameter struct from %s", timePeriodInDatsConfigParameter)
				callableParametersStructPointer.Interface().(registry.StatsMethodParameters).SetStartDate(getStartDateFromPeriod(periodInDays))
				continue
			}
			log.Fatalf("\nParameter %s must be of type integer.", timePeriodInDatsConfigParameter)
		}

		parameterField := callableParametersStruct.FieldByName(configParameter)
		if !parameterField.IsValid() {
			log.Fatalf("Field %s in struct is invalid", configParameter)
		}

		if configValueInt, isInt := configValue.(int); isInt {
			log.Printf("\nSet Int Value %d on field %s", configValueInt, parameterField)
			parameterField.SetInt(int64(configValueInt))
			continue
		}

		if configValueString, isString := configValue.(string); isString {
			log.Printf("\nSet String Value %s on field %s", configValueString, parameterField)
			parameterField.SetString(configValueString)
			continue
		}

		if configValueList, isList := configValue.([]interface{}); isList {
			var configValueStringList []string
			for _, listElement := range configValueList {
				listElementString, isStringable := listElement.(string)
				if !isStringable {
					log.Fatalf("Value %s in configuration list of parameter %s must be a string", listElement, configParameter)
				}
				configValueStringList = append(configValueStringList, listElementString)
			}
			log.Printf("\nSet List Value %v on field %s", configValueStringList, parameterField)
			parameterField.Set(reflect.ValueOf(configValueStringList))
			continue
		}

	}

	//Set chart title
	var concreteTitle string
	if titleTemplate, isString := chartConfig[chartTitleTemplateConfigParameter].(string); isString {
		concreteTitle = formatTemplateString(titleTemplate, callableParametersStructPointer.Interface().(registry.StatsMethodParameters).StartDate())
		log.Printf("\nSet Title of stat to %s", concreteTitle)
	} else {
		log.Fatalf("\nTitle value must be of type string but is: %s", chartConfig[chartTitleTemplateConfigParameter])
	}

	return ChartStatsMethod{
		title:         concreteTitle,
		methodWrapper: callableWrapper.Interface().(func(registry.StatsMethodParameters) outputgen.BarChartable),
		parameters:    callableParametersStructPointer.Interface().(registry.StatsMethodParameters),
	}

}

/*GetAllChartStatsMethods will read the analyst.yaml configuration file
and map the inhere defined statistical method references and their parameters
to the actual implementations of the metods in the registry package.
The methods will be returned in the form of a slice of method container
structs of type ChartStatsMethod. Each method can then be invoked by
calling the .Call() method on the wrapper struct.*/
func getAllChartStatsMethods(registry registry.Registry, config *AnalystConfig, perioToDateConverter func(int) time.Time) []ChartStatsMethod {

	var chartStatsMethods []ChartStatsMethod
	for _, chartConfig := range config.Charts {
		chartStatsMethods = append(chartStatsMethods, getChartStatsMethod(registry, chartConfig))
	}

	return chartStatsMethods

}
