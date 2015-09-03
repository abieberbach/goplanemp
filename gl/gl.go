//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

//OpenGL-Bindings für das Zeichnen der Objekte
package gl

// #cgo darwin  LDFLAGS: -framework OpenGL
// #cgo linux   LDFLAGS: -lGL
// #cgo windows LDFLAGS: -lopengl32
// #ifndef GLAPI
// #define GLAPI extern
// #endif
// #ifndef APIENTRY
// #define APIENTRY
// #endif
// #ifndef APIENTRYP
// #define APIENTRYP APIENTRY *
// #endif
//
// #include <stdlib.h>
// #if defined(__APPLE__)
// 		#include <dlfcn.h>
// #elif defined(_WIN32)
// 		#define WIN32_LEAN_AND_MEAN 1
// 		#include <windows.h>
// #else
// 		#include <GL/glx.h>
// #endif
//
// typedef unsigned int GLbitfield;
// typedef unsigned char GLubyte;
// typedef double GLdouble;
// typedef unsigned int GLenum;
// typedef float GLfloat;
// typedef int GLint;
// typedef int GLsizei;
// typedef unsigned int GLuint;
//
// typedef void  (APIENTRYP GPACTIVETEXTUREARB)(GLenum  texture);
// typedef void  (APIENTRYP GPBEGIN)(GLenum  mode);
// typedef void  (APIENTRYP GPBINDBUFFERARB)(GLenum  target, GLuint  buffer);
// typedef void  (APIENTRYP GPCALLLIST)(GLuint  list);
// typedef void  (APIENTRYP GPCLIENTACTIVETEXTUREARB)(GLenum  texture);
// typedef void  (APIENTRYP GPCOLOR4F)(GLfloat  red, GLfloat  green, GLfloat  blue, GLfloat  alpha);
// typedef void  (APIENTRYP GPDISABLECLIENTSTATE)(GLenum  array);
// typedef void  (APIENTRYP GPDRAWELEMENTS)(GLenum  mode, GLsizei  count, GLenum  type, const void * indices);
// typedef void  (APIENTRYP GPENABLECLIENTSTATE)(GLenum  array);
// typedef void  (APIENTRYP GPEND)();
// typedef void  (APIENTRYP GPENDLIST)();
// typedef GLuint  (APIENTRYP GPGENLISTS)(GLsizei  range);
// typedef GLenum  (APIENTRYP GPGETERROR)();
// typedef void  (APIENTRYP GPGETFLOATV)(GLenum  pname, GLfloat * data);
// typedef void  (APIENTRYP GPGETINTEGERV)(GLenum  pname, GLint * data);
// typedef const GLubyte * (APIENTRYP GPGETSTRING)(GLenum  name);
// typedef void  (APIENTRYP GPLOADIDENTITY)();
// typedef void  (APIENTRYP GPMATRIXMODE)(GLenum  mode);
// typedef void  (APIENTRYP GPNEWLIST)(GLuint  list, GLenum  mode);
// typedef void  (APIENTRYP GPNORMALPOINTER)(GLenum  type, GLsizei  stride, const void * pointer);
// typedef void  (APIENTRYP GPORTHO)(GLdouble  left, GLdouble  right, GLdouble  bottom, GLdouble  top, GLdouble  zNear, GLdouble  zFar);
// typedef void  (APIENTRYP GPPOPCLIENTATTRIB)();
// typedef void  (APIENTRYP GPPOPMATRIX)();
// typedef void  (APIENTRYP GPPUSHCLIENTATTRIB)(GLbitfield  mask);
// typedef void  (APIENTRYP GPPUSHMATRIX)();
// typedef void  (APIENTRYP GPROTATEF)(GLfloat  angle, GLfloat  x, GLfloat  y, GLfloat  z);
// typedef void  (APIENTRYP GPTEXCOORD2F)(GLfloat  s, GLfloat  t);
// typedef void  (APIENTRYP GPTEXCOORDPOINTER)(GLint  size, GLenum  type, GLsizei  stride, const void * pointer);
// typedef void  (APIENTRYP GPTEXENVI)(GLenum  target, GLenum  pname, GLint  param);
// typedef void  (APIENTRYP GPTEXPARAMETERI)(GLenum  target, GLenum  pname, GLint  param);
// typedef void  (APIENTRYP GPTRANSLATEF)(GLfloat  x, GLfloat  y, GLfloat  z);
// typedef void  (APIENTRYP GPVERTEX2F)(GLfloat  x, GLfloat  y);
// typedef void  (APIENTRYP GPVERTEXPOINTER)(GLint  size, GLenum  type, GLsizei  stride, const void * pointer);
//
// static void  glowActiveTextureARB(GPACTIVETEXTUREARB fnptr, GLenum  texture) {
//   (*fnptr)(texture);
// }
// static void  glowBegin(GPBEGIN fnptr, GLenum  mode) {
//   (*fnptr)(mode);
// }
// static void  glowBindBufferARB(GPBINDBUFFERARB fnptr, GLenum  target, GLuint  buffer) {
//   (*fnptr)(target, buffer);
// }
// static void  glowCallList(GPCALLLIST fnptr, GLuint  list) {
//   (*fnptr)(list);
// }
// static void  glowClientActiveTextureARB(GPCLIENTACTIVETEXTUREARB fnptr, GLenum  texture) {
//   (*fnptr)(texture);
// }
// static void  glowColor4f(GPCOLOR4F fnptr, GLfloat  red, GLfloat  green, GLfloat  blue, GLfloat  alpha) {
//   (*fnptr)(red, green, blue, alpha);
// }
// static void  glowDisableClientState(GPDISABLECLIENTSTATE fnptr, GLenum  array) {
//   (*fnptr)(array);
// }
// static void  glowDrawElements(GPDRAWELEMENTS fnptr, GLenum  mode, GLsizei  count, GLenum  type, const void * indices) {
//   (*fnptr)(mode, count, type, indices);
// }
// static void  glowEnableClientState(GPENABLECLIENTSTATE fnptr, GLenum  array) {
//   (*fnptr)(array);
// }
// static void  glowEnd(GPEND fnptr) {
//   (*fnptr)();
// }
// static void  glowEndList(GPENDLIST fnptr) {
//   (*fnptr)();
// }
// static GLuint  glowGenLists(GPGENLISTS fnptr, GLsizei  range) {
//   return (*fnptr)(range);
// }
// static GLenum  glowGetError(GPGETERROR fnptr) {
//   return (*fnptr)();
// }
// static void  glowGetFloatv(GPGETFLOATV fnptr, GLenum  pname, GLfloat * data) {
//   (*fnptr)(pname, data);
// }
// static void  glowGetIntegerv(GPGETINTEGERV fnptr, GLenum  pname, GLint * data) {
//   (*fnptr)(pname, data);
// }
// static const GLubyte * glowGetString(GPGETSTRING fnptr, GLenum  name) {
//   return (*fnptr)(name);
// }
// static void  glowLoadIdentity(GPLOADIDENTITY fnptr) {
//   (*fnptr)();
// }
// static void  glowMatrixMode(GPMATRIXMODE fnptr, GLenum  mode) {
//   (*fnptr)(mode);
// }
// static void  glowNewList(GPNEWLIST fnptr, GLuint  list, GLenum  mode) {
//   (*fnptr)(list, mode);
// }
// static void  glowNormalPointer(GPNORMALPOINTER fnptr, GLenum  type, GLsizei  stride, const void * pointer) {
//   (*fnptr)(type, stride, pointer);
// }
// static void  glowOrtho(GPORTHO fnptr, GLdouble  left, GLdouble  right, GLdouble  bottom, GLdouble  top, GLdouble  zNear, GLdouble  zFar) {
//   (*fnptr)(left, right, bottom, top, zNear, zFar);
// }
// static void  glowPopClientAttrib(GPPOPCLIENTATTRIB fnptr) {
//   (*fnptr)();
// }
// static void  glowPopMatrix(GPPOPMATRIX fnptr) {
//   (*fnptr)();
// }
// static void  glowPushClientAttrib(GPPUSHCLIENTATTRIB fnptr, GLbitfield  mask) {
//   (*fnptr)(mask);
// }
// static void  glowPushMatrix(GPPUSHMATRIX fnptr) {
//   (*fnptr)();
// }
// static void  glowRotatef(GPROTATEF fnptr, GLfloat  angle, GLfloat  x, GLfloat  y, GLfloat  z) {
//   (*fnptr)(angle, x, y, z);
// }
// static void  glowTexCoord2f(GPTEXCOORD2F fnptr, GLfloat  s, GLfloat  t) {
//   (*fnptr)(s, t);
// }
// static void  glowTexCoordPointer(GPTEXCOORDPOINTER fnptr, GLint  size, GLenum  type, GLsizei  stride, const void * pointer) {
//   (*fnptr)(size, type, stride, pointer);
// }
// static void  glowTexEnvi(GPTEXENVI fnptr, GLenum  target, GLenum  pname, GLint  param) {
//   (*fnptr)(target, pname, param);
// }
// static void  glowTexParameteri(GPTEXPARAMETERI fnptr, GLenum  target, GLenum  pname, GLint  param) {
//   (*fnptr)(target, pname, param);
// }
// static void  glowTranslatef(GPTRANSLATEF fnptr, GLfloat  x, GLfloat  y, GLfloat  z) {
//   (*fnptr)(x, y, z);
// }
// static void  glowVertex2f(GPVERTEX2F fnptr, GLfloat  x, GLfloat  y) {
//   (*fnptr)(x, y);
// }
// static void  glowVertexPointer(GPVERTEXPOINTER fnptr, GLint  size, GLenum  type, GLsizei  stride, const void * pointer) {
//   (*fnptr)(size, type, stride, pointer);
// }
// static void* coreGetProcAddress(const char* name) {
// 		#ifdef __APPLE__
// 			return dlsym(RTLD_DEFAULT, name);
// 		#elif _WIN32
// 			return wglGetProcAddress((LPCSTR)name);
// 		#else
// 			return glXGetProcAddress(name);
// 		#endif
// }
import "C"
import (
	"unsafe"
	"errors"
)

var (
	gpActiveTextureARB                                       C.GPACTIVETEXTUREARB
	gpBegin                                                  C.GPBEGIN
	gpBindBufferARB                                          C.GPBINDBUFFERARB
	gpCallList                                               C.GPCALLLIST
	gpClientActiveTextureARB                                 C.GPCLIENTACTIVETEXTUREARB
	gpColor4f                                                C.GPCOLOR4F
	gpDisableClientState                                     C.GPDISABLECLIENTSTATE
	gpDrawElements                                           C.GPDRAWELEMENTS
	gpEnableClientState                                      C.GPENABLECLIENTSTATE
	gpEnd                                                    C.GPEND
	gpEndList                                                C.GPENDLIST
	gpGenLists                                               C.GPGENLISTS
	gpGetError                                               C.GPGETERROR
	gpGetFloatv                                              C.GPGETFLOATV
	gpGetIntegerv                                            C.GPGETINTEGERV
	gpGetString                                              C.GPGETSTRING
	gpLoadIdentity                                           C.GPLOADIDENTITY
	gpMatrixMode                                             C.GPMATRIXMODE
	gpNewList                                                C.GPNEWLIST
	gpNormalPointer                                          C.GPNORMALPOINTER
	gpOrtho                                                  C.GPORTHO
	gpPopClientAttrib                                        C.GPPOPCLIENTATTRIB
	gpPopMatrix                                              C.GPPOPMATRIX
	gpPushClientAttrib                                       C.GPPUSHCLIENTATTRIB
	gpPushMatrix                                             C.GPPUSHMATRIX
	gpRotatef                                                C.GPROTATEF
	gpTexCoord2f                                             C.GPTEXCOORD2F
	gpTexCoordPointer                                        C.GPTEXCOORDPOINTER
	gpTexEnvi                                                C.GPTEXENVI
	gpTexParameteri                                          C.GPTEXPARAMETERI
	gpTranslatef                                             C.GPTRANSLATEF
	gpVertexPointer                                          C.GPVERTEXPOINTER
	gpVertex2f                                               C.GPVERTEX2F
)

func VertexPointer(size int32, xtype uint32, stride int32, pointer unsafe.Pointer) {
	C.glowVertexPointer(gpVertexPointer, (C.GLint)(size), (C.GLenum)(xtype), (C.GLsizei)(stride), pointer)
}

func ClientActiveTextureARB(texture uint32) {
	C.glowClientActiveTextureARB(gpClientActiveTextureARB, (C.GLenum)(texture))
}

func TexCoordPointer(size int32, xtype uint32, stride int32, pointer unsafe.Pointer) {
	C.glowTexCoordPointer(gpTexCoordPointer, (C.GLint)(size), (C.GLenum)(xtype), (C.GLsizei)(stride), pointer)
}

func NormalPointer(xtype uint32, stride int32, pointer unsafe.Pointer) {
	C.glowNormalPointer(gpNormalPointer, (C.GLenum)(xtype), (C.GLsizei)(stride), pointer)
}

func GetIntegerv(pname uint32, data *int32) {
	C.glowGetIntegerv(gpGetIntegerv, (C.GLenum)(pname), (*C.GLint)(unsafe.Pointer(data)))
}
func GetString(name uint32) *uint8 {
	ret := C.glowGetString(gpGetString, (C.GLenum)(name))
	return (*uint8)(ret)
}

func GetError() uint32 {
	ret := C.glowGetError(gpGetError)
	return (uint32)(ret)
}

func TexParameteri(target uint32, pname uint32, param int32) {
	C.glowTexParameteri(gpTexParameteri, (C.GLenum)(target), (C.GLenum)(pname), (C.GLint)(param))
}

func TexEnvi(target uint32, pname uint32, param int32) {
	C.glowTexEnvi(gpTexEnvi, (C.GLenum)(target), (C.GLenum)(pname), (C.GLint)(param))
}

func GetFloatv(pname uint32, data *float32) {
	C.glowGetFloatv(gpGetFloatv, (C.GLenum)(pname), (*C.GLfloat)(unsafe.Pointer(data)))
}

func MatrixMode(mode uint32) {
	C.glowMatrixMode(gpMatrixMode, (C.GLenum)(mode))
}

func PushMatrix() {
	C.glowPushMatrix(gpPushMatrix)
}

func PopMatrix() {
	C.glowPopMatrix(gpPopMatrix)
}

func LoadIdentity() {
	C.glowLoadIdentity(gpLoadIdentity)
}

func Ortho(left float64, right float64, bottom float64, top float64, zNear float64, zFar float64) {
	C.glowOrtho(gpOrtho, (C.GLdouble)(left), (C.GLdouble)(right), (C.GLdouble)(bottom), (C.GLdouble)(top), (C.GLdouble)(zNear), (C.GLdouble)(zFar))
}

func Translatef(x float32, y float32, z float32) {
	C.glowTranslatef(gpTranslatef, (C.GLfloat)(x), (C.GLfloat)(y), (C.GLfloat)(z))
}

func Rotatef(angle float32, x float32, y float32, z float32) {
	C.glowRotatef(gpRotatef, (C.GLfloat)(angle), (C.GLfloat)(x), (C.GLfloat)(y), (C.GLfloat)(z))
}

func GenLists(xrange int32) uint32 {
	ret := C.glowGenLists(gpGenLists, (C.GLsizei)(xrange))
	return (uint32)(ret)
}

func NewList(list uint32, mode uint32) {
	C.glowNewList(gpNewList, (C.GLuint)(list), (C.GLenum)(mode))
}

func EndList() {
	C.glowEndList(gpEndList)
}

func Begin(mode uint32) {
	C.glowBegin(gpBegin, (C.GLenum)(mode))
}

func End() {
	C.glowEnd(gpEnd)
}

func CallList(list uint32) {
	C.glowCallList(gpCallList, (C.GLuint)(list))
}

func Color4f(red float32, green float32, blue float32, alpha float32) {
	C.glowColor4f(gpColor4f, (C.GLfloat)(red), (C.GLfloat)(green), (C.GLfloat)(blue), (C.GLfloat)(alpha))
}

func TexCoord2f(s float32, t float32) {
	C.glowTexCoord2f(gpTexCoord2f, (C.GLfloat)(s), (C.GLfloat)(t))
}

func Vertex2f(x float32, y float32) {
	C.glowVertex2f(gpVertex2f, (C.GLfloat)(x), (C.GLfloat)(y))
}

func ActiveTextureARB(texture uint32) {
	C.glowActiveTextureARB(gpActiveTextureARB, (C.GLenum)(texture))
}

func BindBufferARB(target uint32, buffer uint32) {
	C.glowBindBufferARB(gpBindBufferARB, (C.GLenum)(target), (C.GLuint)(buffer))
}

func PushClientAttrib(mask uint32) {
	C.glowPushClientAttrib(gpPushClientAttrib, (C.GLbitfield)(mask))
}

func PopClientAttrib() {
	C.glowPopClientAttrib(gpPopClientAttrib)
}

func EnableClientState(array uint32) {
	C.glowEnableClientState(gpEnableClientState, (C.GLenum)(array))
}

func DisableClientState(array uint32) {
	C.glowDisableClientState(gpDisableClientState, (C.GLenum)(array))
}

func DrawElements(mode uint32, count int32, xtype uint32, indices unsafe.Pointer) {
	C.glowDrawElements(gpDrawElements, (C.GLenum)(mode), (C.GLsizei)(count), (C.GLenum)(xtype), indices)
}

func GoStr(cstr *uint8) string {
	return C.GoString((*C.char)(unsafe.Pointer(cstr)))
}

func getProcAddress(namea string) unsafe.Pointer {
	cname := C.CString(namea)
	defer C.free(unsafe.Pointer(cname))
	return C.coreGetProcAddress(cname)
}

func InitOpenGL() error {
	gpActiveTextureARB = (C.GPACTIVETEXTUREARB)(getProcAddress("glActiveTextureARB"))
	if gpActiveTextureARB == nil {
		return errors.New("glActiveTextureARB")
	}
	gpBegin = (C.GPBEGIN)(getProcAddress("glBegin"))
	if gpBegin == nil {
		return errors.New("glBegin")
	}

	gpBindBufferARB = (C.GPBINDBUFFERARB)(getProcAddress("glBindBufferARB"))
	if gpBindBufferARB == nil {
		return errors.New("glBindBufferARB")
	}
	gpCallList = (C.GPCALLLIST)(getProcAddress("glCallList"))
	if gpCallList == nil {
		return errors.New("glCallList")
	}
	gpClientActiveTextureARB = (C.GPCLIENTACTIVETEXTUREARB)(getProcAddress("glClientActiveTextureARB"))
	if gpClientActiveTextureARB == nil {
		return errors.New("glClientActiveTextureARB")
	}
	gpColor4f = (C.GPCOLOR4F)(getProcAddress("glColor4f"))
	if gpColor4f == nil {
		return errors.New("glColor4f")
	}
	gpDisableClientState = (C.GPDISABLECLIENTSTATE)(getProcAddress("glDisableClientState"))
	if gpDisableClientState == nil {
		return errors.New("glDisableClientState")
	}
	gpDrawElements = (C.GPDRAWELEMENTS)(getProcAddress("glDrawElements"))
	if gpDrawElements == nil {
		return errors.New("glDrawElements")
	}
	gpEnableClientState = (C.GPENABLECLIENTSTATE)(getProcAddress("glEnableClientState"))
	if gpEnableClientState == nil {
		return errors.New("glEnableClientState")
	}
	gpEnd = (C.GPEND)(getProcAddress("glEnd"))
	if gpEnd == nil {
		return errors.New("glEnd")
	}
	gpEndList = (C.GPENDLIST)(getProcAddress("glEndList"))
	if gpEndList == nil {
		return errors.New("glEndList")
	}
	gpGenLists = (C.GPGENLISTS)(getProcAddress("glGenLists"))
	if gpGenLists == nil {
		return errors.New("glGenLists")
	}
	gpGetError = (C.GPGETERROR)(getProcAddress("glGetError"))
	if gpGetError == nil {
		return errors.New("glGetError")
	}
	gpGetFloatv = (C.GPGETFLOATV)(getProcAddress("glGetFloatv"))
	if gpGetFloatv == nil {
		return errors.New("glGetFloatv")
	}
	gpGetIntegerv = (C.GPGETINTEGERV)(getProcAddress("glGetIntegerv"))
	if gpGetIntegerv == nil {
		return errors.New("glGetIntegerv")
	}
	gpGetString = (C.GPGETSTRING)(getProcAddress("glGetString"))
	if gpGetString == nil {
		return errors.New("glGetString")
	}
	gpLoadIdentity = (C.GPLOADIDENTITY)(getProcAddress("glLoadIdentity"))
	if gpLoadIdentity == nil {
		return errors.New("glLoadIdentity")
	}
	gpMatrixMode = (C.GPMATRIXMODE)(getProcAddress("glMatrixMode"))
	if gpMatrixMode == nil {
		return errors.New("glMatrixMode")
	}
	gpNewList = (C.GPNEWLIST)(getProcAddress("glNewList"))
	if gpNewList == nil {
		return errors.New("glNewList")
	}
	gpNormalPointer = (C.GPNORMALPOINTER)(getProcAddress("glNormalPointer"))
	if gpNormalPointer == nil {
		return errors.New("glNormalPointer")
	}
	gpOrtho = (C.GPORTHO)(getProcAddress("glOrtho"))
	if gpOrtho == nil {
		return errors.New("glOrtho")
	}
	gpPopClientAttrib = (C.GPPOPCLIENTATTRIB)(getProcAddress("glPopClientAttrib"))
	if gpPopClientAttrib == nil {
		return errors.New("glPopClientAttrib")
	}
	gpPopMatrix = (C.GPPOPMATRIX)(getProcAddress("glPopMatrix"))
	if gpPopMatrix == nil {
		return errors.New("glPopMatrix")
	}
	gpPushClientAttrib = (C.GPPUSHCLIENTATTRIB)(getProcAddress("glPushClientAttrib"))
	if gpPushClientAttrib == nil {
		return errors.New("glPushClientAttrib")
	}
	gpPushMatrix = (C.GPPUSHMATRIX)(getProcAddress("glPushMatrix"))
	if gpPushMatrix == nil {
		return errors.New("glPushMatrix")
	}
	gpRotatef = (C.GPROTATEF)(getProcAddress("glRotatef"))
	if gpRotatef == nil {
		return errors.New("glRotatef")
	}
	gpTexCoord2f = (C.GPTEXCOORD2F)(getProcAddress("glTexCoord2f"))
	if gpTexCoord2f == nil {
		return errors.New("glTexCoord2f")
	}
	gpTexCoordPointer = (C.GPTEXCOORDPOINTER)(getProcAddress("glTexCoordPointer"))
	if gpTexCoordPointer == nil {
		return errors.New("glTexCoordPointer")
	}
	gpTexEnvi = (C.GPTEXENVI)(getProcAddress("glTexEnvi"))
	if gpTexEnvi == nil {
		return errors.New("glTexEnvi")
	}
	gpTexParameteri = (C.GPTEXPARAMETERI)(getProcAddress("glTexParameteri"))
	if gpTexParameteri == nil {
		return errors.New("glTexParameteri")
	}
	gpTranslatef = (C.GPTRANSLATEF)(getProcAddress("glTranslatef"))
	if gpTranslatef == nil {
		return errors.New("glTranslatef")
	}
	gpVertexPointer = (C.GPVERTEXPOINTER)(getProcAddress("glVertexPointer"))
	if gpVertexPointer == nil {
		return errors.New("glVertexPointer")
	}
	gpVertex2f = (C.GPVERTEX2F)(getProcAddress("glVertex2f"))
	if gpVertex2f == nil {
		return errors.New("glVertex2f")
	}
	return nil
}

