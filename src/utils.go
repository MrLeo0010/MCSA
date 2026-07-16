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

// containsAny проверяет, содержится ли хотя бы одна фраза из списка в строке str
func ContainsAny(str string, indicators []string) bool {
	for _, indicator := range indicators {
		// Приводим индикатор к нижнему регистру на случай,
		// если в массиве src.RegisterIndicators они записаны в разном регистре
		if strings.Contains(str, strings.ToLower(indicator)) {
			return true
		}
	}
	return false
}
