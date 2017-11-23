// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analytics
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

package parser

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/demonware/harbor-analytics/analyst/registry"
)

const (
	projectCSV    = "../raw/project.csv"
	accessLogCSV  = "../raw/access_log.csv"
	repositoryCSV = "../raw/repository.csv"
	userCSV       = "../raw/user.csv"
)

type csvFieldDescription []string

var projectCSVFields = csvFieldDescription{
	"project_id",
	"owner_id",
	"name",
	"deleted",
	"public",
}

var accessLogCSVFields = csvFieldDescription{
	"log_id",
	"user_id",
	"project_id",
	"repo_name",
	"repo_tag",
	"operation",
	"op_time",
}

var repositoryCSVFields = csvFieldDescription{
	"repository_id",
	"name",
	"project_id",
	"owner_id",
}

var userCSVFields = csvFieldDescription{
	"user_id",
	"username",
}

const (
	pullOperation   = "pull"
	pushOperation   = "push"
	deleteOperation = "delete"
	createOperation = "create"
)

func readCSV(csvFile string, csvFieldDescription csvFieldDescription) ([]map[string]string, error) {

	openFile, err := os.Open(csvFile)
	if err != nil {
		return []map[string]string{}, err
	}

	csvReader := csv.NewReader(bufio.NewReader(openFile))
	csvReader.Comma = ','
	raw, err := csvReader.ReadAll()
	if err != nil {
		return []map[string]string{}, err
	}

	var csvData []map[string]string
	for lineIdx, lineRaw := range raw {
		if lineIdx == 0 {
			continue
		}
		lineParsed := make(map[string]string)
		for fieldIdx, fieldName := range csvFieldDescription {
			lineParsed[fieldName] = lineRaw[fieldIdx]
		}
		csvData = append(csvData, lineParsed)
	}

	return csvData, nil

}

/*CSVsToRegistry converts the raw CSV files
from exported from the harbor database to a registry.Registry struct*/
func CSVsToRegistry() (registry.Registry, error) {

	projectsParsed, err := readCSV(projectCSV, projectCSVFields)
	if err != nil {
		return registry.Registry{}, err
	}

	accessLogsParsed, err := readCSV(accessLogCSV, accessLogCSVFields)
	if err != nil {
		return registry.Registry{}, err
	}

	repositoriesParsed, err := readCSV(repositoryCSV, repositoryCSVFields)
	if err != nil {
		return registry.Registry{}, err
	}

	usersParsed, err := readCSV(userCSV, userCSVFields)
	if err != nil {
		return registry.Registry{}, err
	}

	projects := make(map[int]*registry.Project)
	for _, projectParsed := range projectsParsed {
		projectID, err := strconv.Atoi(projectParsed[projectCSVFields[0]])
		if err != nil {
			return registry.Registry{}, err
		}
		project := registry.Project{
			ID:           projectID,
			Name:         projectParsed[projectCSVFields[2]],
			Repositories: make(map[string]*registry.Repository),
		}
		projects[projectID] = &project
	}

	repositories := make(map[string]*registry.Repository)
	for _, repositoryParsed := range repositoriesParsed {
		repositoryID, err := strconv.Atoi(repositoryParsed[repositoryCSVFields[0]])
		if err != nil {
			return registry.Registry{}, err
		}
		repositoryName := repositoryParsed[repositoryCSVFields[1]]
		repository := registry.Repository{
			ID:   repositoryID,
			Name: repositoryName,
			Tags: map[string]*registry.Tag{},
		}
		repositories[repositoryName] = &repository

		projectID, err := strconv.Atoi(repositoryParsed[repositoryCSVFields[2]])
		if err != nil {
			return registry.Registry{}, err
		}

		log.Printf("Add repsitory %s to project with ID %d\n", repositoryName, projectID)
		projects[projectID].Repositories[repositoryName] = &repository
	}

	users := make(map[int]*registry.User)
	for _, userParsed := range usersParsed {
		userID, err := strconv.Atoi(userParsed[userCSVFields[0]])
		if err != nil {
			return registry.Registry{}, err
		}
		user := registry.User{
			ID:   userID,
			Name: userParsed[userCSVFields[1]],
		}
		users[userID] = &user
	}

	//Iterate over aceess logs to
	//find tags and actions performed on them
	//as well as project creation operations
	for _, accessLogParsed := range accessLogsParsed {

		tagName := accessLogParsed[accessLogCSVFields[4]]
		performedOperation := accessLogParsed[accessLogCSVFields[5]]

		if tagName == "N/A" && performedOperation != createOperation {
			log.Println("Cannot do anything with this access log. Skip.")
			continue
		}

		if performedOperation == deleteOperation {
			log.Println("Ignore delet operations for now. Skip.")
			continue
		}

		//Get the user who performed the logged operation
		var err error
		accessLogUserID, err := strconv.Atoi(accessLogParsed[accessLogCSVFields[1]])
		if err != nil {
			return registry.Registry{}, err
		}
		var ok bool
		accessLogUser, ok := users[accessLogUserID]
		if !ok {
			log.Printf("Failed to find user with ID %d\n", accessLogUserID)
		}

		//What if create operation:
		//A create operation refers to a Project being created
		//We will find the project through its ID and add the
		//creation date. Then leave the accessLog loop because
		//we don't care about tags in the case of project creation
		if performedOperation == createOperation {

			// just to get rid of a false negative go error
			// that complains about shadowing of err
			var projectID int
			projectID, err = strconv.Atoi(accessLogParsed[accessLogCSVFields[2]])
			if err != nil {
				return registry.Registry{}, err
			}

			var project *registry.Project
			if project, ok = projects[projectID]; ok {
				project.CreationDate = accessLogParsed[accessLogCSVFields[6]]
				project.Creator = accessLogUser
			} else {
				log.Printf("Failed to find project with ID %d\n", projectID)
			}

			continue

		}

		tagRepository, ok := repositories[accessLogParsed[accessLogCSVFields[3]]]
		if !ok {
			log.Printf("Could not find repository with name %s. Skip access log parsing.", accessLogParsed[accessLogCSVFields[3]])
			continue
		}

		//Check if tag already exists in our tags list
		//i.e. if we've already parsed an access log
		//if it doesn't create a new one and add it
		if _, ok := tagRepository.Tags[tagName]; !ok {
			tag := registry.Tag{
				Name:   tagName,
				Pulls:  make(map[int]*registry.Pull),
				Pushes: make(map[int]*registry.Push),
			}
			// Add the tag to the associated repository
			tagRepository.Tags[tagName] = &tag
		}

		logID, err := strconv.Atoi(accessLogParsed[accessLogCSVFields[0]])
		if err != nil {
			return registry.Registry{}, err
		}

		logTimestamp, err := time.Parse("2006-01-02 15:04:05", accessLogParsed[accessLogCSVFields[6]])
		if err != nil {
			return registry.Registry{}, err
		}

		accessLog := registry.Log{
			ID:        logID,
			Timestamp: logTimestamp,
			User:      accessLogUser,
		}

		//Switch through pull and push operations and add to tag accordingly
		switch performedOperation {
		case pullOperation:
			tagRepository.Tags[tagName].Pulls[logID] = &registry.Pull{
				Log: accessLog,
			}
		case pushOperation:
			tagRepository.Tags[tagName].Pushes[logID] = &registry.Push{
				Log: accessLog,
			}
		default:
			log.Println("Can't do anything.")
		}

	}

	return registry.Registry{Projects: projects}, nil

}
