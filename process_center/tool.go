package ProcessCenter

import "strings"

func escapeCharacter(s string) string {
	newMap := map[string]string{
		"&amp;": "&",
		"&#91;": "[",
		"&#93;": "]",
		"&#44;": ",",
	}
	for key, value := range newMap {
		s = strings.ReplaceAll(s, key, value)
	}
	return s
}
