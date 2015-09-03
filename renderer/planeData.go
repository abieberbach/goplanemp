//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package renderer
import (
	"github.com/abieberbach/goplanemp"
	"github.com/abieberbach/goplane/xplm/planes"
	"math"
)

func (self *DefaultRenderer) buildPlaneToRender(plane *goplanemp.Plane, x, y, z, distance float32,userAircraftAlt float64,now int) *planeRenderData {
	result := &planeRenderData{}
	result.x = x
	result.y = y
	result.z = z
	result.plane=plane
	result.tcas = showOnTCAS(now,plane,userAircraftAlt)
	result.dist = distance
	result.drawState,result.lightStatus=buildDrawingAndLightState(plane,now)
	return result
}


func showOnTCAS(now int, plane *goplanemp.Plane,userAircraftAlt float64) bool {
	if plane.UpdateData(now, goplanemp.DataType_Radar) != goplanemp.Data_Unavailable {
		if plane.RadarData.Mode == goplanemp.TransponderMode_Standby {
			//TransponderMode = Standby --> keine Anzeige im TCAS
			return false
		}
	}

	// prüfen ob die Höhendifferenz für das Flugzeug eingehalten wird
	alt_diff := math.Abs(plane.PositionData.Elevation - userAircraftAlt)
	if alt_diff > MAX_TCAS_ALTDIFF {
		//Flugzeug ist außerhalb der Höhendifferenz --> keine Anzeige im TCAS
		return false
	}
	return true
}


func buildDrawingAndLightState(plane *goplanemp.Plane,now int) (planes.PlaneDrawState, goplanemp.LightStatus) {
	drawState := planes.PlaneDrawState{}
	lightState := goplanemp.LightStatus{}

	if plane.UpdateData(now,goplanemp.DataType_Surfaces) != goplanemp.Data_Unavailable {
		drawState.GearPosition = plane.SurfacesData.GearPosition;
		drawState.FlapRatio = plane.SurfacesData.FlapRatio;
		drawState.SpoilerRatio = plane.SurfacesData.SpoilerRatio;
		drawState.SpeedBrakeRatio = plane.SurfacesData.SpeedBrakeRatio;
		drawState.SlatRatio = plane.SurfacesData.SlatRatio;
		drawState.WingSweep = plane.SurfacesData.WingSweep;
		drawState.Thrust = plane.SurfacesData.Thrust;
		drawState.YokePitch = plane.SurfacesData.YokePitch;
		drawState.YokeHeading = plane.SurfacesData.YokeHeading;
		drawState.YokeRoll = plane.SurfacesData.YokeRoll;
		lightState = plane.SurfacesData.Lights

	} else {
		drawState.GearPosition = 0.0
		drawState.FlapRatio = 0.0
		if plane.PositionData.Elevation < 70 {
			drawState.GearPosition = 1.0
			drawState.FlapRatio = 1.0
		}
		drawState.SpoilerRatio = 0.0
		drawState.SpeedBrakeRatio = 0.0
		drawState.SlatRatio = 0.0
		drawState.WingSweep = 0.0
		drawState.Thrust = 0.6
		if plane.PositionData.Pitch > 30 {
			drawState.Thrust = 1.0
		}
		drawState.YokePitch = plane.PositionData.Pitch / 90.0
		drawState.YokeHeading = plane.PositionData.Heading / 180.0
		drawState.YokeRoll = plane.PositionData.Roll / 90.0

		// use some smart defaults
		lightState.BeaconLights = true
		lightState.NavLights = true
	}
	//Flugzeug hat kein einfahrbares Fahrwerk --> Fahrwerk ist immer ausgefahren
	if !plane.CslAircraft.HasGear {
		drawState.GearPosition = 1.0;
	}
	return drawState, lightState
}