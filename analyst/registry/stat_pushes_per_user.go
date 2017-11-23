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

/*PushesPerUsers is a slice of pushesPerUser structs
containing the name of a user and the number
of pushes performed by them.*/
type PushesPerUsers struct {
	data  []pushesPerUser
	title string
}

type pushesPerUser struct {
	userName  string
	pushCount int
}

/*GetOrderedBarChartValues for the PushesPerUsers type converts a
PushesPerUsers slice into a map that can be used by a chart generator.
The output slice is guaranteed to be ordered.
This method is a requirement of the outputgen.BarChartable interface.*/
func (p PushesPerUsers) GetOrderedBarChartValues() outputgen.BarChartableValuesOrdered {

	var chartables outputgen.BarChartableValuesOrdered

	for _, pushesPerUser := range p.data {
		chartables = append(chartables, outputgen.BarChartableValue{
			Label: pushesPerUser.userName,
			Value: pushesPerUser.pushCount,
		})
	}

	return chartables

}

/*SetTitle sets the human-readable title of the statistics type.
This method is a requirement of the outputgen.BarChartable interface.*/
func (p *PushesPerUsers) SetTitle(title string) {
	log.Printf("\nPushesPerUsers :: %v .SetTitle %s", p, title)
	p.title = title
	log.Printf("\nNewTitle::%s", p.Title())
}

/*Title returns the human-readable title of the statistics type
and will raise an error if the title is emtystring.
This method is a requirement of the outputgen.BarChartable interface.*/
func (p *PushesPerUsers) Title() string {
	if len(p.title) < 1 {
		log.Fatalf("\nTitle of chartable %v is empty. Abort.", p)
	}
	return p.title
}

/*GetMostPushingUsersParameters is the type
that provides a wrapper for the parameters passed to the
GetMostPushingUsers stats function.
This type implements the StatsMethodParameters interface type.*/
type GetMostPushingUsersParameters struct {
	startDate           time.Time
	MaxNumberOfElements int
	UsersToIgnore       []string
}

/*SetStartDate sets the startDate parameter in the parameters struct,
in this case: GetMostPushingUsersParameters.
This method is required by the StatsMethodParameters interface.*/
func (g *GetMostPushingUsersParameters) SetStartDate(startDate time.Time) {
	log.Printf("\nGetMostPushingUsersParameters.SetStartDate to %s", startDate)
	g.startDate = startDate
}

/*StartDate returns the startDate parameter of the parameters struct,
in this case: GetMostPushingUsersParameters.
This method is required by the StatsMethodParameters interface.*/
func (g *GetMostPushingUsersParameters) StartDate() time.Time {
	return g.startDate
}

/*IsValid check whether all fields in the GetMostPushingUsersParameters
have a valid value. If not valid, false and a reason string is returned.
This method is required by the StatsMethodParameters interface.*/
func (g *GetMostPushingUsersParameters) IsValid() (bool, string) {
	if g.MaxNumberOfElements < 1 {
		return false, "MaxNumberOfElements is less than one"
	}
	return true, ""
}

/*GetMostPushingUsersWrapper is a wrapper method of the
GetMostPushingUsers methods. In contrast to the concrete
GetMostPushingUsers method, it accept interface types
which will then be converted to concrete types and pasded to the
concrete method.

Wrappers are a workaround to make statistical methods of the
registry type generically accessible. The wrappers are accessed in the
configreader module when configuration is mapped to the methods here.
It would not be possibls to access the concrete methods (like
GetMostPushingUsers) without knowing the concrete parameter and
output types.*/
func (registry *Registry) GetMostPushingUsersWrapper(paramsGeneric StatsMethodParameters) outputgen.BarChartable {
	return registry.GetMostPushingUsers(paramsGeneric.(*GetMostPushingUsersParameters))
}

/*GetMostPushingUsers generates a struct containing the top
<MaxNumberOfElements> users who have performed the most pushes to any
repository (since <StartDate>) according to the given CSV data.
Each struct within the list of structs in the data field of the returned
PushesPerUsers struct contains the name of the user
and the number of pushes that have been performed by them ever since <StartDate>.
Users matching a name in the given list of usersToIgnore
will not be included in the returned structure.
Check the GetMostPushingUsersParameters struct for parameters.
This method is wrapped by GetMostPushingUsersWrapper.*/
func (registry *Registry) GetMostPushingUsers(params *GetMostPushingUsersParameters) *PushesPerUsers {

	log.Printf("\nAnalyse :: GetMostPushingUsers :: %v", params)
	if isValid, reason := params.IsValid(); !isValid {
		log.Fatalf("\nGetMostPushingUsers :: params are invalid :: %s", reason)
	}

	allPushesPerUsersMapping := map[string]int{}

	for _, project := range registry.Projects {
		for _, repository := range project.Repositories {
			for _, tag := range repository.Tags {
				for _, push := range tag.Pushes {
					skipUser := false
					for _, userNameToIgnore := range params.UsersToIgnore {
						if userNameToIgnore == push.User.Name {
							skipUser = true
							log.Printf("\nIgnore user %s.", push.User.Name)
							break
						}
					}
					if skipUser {
						continue
					}
					if push.Timestamp.Before(params.StartDate()) {
						log.Printf("\nIgnore push to %s on %s as before relevant time.", repository.Name, push.Timestamp)
						continue
					}
					allPushesPerUsersMapping[push.User.Name] = allPushesPerUsersMapping[push.User.Name] + 1
				}
			}
		}
	}

	var allPushesPerUsers PushesPerUsers
	for username, pushCount := range allPushesPerUsersMapping {
		allPushesPerUsers.data = append(allPushesPerUsers.data, pushesPerUser{
			userName:  username,
			pushCount: pushCount,
		})
	}

	//Sort the elements in the data slace by pushCount descendingly
	sort.Slice(allPushesPerUsers.data, func(idxA, idxB int) bool {
		return allPushesPerUsers.data[idxA].pushCount > allPushesPerUsers.data[idxB].pushCount
	})

	//Trim the output slice to the size defined in MaxNumberOfElements
	//which will determine the number of bars shown in the chart
	if int(params.MaxNumberOfElements) < len(allPushesPerUsers.data) {
		allPushesPerUsers.data = allPushesPerUsers.data[:params.MaxNumberOfElements]
	}

	return &allPushesPerUsers

}
