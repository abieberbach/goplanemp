package texture
import (
	"github.com/abieberbach/goplane/xplm/dataAccess"
	"github.com/abieberbach/goplane/xplm/graphics"
	"errors"
	"fmt"
	"image"
	"os"
	"image/png"
	"github.com/disintegration/imaging"
	"unsafe"
	"strings"
	"path/filepath"
	"github.com/abieberbach/goplanemp/gl"
	"github.com/abieberbach/goplanemp/glu"
)

var TextureManagerInstance = &TextureManager{make(map[string][2]int32), 0, false, false, -1, nil, false}

type TextureManager struct {
	cache                map[string][2]int32
	resolution           int
	mergeLitTextures     bool
	useAnyLitTexture     bool
	lightsTextureNum     int32
	nightLightingDataRef dataAccess.DataRef
	texClampAvail        bool
}

func (self *TextureManager) Init(resolution int, lightsFile string, mergeLitTexture, useAnyLitTexture bool) (err error) {
	self.resolution = resolution
	self.mergeLitTextures = mergeLitTexture
	self.useAnyLitTexture = useAnyLitTexture
	self.lightsTextureNum, _, err = self.LoadTexture(lightsFile, true)
	self.nightLightingDataRef, _ = dataAccess.FindDataRef("sim/graphics/scenery/percent_lights_on");
	extensions := gl.GoStr(gl.GetString(gl.EXTENSIONS))
	self.texClampAvail = strings.Contains(extensions, "GL_SGI_texture_edge_clamp") || strings.Contains(extensions, "GL_SGIS_texture_edge_clamp") || strings.Contains(extensions, "GL_ARB_texture_edge_clamp") || strings.Contains(extensions, "GL_EXT_texture_edge_clamp")
	return
}

func (self *TextureManager) GetLightTextureNumber() int32 {
	return self.lightsTextureNum
}

func (self *TextureManager) UseLitTexture() bool {
	return self.nightLightingDataRef != nil&&dataAccess.GetFloatData(self.nightLightingDataRef) > 0.25
}

func (self *TextureManager) LoadTexture(textureFile string, forceMaxResolution bool) (int32, int32, error) {
	textureNums, found := self.cache[textureFile]
	if found {
		//Texturen wurden bereits geladen --> diese verwenden
		return textureNums[0], textureNums[1], nil
	}
	//Texturen wurden noch nicht geladen
	texImage, err := self.loadTextureFromFile(textureFile, true, forceMaxResolution)
	if err != nil {
		//Haupttextur konnte nicht geladen werden
		return -1, -1, errors.New(fmt.Sprintf("could not load texture \"%v\": %v", textureFile, err))
	}
	textureNums = [2]int32{graphics.GenerateTextureNumbers(1)[0], -1}
	err = self.bindTexture(textureNums[0], texImage)
	if err != nil {
		//Haupttextur konnte nicht gebunden werden
		return -1, -1, errors.New(fmt.Sprintf("could not bind texture \"%v\": %v", textureFile, err))
	}
	litTextureName, hasLitTexture := self.findLitTexture(textureFile)
	if hasLitTexture {
		litTexImage, err := self.loadTextureFromFile(litTextureName, true, forceMaxResolution)
		if err != nil {
			//Nachttextur konnte nicht geladen werden
			return -1, -1, errors.New(fmt.Sprintf("could not load lit texture \"%v\": %v", litTextureName, err))
		}
		if self.mergeLitTextures {
			//Nachttextur mit der normalen Textur mischen, damit es Nachts besser aussieht
			mergeTextures(texImage, litTexImage)
		}
		textureNums[1] = graphics.GenerateTextureNumbers(1)[0]
		err = self.bindTexture(textureNums[1], litTexImage)
		if err != nil {
			//Nachttextur konnte nicht gebunden werden
			return -1, -1, errors.New(fmt.Sprintf("could not bind lit texture \"%v\": %v", litTextureName, err))
		}

	}
	self.cache[textureFile]=textureNums
	return textureNums[0], textureNums[1], nil
}

func (self *TextureManager) loadTextureFromFile(texturePath string, isMagenta, forceMaxResolution bool) (*image.NRGBA, error) {
	textureFile, err := os.Open(texturePath)
	if err != nil {
		return nil, err
	}
	defer textureFile.Close()
	img, err := png.Decode(textureFile)
	if err != nil {
		return nil, err
	}
	//Bild für OpenGL drehen
	texture := imaging.FlipV(img)
	if !forceMaxResolution {
		//es soll nicht die max. Texturauflösung verwendet werden --> Bild entsprechend halbieren
		for res := 5 - self.resolution; res > 0; res-- {
			size := texture.Bounds().Size()
			texture = imaging.Resize(texture, size.X / 2, size.Y / 2, imaging.Linear)
		}
	}
	if isMagenta {
		convertImageToAlpha(texture)
	}
	return texture, nil
}

func (self *TextureManager) bindTexture(textureNum int32, texture *image.NRGBA) error {
	graphics.BindTexture2d(textureNum, 0)
	size := texture.Rect.Size()
	glu.Build2DMipmaps(gl.TEXTURE_2D, 4, size.X, size.Y, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&texture.Pix[0]))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	value := gl.CLAMP
	if self.texClampAvail {
		value = gl.CLAMP_TO_EDGE
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int32(value))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int32(value))
	err := gl.GetError()
	if err > 0 {
		return errors.New(fmt.Sprintf("Texture load got OGL err: %v", err))
	}
	return nil
}

func (self *TextureManager) findLitTexture(texturePath string) (string, bool) {
	pos := strings.LastIndex(strings.ToLower(texturePath), ".png")
	//1. prüfen ob gleichnamige Datei mit _LIT vorhanden ist
	litTexureName := texturePath[:pos] + "_LIT" + texturePath[pos:]
	_, err := os.Stat(litTexureName)
	if err == nil {
		//Textur vorhanden
		return litTexureName, true
	}
	//2. prüfen ob gleichnamige Datei mit LIT vorhanden ist
	litTexureName = texturePath[:pos] + "LIT" + texturePath[pos:]
	_, err = os.Stat(litTexureName)
	if err == nil {
		//Textur vorhanden
		return litTexureName, true
	}
	//3. wenn irgendeine Nachttextur verwendet werden soll
	if self.useAnyLitTexture {
		litTexureName = ""
		dir := filepath.Dir(texturePath)
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir()&&strings.HasSuffix(strings.ToLower(path), "_lit.png") {
				//es wurde eine Nachttextur gefunden
				litTexureName = path
				return filepath.SkipDir
			}
			return nil
		})
		if len(litTexureName) > 0 {
			//es wurde irgendeine Nachttextur gefunden --> diese nehme
			return litTexureName, true
		}
	}
	//es wurde keine Nachttextur gefunden
	return "", false
}

func convertImageToAlpha(img *image.NRGBA) {
	size := img.Rect.Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			i := y * img.Stride + x * 4

			r := img.Pix[i + 0]
			g := img.Pix[i + 1]
			b := img.Pix[i + 2]

			if r == 0xFF&&g == 0&&b == 0xFF {
				img.Pix[i + 0] = 0
				img.Pix[i + 1] = 0
				img.Pix[i + 2] = 0
				img.Pix[i + 3] = 0
			}
		}
	}
}

func mergeTextures(normalTexture, litTexture *image.NRGBA) {
	normSize := normalTexture.Bounds().Size()
	litSize := litTexture.Bounds().Size()
	if normSize.X != litSize.X || normSize.Y != litSize.Y {
		//Größen sind unterschiedlich --> normale Textur auf Größe der Nachttextur skalieren
		normalTexture = imaging.Resize(normalTexture, litSize.X, litSize.Y, imaging.Linear)
	}
	for x := 0; x < litSize.X; x++ {
		for y := 0; y < litSize.Y; y++ {
			i := y * normalTexture.Stride + x * 4
			j := y * litTexture.Stride + x * 4

			litTexture.Pix[j + 0] = mergePixel(normalTexture.Pix[i + 0], litTexture.Pix[j + 0], 1.0)
			litTexture.Pix[j + 1] = mergePixel(normalTexture.Pix[i + 1], litTexture.Pix[j + 1], 1.0)
			litTexture.Pix[j + 2] = mergePixel(normalTexture.Pix[i + 2], litTexture.Pix[j + 2], 0.7)
			litTexture.Pix[j + 3] = normalTexture.Pix[i + 3]
		}
	}
}

func mergePixel(normPixel, litPixel uint8, correction float32) uint8 {
	result := int(normPixel) - int(255 - litPixel)
	if result < 0 {
		result = 0
	}
	result = int(float32(result) * correction)
	return uint8(result)
}