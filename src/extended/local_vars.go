package main

var MCCOwner = "MCCTeam"
var MCCRepo = "Minecraft-Console-Client"

type IndicatorRule struct {
	Indicators []string
	Verdict    string
	Color      string
}

type Release struct {
	Assets []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var IndicatorRules = []IndicatorRule{
	{ErrorIndicators, "LOGIN_ERROR", "red"},
	{RegisterIndicators, "PIRATE_AUTH", "yellow"},
	{LoginIndicators, "NICK_TAKEN", "cyan"},
	{LicenseIndicators, "LICENSE", "red"},
}

var KickIndicators = []string{
	"kicked",
	"banned",
}

var SessionIndicators = []string{
	"failed to verify username",
	"session",
}

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
var UpdateIndicators = []string{
	"Новая версия",
	"Доступна",
	"/upgrade",
}
var SkipTriggers = []string{
	"github.com/MCCTeam",
	"crowdin.com/project/minecraft-console-client",
	"GitHub build",
}
