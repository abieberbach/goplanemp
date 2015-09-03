//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package goplanemp
import (
	"github.com/abieberbach/goplane/xplm/processing"
	"github.com/abieberbach/goplanemp/internal/csl"
)

//Liefert die Anzahl der aktuell verwalteten Flugzeuge zurück
func (self *MultiplayerManager) CountPlanes() int {
	return len(self.planes)
}

// Liefert die ID des übergebenes Index zurück. Der boolsche Wert zeigt an, ob der Index gefunden wurde oder nicht.
func (self *MultiplayerManager) GetNthPlane(index int) (PlaneId, bool) {
	if index < 0&&index >= self.CountPlanes() {
		return PlaneId(-1), false
	}
	i := 0
	self.planesMutex.RLock()
	defer self.planesMutex.RUnlock()
	for id, _ := range self.planes {
		if i == index {
			return id, true
		}
		i++
	}
	return PlaneId(-1), false
}

//Liefert das ICAO-Kürzel und das Livery für das übergebene Flugzeug zurück. Der boolsche Wert zeigt an, ob die ID gültig war.
func (self *MultiplayerManager) GetPlaneICAOAndLivery(id PlaneId) (icao, livery string, found bool) {
	self.planesMutex.RLock()
	plane, exists := self.planes[id]
	self.planesMutex.RUnlock()
	if exists {
		icao = plane.CslAircraft.Icao
		livery = plane.CslAircraft.Livery
	}
	found = exists
	return
}

//Liefert die aktuellen Daten entsprechend des Übergebenen Datentyp für das gewünschte Flugzeug zurück.
//Alle anderen Zeiger sind nil. Zusätzlich liefert die Methode das Alter der Daten und den Rückgabewert des Callback-Moduls zurück
func (self *MultiplayerManager) GetPlaneData(id PlaneId, dataType PlaneDataType) (int, PlaneCallbackResult, *PlanePosition, *PlaneSurfaces, *PlaneRadar) {
	self.planesMutex.RLock()
	plane, exists := self.planes[id]
	self.planesMutex.RUnlock()
	if !exists {
		return -1, Data_Unavailable, nil, nil, nil
	}
	now := processing.GetCycleNumber()
	callbackResult := plane.UpdateData(now, dataType)
	switch dataType {
	case DataType_Position:
		return plane.positionDataAge, callbackResult, plane.PositionData, nil, nil
	case DataType_Surfaces:
		return plane.surfacesDataAge, callbackResult, nil, plane.SurfacesData, nil
	case DataType_Radar:
		return plane.radarDataAge, callbackResult, nil, nil, plane.RadarData
	default:
		return -1, Data_Unavailable, nil, nil, nil
	}
}

//Überprüft ob der übergebene ICAO-Wert gültig ist, d.h. ob es Flugzeuge in den CSL-Paketen mit diesem ICAO gibt
func (self *MultiplayerManager) IsICAOValid(icao string) bool {
	return csl.CslManagerInstance.IsICAOValid(icao)
}

