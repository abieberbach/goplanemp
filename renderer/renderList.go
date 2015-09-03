//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package renderer
import (
	"github.com/abieberbach/goplanemp"
	"github.com/abieberbach/goplane/xplm/planes"
	"github.com/abieberbach/goplanemp/internal/csl"
)


type planeRenderData struct {
	x           float32 // Positional info
	y           float32
	z           float32
	plane       *goplanemp.Plane
	tcas        bool    // Are we visible on TCAS?
	dist        float32
	drawState   planes.PlaneDrawState
	lightStatus goplanemp.LightStatus
}

type renderList struct {
	obj7Planes   []*planeRenderData
	obj8Planes   []*planeRenderData
}

func newRenderList() *renderList {
	result := &renderList{}
	result.obj7Planes = make([]*planeRenderData, 0)
	result.obj8Planes = make([]*planeRenderData, 0)
	return result
}

func (self *renderList) AddPlane(renderData *planeRenderData) {
	switch renderData.plane.CslAircraft.CslType {
	case csl.CSL_TYPE_OBJ7:
		self.obj7Planes = append(self.obj7Planes, renderData)
	case csl.CSL_TYPE_OBJ8:
		self.obj8Planes = append(self.obj8Planes, renderData)
	}

}

func (self *renderList) GetPlanes() ([]*planeRenderData, []*planeRenderData) {
	return self.obj7Planes, self.obj8Planes
}

func (self *renderList) GetAllPlanes() []*planeRenderData {
	result := make([]*planeRenderData, 0, len(self.obj7Planes) + len(self.obj8Planes))
	result = append(result, self.obj7Planes...)
	result = append(result, self.obj8Planes...)
	return result
}

func (self *renderList) Clear() {
	self.obj7Planes = self.obj7Planes[0:0]
	self.obj8Planes = self.obj8Planes[0:0]
}