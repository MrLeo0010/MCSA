package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"minecraft_server_analyser/src"

	"github.com/gookit/color"
)

func printHelp() {
	color.Println("<blue>Справка:</>")
	color.Println("<blue>  help - вывести эту справку</>")
	color.Println("<blue>  exit/quit - выйти из программы</>")
	color.Println("<blue>  update - обновить MCC</>")
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	color.Println("<green>==================================================</>")
	color.Println("<green>        MCSA EXTENDED: ПОЛНЫЙ АНАЛИЗ СЕРВЕРА       </>")
	color.Println("<green>==================================================</>")

	for {
		color.Print("\n<cyan>Введи IP сервера для анализа (help для справки): </>")
		input, err := reader.ReadString('\n')
		if err != nil {
			color.Printf("<red>Ошибка ввода: %v</>\n", err)
			continue
		}

		userInput := strings.TrimSpace(input)
		if userInput == "" {
			continue
		}
		if userInput == "help" {
			printHelp()
			continue
		}
		if userInput == "exit" || userInput == "quit" {
			break
		}
		if userInput == "update" {
			updateMCC()
			continue
		}

		// === ШАГ 1: БАЗОВЫЙ АНАЛИЗ ===
		baseResult := runBaseWorkflow(userInput)

		// === ШАГ 2: ГЛУБОКИЙ АНАЛИЗ (MCC) ===
		mccResult := runMCCWorkflow(userInput)

		// === ШАГ 3: КРАСИВЫЙ СВОДНЫЙ ВЫВОД ===
		printResults(mccResult, baseResult)
	}
}

// printResults собирает всю аналитику воедино и выводит итоговую картину
func printResults(mccResult, baseResult string) {
	color.Println("\n<green>==================================================</>")
	color.Println("<green>               ИТОГОВЫЙ СВОДНЫЙ ОТЧЕТ             </>")
	color.Println("<green>==================================================</>")

	color.Printf("<cyan>Базовый анализ:</> %s\n", baseResult)
	color.Printf("<cyan>Анализ через MCC:</> %s\n", mccResult)
	color.Println("<green>==================================================</>")
}

// runBaseWorkflow пингует сервер, парсит MOTD и онлайн
func runBaseWorkflow(serverIP string) string {
	color.Println("\n<yellow>[1/2] Запуск базового анализа...</>")

	res, err := src.PingServer(serverIP, 5*time.Second)
	if err != nil {
		color.Printf("<red>❌ Ошибка базового подключения: %v</>\n", err)
		return fmt.Sprintf("Ошибка пинга (%v)", err)
	}

	motd := strings.TrimSpace(res.ParseMOTD())

	// Выводим инфу в процессе, чтобы юзер видел промежуточный статус
	color.Println("<green>--------------------------------------------------</>")
	color.Printf("<green>Версия:</> %s (Протокол %d)\n", res.Version.Name, res.Version.Protocol)
	color.Printf("<green>Игроки:</> %d/%d\n", res.Players.Online, res.Players.Max)
	color.Printf("<green>MOTD:</>   %s\n", motd)
	color.Println("<green>--------------------------------------------------</>")

	// Возвращаем компактную строку для финального отчета
	return fmt.Sprintf("Успешно пингуется | Версия: %s | Игроки: %d/%d", res.Version.Name, res.Players.Online, res.Players.Max)
}

// runMCCWorkflow запускает бота и возвращает текстовый вердикт
func runMCCWorkflow(serverIP string) string {
	color.Println("\n<yellow>[2/2] Запуск глубокого анализа через MCC...</>")

	executablePath, err := os.Executable()
	if err != nil {
		return "Ошибка: " + err.Error()
	}
	os.Chdir(filepath.Base(executablePath))
	os.Chdir(".\\MCC")
	mccPath := ".\\MinecraftClient.exe"

	botName := fmt.Sprintf("MCSA_Bot_%d", time.Now().Unix()%10000)
	if len(BotNameVariants) > 0 {
		botNameIndex := rand.IntN(len(BotNameVariants))
		botName = BotNameVariants[botNameIndex]
	}

	cmd := exec.Command(mccPath, botName, "-", serverIP, "Offline")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000010, // CREATE_NEW_CONSOLE
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errMsg := fmt.Sprintf("Ошибка подключения StdoutPipe: %v", err)
		color.Printf("<red>[Ошибка] %s</>\n", errMsg)
		return errMsg
	}

	if err := cmd.Start(); err != nil {
		color.Printf("<red>[Ошибка] Не удалось запустить %s: %v</>\n", mccPath, err)
		return "Ошибка запуска MCC (проверь наличие .exe)"
	}

	color.Println("<gray>[i] Окно MCC открыто. Анализируем лог чата...</>")

	verdict := "UNKNOWN"
	matchedLine := ""

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if scanner.Err() != nil {
			color.Printf("<red>[Ошибка чтения строки] %v</>\n", scanner.Err())
			continue
		}

		var (
			found          bool
			start, end     int
			highlightColor string
		)

		// Основные правила
		for _, rule := range IndicatorRules {
			if found, start, end = src.WhereContainsAny(line, rule.Indicators); found {
				verdict = rule.Verdict
				matchedLine = line
				highlightColor = rule.Color
				break
			}
		}

		if found {
			lineRunes := []rune(line)

			color.Printf(
				"<gray>[MCC]</> %s<%s>%s</>%s\n",
				string(lineRunes[:start]),
				highlightColor,
				string(lineRunes[start:end]),
				string(lineRunes[end:]),
			)
			break
		}

		lineLower := strings.ToLower(line)

		switch {
		case src.ContainsAny(lineLower, UpdateIndicators):
			color.Println("<yellow>[MCSA Warning] ДОСТУПНО ОБНОВЛЕНИЕ MCC</>")
			color.Println("<yellow>[MCSA Warning] ОБНОВИТЬ МОЖНО КОМНДОЙ 'update'</>")
			continue
		case src.ContainsAny(lineLower, SessionIndicators):
			verdict = "LICENSE"
			color.Printf("<gray>[MCC] %s</>\n", line)

		case src.ContainsAny(lineLower, KickIndicators):
			verdict = "KICKED"
			color.Printf("<gray>[MCC] %s</>\n", line)

		case src.ContainsAny(lineLower, SkipTriggers):
			continue

		default:
			color.Printf("<gray>[MCC] %s</>\n", line)
		}
	}

	// Корректно завершаем процесс
	_ = cmd.Process.Kill()
	color.Println("<gray>[Extended] Процесс MCC завершен.</>")

	switch verdict {
	case "PIRATE_AUTH":
		return "❌ ПИРАТКА (Требуется регистрация, найден плагин авторизации) (Совпадение по: " + strings.TrimSpace(matchedLine) + ")"
	case "NICK_TAKEN":
		return "⚠ ПИРАТКА (Ник бота занят, сервер требует войти по логину) (Совпадение по: " + strings.TrimSpace(matchedLine) + ")"
	case "LICENSE":
		return "⚠ ЛИЦЕНЗИЯ (Вход только с аккаунтом Microsoft) (Совпадение по: " + strings.TrimSpace(matchedLine) + ")"
	case "KICKED":
		return "⚠ КИКНУТ / ЗАБАНЕН (Сервер оборвал соединение)"
	case "LOGIN_ERROR":
		return "⚠ ОШИБКА! Бот не смог войти (Совпадение по: " + strings.TrimSpace(matchedLine) + ")"
	default:
		return "❔ СВОБОДНЫЙ ВХОД (Бот зашел без препятствий)"
	}
}
