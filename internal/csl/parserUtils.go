package csl
import "strings"

func trimTokens(tokens []string) []string {
	for idx, token := range tokens {
		tokens[idx] = strings.TrimSpace(token)
	}
	return tokens
}

