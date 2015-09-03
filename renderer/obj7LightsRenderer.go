//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package renderer
import (
	"github.com/abieberbach/goplane/xplm/graphics"
	"github.com/abieberbach/goplane/xplm/camera"
	"math"
	"github.com/abieberbach/goplanemp"
	"github.com/abieberbach/goplane/xplm/processing"
	"github.com/abieberbach/goplanemp/internal/texture"
	"github.com/abieberbach/goplanemp/internal/obj7"
	"github.com/abieberbach/goplanemp/gl"
)

func (self *DefaultRenderer) beginObj7LightDrawing() {
	// Setup OpenGL for the drawing
	lightTexNum := texture.TextureManagerInstance.GetLightTextureNumber()
	if lightTexNum == -1 {
		return
	}
	graphics.SetGraphicsState(true, 1, false, true, true, true, false)
	graphics.BindTexture2d(lightTexNum, 0)
}

var greenNavLightDl uint32 = 0
var redNavLightDl uint32 = 0
var beaconLightDl uint32 = 0
var strobeLightDl uint32 = 0
var landingLightDl uint32 = 0
var taxiLightDl uint32 = 0
var otherLightDl uint32 = 0

func (self *DefaultRenderer) drawObj7Lights(renderData *planeRenderData, fieldOfViewValue float32, cameraPos camera.CameraPosition) {
	gl.MatrixMode(gl.MODELVIEW);
	gl.PushMatrix();
	gl.Translatef(renderData.x, renderData.y, renderData.z);
	gl.Rotatef(renderData.plane.PositionData.Heading, 0.0, -1.0, 0.0);
	gl.Rotatef(renderData.plane.PositionData.Pitch, 1.0, 0.0, 0.0);
	gl.Rotatef(renderData.plane.PositionData.Roll, 0.0, 0.0, -1.0);

	lodInfo := calculateLOD(renderData.plane.CslAircraft.ObjInfo, renderData.dist)
	lights := renderData.lightStatus

	offset := renderData.plane.SurfacesData.Lights.TimeOffset
	// flash frequencies
	baseTime := int(processing.GetElapsedTime() * 1000.0) + offset
	lights.BeaconLights = calculateBeaconFlash(lights.BeaconLights, renderData.plane.SurfacesData.Lights.FlashPattern, baseTime)
	lights.StrobeLights = calculateStrobeFlash(lights.StrobeLights, renderData.plane.SurfacesData.Lights.FlashPattern, baseTime)

	// Find our distance from the camera
	dx := cameraPos.X - renderData.x
	dy := cameraPos.Y - renderData.y
	dz := cameraPos.Z - renderData.z
	distance := float32(math.Sqrt(float64((dx * dx) + (dy * dy) + (dz * dz))))

	// Convert to NM
	distance *= MetersToNM

	// Scale based on our FOV and Zoom. I did my initial
	// light adjustments at a FOV of 60 so thats why
	// I divide our current FOV by 60 to scale it appropriately.
	distance *= fieldOfViewValue / 60.0
	distance /= cameraPos.Zoom

	// Calculate our light size. This is piecewise linear. I noticed
	// that light size changed more rapidly when closer than 3nm so
	// I have a separate equation for that.
	var size float32
	if (distance <= 3.6) {
		size = (10 * distance) + 1
	} else {
		size = (6.7 * distance) + 12
	}

	for _, currentLight := range lodInfo.Lights {
		gl.MatrixMode(gl.MODELVIEW)
		gl.PushMatrix()
		// First we translate to our coordinate system and move the origin
		// to the center of our lights.
		gl.Translatef(currentLight.XYZ[0], currentLight.XYZ[1], currentLight.XYZ[2])

		// Now we undo the rotation of the plane
		gl.Rotatef(-renderData.plane.PositionData.Pitch, 1.0, 0.0, 0.0)
		gl.Rotatef(-renderData.plane.PositionData.Heading, 0.0, -1.0, 0.0)
		gl.Rotatef(-renderData.plane.PositionData.Roll, 0.0, 0.0, -1.0)

		// Now we undo the rotation of the camera
		// NOTE: The order and sign of the camera is backwards
		// from what we'd expect (the plane rotations) because
		// the camera works backwards. If you pan right, everything
		// else moves left!
		gl.Rotatef(cameraPos.Pitch, 1.0, 0.0, 0.0)
		gl.Rotatef(cameraPos.Heading, 0.0, -1.0, 0.0)
		gl.Rotatef(cameraPos.Roll, 0.0, 0.0, -1.0)

		drawSingleLight(currentLight, size, distance, &lights)
		// Put OpenGL back how we found it
		gl.PopMatrix()
	}
	gl.PopMatrix()
}

func drawSingleLight(light *obj7.LightInfo, size float32, distance float32, lights *goplanemp.LightStatus) {
	// Finally we can draw our lights
	switch light.LightType {
	case obj7.LightType_RedNavigation:
		if lights.NavLights {
			if redNavLightDl == 0 {
				redNavLightDl = gl.GenLists(1)
				gl.NewList(redNavLightDl, gl.COMPILE)
				gl.Begin(gl.QUADS)
				gl.Color4f(light.RGB_Float[0], light.RGB_Float[1], light.RGB_Float[2], light.RGB_Float[3])
				gl.TexCoord2f(0, 0.5)
				gl.Vertex2f(-(size / 2.0), -(size / 2.0))
				gl.TexCoord2f(0, 1.0)
				gl.Vertex2f(-(size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.25, 1.0)
				gl.Vertex2f((size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.25, 0.5)
				gl.Vertex2f((size / 2.0), -(size / 2.0))
				gl.End()
				gl.EndList()
			}
			gl.CallList(redNavLightDl)
		}
	case obj7.LightType_GreenNavigation:
		if lights.NavLights {
			if greenNavLightDl == 0 {
				greenNavLightDl = gl.GenLists(1)
				gl.NewList(greenNavLightDl, gl.COMPILE)
				gl.Begin(gl.QUADS)
				gl.Color4f(light.RGB_Float[0], light.RGB_Float[1], light.RGB_Float[2], light.RGB_Float[3])
				gl.TexCoord2f(0, 0.5)
				gl.Vertex2f(-(size / 2.0), -(size / 2.0))
				gl.TexCoord2f(0, 1.0)
				gl.Vertex2f(-(size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.25, 1.0)
				gl.Vertex2f((size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.25, 0.5)
				gl.Vertex2f((size / 2.0), -(size / 2.0))
				gl.End()
				gl.EndList()
			}
			gl.CallList(greenNavLightDl)
		}
	case obj7.LightType_Beacon:
		if lights.BeaconLights {
			if beaconLightDl == 0 {
				beaconLightDl = gl.GenLists(1)
				gl.NewList(beaconLightDl, gl.COMPILE)
				gl.Begin(gl.QUADS)
				gl.Color4f(light.RGB_Float[0], light.RGB_Float[1], light.RGB_Float[2], light.RGB_Float[3])
				gl.TexCoord2f(0, 0.5)
				gl.Vertex2f(-(size / 2.0), -(size / 2.0))
				gl.TexCoord2f(0, 1.0)
				gl.Vertex2f(-(size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.25, 1.0)
				gl.Vertex2f((size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.25, 0.5)
				gl.Vertex2f((size / 2.0), -(size / 2.0))
				gl.End()
				gl.EndList()
			}
			gl.CallList(beaconLightDl)
		}
	case obj7.LightType_Strobe:
		if lights.StrobeLights {
			if strobeLightDl == 0 {
				strobeLightDl = gl.GenLists(1)
				gl.NewList(strobeLightDl, gl.COMPILE)
				gl.Begin(gl.QUADS)
				gl.Color4f(light.RGB_Float[0], light.RGB_Float[1], light.RGB_Float[2], light.RGB_Float[3])
				gl.TexCoord2f(0.25, 0.0)
				gl.Vertex2f(-(size / 1.5), -(size / 1.5))
				gl.TexCoord2f(0.25, 0.5)
				gl.Vertex2f(-(size / 1.5), (size / 1.5))
				gl.TexCoord2f(0.50, 0.5)
				gl.Vertex2f((size / 1.5), (size / 1.5))
				gl.TexCoord2f(0.50, 0.0)
				gl.Vertex2f((size / 1.5), -(size / 1.5))
				gl.End()
				gl.EndList()
			}
			gl.CallList(strobeLightDl)
		}
	case obj7.LightType_Landing:
		if lights.LandingLights {
			// BEN SEZ: modulate the _alpha to make this dark, not
			// the light color.  Otherwise if the sky is fairly light the light
			// will be darker than the sky, which looks f---ed during the day.
			color := light.RGB_Float
			if (color[0] < 0.0) {
				color[0] = 0.0
			}
			if (color[0] < 0.0) {
				color[0] = 0.0
			}
			if (color[0] < 0.0) {
				color[0] = 0.0
			}
			color[3] *= (distance * -0.05882) + 1.1764
			gl.Color4f(color[0], color[1], color[2], color[3])
			if landingLightDl == 0 {
				landingLightDl = gl.GenLists(1)
				gl.NewList(landingLightDl, gl.COMPILE)
				gl.Begin(gl.QUADS)
				gl.TexCoord2f(0.25, 0.0)
				gl.Vertex2f(-(size / 2.0), -(size / 2.0))
				gl.TexCoord2f(0.25, 0.5)
				gl.Vertex2f(-(size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.50, 0.5)
				gl.Vertex2f((size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.50, 0.0)
				gl.Vertex2f((size / 2.0), -(size / 2.0))
				gl.End()
				gl.EndList()
			}
			gl.CallList(landingLightDl)
		}
	case obj7.LightType_Taxi:
		if lights.LandingLights {
			// BEN SEZ: modulate the _alpha to make this dark, not
			// the light color.  Otherwise if the sky is fairly light the light
			// will be darker than the sky, which looks f---ed during the day.
			color := light.RGB_Float
			if (color[0] < 0.0) {
				color[0] = 0.0
			}
			if (color[0] < 0.0) {
				color[0] = 0.0
			}
			if (color[0] < 0.0) {
				color[0] = 0.0
			}
			color[3] *= (distance * -0.05882) + 1.1764
			gl.Color4f(color[0], color[1], color[2], color[3])
			if taxiLightDl == 0 {
				taxiLightDl = gl.GenLists(1)
				gl.NewList(taxiLightDl, gl.COMPILE)
				gl.Begin(gl.QUADS)
				gl.TexCoord2f(0.25, 0.0)
				gl.Vertex2f(-(size / 2.0), -(size / 2.0))
				gl.TexCoord2f(0.25, 0.5)
				gl.Vertex2f(-(size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.50, 0.5)
				gl.Vertex2f((size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.50, 0.0)
				gl.Vertex2f((size / 2.0), -(size / 2.0))
				gl.End()
				gl.EndList()
			}
			gl.CallList(taxiLightDl)
		}
	default:
		if (lights.NavLights) {
			gl.Color4f(light.RGB_Float[0], light.RGB_Float[1], light.RGB_Float[2], light.RGB_Float[3])
			if otherLightDl == 0 {
				otherLightDl = gl.GenLists(1)
				gl.NewList(otherLightDl, gl.COMPILE)
				gl.Begin(gl.QUADS)
				gl.TexCoord2f(0, 0.5)
				gl.Vertex2f(-(size / 2.0), -(size / 2.0))
				gl.TexCoord2f(0, 1.0)
				gl.Vertex2f(-(size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.25, 1.0)
				gl.Vertex2f((size / 2.0), (size / 2.0))
				gl.TexCoord2f(0.25, 0.5)
				gl.Vertex2f((size / 2.0), -(size / 2.0))
				gl.End()
				gl.EndList()
			}
			gl.CallList(otherLightDl)
		}
	}
}

func calculateBeaconFlash(bcnLights bool, flashPattern goplanemp.LightsPattern, baseTime int) bool {
	if (bcnLights) {
		bcnLights = false
		x := baseTime % 1200
		switch flashPattern {
		case goplanemp.LightsPattern_EADS:
			// EADS pattern: two flashes every 1.2 seconds
			bcnLights = x < 120 || (x > 240 && x < 360)
		case goplanemp.LightsPattern_GA:
			// GA pattern: 900ms / 1200ms
			bcnLights = (baseTime % 2100) < 900
		case goplanemp.LightsPattern_Default:
			fallthrough
		default:
			// default pattern: one flash every 1.2 seconds
			bcnLights = x < 120
		}
	}
	return bcnLights
}

func calculateStrobeFlash(strbLights bool, flashPattern goplanemp.LightsPattern, baseTime int) bool {
	if strbLights {
		strbLights = false
		x := baseTime % 1700
		switch flashPattern {
		case goplanemp.LightsPattern_EADS:
			strbLights = x < 80 || (x > 260 && x < 340)
		case goplanemp.LightsPattern_GA:
			// similar to the others.. but a little different frequency :)
			strbLights = (baseTime % 1900) < 100
		case goplanemp.LightsPattern_Default:
			fallthrough
		default:
			strbLights = x < 80
		}
	}
	return strbLights
}