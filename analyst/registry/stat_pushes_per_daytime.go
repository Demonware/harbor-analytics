// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analyst
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

package registry

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/demonware/harbor-analytics/analyst/outputgen"
)

/*PushesPerDaytimes is a slice of pushesPerDaytime structs
containing the hour of a day and the number
of pushes to the registry within this hour.*/
type PushesPerDaytimes struct {
	data  []pushesPerDaytime
	title string
}

type pushesPerDaytime struct {
	hourOfDay int
	pushCount int
}

/*GetOrderedBarChartValues for the PushesPerDaytimes type converts a
PushesPerDaytimes slice into a map that can be used by a chart generator.
The output slice is guaranteed to be ordered.
This method is a requirement of the outputgen.BarChartable interface.*/
func (p *PushesPerDaytimes) GetOrderedBarChartValues() outputgen.BarChartableValuesOrdered {

	var chartables outputgen.BarChartableValuesOrdered

	for _, pushesPerDaytime := range p.data {
		chartables = append(chartables, outputgen.BarChartableValue{
			Label: fmt.Sprintf("%d:00", pushesPerDaytime.hourOfDay),
			Value: pushesPerDaytime.pushCount,
		})
	}

	return chartables

}

/*SetTitle sets the human-readable title of the statistics type.
This method is a requirement of the outputgen.BarChartable interface.*/
func (p *PushesPerDaytimes) SetTitle(title string) {
	log.Printf("\nPushesPerDaytimes :: %v .SetTitle %s", p, title)
	p.title = title
	log.Printf("\nNewTitle::%s", p.Title())
}

/*Title returns the human-readable title of the statistics type
and will raise an error if the title is emtystring.
This method is a requirement of the outputgen.BarChartable interface.*/
func (p *PushesPerDaytimes) Title() string {
	if len(p.title) < 1 {
		log.Fatalf("\nTitle of chartable %v is empty. Abort.", p)
	}
	return p.title
}

/*GetPushesPerDaytimesParameters is the type
that provides a wrapper for the parameters passed to the
GetPushesPerDaytimes stats function.
This type implements the StatsMethodParameters interface type.*/
type GetPushesPerDaytimesParameters struct {
	startDate time.Time
}

/*SetStartDate sets the startDate parameter in the parameters struct,
in this case: GetPushesPerDaytimesParameters.
This method is required by the StatsMethodParameters interface.*/
func (g *GetPushesPerDaytimesParameters) SetStartDate(startDate time.Time) {
	log.Printf("\nGetPushesPerDaytimesParameters.SetStartDate to %s", startDate)
	g.startDate = startDate
}

/*StartDate returns the startDate parameter of the parameters struct,
in this case: GetPushesPerDaytimesParameters.
This method is required by the StatsMethodParameters interface.*/
func (g *GetPushesPerDaytimesParameters) StartDate() time.Time {
	return g.startDate
}

/*IsValid check whether all fields in the GetPushesPerDaytimesParameters
have a valid value. If not valid, false and a reason string is returned.
This method is required by the StatsMethodParameters interface.*/
func (g *GetPushesPerDaytimesParameters) IsValid() (bool, string) {
	return true, ""
}

/*GetPushesPerDaytimesWrapper is a wrapper method of the
GetPushesPerDaytimes methods. In contrast to the concrete
GetPushesPerDaytimes method, it accept interface types
which will then be converted to concrete types and pasded to the
concrete method.

Wrappers are a workaround to make statistical methods of the
registry type generically accessible. The wrappers are accessed in the
configreader module when configuration is mapped to the methods here.
It would not be possibls to access the concrete methods (like
GetPushesPerDaytimes) without knowing the concrete parameter and
output types.*/
func (registry *Registry) GetPushesPerDaytimesWrapper(paramsGeneric StatsMethodParameters) outputgen.BarChartable {
	return registry.GetPushesPerDaytimes(paramsGeneric.(*GetPushesPerDaytimesParameters))
}

/*GetPushesPerDaytimes generates a struct containing the
number of pushed performed for a given hour of the day
accumulated over the timeperiod since <StartDate>
according to the given CSV data.

Each struct within the list of structs in the data field of the returned
PushesPerDaytimes struct contains
the hour of a day and the number of pushes that
have been performed within this hour accumulated ever since <StartDate>.

Check the GetPushesPerDaytimesParameters struct for parameters.
This method is wrapped by GetPushesPerDaytimesWrapper.*/
func (registry *Registry) GetPushesPerDaytimes(params *GetPushesPerDaytimesParameters) *PushesPerDaytimes {

	log.Printf("\nAnalyse :: GetPushesPerDaytimes :: %v", params)
	if isValid, reason := params.IsValid(); !isValid {
		log.Fatalf("\nGetPushesPerDaytimes :: params are invalid :: %s", reason)
	}

	pushesPerDayimeMapping := map[int]int{}

	//Go through all repositories and sum up the pushes
	//performed to any tag in the repository
	for _, project := range registry.Projects {
		for _, repository := range project.Repositories {
			for _, tag := range repository.Tags {
				for _, push := range tag.Pushes {
					if push.Timestamp.Before(params.StartDate()) {
						log.Printf("\nIgnore push to %s on %s as before relevant time.", repository.Name, push.Timestamp)
						continue
					}
					pushesPerDayimeMapping[push.Timestamp.Hour()] = pushesPerDayimeMapping[push.Timestamp.Hour()] + 1
				}
			}
		}
	}

	var allPushesPerDaytimes PushesPerDaytimes
	for hourOfDay, pushCount := range pushesPerDayimeMapping {
		allPushesPerDaytimes.data = append(allPushesPerDaytimes.data, pushesPerDaytime{
			hourOfDay: hourOfDay,
			pushCount: pushCount,
		})
	}

	//Sort the elements in the data slace by pushCount descendingly
	sort.Slice(allPushesPerDaytimes.data, func(idxA, idxB int) bool {
		return allPushesPerDaytimes.data[idxA].hourOfDay < allPushesPerDaytimes.data[idxB].hourOfDay
	})

	return &allPushesPerDaytimes

}
