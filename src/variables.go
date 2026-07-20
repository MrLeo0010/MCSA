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
	PirateStatus       string `json:"-"`
	PirateStatusReason string `json:"-"`
}

// Регулярка для поиска следов плагинов авторизации в тексте кика
var AuthKeywordsRegex = regexp.MustCompile(`(?i)(reg|login|auth|log in|register|войти|вход|пароль|авториз)`)

// Список вариантов имен ботов
var BotNameVariants = []string{
	"Alex228",
	"Den4ik",
	"KirillPlay",
	"Vladik",
	"JustPlayer",
	"NoobMaster",
	"DarkFox",
	"Ghost",
	"StevePro",
	"CraftBoy",
	"PlayerOne",
	"MineCraftik",
	"DragonSlayer",
	"Kotik",
	"DefinitelyHuman",
	"TotallyNotBot",
	"AFK_Player",
	"NoCheats",
	"TrustMeBro",
	"LoadingChunks",
	"Respawning",
	"StoneEnjoyer",
	"HerobrineFan",
	"Join_server",
}

var LicenseIndicators = []string{
	"Не удалось проверить имя пользователя",
	"Failed to verify username",
	"Microsoft",
	"Minecraft account",
}
var RegisterIndicators = []string{
	"/reg",
	"/register",
	"зарегистрируйтесь",
	"regicter",
	"confirmpassword",
	"парол",
}
var LoginIndicators = []string{
	"/login",
	"вход",
	"/log",
	"войдите",
	"login",
}
var ErrorIndicators = []string{
	"Не удалось выполнить ping этого IP",
	"System.Net.Sockets.SocketException",
	"эта версия не поддерживается",
	"Не удается подключиться к серверу",
	"[ERROR]",
	"IOException",
	"at System.",
}
