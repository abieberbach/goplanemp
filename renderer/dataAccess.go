package renderer
import (
	"github.com/abieberbach/goplanemp"
	"github.com/abieberbach/goplane/xplm/dataAccess"
	"github.com/abieberbach/goplane/extra"
)

type dataType int

const (
	gear_ratio dataType = 0
	flap_ratio dataType = 1
	spoiler_ratio dataType = 2
	speedbreak_ratio dataType = 3
	slat_ratio dataType = 4
	wing_sweep_ratio dataType = 5
	thrust dataType = 6
	pitch dataType = 7
	heading dataType = 8
	roll dataType = 9
	landing_light_on dataType = 10
	beacon_light_on dataType = 11
	strobe_light_on dataType = 12
	nav_light_on dataType = 13
	taxi_light_on dataType = 14
)

var currentAircraft *goplanemp.Plane

func getFloatValue(ref interface{}) float32 {
	if currentAircraft == nil {
		return 0.0
	}

	switch ref.(dataType) {
	case gear_ratio:
		return currentAircraft.SurfacesData.GearPosition
	case flap_ratio:
		return currentAircraft.SurfacesData.FlapRatio
	case spoiler_ratio:
		return currentAircraft.SurfacesData.SpoilerRatio
	case speedbreak_ratio:
		return currentAircraft.SurfacesData.SpeedBrakeRatio
	case slat_ratio:
		return currentAircraft.SurfacesData.SlatRatio
	case wing_sweep_ratio:
		return currentAircraft.SurfacesData.WingSweep
	case thrust:
		return currentAircraft.SurfacesData.Thrust
	case pitch:
		return currentAircraft.SurfacesData.YokePitch
	case heading:
		return currentAircraft.SurfacesData.YokeHeading
	case roll:
		return currentAircraft.SurfacesData.YokeRoll

	case landing_light_on:
		if currentAircraft.SurfacesData.Lights.LandingLights {
			return 1.0
		}
		return 0.0
	case beacon_light_on:
		if currentAircraft.SurfacesData.Lights.BeaconLights {
			return 1.0
		}
		return 0.0
	case strobe_light_on:
		if currentAircraft.SurfacesData.Lights.StrobeLights {
			return 1.0
		}
		return 0.0
	case nav_light_on:
		if currentAircraft.SurfacesData.Lights.NavLights {
			return 1.0
		}
		return 0.0
	case taxi_light_on:
		if currentAircraft.SurfacesData.Lights.TaxiLights {
			return 1.0
		}
		return 0.0
	default:
		return 0.0

	}
}

func registerMultiplayerDataRefs() {
	accessors := dataAccess.DataRefAccessors{}
	accessors.ReadFloat = getFloatValue
	dataAccess.RegisterDataAccessor("goplanemp/controls/gear_ratio", dataAccess.TypeFloat, false, accessors, gear_ratio, nil)
	dataAccess.RegisterDataAccessor("goplanemp/controls/flap_ratio", dataAccess.TypeFloat, false, accessors, flap_ratio, nil)
	dataAccess.RegisterDataAccessor("goplanemp/controls/spoiler_ratio", dataAccess.TypeFloat, false, accessors, spoiler_ratio, nil)
	dataAccess.RegisterDataAccessor("goplanemp/controls/speed_brake_ratio", dataAccess.TypeFloat, false, accessors, speedbreak_ratio, nil)
	dataAccess.RegisterDataAccessor("goplanemp/controls/slat_ratio", dataAccess.TypeFloat, false, accessors, slat_ratio, nil)
	dataAccess.RegisterDataAccessor("goplanemp/controls/wing_sweep_ratio", dataAccess.TypeFloat, false, accessors, wing_sweep_ratio, nil)
	dataAccess.RegisterDataAccessor("goplanemp/controls/thrust_ratio", dataAccess.TypeFloat, false, accessors, thrust, nil)
	dataAccess.RegisterDataAccessor("goplanemp/controls/yoke_pitch_ratio", dataAccess.TypeFloat, false, accessors, pitch, nil)
	dataAccess.RegisterDataAccessor("goplanemp/controls/yoke_heading_ratio", dataAccess.TypeFloat, false, accessors, heading, nil)
	dataAccess.RegisterDataAccessor("goplanemp/controls/yoke_roll_ratio", dataAccess.TypeFloat, false, accessors, roll, nil)
	dataAccess.RegisterDataAccessor("goplanemp/lights/landing_lights_on", dataAccess.TypeFloat, false, accessors, landing_light_on, nil)
	dataAccess.RegisterDataAccessor("goplanemp/lights/beacon_lights_on", dataAccess.TypeFloat, false, accessors, beacon_light_on, nil)
	dataAccess.RegisterDataAccessor("goplanemp/lights/strobe_lights_on", dataAccess.TypeFloat, false, accessors, strobe_light_on, nil)
	dataAccess.RegisterDataAccessor("goplanemp/lights/nav_lights_on", dataAccess.TypeFloat, false, accessors, nav_light_on, nil)
	dataAccess.RegisterDataAccessor("goplanemp/lights/taxi_lights_on", dataAccess.TypeFloat, false, accessors, taxi_light_on, nil)

	extra.RegisterDataRefToDataRefEditor("goplanemp/controls/gear_ratio", "goplanemp/controls/flap_ratio", "goplanemp/controls/spoiler_ratio", "goplanemp/controls/speed_brake_ratio", "goplanemp/controls/slat_ratio", "goplanemp/controls/wing_sweep_ratio", "goplanemp/controls/thrust_ratio", "goplanemp/controls/yoke_pitch_ratio", "goplanemp/controls/yoke_heading_ratio", "goplanemp/controls/yoke_roll_ratio", "goplanemp/lights/landing_lights_on", "goplanemp/lights/beacon_lights_on", "goplanemp/lights/strobe_lights_on", "goplanemp/lights/nav_lights_on", "goplanemp/lights/taxi_lights_on")

}