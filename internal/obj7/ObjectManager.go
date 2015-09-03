//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package obj7
import (
	"sync"
)

var ObjectManagerInstance = &ObjectManager{make(map[string]*ObjectInfo), sync.RWMutex{}}

type ObjectManager struct {
	objList    map[string]*ObjectInfo
	objRwMutex sync.RWMutex
}

func (self *ObjectManager) GetObject(path string) (objInfo *ObjectInfo, err error) {
	self.objRwMutex.RLock()
	objInfo, found := self.objList[path]
	self.objRwMutex.RUnlock()
	if !found {
		objInfo, err = NewObjectInfo(path)
		if err != nil {
			return nil, err
		}
		self.objRwMutex.Lock()
		self.objList[path] = objInfo
		self.objRwMutex.Unlock()
	}
	return objInfo, nil
}
