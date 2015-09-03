//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package goplanemp
import (
	"github.com/abieberbach/goplane/xplm/display"
	"github.com/abieberbach/goplane/xplm/dataAccess"
	"github.com/abieberbach/goplane/extra/logging"
	"time"
)

//Setzt die Renderer-Implementierung die genutzt werden soll.
//Soll der Standardrenderer genutzt werden, dann sieht der Aufruf wie folgt aus:
//	mpManager.SetPlaneRenderer(renderer.NewDefaultRenderer(config), nil)
func (self *MultiplayerManager) SetPlaneRenderer(renderer Renderer, ref interface{}) {
	self.renderer = renderer
	self.rendererRef = ref
}

//Dumpt alle Informationen eines Zykluses in die Log-Datei von X-Plane.
func (self *MultiplayerManager) DumpOneCycle() {
	self.dumpDebugInfo = true
}

func (self *MultiplayerManager) drawAirplanes(phase display.DrawingPhase, isBefore bool, ref interface{}) int {
	defer func() {
		self.dumpDebugInfo = false
	}()
	if len(self.planes) == 0 {
		if self.dumpDebugInfo {
			logging.Debug("No planes this cycle.")
		}
		return 1
	}
	is_shadow := self.wrtDataRef != nil && dataAccess.GetIntData(self.wrtDataRef) != 0;
	if self.prtDataRef != nil {
		self.isBlend = dataAccess.GetIntData(self.prtDataRef) == 2;
	}
	var blendValue bool
	if is_shadow {
		blendValue = false
	} else {
		blendValue = self.isBlend
	}
	var startTime time.Time
	if self.dumpDebugInfo {
		startTime = time.Now()
	}
	self.planesMutex.RLock()
	self.renderer.Render(blendValue, self.renderAircraftLabels, self.planes, self.rendererRef)
	self.planesMutex.RUnlock()
	if (!is_shadow) {
		self.isBlend = !self.isBlend;
	}
	if self.dumpDebugInfo {
		logging.Debugf("Drawing time: %v", time.Now().Sub(startTime))
		logging.Debug("Dumping one cycle map of planes")
		self.planesMutex.RLock()
		for _, plane := range self.planes {
			logging.Debugf("[%v] - icao=%v, airline=%v, livery:=%v", plane.PlaneId, plane.CslAircraft.Icao, plane.CslAircraft.Airline, plane.CslAircraft.Livery)
			if plane.PositionData == nil {
				logging.Debugf("Drawing plane: %v at <no position data>", plane.CslAircraft.ObjFile)
			} else {
				logging.Debugf("Drawing plane: %v at %v,%v,%v (%vx%vx%v full=%v)", plane.CslAircraft.ObjFile, plane.PositionData.Lat, plane.PositionData.Lon, plane.PositionData.Elevation, plane.PositionData.Pitch, plane.PositionData.Roll, plane.PositionData.Heading, plane.RenderFull)
			}
		}
		self.planesMutex.RUnlock()
		logging.Debug("End of cycle dump.")
	}
	return 1
}

//Aktiviert das Zeichnen der Flugzeugbeschriftungen.
func (self *MultiplayerManager) EnableAircraftLabels() {
	self.renderAircraftLabels = true
}

//Deaktiviert das Zeichnen der Flugzeugbeschriftungen.
func (self *MultiplayerManager) DisableAircraftLabels() {
	self.renderAircraftLabels = false
}

//Gibt an, ob die Flugzeugbeschriftungen gezeichnet werden.
func (self *MultiplayerManager) DrawingAircraftLabels() bool {
	return self.renderAircraftLabels
}