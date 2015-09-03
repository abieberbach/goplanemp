//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

//OpenGL GLU-Bindings für das Zeichnen der Objekte
package glu

// #cgo darwin LDFLAGS: -framework Carbon -framework OpenGL -framework GLUT
// #cgo linux LDFLAGS: -lGLU
// #cgo windows LDFLAGS: -lglu32
//
// #ifdef __APPLE__
//   #include <OpenGL/glu.h>
// #else
//   #include <GL/glu.h>
// #endif
import "C"
import "unsafe"

func Build2DMipmaps(target uint32, internalFormat int, width, height int, format, typ uint32, data unsafe.Pointer) int {
	return int(C.gluBuild2DMipmaps(
		C.GLenum(target),
		C.GLint(internalFormat),
		C.GLsizei(width),
		C.GLsizei(height),
		C.GLenum(format),
		C.GLenum(typ),
		data))
}
