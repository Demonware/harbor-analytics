// Copyright (C) Activision Publishing, Inc. 2017
// https://github.com/Demonware/harbor-analyst
// Author: David Rieger
// Licensed under the 3-Clause BSD License (the "License");
// you may not use this file except in compliance with the License.

package registry

import "time"

/*StatsMethodParameters is an interface type that needs to be implemented
by every type that is to be used as a method/function parameter of a
StatsMethod i.e. a method that is used to process registry data to generate
statistical output.*/
type StatsMethodParameters interface {
	SetStartDate(time.Time)
	StartDate() time.Time
	IsValid() (bool, string)
}
