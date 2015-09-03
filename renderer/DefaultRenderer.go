//Enthält die Standardimplementierung für das Zeichnen der Objekte
package renderer
import (
	"github.com/abieberbach/goplane/xplm/dataAccess"
	"github.com/abieberbach/goplane/extra/logging"
	"github.com/abieberbach/goplane/extra"
	"github.com/abieberbach/goplane/xplm/camera"
	"github.com/abieberbach/goplanemp"
	"github.com/abieberbach/goplane/xplm/graphics"
	"math"
	"github.com/abieberbach/goplane/xplm/processing"
	"github.com/abieberbach/goplane/xplm/scenery"
	"github.com/abieberbach/goplanemp/internal/csl"
	"time"
)

const (
	MAX_TCAS_ALTDIFF = 5000 //max. Höhendifferenz für TCAS
	MAX_LABEL_DIST = 5000.0 //max. Sichtweite der Labels in m
	MAX_TCAS_DIST = 40.0 * NMToMeters //max. TCAS-Entfernung (40 NM)

	FtToMeters = 0.3048 //Umrechnungsfaktor Fuss -> Meter
	MetersToNM = 0.000539956803 //Umrechnungsfaktor Meter -> NM
	NMToMeters = 1 / MetersToNM //Umrechnungsfaktor NM -> Meter
)

//Standardimplementierung der Renderer-Schnittstelle
type DefaultRenderer struct {
	configuration   *goplanemp.Configuration

	visDataRef      dataAccess.DataRef
	altitudeDataRef dataAccess.DataRef
	fovDataRef      dataAccess.DataRef

	totalPlanes     int
	obj7Planes      int
	obj8Planes      int
	drawTime        float64
	renderPlanes    *renderList
	greenNavRef     scenery.ObjectRef
	redNavRef       scenery.ObjectRef
	strobeRef       scenery.ObjectRef
	beaconRef       scenery.ObjectRef
	landingRef      scenery.ObjectRef
	taxiRef         scenery.ObjectRef
	otherRef        scenery.ObjectRef


}

//Erzeugt einen neuen Standardrenderer.
func NewDefaultRenderer(config *goplanemp.Configuration) goplanemp.Renderer {
	result := &DefaultRenderer{}
	result.renderPlanes = newRenderList()
	result.configuration = config
	return result
}

func getRenderStat(ref interface{}) int32 {
	return int32(*(ref.(*int)))
}

func (self *DefaultRenderer) Init() {
	var dataRefFound bool
	self.visDataRef, dataRefFound = dataAccess.FindDataRef("sim/graphics/view/visibility_effective_m")
	if !dataRefFound {
		self.visDataRef, dataRefFound = dataAccess.FindDataRef("sim/weather/visibility_effective_m")
		if !dataRefFound {
			logging.Warning("Default renderer could not find effective visibility in the sim.")
		}
	}
	self.altitudeDataRef, _ = dataAccess.FindDataRef("sim/flightmodel/position/elevation")
	self.fovDataRef, _ = dataAccess.FindDataRef("sim/graphics/view/field_of_view_deg");
	if self.configuration.EnableRenderStatistics {
		accessors := dataAccess.DataRefAccessors{}
		accessors.ReadInt = getRenderStat
		dataAccess.RegisterDataAccessor("hack/renderer/planes", dataAccess.TypeInt, false, accessors, &self.totalPlanes, nil)
		dataAccess.RegisterDataAccessor("hack/renderer/objects7", dataAccess.TypeInt, false, accessors, &self.obj7Planes, nil)
		dataAccess.RegisterDataAccessor("hack/renderer/objects8", dataAccess.TypeInt, false, accessors, &self.obj8Planes, nil)
		dataAccess.RegisterDataAccessor("hack/renderer/drawTime", dataAccess.TypeDouble, false, accessors, &self.drawTime, nil)
		extra.RegisterDataRefToDataRefEditor("hack/renderer/planes", "hack/renderer/objects7", "hack/renderer/objects8", "hack/renderer/drawTime")
	}
	registerMultiplayerDataRefs()
}

func (self *DefaultRenderer) Render(isBlend, renderAircraftLabels bool, aircrafts map[goplanemp.PlaneId]*goplanemp.Plane, ref interface{}) {
	var startTime time.Time
	if self.configuration.EnableRenderStatistics {
		startTime = time.Now()
	}
	self.totalPlanes = len(aircrafts)
	self.obj7Planes, self.obj8Planes = 0, 0
	now := processing.GetCycleNumber()
	cameraPos := camera.XPLMReadCameraPosition()
	cameraZoom := cameraPos.Zoom
	maxDist := dataAccess.GetFloatData(self.visDataRef)
	labelDist := float32(math.Min(float64(maxDist), MAX_LABEL_DIST)) * cameraZoom
	fullPlaneDist := cameraZoom * (5280.0 / 3.2) * self.configuration.FullDistance * 0.6213712
	cullInfo := newCullInfo()
	if now % self.configuration.CalculationMod == 0 {
		self.renderPlanes.Clear()
		userAircraftAlt := dataAccess.GetDoubleData(self.altitudeDataRef) / FtToMeters
		for _, plane := range aircrafts {
			//neue Liste erzeugen
			if plane.CslAircraft.CslType == csl.CSL_TYPE_OBJ7&&plane.CslAircraft.ObjInfo == nil {
				//es wurde kein Objekt für ein OBJ7-Flugzeug geladen --> aktuelles Flugzeug ignorieren
				continue
			}
			if plane.UpdateData(now, goplanemp.DataType_Position) == goplanemp.Data_Unavailable {
				//keine Positionsdaten vorhanden --> aktuelles Flugzeug ignorieren
				continue
			}
			//Entfernung des Flugzeugs bestimmen
			x64, y64, z64 := graphics.WorldToLocal(plane.PositionData.Lat, plane.PositionData.Lon, plane.PositionData.Elevation * FtToMeters)
			x, y, z := float32(x64), float32(y64), float32(z64)

			//prüfen ob das Flugzeug angezeigt werden soll
			cull, distance := self.isCulling(plane, cullInfo, x, y, z, maxDist)
			//Flugzeug ist außerhalb des TCAS-Bereichs --> ignorieren
			if distance > MAX_TCAS_DIST {
				continue
			}
			plane.RenderFull = distance < fullPlaneDist
			if !cull {
				self.renderPlanes.AddPlane(self.buildPlaneToRender(plane, x, y, z, distance, userAircraftAlt, now))
			}
		}
	}

	//und Flugzeuge zeichen
	obj7Planes, obj8Planes := self.renderPlanes.GetPlanes()
	if len(obj7Planes) > 0 {
		if isBlend {
			// Durchlauf 1 - OBJ7-Flugzeuge
			// Blend for solid OBJ7s?  YES!  First, in HDR mode, they DO NOT draw to the gbuffer properly -
			// they splat their livery into the normal map, which is terrifying and stupid.  Then they are also
			// pre-lit...the net result is surprisingly not much worse than regular rendering considering how many
			// bad things have happened, but for all I know we're getting NaNs somewhere.
			//
			// Blending isn't going to hurt things in NON-HDR because our rendering is so stupid for old objs - there's
			// pretty much never translucency so we aren't going to get Z-order fails.  So f--- it...always draw blend.<

			for _, current := range obj7Planes {
				if current.dist > fullPlaneDist {
					//max. Entfernung wurde erreicht --> aktuelles Flugzeug nicht zeichnen
					continue
				}
				self.drawObj7Model(current)
				self.obj7Planes++
			}
			// Durchlauf 2- OBJ7-Lichter
			fov := dataAccess.GetFloatData(self.fovDataRef)
			self.beginObj7LightDrawing()
			for _, current := range obj7Planes {
				if !self.configuration.UseObj8Lights || current.dist > fullPlaneDist {
					self.drawObj7Lights(current, fov, cameraPos)
				}
			}
		} else {
			if self.configuration.UseObj8Lights {
				for _, current := range obj7Planes {
					if current.dist <= fullPlaneDist {
						self.drawObj8Lights(current)
					}
				}
			}
		}
	}

	//Durchlauf 3 - OBJ8-Flugzeuge zeichnen
	if len(obj8Planes) > 0 {

		if isBlend {
			for _, plane := range obj8Planes {
				currentAircraft = plane.plane
				self.drawObj8Tanslucent(plane)
			}
		} else {
			for _, plane := range obj8Planes {
				currentAircraft = plane.plane
				self.drawObj8Solid(plane)
			}
		}
		self.obj8Planes = len(obj8Planes)
	}
	// Durchlauf 4 - Beschriftungen zeichnen
	if isBlend && renderAircraftLabels {
		self.drawPlaneLabels(self.renderPlanes, cullInfo, labelDist)
	}
	if self.configuration.EnableRenderStatistics {
		self.drawTime = time.Now().Sub(startTime).Seconds() / 1000.0
	}
}

func (self *DefaultRenderer) isCulling(plane *goplanemp.Plane, cullInfo *cullInfo, x, y, z, maxDist float32) (bool, float32) {
	//prüfen ob Flugzeug außerhalb des Anzeigebereichs ist
	distance := float32(math.Sqrt(float64(cullInfo.sphereDistanceSqr(x, y, z))))
	result := distance > maxDist
	// Calculate the angles between the camera angles and the real angles.
	// Cull if we exceed half the FOV.
	if !result && !cullInfo.sphereIsVisible(float32(x), float32(y), float32(z), 50.0) {
		result = true
	}
	return result, distance
}


