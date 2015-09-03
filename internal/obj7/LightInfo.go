package obj7

type LightType int

const (
	LightType_Other LightType = 0
	LightType_RedNavigation LightType = 1
	LightType_GreenNavigation LightType = 2
	LightType_Beacon LightType = 3
	LightType_Strobe LightType = 4
	LightType_Landing LightType = 5
	LightType_Taxi LightType = 6
)

type LightInfo struct {
	XYZ       [3]float32
	RGB       [3]int
	LightType LightType
	RGB_Float [4]float32
}

var (
	NavLightRedColor = [4]float32{1.0, 0.0, 0.2, 0.5}
	NavLightGreenColor = [4]float32{0.0, 1.0, 0.3, 0.5}
	LandingLightColor = [4]float32{1.0, 1.0, 0.7, 0.6}
	StrobeLightColor = [4]float32{1.0, 1.0, 1.0, 0.7}
)

func NewLightInfo(vrgb VecRGB) *LightInfo {
	result := &LightInfo{}
	result.XYZ[0] = vrgb.V[0]
	result.XYZ[1] = vrgb.V[1]
	result.XYZ[2] = vrgb.V[2]
	result.RGB[0] = int(vrgb.RGB[0])
	result.RGB[1] = int(vrgb.RGB[1])
	result.RGB[2] = int(vrgb.RGB[2])
	result.LightType = result.buildType()
	result.RGB_Float = result.buildFloatColor()
	return result
}

func (self *LightInfo) buildType() LightType {
	if self.RGB[0] == 11 && self.RGB[1] == 11 && self.RGB[2] == 11 {
		return LightType_RedNavigation
	}
	if self.RGB[0] == 22 && self.RGB[1] == 22 && self.RGB[2] == 22 {
		return LightType_GreenNavigation
	}
	if self.RGB[0] == 33 && self.RGB[1] == 33 && self.RGB[2] == 33 {
		return LightType_Beacon
	}
	if self.RGB[0] == 44 && self.RGB[1] == 44 && self.RGB[2] == 44 {
		return LightType_Strobe
	}
	if self.RGB[0] == 55 && self.RGB[1] == 55 && self.RGB[2] == 55 {
		return LightType_Landing
	}
	if self.RGB[0] == 66 && self.RGB[1] == 66 && self.RGB[2] == 66 {
		return LightType_Taxi
	}
	return LightType_Other
}


func (self *LightInfo) buildFloatColor() [4]float32 {
	switch self.LightType {
	case LightType_RedNavigation, LightType_Beacon:
		return NavLightRedColor
	case LightType_GreenNavigation:
		return NavLightGreenColor
	case LightType_Strobe:
		return StrobeLightColor
	case LightType_Landing:
		return LandingLightColor
	default:
		return [4]float32{float32(self.RGB[0]) * 0.1, float32(self.RGB[1]) * 0.1, float32(self.RGB[2]) * 0.1, 0.0}
	}
}