package renderer
import (
	"github.com/abieberbach/goplane/xplm/scenery"
	"os"
	"path/filepath"
	"github.com/abieberbach/goplane/extra/logging"
	"bufio"
	"fmt"
	"strings"
	"github.com/abieberbach/goplanemp/internal/csl"
)

func (self *DefaultRenderer) drawObj8Tanslucent(renderData *planeRenderData) {
	self.drawAttachments(renderData, true)
}

func (self *DefaultRenderer) drawObj8Solid(renderData *planeRenderData) {
	self.drawAttachments(renderData, false)
}

func (self *DefaultRenderer) drawAttachments(renderData *planeRenderData, translucent bool) {
	currentDrawInfo := scenery.NewDrawInfoLocal(renderData.x, renderData.y, renderData.z, renderData.plane.PositionData.Pitch, renderData.plane.PositionData.Heading, renderData.plane.PositionData.Roll)
	for _, currentAttachment := range renderData.plane.CslAircraft.Attachments {
		if currentAttachment.ObjectReference == nil {
			go self.loadAttachment(currentAttachment)
		} else if (translucent && currentAttachment.DrawType == csl.CSL_DRAW_GLASS) || (!translucent && currentAttachment.DrawType != csl.CSL_DRAW_GLASS) {
			scenery.DrawObjects(currentAttachment.ObjectReference, 1, []scenery.DrawInfo{currentDrawInfo}, true, false)
		}

	}
}

func (self *DefaultRenderer) loadAttachment(attachment *csl.CslAircraftAttachment) {
	if attachment.ObjectLoadRequested {
		//Objekt wird bereits geladen --> ignorieren
		return
	}
	attachment.ObjectLoadRequested = true
	if self.configuration.DisableObj8LightSpills || len(attachment.ReplaceTexture) > 0 {
		//es sollen die OBJ8-LightSpills ausgeblendet werden bzw. es gibt eine Austauschtextur --> Objekt ersetzen
		var err error
		attachment.ObjFile, err = self.replaceObject(attachment.ObjFile, attachment.ReplaceTexture)
		if err != nil {
			logging.Errorf("could not replace object: %v (obj file = %v)", err, attachment.ObjFile)
		}
	}
	scenery.LoadObjectAsync(attachment.ObjFile, self.objectLoaded, attachment)
}

func (self *DefaultRenderer) replaceObject(objFilePath, replaceTexturePath string) (string, error) {
	objFile, err := os.Open(objFilePath)
	if err != nil {
		return "", err
	}
	defer objFile.Close()
	replaceObjFile, err := os.Create(filepath.Join(filepath.Dir(objFilePath), "alt_" + filepath.Base(objFilePath)))
	if err != nil {
		return "", err
	}
	defer replaceObjFile.Close()
	scanner := bufio.NewScanner(objFile)
	writer := bufio.NewWriter(replaceObjFile)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return "", nil
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "TEXTURE ") &&len(replaceTexturePath) > 0 {
			//Zeile mit der Texturangabe gefunden --> neue Textur setzen
			line = "TEXTURE " + filepath.Base(replaceTexturePath)
		}
		if strings.HasPrefix(line, "LIGHT_SPILL_CUSTOM ") &&self.configuration.DisableObj8LightSpills {
			//Light-Spills sollen ignoriert werden --> aktuelle Zeile ignorieren
			continue
		}
		_, err = fmt.Fprintln(writer, line)
		if err != nil {
			return "", nil
		}
	}
	writer.Flush()
	return replaceObjFile.Name(), nil
}

func (self *DefaultRenderer) objectLoaded(objRef scenery.ObjectRef, ref interface{}) {
	attachment := ref.(*csl.CslAircraftAttachment)
	attachment.ObjectReference = objRef
	if len(attachment.ReplaceTexture) > 0 || self.configuration.DisableObj8LightSpills {
		os.Remove(attachment.ObjFile)
	}
}