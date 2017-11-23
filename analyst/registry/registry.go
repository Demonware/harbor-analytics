// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analytics
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

package registry

import "time"

/*User represnets a user of the Harbor registry.
A user is identified by a unique name and a unique ID.*/
type User struct {
	ID   int
	Name string
}

/*Log is a generic structure which represents
a log on a image tag on the Harbor registry.
This might be the log of a pull or a push action.

A log holds (apart from its ID) the time
at which the logged action was performed and
a reference to the user who performed the action.

Log structs are not held directly by any other
data structure of the inhere designed registry,
there is not reason to ever directly instantiate a
log type. Instead use the concrete types that
inherit from the Log struct and are useg in the
Tag struct, such as Push and Pull.*/
type Log struct {
	ID        int
	Timestamp time.Time
	User      *User
}

/*Pull represents a pull action of a certain docker image
from the Harbor registry.
Pull inherits the fields of the generic structure "Log"

Each Tag can hold any number of pull logs.*/
type Pull struct {
	Log
}

/*Push represents a push action of a certain docker image
to the Harbor registry.
Push inherits the fields of the generic structure "Log"

Each Tag can hold any number of push logs.*/
type Push struct {
	Log
}

/*Tag is the structure for a image tag
inside a repository of a project on the registry.
A repository can hold any number of unique tags.
A tag can hold any number of Pulls and Pushes.

Looking at a full qualified docker image name, the tag name
is the name after the colon.

e.g. for the image name "docker-registry.comany.com/coreapp/base:0.2.1"
the tag is "0.2.1".*/
type Tag struct {
	Name   string
	Pulls  map[int]*Pull
	Pushes map[int]*Push
}

/*Repository is the structure for a repository
inside a project on the registry.
A project can hold any number of unique repositories.
Each repository hold any number of unique tags.

Looking at a full qualified docker image name, the repository name
is the name after the docker registry URL and project name and before
the colon.

e.g. for the image name "docker-registry.comany.com/coreapp/base:0.2.1"
the repository is "base".*/
type Repository struct {
	ID   int
	Name string
	Tags map[string]*Tag
}

/*Project is the structure for a project
on the Harbor registry.
A registry can hold any number of unique projects.
Each project can hold any number of unique repositories.
Each project also has a creator which is a reference
to a User and a date of creation.

Looking at a full qualified docker image name, the project name
is the first name after the docker registry URL and before
the next slash.

e.g. for the image name "docker-registry.comany.com/coreapp/base:0.2.1"
the project is "coreapp".*/
type Project struct {
	ID           int
	Name         string
	CreationDate string
	Creator      *User
	Repositories map[string]*Repository
}

/*Registry represents the structure for all
the date on the Harbor docker registry.
It contains all projects of the registry mapped to
their IDs.*/
type Registry struct {
	Projects map[int]*Project
}
