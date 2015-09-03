//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package csl
import (
	"bufio"
	"os"
	"github.com/abieberbach/goplane/extra/logging"
	"strings"
)

func parseRelatedFile(relatedFilePath string) map[string][]string {
	relatedMap := make(map[string][]string)
	relatedFile, err := os.Open(relatedFilePath)
	if err != nil {
		logging.Errorf("could not open related file \"%v\": %v", relatedFilePath, err)
		return relatedMap
	}
	scanner := bufio.NewScanner(relatedFile)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if len(text) == 0 || strings.HasPrefix(text, ";") {
			//Zeile enthält keinen Text bzw. nur Kommentare --> ignorieren
			continue
		}
		tokens := strings.Split(text, " ")
		tokens = trimTokens(tokens)
		//für jeden ICAO-Code die Liste der verwandten ICAOs merken
		for _, token := range tokens {
			relatedMap[token] = tokens
		}
	}
	return relatedMap
}