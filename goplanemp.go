package goplanemp

//Positionsdaten des Objektes
type PlanePosition    struct {
	Lat       float64 //Breitengrad
	Lon       float64 //Längengrad
	Elevation float64 //Höhe MSL (in Fuss)
	Pitch     float32 //aktueller Pitch-Winkel
	Roll      float32 //aktueller Roll-Winkel
	Heading   float32 //aktueller Heading-Winkel
	Label     string  //Bezeichnung
}


// Muster für die Strobe- und Beacon-Lichter
type LightsPattern int

const (
	LightsPattern_Default    LightsPattern = 0    // Jets: one strobe flash, short beacon (-*---*---*---)
	LightsPattern_EADS       LightsPattern = 1    // Airbus+EADS: strobe flashes twice (-*-*-----*-*--), short beacon
	LightsPattern_GA         LightsPattern = 2    // GA: one strobe flash, long beacon (-*--------*---)
)

//Enthält die Informationen über die Beleuchtung
type LightStatus struct {
	TimeOffset    int           //Offset für die Berechnung der Strobe und Beacon-Lichter
	LandingLights bool          //Zustand der Landelichter
	BeaconLights  bool          //Zustand der Beacon-Lights
	StrobeLights  bool          //Zustand der Strobe-Lights
	NavLights     bool          //Zustand der Navigationslichter
	TaxiLights    bool          //Zustand der Taxi-Lights
	FlashPattern  LightsPattern //Muster für die Beacon- und Strobe-Lights
}

//Enthält die Daten für die Steuerflächen
type PlaneSurfaces struct {
	GearPosition    float32     //Position des Fahrwerks [0 (eingefahren) .. 1 (ausgefahren)]
	FlapRatio       float32     //Position der Flaps [0 (eingefahren) .. 1 (ausgefahren)]
	SpoilerRatio    float32     //Position der Spoiler [0 (eingefahren) .. 1 (ausgefahren)]
	SpeedBrakeRatio float32     //Position der Speedbrakes [0 (eingefahren) .. 1 (ausgefahren)]
	SlatRatio       float32     //Position des Slats [0 (eingefahren) .. 1 (ausgefahren)]
	WingSweep       float32     //Position der Flügel
	Thrust          float32     //aktueller Schub [-1 (voller Umkehrschub) .. [0] (kein Schub) .. 1 (voller Schub)]
	YokePitch       float32     //aktuller Pitch-Winkel des Yoke
	YokeHeading     float32     //aktuller Heading-Winkel des Yoke
	YokeRoll        float32     //aktuller Roll-Winkel des Yoke
	Lights          LightStatus //aktueller Zustand der Lichter
}

//Transponder Modus
type TransponderMode int

const (
	TransponderMode_Standby    TransponderMode = 0 //Modus: Standby
	TransponderMode_Mode3A        TransponderMode = 1 //Modus: 3a
	TransponderMode_ModeC        TransponderMode = 2 //Modus: C
	TransponderMode_ModeC_Low    TransponderMode = 3 //Modus: C Low
	TransponderMode_ModeC_Ident    TransponderMode = 4 //Modus: C Ident
)

//Enthält die Radar-Informationen
type PlaneRadar struct {
	Code int
	Mode TransponderMode
}

// Ergebnis des Callback Aufrufs
type PlaneCallbackResult int

const (
	Data_Unavailable PlaneCallbackResult = 0 // Information kann nicht ermittelt werden
	Data_Unchanged PlaneCallbackResult = 1 // Information hat sich seit dem letzten Aufruf nicht geändert
	Data_NewData PlaneCallbackResult = 2 // Informationen haben sich seit dem letzten Aufruf geändert
)

//Art der abgefragten Daten
type PlaneDataType int

const (
	DataType_Position    PlaneDataType = 0 //Positionsdaten des Objekts
	DataType_Surfaces    PlaneDataType = 1 //Steuerinformationen des Objekts
	DataType_Radar       PlaneDataType = 2 //Radarinformationen des Objekts
)

// ID des Flugzeugs
type PlaneId int

// Callback-Funktion zur Ermittlung der gewünschten Daten. Alle nicht abgefragten Daten haben als Ergebnis "nil".
type PlaneDataFunc func(id PlaneId, dataType PlaneDataType, ref interface{}) (PlaneCallbackResult, *PlanePosition, *PlaneSurfaces, *PlaneRadar)

//Enthält die Einstellungen für die Multiplayer-Bibliothek
type Configuration struct {
	FullDistance           float32 //Entfernung in km bis zu der das Objekt voll gezeichnet wird, danach nur noch die Lichter
	Resolution             int //Auflösung für OBJ7-Texturen (0[schlechteste Qualität] .. 5 [beste Qualität])
	CalculationMod         int //Modulo-Wert für die Abfrage der Daten (je niedriger der Wert umso öfter werden die Daten abgerufen)
	EnableRenderStatistics bool //Flag ob Statistiken zum DefaultRenderer als DataRef zur Verfügung gestellt werden sollen
	RenderAircraftLabels   bool //Flag ob Labels gezeichnet werden sollen
	UseAnyLitTexture       bool //Flag ob bei fehlender Nachttextur eine beliebige andere Nachttextur des Models genommen werden soll
	MergeLitTextures       bool //Flag ob die Nachttexturen mit der Tagtextur gemischt werden soll (besseres Aussehen)
	UseObj8Lights          bool //Flag ob bei OBJ7-Objekten anstatt der eigenen Lichter die neuen OBJ8 Lichter genutzt werden sollen
	DisableObj8LightSpills bool //Flag ob bei OBJ8-Objekten die LightSpills entfernt werden sollen (bessere Performance)
}

// Erzeugt eine neue Standardkonfiguration
//	FullDistance           = 5.0 km
//	Resolution             = 3
//	CalculationMod         = 50
//	EnableRenderStatistics = false
//	RenderAircraftLabels   = false
//	UseAnyLitTexture       = false
//	MergeLitTextures       = false
//	UseObj8Lights          = false
//	DisableObj8LightSpills = true
//
func NewDefaultConfiguration() *Configuration {
	return &Configuration{5.0, 3, 50,  false, false, false, false, false, true}
}

