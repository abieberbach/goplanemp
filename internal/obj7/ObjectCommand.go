//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package obj7

type CommandType int
type CommandId int

const (
	CommandType_None CommandType = iota
	CommandType_PtLine
	CommandType_Poly
	CommandType_Attr
)

const (
	CommandId_ObjEnd CommandId = iota
	CommandId_ObjLight
	CommandId_ObjLine
	CommandId_ObjTri
	CommandId_ObjQuad
	CommandId_ObjQuadHard
	CommandId_ObjQuadCockpit
	CommandId_ObjSmokeBlack
	CommandId_ObjSmokeWhite
	CommandId_ObjMovie
	CommandId_ObjPolygon
	CommandId_ObjQuadStrip
	CommandId_ObjTriStrip
	CommandId_ObjTriFan
	CommandId_AttrShadeFlat
	CommandId_AttrShadeSmooth
	CommandId_AttrAmbientRGB
	CommandId_AttrDiffuseRGB
	CommandId_AttrEmissionRGB
	CommandId_AttrSpecularRGB
	CommandId_AttrShinyRat
	CommandId_AttrNoDepth
	CommandId_AttrDepth
	CommandId_AttrLOD
	CommandId_AttrReset
	CommandId_AttrCull
	CommandId_AttrNoCull
	CommandId_AttrOffset
	CommandId_AttrMax
)

type VecTex struct {
	V  [3]float32
	St [2]float32
}

type VecRGB struct {
	V   [3]float32
	RGB [3]float32
}


type ObjectCommand struct {
	CmdId      CommandId
	CmdType    CommandType
	Attributes []float32
	St         []VecTex
	RGB        []VecRGB
}

func NewObjectCommand(cmdId CommandId, cmdType CommandType) *ObjectCommand {
	return &ObjectCommand{cmdId, cmdType, make([]float32, 0), make([]VecTex, 0), make([]VecRGB, 0)}
}