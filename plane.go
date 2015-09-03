package goplanemp
import (
	"github.com/abieberbach/goplane/extra/logging"
	"strings"
	"github.com/abieberbach/goplanemp/internal/csl"
	"github.com/abieberbach/goplanemp/internal/texture"
)

//Enthält die notwendigen Daten für die Verwaltung eines Flugzeugs
//BUG(abieberbach) Unterstützung für ACF-Flugzeuge wurde entfernt
type Plane struct {
	PlaneId         PlaneId
	CslAircraft     *csl.CslAircraft //zugeordnete CSl-Definition
	TextureNum      int32            //Nummer der Textur für den Tag
	LitTextureNum   int32            //Nummer der Textur für die Nacht
	PositionData    *PlanePosition   //aktuelle Positionsdaten
	SurfacesData    *PlaneSurfaces   //aktuelle Steuerungsdaten
	RadarData       *PlaneRadar      //aktuelle Radardaten
	RenderFull      bool             //Flag ob das Flugzeug voll oder nur die Lichter gezeichnet werden sollen
	dataFunc        PlaneDataFunc    //Funktion zur Ermittlung der Daten
	refData         interface{}      //zusätzliche Referenzdaten für die PlaneDataFunc
	positionDataAge int              //Alter der Positionsdaten
	surfacesDataAge int              //Alter der Steuerungsdaten
	radarDataAge    int              //Alter der Radardaten
}

func newPlane(id PlaneId, dataFunc PlaneDataFunc, refData interface{}) *Plane {
	newAircraft := &Plane{}
	newAircraft.PlaneId = id
	newAircraft.dataFunc = dataFunc
	newAircraft.refData = refData
	newAircraft.RenderFull = true
	return newAircraft
}

func (self *Plane) UpdateData(now int, dataType PlaneDataType) (PlaneCallbackResult) {
	switch dataType {
	case DataType_Position:
		if self.positionDataAge != now {
			callbackResult, newPosData, _, _ := self.dataFunc(self.PlaneId, dataType, self.refData)
			if callbackResult == Data_NewData {
				self.PositionData = newPosData
				self.positionDataAge = now
			}
			return callbackResult
		}
		return Data_Unchanged
	case DataType_Surfaces:
		if self.surfacesDataAge != now {
			callbackResult, _, newSurfacesData, _ := self.dataFunc(self.PlaneId, dataType, self.refData)
			if callbackResult == Data_NewData {
				self.SurfacesData = newSurfacesData
				self.surfacesDataAge = now
			}
			return callbackResult
		}
		return Data_Unchanged
	case DataType_Radar:
		if self.radarDataAge != now {
			callbackResult, _, _, newRadarData := self.dataFunc(self.PlaneId, dataType, self.refData)
			if callbackResult == Data_NewData {
				self.RadarData = newRadarData
				self.radarDataAge = now
			}
			return callbackResult
		}
		return Data_Unchanged
	}
	return Data_Unavailable
}

func (self *Plane) updateAircraftData() {
	if self.CslAircraft.CslType == csl.CSL_TYPE_OBJ7 {
		//für OBJ7-Flugzeuge das Model parsen und die Texturen laden
		self.TextureNum, self.LitTextureNum = self.loadTextures()
	}
}

func (self *Plane)loadTextures() (int32, int32) {
	if self.CslAircraft.ObjInfo == nil {
		return -1, -1
	}
	//erstmal Standardtextur nehmen
	textureFile := self.CslAircraft.ObjInfo.DefaultTexture
	if !strings.HasSuffix(textureFile, ".png") {
		textureFile += ".png"
	}
	if len(self.CslAircraft.TextureFile) != 0 {
		//es gibt eine Textur aus der CSL-Definition --> diese nehmen
		textureFile = self.CslAircraft.TextureFile
	}
	textureNum, litTextureNum, err := texture.TextureManagerInstance.LoadTexture(textureFile, false)

	if err != nil {
		logging.Errorf("aircraft (icao=%v, airline=%v, livery=%v): %v", self.CslAircraft.Icao, self.CslAircraft.Airline, self.CslAircraft.Livery, err)
		return -1, -1
	}
	//normale und LIT-Textur existieren --> beide Nummern zurückgeben
	return textureNum, litTextureNum
}

