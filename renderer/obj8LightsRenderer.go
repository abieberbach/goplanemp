//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package renderer
import (
	"github.com/abieberbach/goplane/xplm/scenery"
	"github.com/abieberbach/goplane/xplm/plugins"
	"path/filepath"
	"github.com/abieberbach/goplanemp/internal/obj7"
)




func (self *DefaultRenderer) drawObj8Lights(renderData *planeRenderData) {
	if self.greenNavRef == nil {
		self.initObj8Lights()
	}
	lodInfo := calculateLOD(renderData.plane.CslAircraft.ObjInfo, renderData.dist)
	greenNavDrawInfos := make([]scenery.DrawInfo, 0)
	redNavDrawInfos := make([]scenery.DrawInfo, 0)
	strobeDrawInfos := make([]scenery.DrawInfo, 0)
	beaconDrawInfos := make([]scenery.DrawInfo, 0)
	landingDrawInfos := make([]scenery.DrawInfo, 0)
	taxiDrawInfos := make([]scenery.DrawInfo, 0)
	otherDrawInfos := make([]scenery.DrawInfo, 0)
	for _, currentLight := range lodInfo.Lights {
		switch currentLight.LightType {
		case obj7.LightType_GreenNavigation:
			if renderData.lightStatus.NavLights {
				greenNavDrawInfos = appendLight(renderData, currentLight, greenNavDrawInfos)
			}
		case obj7.LightType_RedNavigation:
			if renderData.lightStatus.NavLights {
				redNavDrawInfos = appendLight(renderData, currentLight, redNavDrawInfos)
			}
		case obj7.LightType_Strobe:
			if renderData.lightStatus.StrobeLights {
				strobeDrawInfos = appendLight(renderData, currentLight, strobeDrawInfos)
			}
		case obj7.LightType_Beacon:
			if renderData.lightStatus.BeaconLights {
				beaconDrawInfos = appendLight(renderData, currentLight, beaconDrawInfos)
			}
		case obj7.LightType_Landing:
			if renderData.lightStatus.LandingLights {
				landingDrawInfos = appendLight(renderData, currentLight, landingDrawInfos)
			}
		case obj7.LightType_Taxi:
			if renderData.lightStatus.TaxiLights {
				taxiDrawInfos = appendLight(renderData, currentLight, taxiDrawInfos)
			}
		case obj7.LightType_Other:
			if renderData.lightStatus.NavLights {
				otherDrawInfos = appendLight(renderData, currentLight, otherDrawInfos)
			}
		default:
			otherDrawInfos = appendLight(renderData, currentLight, otherDrawInfos)
		}
	}
	scenery.DrawObjects(self.greenNavRef, len(greenNavDrawInfos), greenNavDrawInfos, true, false)
	scenery.DrawObjects(self.redNavRef, len(redNavDrawInfos), redNavDrawInfos, true, false)
	scenery.DrawObjects(self.strobeRef, len(strobeDrawInfos), strobeDrawInfos, true, false)
	scenery.DrawObjects(self.beaconRef, len(beaconDrawInfos), beaconDrawInfos, true, false)
	scenery.DrawObjects(self.landingRef, len(landingDrawInfos), landingDrawInfos, true, false)
	scenery.DrawObjects(self.taxiRef, len(taxiDrawInfos), taxiDrawInfos, true, false)
	scenery.DrawObjects(self.otherRef, len(otherDrawInfos), otherDrawInfos, true, false)
}

func appendLight(renderData *planeRenderData, currentLight *obj7.LightInfo, drawInfos []scenery.DrawInfo) []scenery.DrawInfo {
	return append(drawInfos, scenery.NewDrawInfoLocal(renderData.x + currentLight.XYZ[0], renderData.y + currentLight.XYZ[1], renderData.z + currentLight.XYZ[2], renderData.plane.PositionData.Pitch, renderData.plane.PositionData.Heading, renderData.plane.PositionData.Roll))

}

func (self *DefaultRenderer) initObj8Lights() {
	_, path, _, _ := plugins.GetPluginInfo(plugins.GetMyId())
	//aktueller Pfad ist die XPL-Datei --> noch eins nach oben
	path = filepath.Dir(path)
	//aktueller Pfad ist jetzt der 64-Ordner --> noch eins nach oben
	path = filepath.Dir(path)
	if self.configuration.DisableObj8LightSpills {
		path = filepath.Join(path,"Resources", "obj8Lights")
	} else {
		path = filepath.Join(path,"Resources", "obj8Lights", "spill")
	}
	self.greenNavRef = scenery.LoadObject(filepath.Join(path, "greenNavLight.obj"))
	self.redNavRef = scenery.LoadObject(filepath.Join(path, "redNavLight.obj"))
	self.strobeRef = scenery.LoadObject(filepath.Join(path, "strobeLight.obj"))
	self.beaconRef = scenery.LoadObject(filepath.Join(path, "beaconLight.obj"))
	self.landingRef = scenery.LoadObject(filepath.Join(path, "landingLight.obj"))
	self.taxiRef = scenery.LoadObject(filepath.Join(path, "taxiLight.obj"))
	self.otherRef = scenery.LoadObject(filepath.Join(path, "otherLight.obj"))

}