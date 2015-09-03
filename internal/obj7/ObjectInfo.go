package obj7
import (
	"errors"
	"fmt"
)

type ObjectInfo struct {
	Path           string
	DefaultTexture string
	Lods           []*LodObjectInfo
}

func NewObjectInfo(objFile string) (*ObjectInfo, error) {
	info := &ObjectInfo{objFile, "", make([]*LodObjectInfo, 0)}
	//Objektdefinition laden
	parser, err := NewObjectParser(objFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not create object parser for \"%v\": %v", objFile, err))
	}
	parser.Parse()
	if parser.Err() != nil {
		return nil, errors.New(fmt.Sprintf("could not load object \"%v\": %v (line %v)", objFile, parser.Err(), parser.CurrentLine()))
	}
	info.DefaultTexture = parser.GetTexture()
	info.evaluateCommands(parser.GetCommands())
	info.calculateLodNormals()
	return info, nil
}

func (self *ObjectInfo) evaluateCommands(commands []*ObjectCommand) {
	var currentLod *LodObjectInfo
	if !hasLods(commands) {
		//es wurden keine LODs gefunden --> Standard LOD hinzufügen
		currentLod = NewLodObjectInfo(0, 40000)
		self.Lods = append(self.Lods, currentLod)
	}
	//über alle Kommandos iterieren und die Polygone und Lichter herausfiltern
	for _, command := range commands {
		switch command.CmdType {
		case CommandType_Attr:
			if (command.CmdId == CommandId_AttrLOD) {
				// wir haben einen neuen LOD-Abschnitt gefunden, ab jetzt diesen verwenden
				currentLod = NewLodObjectInfo(command.Attributes[0], command.Attributes[1])
				self.Lods = append(self.Lods, currentLod)
			}
		case CommandType_PtLine:
			if command.CmdId == CommandId_ObjLight {
				// For each light we've found, copy the data into our
				// own light vector
				for _, vrgb := range command.RGB {
					currentLod.Lights = append(currentLod.Lights, NewLightInfo(vrgb))
				}
			}
		case CommandType_Poly:
			indexes := make([]int32, 0)
			// alle Punkte im PointPool registrieren
			for _, currentVecTex := range command.St {
				index := currentLod.PointPool.AddPoint(currentVecTex.V, currentVecTex.St)
				indexes = append(indexes, index)
			}
			switch command.CmdId {
			case CommandId_ObjTri:
				for _, idx := range indexes {
					currentLod.TriangleList = append(currentLod.TriangleList, idx)
				}
			case CommandId_ObjTriFan:
				for n := 2; n < len(indexes); n++ {
					currentLod.TriangleList = append(currentLod.TriangleList, indexes[0])
					currentLod.TriangleList = append(currentLod.TriangleList, indexes[n - 1])
					currentLod.TriangleList = append(currentLod.TriangleList, indexes[n])
				}
			case CommandId_ObjTriStrip, CommandId_ObjQuadStrip:
				for n := 2; n < len(indexes); n++ {
					if ((n % 2) == 1) {
						currentLod.TriangleList = append(currentLod.TriangleList, indexes[n - 2])
						currentLod.TriangleList = append(currentLod.TriangleList, indexes[n])
						currentLod.TriangleList = append(currentLod.TriangleList, indexes[n - 1])
					} else {
						currentLod.TriangleList = append(currentLod.TriangleList, indexes[n - 2])
						currentLod.TriangleList = append(currentLod.TriangleList, indexes[n - 1])
						currentLod.TriangleList = append(currentLod.TriangleList, indexes[n])
					}
				}
			case CommandId_ObjQuad:
				for n := 3; n < len(indexes); n += 4 {
					currentLod.TriangleList = append(currentLod.TriangleList, indexes[n - 3])
					currentLod.TriangleList = append(currentLod.TriangleList, indexes[n - 2])
					currentLod.TriangleList = append(currentLod.TriangleList, indexes[n - 1])
					currentLod.TriangleList = append(currentLod.TriangleList, indexes[n - 3])
					currentLod.TriangleList = append(currentLod.TriangleList, indexes[n - 1])
					currentLod.TriangleList = append(currentLod.TriangleList, indexes[n])
				}
			}
		}
	}
}

func (self *ObjectInfo) calculateLodNormals() {
	for _, currentLod := range self.Lods {
		for n := 0; n < len(currentLod.TriangleList); n += 3 {
			currentLod.PointPool.CalcTriNormal(currentLod.TriangleList[n], currentLod.TriangleList[n + 1], currentLod.TriangleList[n + 2])
		}
		currentLod.PointPool.NormalizeNormals()
		currentLod.Dl = 0
	}
}

func hasLods(commands []*ObjectCommand) bool {
	for _, command := range commands {
		if command.CmdType == CommandType_Attr && command.CmdId == CommandId_AttrLOD {
			return true
		}
	}
	return false
}
