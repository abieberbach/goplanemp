package extra
import (
	"fmt"
	"github.com/abieberbach/goplane/xplm/dataAccess"
	"github.com/abieberbach/goplane/extra/logging"
	"github.com/abieberbach/goplane/xplm/processing"
	"github.com/abieberbach/goplanemp"
	"github.com/abieberbach/goplanemp/internal/csl"
)

//Testet alle gefundenen Flugzeuge und zeichnet sie in X-Plane
type AircraftTest struct {
	multiplayerManager *goplanemp.MultiplayerManager
	currentTestPlaneId goplanemp.PlaneId
	allAircrafts       []*csl.CslAircraft
	currentIndex       int
	oldConfigMod       int
	latDataRef         dataAccess.DataRef
	lonDataRef         dataAccess.DataRef
	elevDataRef        dataAccess.DataRef
}

//Erzeugt einen neuen Test
func NewAircraftTest(multiplayerManager *goplanemp.MultiplayerManager) *AircraftTest {
	result := &AircraftTest{}
	result.multiplayerManager = multiplayerManager
	result.currentTestPlaneId = -1
	result.currentIndex = 0
	result.oldConfigMod = 0
	result.latDataRef, _ = dataAccess.FindDataRef("sim/flightmodel/position/latitude")
	result.lonDataRef, _ = dataAccess.FindDataRef("sim/flightmodel/position/longitude")
	result.elevDataRef, _ = dataAccess.FindDataRef("sim/flightmodel/position/elevation")
	return result
}

//Führt den Test aus.
func (self *AircraftTest) TestAllAirplanes() {
	self.allAircrafts = csl.CslManagerInstance.GetAllAircrafts()
	self.oldConfigMod = self.multiplayerManager.GetConfiguration().CalculationMod
	self.multiplayerManager.GetConfiguration().CalculationMod = 1
	processing.RegisterFlightLoopCallback(self.testFlightLoop, 0.125, nil)
}

func (self *AircraftTest) testFlightLoop(elapsedSinceLastCall, elapsedTimeSinceLastFlightLoop float32, counter int, ref interface{}) float32 {
	if self.currentTestPlaneId > 0 {
		self.multiplayerManager.DestroyPlane(self.currentTestPlaneId)
	}
	if self.currentIndex >= len(self.allAircrafts) {
		processing.UnregisterFlightLoopCallback(self.testFlightLoop, nil)
		self.multiplayerManager.GetConfiguration().CalculationMod = self.oldConfigMod
		return 0
	}
	currentAircraft := self.allAircrafts[self.currentIndex]
	logging.Debugf("testing Plane: index=%v/%v icao=%v airline=%v livery=%v", self.currentIndex, len(self.allAircrafts), currentAircraft.Icao, currentAircraft.Airline, currentAircraft.Livery)
	self.currentTestPlaneId = self.multiplayerManager.CreatePlane(currentAircraft.Icao, currentAircraft.Airline, currentAircraft.Livery, self.testDataFunc, currentAircraft)
	self.currentIndex += 1
	return 0.125

}

func (self *AircraftTest) testDataFunc(id goplanemp.PlaneId, dataType goplanemp.PlaneDataType, ref interface{}) (goplanemp.PlaneCallbackResult, *goplanemp.PlanePosition, *goplanemp.PlaneSurfaces, *goplanemp.PlaneRadar) {
	currentPlane := ref.(*csl.CslAircraft)
	switch dataType {
	case goplanemp.DataType_Position:
		posData := &goplanemp.PlanePosition{}
		posData.Lat = dataAccess.GetDoubleData(self.latDataRef)
		posData.Lon = dataAccess.GetDoubleData(self.lonDataRef)
		posData.Elevation = dataAccess.GetDoubleData(self.elevDataRef) + 25.0 //25 Fuß höher damit es gut zu sehen ist
		posData.Heading = 0.0
		posData.Pitch = 0.0
		posData.Roll = 0.0
		posData.Label = fmt.Sprintf("%v - %v - %v", currentPlane.Icao, currentPlane.Airline, currentPlane.Livery)
		return goplanemp.Data_NewData, posData, nil, nil
	case goplanemp.DataType_Surfaces:
		surfaceData := &goplanemp.PlaneSurfaces{}
		surfaceData.Lights = goplanemp.LightStatus{}
		surfaceData.Lights.BeaconLights = true
		surfaceData.Lights.LandingLights = true
		surfaceData.Lights.NavLights = true
		surfaceData.Lights.StrobeLights = true
		surfaceData.Lights.FlashPattern = goplanemp.LightsPattern_Default
		return goplanemp.Data_NewData, nil, surfaceData, nil
	case goplanemp.DataType_Radar:
		return goplanemp.Data_Unavailable, nil, nil, nil
	}
	return goplanemp.Data_Unavailable, nil, nil, nil
}
