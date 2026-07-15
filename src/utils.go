package src

import "regexp"

// CleanMOTD очищает текст от параграфов (§) разметки Minecraft
func CleanMOTD(motd string) string {
	re := regexp.MustCompile(`§[0-9a-fk-orA-FK-OR]`)
	return re.ReplaceAllString(motd, "")
}
