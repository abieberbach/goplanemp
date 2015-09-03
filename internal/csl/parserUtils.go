//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package csl
import "strings"

func trimTokens(tokens []string) []string {
	for idx, token := range tokens {
		tokens[idx] = strings.TrimSpace(token)
	}
	return tokens
}

