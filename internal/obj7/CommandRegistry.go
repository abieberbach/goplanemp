//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package obj7

type cmdInfo struct {
	cmdId        CommandId
	cmdType      CommandType
	elementCount int
}

var registeredCommands = make(map[string]cmdInfo)

func init() {

	registeredCommands["end"] = cmdInfo{CommandId_ObjEnd, CommandType_None, 0 }
	registeredCommands["light"] = cmdInfo{CommandId_ObjLight, CommandType_PtLine, 1 }
	registeredCommands["line"] = cmdInfo{CommandId_ObjLine, CommandType_PtLine, 2 }
	registeredCommands["tri"] = cmdInfo{CommandId_ObjTri, CommandType_Poly, 3 }
	registeredCommands["quad"] = cmdInfo{CommandId_ObjQuad, CommandType_Poly, 4 }
	registeredCommands["quad_cockpit"] = cmdInfo{CommandId_ObjQuadCockpit, CommandType_Poly, 4 }
	registeredCommands["quad_hard"] = cmdInfo{CommandId_ObjQuadHard, CommandType_Poly, 4 }
	registeredCommands["smoke_black"] = cmdInfo{CommandId_ObjSmokeBlack, CommandType_Attr, 4 }
	registeredCommands["smoke_white"] = cmdInfo{CommandId_ObjSmokeWhite, CommandType_Attr, 4 }
	registeredCommands["quad_movie"] = cmdInfo{CommandId_ObjMovie, CommandType_Poly, 4 }
	registeredCommands["polygon"] = cmdInfo{CommandId_ObjPolygon, CommandType_Poly, 0 }
	registeredCommands["quad_strip"] = cmdInfo{CommandId_ObjQuadStrip, CommandType_Poly, 0 }
	registeredCommands["tri_strip"] = cmdInfo{CommandId_ObjTriStrip, CommandType_Poly, 0 }
	registeredCommands["tri_fan"] = cmdInfo{CommandId_ObjTriFan, CommandType_Poly, 0 }
	registeredCommands["ATTR_shade_flat"] = cmdInfo{CommandId_AttrShadeFlat, CommandType_Attr, 0 }
	registeredCommands["ATTR_shade_smooth"] = cmdInfo{CommandId_AttrShadeSmooth, CommandType_Attr, 0 }
	registeredCommands["shade_flat"] = cmdInfo{CommandId_AttrShadeFlat, CommandType_Attr, 0 }
	registeredCommands["shade_smooth"] = cmdInfo{CommandId_AttrShadeSmooth, CommandType_Attr, 0 }
	registeredCommands["ATTR_ambient_rgb"] = cmdInfo{CommandId_AttrAmbientRGB, CommandType_Attr, 3 }
	registeredCommands["ATTR_difuse_rgb"] = cmdInfo{CommandId_AttrDiffuseRGB, CommandType_Attr, 3 }
	registeredCommands["ATTR_emission_rgb"] = cmdInfo{CommandId_AttrEmissionRGB, CommandType_Attr, 3 }
	registeredCommands["ATTR_specular_rgb"] = cmdInfo{CommandId_AttrSpecularRGB, CommandType_Attr, 3 }
	registeredCommands["ATTR_shiny_rat"] = cmdInfo{CommandId_AttrShinyRat, CommandType_Attr, 1 }
	registeredCommands["ATTR_no_depth"] = cmdInfo{CommandId_AttrNoDepth, CommandType_Attr, 0 }
	registeredCommands["ATTR_depth"] = cmdInfo{CommandId_AttrDepth, CommandType_Attr, 0 }
	registeredCommands["ATTR_LOD"] = cmdInfo{CommandId_AttrLOD, CommandType_Attr, 2 }
	registeredCommands["ATTR_reset"] = cmdInfo{CommandId_AttrReset, CommandType_Attr, 0 }
	registeredCommands["ATTR_cull"] = cmdInfo{CommandId_AttrCull, CommandType_Attr, 0 }
	registeredCommands["ATTR_no_cull"] = cmdInfo{CommandId_AttrNoCull, CommandType_Attr, 0 }
	registeredCommands["ATTR_poly_os"] = cmdInfo{CommandId_AttrOffset, CommandType_Attr, 1 }
}

func FindObjectCommand(token string) (*ObjectCommand, int) {
	info, found := registeredCommands[token]
	if !found {
		return nil, 0
	}
	return NewObjectCommand(info.cmdId, info.cmdType), info.elementCount
}