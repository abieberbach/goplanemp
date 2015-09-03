//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package goplanemp

//Schnittstelle für eine eigene Renderer-Implementierung
type Renderer interface {
	//Initialisiert den Renderer.
	Init()
	//Zeichnet alle übergebenen Flugzeuge.
	Render(isBlend, renderAircraftLabels bool, aircrafts map[PlaneId]*Plane, ref interface{})
}