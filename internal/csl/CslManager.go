//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package csl
import (
	"path/filepath"
	"os"
	"github.com/abieberbach/goplane/extra/logging"
	"time"
	"sync"
)

var CslManagerInstance = &CslManager{"", "", "", make([]*CslAircraft, 0), make(map[string][]*CslAircraft)}

type CslManager struct {
	cslPath      string
	defaultICAO  string
	relatedFile  string
	allAircrafts []*CslAircraft
	aircraftMap  map[string][]*CslAircraft //Map mit allen gefunden Flugzeugen inkl. Mapping auf verwandte ICAOs (Key: ICAO-Code)
}


func (self *CslManager) Init(cslPath, relatedFile, defaultICAO string) {
	self.cslPath = cslPath
	self.relatedFile = relatedFile
	self.defaultICAO = defaultICAO
}

func (self *CslManager) LoadAndValidatePackages() {
	go func() {
		packages := self.loadPackages()
		self.validatePackages(packages)
		relatedICAOs := parseRelatedFile(self.relatedFile)
		self.createAircraftMap(packages, relatedICAOs)
	}()
}

func (self *CslManager) loadPackages() *CslPackages {
	logging.Info("loading CSL packages ")
	startTime := time.Now()
	packages := NewCslPackages()
	filepath.Walk(self.cslPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			cslPackage := self.loadOnePackage(path)
			if cslPackage != nil {
				packages.AddPackage(cslPackage)
				logging.Infof("found package %v with %v aircrafts (directory: %v)", cslPackage.Name, len(cslPackage.Aircrafts), cslPackage.BaseDirectory)
				self.allAircrafts = append(self.allAircrafts, cslPackage.Aircrafts...)
			}
		}
		return nil
	})
	logging.Infof("CSL packages loaded (time: %v)", time.Now().Sub(startTime))
	return packages
}

func (self *CslManager) loadOnePackage(path string) *CslPackage {
	file, err := os.Open(filepath.Join(path, "xsb_aircraft.txt"))
	defer file.Close()
	if err != nil {
		if os.IsNotExist(err) {
			logging.Debugf("ignore folder: %v", path)
		} else {
			logging.Warningf("could not open CSL package \"%v\": %v", path, err)
		}
		return nil
	}
	cslPackage, err := parseCslPackage(file, path)
	if err != nil {
		logging.Warningf("could not read CSL aircraft file \"%v\": %v", path, err)
		return nil
	}
	return cslPackage
}

func (self *CslManager) validatePackages(packages *CslPackages) {
	logging.Info("validating CSL packages ")
	startTime := time.Now()
	waitGroup := sync.WaitGroup{}
	for _, cslPackage := range packages.GetAllPackages() {
		waitGroup.Add(1)
		go func(currentPackage *CslPackage) {
			currentPackage.validate(packages)
			waitGroup.Done()
		}(cslPackage)
	}
	waitGroup.Wait()
	logging.Infof("CSL packages validated (time: %v)", time.Now().Sub(startTime))
}

func (self *CslManager) createAircraftMap(packages *CslPackages, relatedICAOMap map[string][]string) {
	logging.Info("creating aircraft map")
	startTime := time.Now()
	totalCount := 0
	for _, cslPackage := range packages.GetAllPackages() {
		for _, cslAircraft := range cslPackage.Aircrafts {
			totalCount += self.appendAircraftToMap(cslAircraft, relatedICAOMap)
		}
	}
	logging.Infof("aircraft map created: %v aircrafts registerd (time: %v)", totalCount, time.Now().Sub(startTime))
}

func (self *CslManager) appendAircraftToMap(cslAircraft *CslAircraft, relatedICAOMap map[string][]string) int {
	relatedICAOs, found := relatedICAOMap[cslAircraft.Icao]
	if !found {
		relatedICAOs = []string{cslAircraft.Icao}
	}
	//für jeden ICAO-Code das Flugzeug zur Map hinzufügen
	for _, currentIcao := range relatedICAOs {
		aircrafts, exists := self.aircraftMap[currentIcao]
		if !exists {
			aircrafts = make([]*CslAircraft, 0)
		}
		aircrafts = append(aircrafts, cslAircraft)
		self.aircraftMap[currentIcao] = aircrafts
		logging.Debugf("adding aircraft to map: icao=%v airline=%v livery=%v", currentIcao, cslAircraft.Airline, cslAircraft.Livery)
	}
	return len(relatedICAOs)
}

/*
Ermittelt anhand der Daten das gewünschte Flugzeug. Bei der Ermittlung gilt folgende Priorität
1. ICAO, Airline und Livery stimmen überein
2. ICAO und Airline stimmen überein
3. verwandter ICAO, Airline und Livery stimmen überein
4. verwandter ICAO und Airline stimmen überein
5. ICAO stimmt überein
6. verwandter ICAO stimmt überein
7. Standard-ICAO stimmt überein
8. verwandter Standard-ICAO stimmt überein
9. beliebiges Flugzeug nehmen
 */
func (self *CslManager) GetAircraft(icao, airline, livery string) *CslAircraft {
	logging.Debugf("request aircraft: icao=%v, airline=%v, livery=%v", icao, airline, livery)
	startTime := time.Now()
	aircraftsWithIcao, found := self.aircraftMap[icao]
	defaultIcao := false
	if !found {
		//es gibt keine Flugzeuge mit dem gewünschten ICAO --> Standard-ICAO versuchen
		aircraftsWithIcao, found = self.aircraftMap[self.defaultICAO]
		icao = self.defaultICAO
		airline = ""
		livery = ""
		defaultIcao = true
		if !found {
			//es gibt auch keine Flugzeuge mit dem Standard-ICAO --> Punkt 9 beliebiges Flugzeug nehmen
			for _, aircrafts := range self.aircraftMap {
				logging.Debugf("no matching aircraft found, use: icao=%v, airline=%v, livery=%v, priority=9 (time: %v)", aircrafts[0].Icao, aircrafts[0].Airline, aircrafts[0].Livery, time.Now().Sub(startTime))
				return aircrafts[0]
			}
		}
	}
	matchingResult := matchPlane(aircraftsWithIcao, icao, airline, livery)
	var result *CslAircraft
	var priority int
	if matchingResult.exactIcaoAirlineLiveryMatch != nil {
		//ICAO, Airline und Livery stimmen überein
		result = matchingResult.exactIcaoAirlineLiveryMatch
		priority = 1
	} else if matchingResult.exactIcaoAirlineMatch != nil {
		//ICAO und Airline stimmen überein
		result = matchingResult.exactIcaoAirlineMatch
		priority = 2
	} else if matchingResult.relatedIcaoAirlineLiveryMatch != nil {
		//verwandter ICAO, Airline und Livery stimmen überein
		result = matchingResult.relatedIcaoAirlineLiveryMatch
		priority = 3
	} else if matchingResult.relatedIcaoAirlineMatch != nil {
		//verwandter ICAO und Airline stimmen überein
		result = matchingResult.relatedIcaoAirlineMatch
		priority = 4
	} else if matchingResult.exactIcaoMatch != nil {
		//ICAO stimmt überein
		result = matchingResult.exactIcaoMatch
		priority = 5
		if defaultIcao {
			priority = 7
		}
	} else {
		//verwandter ICAO stimmt überein
		result = matchingResult.relatedIcaoMatch
		priority = 6
		if defaultIcao {
			priority = 8
		}
	}
	logging.Debugf("matching aircraft found, use: icao=%v, airline=%v, livery=%v, priority=%v (time: %v)", result.Icao, result.Airline, result.Livery, priority, time.Now().Sub(startTime))
	return result
}

func (self *CslManager) IsICAOValid(icao string) bool {
	_, found := self.aircraftMap[icao]
	return found
}


func (self *CslManager) SetDefaultICAO(defaultICAO string) {
	self.defaultICAO = defaultICAO
}

func (self *CslManager) GetAllAircrafts() []*CslAircraft {
	return self.allAircrafts
}