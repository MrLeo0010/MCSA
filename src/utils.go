package src

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// CleanMOTD очищает текст от параграфов (§) разметки Minecraft
func CleanMOTD(motd string) string {
	re := regexp.MustCompile(`§[0-9a-fk-orA-FK-OR]`)
	return re.ReplaceAllString(motd, "")
}

// ParseMOTD безопасно извлекает текстовый MOTD
func (s *StatusResponse) ParseMOTD() string {
	if s.Description == nil {
		return ""
	}
	switch v := s.Description.(type) {
	case string:
		return CleanMOTD(v)
	case map[string]any:
		if text, ok := v["text"].(string); ok {
			return CleanMOTD(text)
		}
	}
	return ""
}

func WhereContainsAny(str string, indicators []string) (bool, int, int) {
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

func ContainsAny(s string, words []string) bool {
	s = strings.ToLower(s)

	for _, word := range words {
		if strings.Contains(s, strings.ToLower(word)) {
			return true
		}
	}
	return false
}

// Вспомогательная функция для поиска одного слайса рун в другом
func findRuneSubslice(haystack, needle []rune) int {
	if len(needle) == 0 {
		return -1
	}
	for i := 0; i <= len(haystack)-len(needle); i++ {
		match := true
		for j := range needle {
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

func BuildGitHubURL(owner, repo string) string {
	url_str := strings.ReplaceAll(BaseGitHubURL, "{OWNER}", owner)
	url_str = strings.ReplaceAll(url_str, "{REPO}", repo)
	return url_str
}

func BuildGitHubAPIURL(owner, repo string) string {
	url_str := strings.ReplaceAll(BaseGitHubAPIURL, "{OWNER}", owner)
	return strings.ReplaceAll(url_str, "{REPO}", repo)
}

func (p *ProgressReader) Read(buf []byte) (int, error) {
	n, err := p.Reader.Read(buf)
	p.Downloaded += int64(n)

	if p.Total > 0 {
		percent := float64(p.Downloaded) / float64(p.Total) * 100

		barWidth := 30
		filled := int(percent / 100 * float64(barWidth))

		bar := strings.Repeat("█", filled) +
			strings.Repeat("░", barWidth-filled)

		fmt.Printf("\r%s[%s] %.1f%%", p.Prefix, bar, percent)
	}

	return n, err
}

func DownloadFile(url, filename string, silent bool, prefix string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	var reader io.Reader = resp.Body

	if !silent {
		reader = &ProgressReader{
			Reader: resp.Body,
			Total:  resp.ContentLength,
			Prefix: prefix,
		}
	}

	_, err = io.Copy(out, reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !silent {
		fmt.Println()
	}
}
