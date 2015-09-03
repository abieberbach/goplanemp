//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package csl
import "sync"

type CslPackages struct {
	packages map[string]*CslPackage
	rwMutex  sync.RWMutex
}

func NewCslPackages() *CslPackages {
	return &CslPackages{make(map[string]*CslPackage), sync.RWMutex{}}
}

func (self *CslPackages) AddPackage(cslPackage *CslPackage) {
	self.rwMutex.Lock()
	self.packages[cslPackage.Name] = cslPackage
	self.rwMutex.Unlock()
}

func (self *CslPackages) GetPackage(name string) (cslPackage *CslPackage, found bool) {
	self.rwMutex.RLock()
	cslPackage, found = self.packages[name]
	self.rwMutex.RUnlock()
	return
}

func (self *CslPackages) GetAllPackages() (result []*CslPackage) {
	self.rwMutex.RLock()
	result = make([]*CslPackage, 0, len(self.packages))
	for _, currentPackage := range self.packages {
		result = append(result, currentPackage)
	}
	self.rwMutex.RUnlock()
	return
}


