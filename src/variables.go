package src

import "regexp"

// Базовая структура ответа Server List Ping
type StatusResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int            `json:"max"`
		Online int            `json:"online"`
		Sample []PlayerSample `json:"sample"`
	} `json:"players"`
	Description interface{} `json:"description"` // Может быть строкой или объектом Chat

	// Наши кастомные поля (Go пропустит их при Unmarshal из JSON, заполним сами)
	IsPirate     bool   `json:"-"`
	PirateReason string `json:"-"`
}

// Регулярка для поиска следов плагинов авторизации в тексте кика
var AuthKeywordsRegex = regexp.MustCompile(`(?i)(reg|login|auth|log in|register|войти|вход|пароль|авториз)`)

// Список известных плагинов авторизации для проверки через Query
var AuthPlugins = []string{"authme", "loginsecurity", "advancedlogin", "xauth", "userconn", "fastlogin"}

// Список вариантов имен ботов
var BotNameVariants = []string{"I_am_player", "Real_player", "Bratanchik228", "MrBratik"}

var LicenseIndicators = []string{"Не удалось проверить имя пользователя", "Failed to verify username", "Microsoft"}
var RegisterIndicators = []string{"/reg", "/register", "зарегистрируйтесь", "regicter", "confirmpassword", "пароля"}
var LoginIndicators = []string{"/login", "вход", "/log", "войдите", "login"}
var ErrorIndicators = []string{"Не удалось выполнить ping этого IP", "System.Net.Sockets.SocketException"}
