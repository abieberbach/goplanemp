//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package renderer
import "github.com/abieberbach/goplanemp/gl"


type cullInfo struct {
	model_view [16]float32 // The model view matrix, to get from local OpenGL to eye coordinates.
	proj       [16]float32 // Proj matrix - this is just a hack to use for gluProject.
	nea_clip   [4]float32  // Four clip planes in the form of Ax + By + Cz + D = 0 (ABCD are in the array.)
	far_clip   [4]float32  // They are oriented so the positive side of the clip plane is INSIDE the view volume.
	lft_clip   [4]float32
	rgt_clip   [4]float32
	bot_clip   [4]float32
	top_clip   [4]float32
}

func newCullInfo() *cullInfo {
	info := &cullInfo{}
	// First, just read out the current OpenGL matrices...do this once at setup because it's not the fastest thing to do.
	gl.GetFloatv(gl.MODELVIEW_MATRIX, &(info.model_view[0]))
	gl.GetFloatv(gl.PROJECTION_MATRIX, &(info.proj[0]))

	// Now...what the heck is this?  Here's the deal: the clip planes have values in "clip" coordinates of: Left = (1,0,0,1)
	// Right = (-1,0,0,1), Bottom = (0,1,0,1), etc.  (Clip coordinates are coordinates from -1 to 1 in XYZ that the driver
	// uses.  The projection matrix converts from eye to clip coordinates.)
	//
	// How do we convert a plane backward from clip to eye coordinates?  Well, we need the transpose of the inverse of the
	// inverse of the projection matrix.  (Transpose of the inverse is needed to transform a plane, and the inverse of the
	// projection is the matrix that goes clip -> eye.)  Well, that cancels out to the transpose of the projection matrix,
	// which is nice because it means we don't need a matrix inversion in this bit of sample code.

	// So this nightmare down here is simply:
	// clip plane * transpose (proj_matrix)
	// worked out for all six clip planes.  If you squint you can see the patterns:
	// L:  1  0 0 1
	// R: -1  0 0 1
	// B:  0  1 0 1
	// T:  0 -1 0 1
	// etc.

	info.lft_clip[0] = info.proj[0] + info.proj[3]; info.lft_clip[1] = info.proj[4] + info.proj[7]; info.lft_clip[2] = info.proj[8] + info.proj[11]; info.lft_clip[3] = info.proj[12] + info.proj[15];
	info.rgt_clip[0] = -info.proj[0] + info.proj[3]; info.rgt_clip[1] = -info.proj[4] + info.proj[7]; info.rgt_clip[2] = -info.proj[8] + info.proj[11]; info.rgt_clip[3] = -info.proj[12] + info.proj[15];

	info.bot_clip[0] = info.proj[1] + info.proj[3]; info.bot_clip[1] = info.proj[5] + info.proj[7]; info.bot_clip[2] = info.proj[9] + info.proj[11]; info.bot_clip[3] = info.proj[13] + info.proj[15];
	info.top_clip[0] = -info.proj[1] + info.proj[3]; info.top_clip[1] = -info.proj[5] + info.proj[7]; info.top_clip[2] = -info.proj[9] + info.proj[11]; info.top_clip[3] = -info.proj[13] + info.proj[15];

	info.nea_clip[0] = info.proj[2] + info.proj[3]; info.nea_clip[1] = info.proj[6] + info.proj[7]; info.nea_clip[2] = info.proj[10] + info.proj[11]; info.nea_clip[3] = info.proj[14] + info.proj[15];
	info.far_clip[0] = -info.proj[2] + info.proj[3]; info.far_clip[1] = -info.proj[6] + info.proj[7]; info.far_clip[2] = -info.proj[10] + info.proj[11]; info.far_clip[3] = -info.proj[14] + info.proj[15];

	return info
}

func (self *cullInfo) sphereDistanceSqr(x, y, z float32) float32 {
	xp, yp, zp := self.transformIntoEyeCoordinates(x, y, z)
	return xp * xp + yp * yp + zp * zp;
}

func (self *cullInfo) transformIntoEyeCoordinates(x, y, z float32) (xp, yp, zp float32) {
	xp = x * self.model_view[0] + y * self.model_view[4] + z * self.model_view[ 8] + self.model_view[12];
	yp = x * self.model_view[1] + y * self.model_view[5] + z * self.model_view[ 9] + self.model_view[13];
	zp = x * self.model_view[2] + y * self.model_view[6] + z * self.model_view[10] + self.model_view[14];
	return
}

func (self *cullInfo) sphereIsVisible(x, y, z, r float32) bool {
	// First: we transform our coordinate into eye coordinates from model-view.
	xp, yp, zp := self.transformIntoEyeCoordinates(x, y, z)

	// Now - we apply the "plane equation" of each clip plane to see how far from the clip plane our point is.
	// The clip planes are directed: positive number distances mean we are INSIDE our viewing area by some distance;
	// negative means outside.  So ... if we are outside by less than -r, the ENTIRE sphere is out of bounds.
	// We are not visible!  We do the near clip plane, then sides, then far, in an attempt to try the planes
	// that will eliminate the most geometry first...half the world is behind the near clip plane, but not much is
	// behind the far clip plane on sunny day.
	if (xp * self.nea_clip[0] + yp * self.nea_clip[1] + zp * self.nea_clip[2] + self.nea_clip[3] + r) < 0 {
		return false
	}
	if (xp * self.bot_clip[0] + yp * self.bot_clip[1] + zp * self.bot_clip[2] + self.bot_clip[3] + r) < 0 {
		return false
	}
	if (xp * self.top_clip[0] + yp * self.top_clip[1] + zp * self.top_clip[2] + self.top_clip[3] + r) < 0 {
		return false
	}
	if (xp * self.lft_clip[0] + yp * self.lft_clip[1] + zp * self.lft_clip[2] + self.lft_clip[3] + r) < 0 {
		return false
	}
	if (xp * self.rgt_clip[0] + yp * self.rgt_clip[1] + zp * self.rgt_clip[2] + self.rgt_clip[3] + r) < 0 {
		return false
	}
	if (xp * self.far_clip[0] + yp * self.far_clip[1] + zp * self.far_clip[2] + self.far_clip[3] + r) < 0 {
		return false
	}
	return true
}

func (self *cullInfo)convertTo2d(viewPoint [4]float32, x, y, z, w float32) (float32, float32) {
	xe := x * self.model_view[0] + y * self.model_view[4] + z * self.model_view[ 8] + w * self.model_view[12];
	ye := x * self.model_view[1] + y * self.model_view[5] + z * self.model_view[ 9] + w * self.model_view[13];
	ze := x * self.model_view[2] + y * self.model_view[6] + z * self.model_view[10] + w * self.model_view[14];
	we := x * self.model_view[3] + y * self.model_view[7] + z * self.model_view[11] + w * self.model_view[15];

	xc := xe * self.proj[0] + ye * self.proj[4] + ze * self.proj[ 8] + we * self.proj[12];
	yc := xe * self.proj[1] + ye * self.proj[5] + ze * self.proj[ 9] + we * self.proj[13];
	wc := xe * self.proj[3] + ye * self.proj[7] + ze * self.proj[11] + we * self.proj[15];

	xc /= wc;
	yc /= wc;

	return viewPoint[0] + (1.0 + xc) * viewPoint[2] / 2.0, viewPoint[1] + (1.0 + yc) * viewPoint[3] / 2.0
}


