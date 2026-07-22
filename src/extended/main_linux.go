package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"minecraft_server_analyser/src"

	"github.com/gookit/color"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	color.Println("<green>==================================================</>")
	color.Println("<green>        MCSA EXTENDED: ПОЛНЫЙ АНАЛИЗ СЕРВЕРА (LINUX)</>")
	color.Println("<green>==================================================</>")

	for {
		color.Print("\n<cyan>Введи IP сервера для анализа (exit для выхода): </>")
		input, err := reader.ReadString('\n')
		if err != nil {
			color.Printf("<red>Ошибка ввода: %v</>\n", err)
			continue
		}

		target := strings.TrimSpace(input)
		if target == "" {
			continue
		}
		if target == "exit" || target == "quit" {
			break
		}

		// === ШАГ 1: БАЗОВЫЙ АНАЛИЗ ===
		baseResult := runBaseWorkflow(target)

		// === ШАГ 2: ГЛУБОКИЙ АНАЛИЗ (MCC) ===
		mccResult := runMCCWorkflow(target)

		// === ШАГ 3: КРАСИВЫЙ СВОДНЫЙ ВЫВОД ===
		printResults(mccResult, baseResult)
	}
}

func printResults(mccResult, baseResult string) {
	color.Println("\n<green>==================================================</>")
	color.Println("<green>               ИТОГОВЫЙ СВОДНЫЙ ОТЧЕТ             </>")
	color.Println("<green>==================================================</>")

	color.Printf("<cyan>Базовый анализ:</> %s\n", baseResult)
	color.Printf("<cyan>Анализ через MCC:</> %s\n", mccResult)
	color.Println("<green>==================================================</>")
}

func runBaseWorkflow(serverIP string) string {
	color.Println("\n<yellow>[1/2] Запуск базового анализа...</>")

	res, err := src.PingServer(serverIP, 5*time.Second)
	if err != nil {
		color.Printf("<red>❌ Ошибка базового подключения: %v</>\n", err)
		return fmt.Sprintf("Ошибка пинга (%v)", err)
	}

	motd := strings.TrimSpace(res.ParseMOTD())

	color.Println("<green>--------------------------------------------------</>")
	color.Printf("<green>Версия:</> %s (Протокол %d)\n", res.Version.Name, res.Version.Protocol)
	color.Printf("<green>Игроки:</> %d/%d\n", res.Players.Online, res.Players.Max)
	color.Printf("<green>MOTD:</>   %s\n", motd)
	color.Println("<green>--------------------------------------------------</>")

	return fmt.Sprintf("Успешно пингуется | Версия: %s | Игроки: %d/%d", res.Version.Name, res.Players.Online, res.Players.Max)
}

func runMCCWorkflow(serverIP string) string {
	color.Println("\n<yellow>[2/2] Запуск глубокого анализа через MCC...</>")

	botName := fmt.Sprintf("MCSA_Bot_%d", time.Now().Unix()%10000)
	if len(src.BotNameVariants) > 0 {
		botNameIndex := rand.IntN(len(src.BotNameVariants))
		botName = src.BotNameVariants[botNameIndex]
	}

	// Указываем путь к бинарнику через слэш от текущей папки
	mccPath := "./MinecraftClient" //filepath.Join(".", "MCC", "MinecraftClient")

	// Передаем путь с явным указанием папки
	cmd := exec.Command(mccPath, botName, "-", serverIP, "Offline")

	// Обязательно задаем Dir, чтобы MCC нашел свой файл конфигурации MinecraftClient.ini
	os.Chdir("MCC")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errMsg := fmt.Sprintf("Ошибка подключения StdoutPipe: %v", err)
		color.Printf("<red>[Ошибка] %s</>\n", errMsg)
		return errMsg
	}

	if err := cmd.Start(); err != nil {
		color.Printf("<red>[Ошибка] Не удалось запустить %s: %v</>\n", mccPath, err)
		return "Ошибка запуска MCC (проверь наличие бинарника и права +x)"
	}

	color.Println("<gray>[i] MCC запущен в фоновом режиме. Анализируем лог чата...</>")

	verdict := "UNKNOWN"
	matchedLine := ""

	scanner := bufio.NewScanner(stdout)
	if err := scanner.Err(); err != nil {
		color.Printf("<red>[Ошибка] %v</>\n", err)
		return "Ошибка чтения лога"
	}
	for scanner.Scan() {
		line := scanner.Text()

		var (
			found          bool
			start, end     int
			highlightColor string
		)

		for _, rule := range src.IndicatorRules {
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
		case src.ContainsAny(lineLower, src.SessionIndicators):
			verdict = "LICENSE"
			color.Printf("<gray>[MCC] %s</>\n", line)

		case src.ContainsAny(lineLower, src.KickIndicators):
			verdict = "KICKED"
			color.Printf("<gray>[MCC] %s</>\n", line)

		case src.ContainsAny(lineLower, src.SkipTriggers):
			continue

		default:
			color.Printf("<gray>[MCC] %s</>\n", line)
		}
	}

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
