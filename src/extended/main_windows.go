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

	// Импортируем твой общий пакет
	"minecraft_server_analyser/src"

	"github.com/gookit/color"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	color.Println("<green>==================================================</>")
	color.Println("<green>        MCSA EXTENDED: ПОЛНЫЙ АНАЛИЗ СЕРВЕРА       </>")
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

	os.Chdir(".\\MCC")
	mccPath := ".\\MinecraftClient.exe"

	botName := fmt.Sprintf("MCSA_Bot_%d", time.Now().Unix()%10000)
	if len(src.BotNameVariants) > 0 {
		botNameIndex := rand.IntN(len(src.BotNameVariants))
		botName = src.BotNameVariants[botNameIndex]
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

		// Печатаем логи MCC красивым приглушенным серым
		color.Printf("<gray>[MCC] %s</>\n", line)

		lineLower := strings.ToLower(line)

		// Проверяем индикаторы ошибок
		if src.ContainsAny(lineLower, src.ErrorIndicators) {
			verdict = "ERROR"
			matchedLine = line
			break
		}
		if src.ContainsAny(lineLower, src.RegisterIndicators) {
			verdict = "PIRATE_AUTH"
			matchedLine = line
			break
		}
		// Проверяем индикаторы авторизации (логина)
		if src.ContainsAny(lineLower, src.LoginIndicators) {
			verdict = "NICK_TAKEN"
			matchedLine = line
			break
		}
		if strings.Contains(lineLower, "failed to verify username") || strings.Contains(lineLower, "session") {
			verdict = "LICENSE"
			break
		}
		if strings.Contains(lineLower, "kicked") || strings.Contains(lineLower, "banned") {
			verdict = "KICKED"
			break
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
		return "✅ ЛИЦЕНЗИЯ (Вход только с аккаунтом Microsoft)"
	case "KICKED":
		return "⚠ КИКНУТ / ЗАБАНЕН (Сервер оборвал соединение)"
	case "ERROR":
		return "⚠ ОШИБКА! Бот не смог войти (Совпадение по: " + strings.TrimSpace(matchedLine) + ")"
	default:
		return "❔ СВОБОДНЫЙ ВХОД (Бот зашел без препятствий)"
	}
}
