package api

import "strings"

func replaceWords(original, replacement string, toReplace []string) string {
	output := []string{}
	words := strings.FieldsSeq(original)
	for word := range words {
		replaced := false
		for _, x := range toReplace {
			if strings.EqualFold(word, x) {
				output = append(output, replacement)
				replaced = true
				break
			}
		}

		if replaced {
			continue
		}
		output = append(output, word)
	}
	return strings.Join(output, " ")
}
