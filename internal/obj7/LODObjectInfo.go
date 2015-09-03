//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package obj7

type LodObjectInfo struct {
	NearDist     float32 // The visible range
	FarDist      float32 // of this LOD
	TriangleList []int32
	Lights       []*LightInfo
	PointPool    *ObjectPointPool
	Dl           uint32
}

func NewLodObjectInfo(nearDist, farDist float32) *LodObjectInfo {
	result := &LodObjectInfo{}
	result.NearDist = nearDist
	result.FarDist = farDist
	result.TriangleList = make([]int32, 0)
	result.Lights = make([]*LightInfo, 0)
	result.PointPool = NewObjectPointPool()
	result.Dl = 0
	return result
}