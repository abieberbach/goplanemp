//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package obj7
import (
	"unsafe"
	"math"
	"github.com/abieberbach/goplanemp/gl"
)

type ObjectPointPool struct {
	pool []float32
}

func NewObjectPointPool() *ObjectPointPool {
	return &ObjectPointPool{make([]float32, 0)}
}

func (self *ObjectPointPool) AddPoint(xyz [3]float32, st [2]float32) int32 {
	//prüfen ob der Punkt ggf. schon in der Liste ist
	for n := 0; n < len(self.pool); n += 8 {
		if (xyz[0] == self.pool[n]) &&
		(xyz[1] == self.pool[n + 1]) &&
		(xyz[2] == self.pool[n + 2]) &&
		(st[0] == self.pool[n + 3]) &&
		(st[1] == self.pool[n + 4]) {
			//ja --> Position des Punktes zurückgeben
			return int32(n / 8)
		}
	}
	// Punkt wurde nicht gefunden --> Daten in die Liste aufnehmen
	// XYZ hinzufügen
	self.pool = append(self.pool, xyz[0])
	self.pool = append(self.pool, xyz[1])
	self.pool = append(self.pool, xyz[2])
	// ST hinzufügen
	self.pool = append(self.pool, st[0])
	self.pool = append(self.pool, st[1])
	//zusätzliche Were für die Normalen reservieren
	self.pool = append(self.pool, 0.0)
	self.pool = append(self.pool, 0.0)
	self.pool = append(self.pool, 0.0)
	return int32((len(self.pool) / 8) - 1)

}

func (self *ObjectPointPool) PreparePoolToDraw() {
	// Setup our triangle data (20 represents 5 elements of 4 bytes each
	// namely s,t,xn,yn,zn)
	gl.VertexPointer(3, gl.FLOAT, 32, unsafe.Pointer(&self.pool[0]))
	// Set our texture data (24 represents 6 elements of 4 bytes each
	// namely xn, yn, zn, x, y, z. We start 3 from the beginning to skip
	// over x, y, z initially.
	gl.ClientActiveTextureARB(gl.TEXTURE1)
	gl.TexCoordPointer(2, gl.FLOAT, 32, unsafe.Pointer(&self.pool[3]))

	gl.ClientActiveTextureARB(gl.TEXTURE0)
	gl.TexCoordPointer(2, gl.FLOAT, 32, unsafe.Pointer(&self.pool[3]))
	// Set our normal data...
	gl.NormalPointer(gl.FLOAT, 32, unsafe.Pointer(&self.pool[5]))
}

func (self *ObjectPointPool) CalcTriNormal(idx1, idx2, idx3 int32) {
	if self.pool[idx1 * 8  ] == self.pool[idx2 * 8  ]&&
	self.pool[idx1 * 8 + 1] == self.pool[idx2 * 8 + 1]&&
	self.pool[idx1 * 8 + 2] == self.pool[idx2 * 8 + 2] {
		return
	}
	if (self.pool[idx1 * 8  ] == self.pool[idx3 * 8  ]&&
	self.pool[idx1 * 8 + 1] == self.pool[idx3 * 8 + 1]&&
	self.pool[idx1 * 8 + 2] == self.pool[idx3 * 8 + 2]) {
		return
	}
	if (self.pool[idx2 * 8  ] == self.pool[idx3 * 8  ]&&
	self.pool[idx2 * 8 + 1] == self.pool[idx3 * 8 + 1]&&
	self.pool[idx2 * 8 + 2] == self.pool[idx3 * 8 + 2]) {
		return
	}

	// idx2->idx1 cross idx1->idx3 = normal product
	var v1 [3]float32
	var v2 [3]float32
	v1[0] = self.pool[idx2 * 8  ] - self.pool[idx1 * 8  ];
	v1[1] = self.pool[idx2 * 8 + 1] - self.pool[idx1 * 8 + 1];
	v1[2] = self.pool[idx2 * 8 + 2] - self.pool[idx1 * 8 + 2];

	v2[0] = self.pool[idx2 * 8  ] - self.pool[idx3 * 8  ];
	v2[1] = self.pool[idx2 * 8 + 1] - self.pool[idx3 * 8 + 1];
	v2[2] = self.pool[idx2 * 8 + 2] - self.pool[idx3 * 8 + 2];

	// We do NOT normalize the cross product; we want larger triangles
	// to make bigger normals.  When we blend them, bigger sides will
	// contribute more to the normals.  We'll normalize the normals
	// after the blend is done.
	n := crossVec(v1, v2)
	self.pool[idx1 * 8 + 5] += n[0];
	self.pool[idx1 * 8 + 6] += n[1];
	self.pool[idx1 * 8 + 7] += n[2];

	self.pool[idx2 * 8 + 5] += n[0];
	self.pool[idx2 * 8 + 6] += n[1];
	self.pool[idx2 * 8 + 7] += n[2];

	self.pool[idx3 * 8 + 5] += n[0];
	self.pool[idx3 * 8 + 6] += n[1];
	self.pool[idx3 * 8 + 7] += n[2];
}

func (self *ObjectPointPool) NormalizeNormals() {
	for n := 0; n < len(self.pool); n += 8 {
		for m := 0; m < len(self.pool); m += 8 {
			if self.pool[n  ] == self.pool[m  ] &&
			self.pool[n + 1] == self.pool[m + 1] &&
			self.pool[n + 2] == self.pool[m + 2] &&
			m != n {
				self.pool[n + 5], self.pool[m + 5] = swappedAdd(self.pool[n + 5], self.pool[m + 5])
				self.pool[n + 6], self.pool[m + 6] = swappedAdd(self.pool[n + 6], self.pool[m + 6])
				self.pool[n + 7], self.pool[m + 7] = swappedAdd(self.pool[n + 7], self.pool[m + 7])
			}
		}
	}
	for n := 5; n < len(self.pool); n += 8 {
		self.pool[n], self.pool[n + 1], self.pool[n + 2] = normalizeVec(self.pool[n], self.pool[n + 1], self.pool[n + 2])
	}
}

func (self *ObjectPointPool) Purge() {
	self.pool = make([]float32, 0)
}

func (self *ObjectPointPool) Size() int {
	return len(self.pool)
}

func crossVec(a [3]float32, b [3]float32) (dst [3]float32) {
	dst[0] = a[1] * b[2] - a[2] * b[1]
	dst[1] = a[2] * b[0] - a[0] * b[2]
	dst[2] = a[0] * b[1] - a[1] * b[0]
	return
}

func swappedAdd(a float32, b float32) (float32, float32) {
	return a + b, b + a
}

func normalizeVec(x, y, z float32) (float32, float32, float32) {
	len := float32(math.Sqrt(float64(x * x + y * y + z * z)))
	if (len > 0.0) {
		len = 1.0 / len
		x *= len
		y *= len
		z *= len
	}
	return x, y, z
}
