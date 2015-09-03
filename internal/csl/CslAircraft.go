package csl
import (
	"fmt"
	"strings"
	"path/filepath"
	"errors"
	"os"
	"github.com/abieberbach/goplane/xplm/scenery"
	"github.com/abieberbach/goplanemp/internal/obj7"
)

type CslType int

const (
	CSL_TYPE_AUSTIN CslType = 0
	CSL_TYPE_OBJ7 CslType = 1
	CSL_TYPE_OBJ8 CslType = 2
)

type CslDrawType int

const (
	CSL_DRAW_LIGHTS CslDrawType = 0
	CSL_DRAW_LOW_LOD CslDrawType = 1
	CSL_DRAW_SOLID CslDrawType = 2
	CSL_DRAW_GLASS CslDrawType = 3
)

type CslAircraftAttachment struct {
	ObjFile             string
	ObjectReference     scenery.ObjectRef
	DrawType            CslDrawType
	NeedsAnimation      bool
	ObjectLoadRequested bool
	ReplaceTexture      string
}

type CslAircraft struct {
	AcfFile     string
	ObjFile     string
	HasGear     bool
	TextureFile string
	Icao        string
	Airline     string
	Livery      string
	CslType     CslType
	ObjInfo     *obj7.ObjectInfo
	Attachments []*CslAircraftAttachment
}

func (self *CslAircraft) validate(allPackages *CslPackages) (bool, string) {
	//1. Pflichtangaben prüfen
	if (self.CslType == CSL_TYPE_AUSTIN&&len(self.AcfFile) == 0) || (self.CslType == CSL_TYPE_OBJ7&&len(self.ObjFile) == 0) {
		return false, "aircraft file or object file is mandatory"
	}
	var err error
	//2. ACF-Datei prüfen
	self.AcfFile, err = replaceFileName(self.AcfFile, allPackages, true)
	if err != nil {
		return false, fmt.Sprintf("aircraft file \"%v\" is not correct: %v", self.AcfFile, err.Error())
	}
	//3. OBJ-Datei prüfen
	self.ObjFile, err = replaceFileName(self.ObjFile, allPackages, true)
	if err != nil {
		return false, fmt.Sprintf("object file \"%v\" is not correct: %v", self.ObjFile, err.Error())
	}
	//4. bei OBJ7-Modellen die OBJ-Datei einlesen
	if self.CslType == CSL_TYPE_OBJ7 {
		self.ObjInfo, err = obj7.ObjectManagerInstance.GetObject(self.ObjFile)
		if err != nil {
			return false, fmt.Sprintf("obj file \"%v\" is not correct: %v", self.ObjFile, err.Error())
		}
	}

	//5. bei OBJ8-Modellen die OBJ-Anhänge einlesen
	if self.CslType == CSL_TYPE_OBJ8 {
		for _, currentAttachment := range self.Attachments {
			currentAttachment.ObjFile, err = replaceFileName(currentAttachment.ObjFile, allPackages, true)
			if err != nil {
				return false, fmt.Sprintf("object attachment file \"%v\" is not correct: %v", currentAttachment.ObjFile, err.Error())
			}
			currentAttachment.ReplaceTexture, err = replaceFileName(currentAttachment.ReplaceTexture, allPackages, true)
			if err != nil {
				return false, fmt.Sprintf("object replace texture file \"%v\" is not correct: %v", currentAttachment.ReplaceTexture, err.Error())
			}
		}
	}

	//6. Textur-Datei prüfen
	self.TextureFile, err = replaceFileName(self.TextureFile, allPackages, true)
	if err != nil {
		return false, fmt.Sprintf("texture file \"%v\" is not correct: %v", self.TextureFile, err.Error())
	}
	return true, ""
}

func replaceFileName(fileName string, allPackages *CslPackages, correctName bool) (string, error) {
	if len(fileName) == 0 {
		return fileName, nil
	}
	//Pfade normalisieren
	fileName = strings.Replace(fileName, "\\", "/", -1)
	fileName = strings.Replace(fileName, ":", "/", -1)
	pathTokens := strings.Split(fileName, "/")
	cslPackage, found := allPackages.GetPackage(pathTokens[0])
	if !found {
		return fileName, errors.New(fmt.Sprintf("package \"%v\" not found", pathTokens[0]))
	}
	pathTokens[0] = cslPackage.BaseDirectory
	fileName = filepath.Join(pathTokens...)
	_, err := os.Stat(fileName)
	if correctName&&os.IsNotExist(err) {
		//Datei ist nicht vorhanden --> prüfen ob evt. nur die Schreibweise falsch war
		fileName = correctFileName(fileName)
		//nochmal prüfen
		_, err = os.Stat(fileName)
	}
	return fileName, err
}

func correctFileName(fileName string) string {
	baseDir := filepath.Dir(fileName)
	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir()&&path != baseDir {
			return filepath.SkipDir
		}
		if strings.EqualFold(fileName, path) {
			fileName = path
		}
		return nil
	})
	return fileName
}
