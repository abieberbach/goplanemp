package renderer
import (
	"github.com/abieberbach/goplane/xplm/graphics"
	"image/color"
	"github.com/abieberbach/goplanemp/gl"
)

func (self *DefaultRenderer) drawPlaneLabels(renderList *renderList, cullInfo *cullInfo, labelDist float32) {
	var vp [4] float32;
	gl.GetFloatv(gl.VIEWPORT, &vp[0]);

	gl.MatrixMode(gl.PROJECTION);
	gl.PushMatrix();
	gl.LoadIdentity();
	gl.Ortho(0, float64(vp[2]), 0, float64(vp[3]), -1, 1);

	gl.MatrixMode(gl.MODELVIEW);
	gl.PushMatrix();
	gl.LoadIdentity();

	c := []float32{1, 1, 0, 1 }
	for _, plane := range renderList.GetAllPlanes() {
		if plane.dist < labelDist {
			x, y := cullInfo.convertTo2d(vp, plane.x, plane.y, plane.z, 1.0);
			rat := 1.0 - (plane.dist / MAX_LABEL_DIST);
			c[0] = 0.5 + 0.5 * rat
			c[1] = c[0]
			c[2] = 0.5 - 0.5 * rat; // gray -> yellow - no alpha in the SDK - foo!
			graphics.DrawString(color.RGBA{uint8(255.0 * c[0]), uint8(255.0 * c[1]), uint8(255.0 * c[2]), 0}, int(x), int(y + 10), plane.plane.PositionData.Label, []int32{}, graphics.Font_Basic)
		}
	}
	gl.MatrixMode(gl.PROJECTION);
	gl.PopMatrix();
	gl.MatrixMode(gl.MODELVIEW);
	gl.PopMatrix();
}