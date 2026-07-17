package src

import (
	"regexp"
	"strings"
)

// CleanMOTD очищает текст от параграфов (§) разметки Minecraft
func CleanMOTD(motd string) string {
	re := regexp.MustCompile(`§[0-9a-fk-orA-FK-OR]`)
	return re.ReplaceAllString(motd, "")
}

func ContainsAny(str string, indicators []string) (bool, int, int) {
	// Переводим исходную строку в нижний регистр и сразу в руны
	strRunes := []rune(strings.ToLower(str))

	for _, indicator := range indicators {
		indicatorRunes := []rune(strings.ToLower(indicator))

		// Ищем совпадение слайса рун indicatorRunes внутри strRunes
		start := findRuneSubslice(strRunes, indicatorRunes)
		if start != -1 {
			end := start + len(indicatorRunes)
			return true, start, end
		}
	}
	return false, -1, -1
}

// Вспомогательная функция для поиска одного слайса рун в другом
func findRuneSubslice(haystack, needle []rune) int {
	if len(needle) == 0 {
		return -1
	}
	for i := 0; i <= len(haystack)-len(needle); i++ {
		match := true
		for j := 0; j < len(needle); j++ {
			if haystack[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
