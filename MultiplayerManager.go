package goplanemp
import (
	"github.com/abieberbach/goplane/xplm/display"
	"github.com/abieberbach/goplane/extra/logging"
	"fmt"
	"sync/atomic"
	"github.com/abieberbach/goplane/xplm/dataAccess"
	"errors"
	"runtime"
	"sync"
	"github.com/abieberbach/goplanemp/internal/csl"
	"github.com/abieberbach/goplanemp/internal/texture"
	"github.com/abieberbach/goplanemp/gl"
)


//Schnittstelle zur Multiplayer-Bibliothek
type MultiplayerManager struct {
	configuration        *Configuration
	renderer             Renderer
	rendererRef          interface{}
	dumpDebugInfo        bool
	renderAircraftLabels bool
	counter              uint32
	planes               map[PlaneId]*Plane
	isBlend              bool
	wrtDataRef           dataAccess.DataRef
	prtDataRef           dataAccess.DataRef
	planesMutex          sync.RWMutex
}

//Erzeugt einen neuen MultiplayerManager
func New(cslFolder, relatedFile, lightsFile, defaultICAO string, config *Configuration) (*MultiplayerManager, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//OpenGL initialisieren
	err := gl.InitOpenGL()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not initialize OpenGL: %v", err))
	}
	manager := &MultiplayerManager{}
	manager.planes = make(map[PlaneId]*Plane)
	manager.isBlend = false
	manager.wrtDataRef = nil
	manager.prtDataRef = nil
	manager.renderAircraftLabels = config.RenderAircraftLabels
	manager.configuration = config
	manager.planesMutex = sync.RWMutex{}
	csl.CslManagerInstance.Init(cslFolder, relatedFile, defaultICAO)
	err = texture.TextureManagerInstance.Init(config.Resolution, lightsFile, config.MergeLitTextures, config.UseAnyLitTexture)
	if err != nil {
		logging.Error(fmt.Sprintf("Could not load lights texture: %v", err))
		return nil, err
	}
	return manager, nil
}

//Aktiviert den MultiplayerManager. Erst nach diesem Aufruf können neue Objekte erzeugt werden.
func (self *MultiplayerManager) Enable() error {
	csl.CslManagerInstance.LoadAndValidatePackages()
	if self.renderer == nil {
		return errors.New("no renderer available")
	}
	self.wrtDataRef, _ = dataAccess.FindDataRef("sim/graphics/view/world_render_type")
	self.prtDataRef, _ = dataAccess.FindDataRef("sim/graphics/view/plane_render_type")
	self.renderer.Init()
	display.RegisterDrawCallback(self.drawAirplanes, display.Phase_Airplanes, false, nil)
	return nil
}

// Fügt ein neues Flugzeug dem Manager hinzu. Dazu wird versucht entsprechend der Angaben das genaue Flugzeug zu finden.
// Die Datenbereitstellung erfolgt über die PlaneDataFunc. Die zusätzlichen refData werden an die PlaneDataFunc übergeben.
func (self *MultiplayerManager) CreatePlane(icao, airline, livery string, dataFunc PlaneDataFunc, refData interface{}) PlaneId {
	id := PlaneId(atomic.AddUint32(&self.counter, 1))
	plane := newPlane(id, dataFunc, refData)
	plane.CslAircraft = csl.CslManagerInstance.GetAircraft(icao, airline, livery)
	plane.updateAircraftData()
	self.planesMutex.Lock()
	self.planes[id] = plane
	self.planesMutex.Unlock()
	return id
}

// Entfernt das übergebene Flugzeug aus der Überwachung des MultiplayerManagers. Nach diesem Aufruf wird das Flugzeug nicht mehr gezeichnet.
func (self *MultiplayerManager) DestroyPlane(id PlaneId) {
	self.planesMutex.Lock()
	delete(self.planes, id)
	self.planesMutex.Unlock()

}

// Ändert für ein Flugzeug das Aussehen. Dabei kommen die gleichen Regeln wie bei CreatePlane zum Einsatz.
func (self *MultiplayerManager) ChangePlaneModel(id PlaneId, icao, airline, livery string) {
	self.planesMutex.RLock()
	plane, found := self.planes[id]
	self.planesMutex.RUnlock()
	if found {
		plane.CslAircraft = csl.CslManagerInstance.GetAircraft(icao, airline, livery)
		plane.updateAircraftData()
	}
}

// Ändert das Standard-ICAO-Kürzel für die Ermittlung des Flugzeugs, wenn das gewünschte ICAO nicht registriert ist.
func (self *MultiplayerManager) SetDefaultPlaneICAO(defaultICAO string) {
	csl.CslManagerInstance.SetDefaultICAO(defaultICAO)
}

//Liefert die aktuelle Konfiguration
func (self *MultiplayerManager) GetConfiguration() *Configuration {
	return self.configuration
}