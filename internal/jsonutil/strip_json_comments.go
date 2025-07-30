package jsonutil

import (
	"bufio"
	"bytes"
	"strings"
)

func StripJSONComments(data []byte) []byte {
	var result bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader(data))

	inString := false
	inMultiLineComment := false

	for scanner.Scan() {
		line := scanner.Text()
		var cleanLine strings.Builder

		i := 0
		for i < len(line) {
			char := line[i]

			// Handle multi-line comment end
			if inMultiLineComment {
				if i < len(line)-1 && char == '*' && line[i+1] == '/' {
					inMultiLineComment = false
					i += 2
					continue
				}
				i++
				continue
			}

			// Handle string literals (don't process comments inside strings)
			if char == '"' && (i == 0 || line[i-1] != '\\') {
				inString = !inString
				cleanLine.WriteByte(char)
				i++
				continue
			}

			// Skip comment processing if we're inside a string
			if inString {
				cleanLine.WriteByte(char)
				i++
				continue
			}

			// Handle single-line comments
			if i < len(line)-1 && char == '/' && line[i+1] == '/' {
				break // Skip rest of line
			}

			// Handle multi-line comment start
			if i < len(line)-1 && char == '/' && line[i+1] == '*' {
				inMultiLineComment = true
				i += 2
				continue
			}

			cleanLine.WriteByte(char)
			i++
		}

		// Add the cleaned line if it's not empty or just whitespace
		cleanedLine := strings.TrimSpace(cleanLine.String())
		if cleanedLine != "" {
			result.WriteString(cleanedLine)
			result.WriteByte('\n')
		}
	}

	return result.Bytes()
}
