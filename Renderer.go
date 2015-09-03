package goplanemp

//Schnittstelle für eine eigene Renderer-Implementierung
type Renderer interface {
	//Initialisiert den Renderer.
	Init()
	//Zeichnet alle übergebenen Flugzeuge.
	Render(isBlend, renderAircraftLabels bool, aircrafts map[PlaneId]*Plane, ref interface{})
}