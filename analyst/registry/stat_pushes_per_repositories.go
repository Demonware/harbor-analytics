// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analytics
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

package registry

import (
	"log"
	"sort"
	"time"

	"github.com/demonware/harbor-analytics/analyst/outputgen"
)

/*PushesPerRepositories is a slice of pushesPerRepository structs
containing the name of a repository and the number
of pushes to it.*/
type PushesPerRepositories struct {
	data  []pushesPerRepository
	title string
}

type pushesPerRepository struct {
	repositoryName string
	pushCount      int
}

/*GetOrderedBarChartValues for the PushesPerRepositories type converts a
PushesPerRepositories slice into a map that can be used by a chart generator.
The output slice is guaranteed to be ordered.
This method is a requirement of the outputgen.BarChartable interface.*/
func (p PushesPerRepositories) GetOrderedBarChartValues() outputgen.BarChartableValuesOrdered {

	var chartables outputgen.BarChartableValuesOrdered

	for _, pushesPerRepository := range p.data {
		chartables = append(chartables, outputgen.BarChartableValue{
			Label: pushesPerRepository.repositoryName,
			Value: pushesPerRepository.pushCount,
		})
	}

	return chartables

}

/*SetTitle sets the human-readable title of the statistics type.
This method is a requirement of the outputgen.BarChartable interface.*/
func (p *PushesPerRepositories) SetTitle(title string) {
	log.Printf("\nPushesPerRepositories :: %v .SetTitle %s", p, title)
	p.title = title
	log.Printf("\nNewTitle::%s", p.Title())
}

/*Title returns the human-readable title of the statistics type
and will raise an error if the title is emtystring.
This method is a requirement of the outputgen.BarChartable interface.*/
func (p *PushesPerRepositories) Title() string {
	if len(p.title) < 1 {
		log.Fatalf("\nTitle of chartable %v is empty. Abort.", p)
	}
	return p.title
}

/*GetMostPushedToRepositoriesParameters is the type
that provides a wrapper for the parameters passed to the
GetMostPushedToRepositories stats function.
This type implements the StatsMethodParameters interface type.*/
type GetMostPushedToRepositoriesParameters struct {
	startDate            time.Time
	MaxNumberOfElements  int
	RepositoriesToIgnore []string
}

/*SetStartDate sets the startDate parameter in the parameters struct,
in this case: GetMostPushedToRepositoriesParameters.
This method is required by the StatsMethodParameters interface.*/
func (g *GetMostPushedToRepositoriesParameters) SetStartDate(startDate time.Time) {
	log.Printf("\nGetMostPushedToRepositoriesParameters.SetStartDate to %s", startDate)
	g.startDate = startDate
}

/*StartDate returns the startDate parameter of the parameters struct,
in this case: GetMostPushedToRepositoriesParameters.
This method is required by the StatsMethodParameters interface.*/
func (g *GetMostPushedToRepositoriesParameters) StartDate() time.Time {
	return g.startDate
}

/*IsValid check whether all fields in the GetMostPushedToRepositoriesParameters
have a valid value. If not valid, false and a reason string is returned.
This method is required by the StatsMethodParameters interface.*/
func (g *GetMostPushedToRepositoriesParameters) IsValid() (bool, string) {
	if g.MaxNumberOfElements < 1 {
		return false, "MaxNumberOfElements is less than one"
	}
	return true, ""
}

/*GetMostPushedToRepositoriesWrapper is a wrapper method of the
GetMostPushedToRepositories methods. In contrast to the concrete
GetMostPushedToRepositories method, it accept interface types
which will then be converted to concrete types and pasded to the
concrete method.

Wrappers are a workaround to make statistical methods of the
registry type generically accessible. The wrappers are accessed in the
configreader module when configuration is mapped to the methods here.
It would not be possibls to access the concrete methods (like
GetMostPushedToRepositories) without knowing the concrete parameter and
output types.*/
func (registry *Registry) GetMostPushedToRepositoriesWrapper(paramsGeneric StatsMethodParameters) outputgen.BarChartable {
	return registry.GetMostPushedToRepositories(paramsGeneric.(*GetMostPushedToRepositoriesParameters))
}

/*GetMostPushedToRepositories generates a struct containing the
<MaxNumberOfElements> most pushed-to repositories (since <StartDate>)
according to the given CSV data.

Each struct within the list of structs in the data field of the returned
PushesPerRepositories struct contains
the name of the repository and the number of pushes that
have been performed to it ever since <StartDate>.
Repositories matching a name in the given list of repositoriesToIgnore
will not be included in the returned structure.
Check the GetMostPushedToRepositoriesParameters struct for parameters.
This method is wrapped by GetMostPushedToRepositoriesWrapper.*/
func (registry *Registry) GetMostPushedToRepositories(params *GetMostPushedToRepositoriesParameters) *PushesPerRepositories {

	log.Printf("\nAnalyse :: GetMostPushedToRepositories :: %v", params)
	if isValid, reason := params.IsValid(); !isValid {
		log.Fatalf("\nGetMostPushedToRepositories :: params are invalid :: %s", reason)
	}

	var allPushesPerRepositories PushesPerRepositories

	//Go through all repositories and sum up the pushes
	//performed to any tag in the repository
	for _, project := range registry.Projects {
		for _, repository := range project.Repositories {

			skipRepository := false
			for _, repositoryNameToIgnore := range params.RepositoriesToIgnore {
				if repositoryNameToIgnore == repository.Name {
					skipRepository = true
					log.Printf("\nIgnore repository %s.", repository.Name)
					break
				}
			}

			if skipRepository {
				continue
			}

			totalPushes := 0
			for _, tag := range repository.Tags {
				pushes := 0
				for _, push := range tag.Pushes {
					if push.Timestamp.Before(params.StartDate()) {
						log.Printf("\nIgnore push to %s on %s as before relevant time.", repository.Name, push.Timestamp)
						continue
					}
					pushes++
				}
				totalPushes += pushes
			}
			allPushesPerRepositories.data = append(allPushesPerRepositories.data, pushesPerRepository{
				repositoryName: repository.Name,
				pushCount:      totalPushes,
			})
		}
	}

	//Sort the elements in the data slace by pushCount descendingly
	sort.Slice(allPushesPerRepositories.data, func(idxA, idxB int) bool {
		return allPushesPerRepositories.data[idxA].pushCount > allPushesPerRepositories.data[idxB].pushCount
	})

	//Trim the output slice to the size defined in MaxNumberOfElements
	//which will determine the number of bars shown in the chart
	if int(params.MaxNumberOfElements) < len(allPushesPerRepositories.data) {
		allPushesPerRepositories.data = allPushesPerRepositories.data[:params.MaxNumberOfElements]
	}

	return &allPushesPerRepositories

}
