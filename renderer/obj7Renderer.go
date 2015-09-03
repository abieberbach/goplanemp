package renderer
import (
	"github.com/abieberbach/goplane/xplm/graphics"
	"unsafe"
	"github.com/abieberbach/goplanemp/internal/texture"
	"github.com/abieberbach/goplanemp/internal/obj7"
	"github.com/abieberbach/goplanemp/gl"
)

func (self *DefaultRenderer) drawObj7Model(renderData *planeRenderData) {
	if renderData.plane.CslAircraft.ObjInfo == nil {
		//es wurde kein Objekt geladen --> nichts zeichnen
		return
	}
	// Find out what LOD we need to draw
	lodInfo := calculateLOD(renderData.plane.CslAircraft.ObjInfo, renderData.dist)
	if lodInfo == nil {
		//es wurde kein gutes LOD gefunden --> nichts zeichnen
		return
	}
	//PointPool ist leer
	if lodInfo.PointPool.Size() == 0 && lodInfo.Dl == 0 {
		return
	}
	gl.MatrixMode(gl.MODELVIEW);
	gl.PushMatrix();
	gl.Translatef(renderData.x, renderData.y, renderData.z);
	gl.Rotatef(renderData.plane.PositionData.Heading, 0.0, -1.0, 0.0);
	gl.Rotatef(renderData.plane.PositionData.Pitch, 1.0, 0.0, 0.0);
	gl.Rotatef(renderData.plane.PositionData.Roll, 0.0, 0.0, -1.0);

	textureId, litTextureId := renderData.plane.TextureNum, renderData.plane.LitTextureNum
	useNight := texture.TextureManagerInstance.UseLitTexture()&&litTextureId != -1
	texUnits := 1
	if useNight {
		texUnits++
	} else {
		litTextureId = -1
	}
	graphics.SetGraphicsState(true, texUnits, true, true, true, true, true)
	graphics.BindTexture2d(textureId, 0)
	if litTextureId != -1 {
		graphics.BindTexture2d(litTextureId, 1)
	}
	gl.TexEnvi(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)
	if litTextureId != -1 {
		gl.ActiveTextureARB(gl.TEXTURE1)
		gl.TexEnvi(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.ADD)
		gl.ActiveTextureARB(gl.TEXTURE0)
	}

	if lodInfo.Dl == 0 {
		lodInfo.Dl = gl.GenLists(1)

		var xpBuffer int32
		// See if the card even has VBO. If it does, save xplane's pointer
		// and bind to 0 for us.
		gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING_ARB, &xpBuffer)
		gl.BindBufferARB(gl.ARRAY_BUFFER_ARB, 0)
		// Save XPlanes OpenGL state
		gl.PushClientAttrib(gl.CLIENT_ALL_ATTRIB_BITS)
		// Setup OpenGL pointers to our pool
		lodInfo.PointPool.PreparePoolToDraw()
		// Enable vertex data sucking
		gl.EnableClientState(gl.VERTEX_ARRAY)
		// Enable normal array sucking
		gl.EnableClientState(gl.NORMAL_ARRAY)
		// Enable texture coordinate data sucking
		gl.ClientActiveTextureARB(gl.TEXTURE1)
		gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
		gl.ClientActiveTextureARB(gl.TEXTURE0)
		gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
		// Disable colors - maybe x-plane left it around.
		gl.DisableClientState(gl.COLOR_ARRAY)

		gl.NewList(lodInfo.Dl, gl.COMPILE)
		// Kick OpenGL and draw baby!
		gl.DrawElements(gl.TRIANGLES, int32(len(lodInfo.TriangleList)), gl.UNSIGNED_INT, unsafe.Pointer(&(lodInfo.TriangleList[0])))
		gl.EndList()

		// Disable vertex data sucking
		gl.DisableClientState(gl.VERTEX_ARRAY)
		// Disable texture coordinate data sucking
		gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
		// Disable normal array sucking
		gl.DisableClientState(gl.NORMAL_ARRAY)

		// Restore Xplane's OpenGL State
		gl.PopClientAttrib()

		// If we bound before, we need to put xplane back where it was
		gl.BindBufferARB(gl.ARRAY_BUFFER_ARB, uint32(xpBuffer))

		lodInfo.TriangleList = make([]int32, 0)
		lodInfo.PointPool.Purge()
	}
	gl.CallList(lodInfo.Dl)
	gl.PopMatrix()
}

func calculateLOD(objInfo *obj7.ObjectInfo, distance float32) *obj7.LodObjectInfo {
	for _, lodInfo := range objInfo.Lods {
		if distance >= lodInfo.NearDist && distance <= lodInfo.FarDist {
			return lodInfo
		}
	}
	return nil
}