//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package csl
import (
	"github.com/abieberbach/goplane/extra/logging"
	"fmt"
)

type CslPackage struct {
	Name         string
	Dependencies []string
	BaseDirectory string
	Aircrafts    []*CslAircraft
	Valid        bool
}

func (self *CslPackage) validate(allPackages *CslPackages) {
	//1. prüfen ob Abhängig
	self.checkDependencies(allPackages)
	if !self.Valid {
		return
	}
	//2. prüfen ob alle benötigten ACF-Dateien vorhanden sind
	self.checkAircrafts(allPackages)
	if !self.Valid {
		return
	}
	logging.Infof("package \"%v\" is valid", self.Name)
}

func (self *CslPackage) checkDependencies(allPackages *CslPackages) {
	for _, dep := range self.Dependencies {
		_, found := allPackages.GetPackage(dep)
		if !found {
			self.invalidate("missing dependency package \"%v\"", dep)
		}
	}
}

func (self *CslPackage) checkAircrafts(allPackages *CslPackages) {
	for _, currentAircraft := range self.Aircrafts {
		ok,validationMessage := currentAircraft.validate(allPackages)
		if !ok {
			self.invalidate(validationMessage)
		}
	}
}

func (self *CslPackage) invalidate(msg string, params... interface{}) {
	logging.Warningf("invalid package \"%v\", reason: %v [base directory: %v]", self.Name, fmt.Sprintf(msg, params...),self.BaseDirectory)
	self.Valid = false
}