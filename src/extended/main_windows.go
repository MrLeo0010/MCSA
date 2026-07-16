package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	// Импортируем твой общий пакет (замени имя модуля на свое, если оно другое в go.mod)
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

		// === ШАГ 1: БЫСТРЫЙ БАЗОВЫЙ ПИНГ ===
		color.Println("\n<yellow>[1/2] Запуск базового анализа...</>")

		// Используем твою общую функцию пинга из пакета src
		res, err := src.PingServer(target, 5*time.Second)
		if err != nil {
			color.Printf("<red>❌ Ошибка базового подключения: %v</>\n", err)
			color.Println("<gray>Пробуем запустить MCC напрямую, возможно сервер блокирует обычный пинг...</gray>")
			// Если обычный пинг упал, всё равно даем шанс MCC
			runMCCWorkflow(target)
			continue
		}

		motd := strings.TrimSpace(res.ParseMOTD())
		// Выводим базовую инфу красиво (как в твоей base-версии)
		color.Println("<green>--------------------------------------------------</>")
		color.Printf("<green>Версия:</> %s (Протокол %d)\n", res.Version.Name, res.Version.Protocol)
		color.Printf("<green>Игроки:</> %d/%d\n", res.Players.Online, res.Players.Max)
		// Чистим MOTD от параграфов параграфов цвета (§) перед выводом, если у тебя есть для этого функция
		// Берем готовый метод прямо из твоих моделей!
		color.Printf("<green>MOTD:</>   %s\n", motd)
		color.Println("<green>--------------------------------------------------</>")

		// === ШАГ 2: ГЛУБОКИЙ АНАЛИЗ ЧЕРЕЗ MCC ===
		color.Println("\n<yellow>[2/2] Запуск глубокого анализа через MCC...</>")
		runMCCWorkflow(target)
	}
}

func runMCCWorkflow(serverIP string) {
	mccPath := ".\\MinecraftClient.exe"
	botName := fmt.Sprintf("MCSA_Bot_%d", time.Now().Unix()%10000)

	cmd := exec.Command(mccPath, botName, "-", serverIP, "Offline")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000010, // CREATE_NEW_CONSOLE
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		color.Printf("<red>[Ошибка] Не удалось подключить StdoutPipe: %v</>\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		color.Printf("<red>[Ошибка] Не удалось запустить MinecraftClient.exe: %v</>\n", err)
		color.Println("<yellow>Убедись, что MinecraftClient.exe лежит в одной папке со сканером!</>")
		return
	}

	color.Println("<gray>[i] Окно MCC открыто. Парсим логи чата в реальном времени...</>")

	// Переменная для финального вердикта
	verdict := "UNKNOWN"

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		// Теперь логи выводятся красиво серым цветом, без сырых тегов!
		color.Printf("<gray>[MCC] %s</>\n", line)

		lineLower := strings.ToLower(line)

		// Логика определения вердикта
		if strings.Contains(lineLower, "/register") || strings.Contains(lineLower, "зарегистрируйтесь") {
			verdict = "PIRATE_AUTH"
			break
		}
		if strings.Contains(lineLower, "/login") {
			verdict = "NICK_TAKEN"
			break
		}
		if strings.Contains(lineLower, "failed to verify username") || strings.Contains(lineLower, "session") {
			verdict = "LICENSE"
			break
		}
		// Если бота кикнуло по другой причине (например, спам-фильтр или вайтлист)
		if strings.Contains(lineLower, "kicked") || strings.Contains(lineLower, "banned") {
			verdict = "KICKED"
			break
		}
	}

	// Убиваем процесс MCC
	_ = cmd.Process.Kill()
	color.Println("<gray>[Extended] Процесс MCC завершен.</>")

	// === ШАГ 3: ФИНАЛЬНЫЙ СВОДНЫЙ ОТЧЕТ ===
	color.Println("\n<green>==================================================</>")
	color.Println("<green>               ИТОГОВЫЙ ВЕРДИКТ                   </>")
	color.Println("<green>==================================================</>")

	switch verdict {
	case "PIRATE_AUTH":
		color.Println("<red>❌ Статус авторизации: ТРЕБУЕТСЯ РЕГИСТРАЦИЯ</>")
		color.Println("<gray>Обнаружен плагин авторизации (AuthMe / LogIt / аналоги).</ gray>")
	case "NICK_TAKEN":
		color.Println("<yellow>⚠ Статус авторизации: НИК ЗАНЯТ</>")
		color.Println("<gray>Сервер пиратский, но ник нашего бота уже зарегистрирован в базе.</>")
	case "LICENSE":
		color.Println("<green>✅ Статус авторизации: ЛИЦЕНЗИЯ (Online Mode)</>")
		color.Println("<gray>Вход без лицензионного аккаунта Microsoft невозможен.</>")
	case "KICKED":
		color.Println("<yellow>⚠ Статус авторизации: БОТ БЫЛ КИКНУТ / ЗАБАНЕН</> ")
		color.Println("<gray>Сервер прервал соединение до того, как мы успели прочитать чат.</ gray>")
	default:
		color.Println("<white>❔ Статус авторизации: СВОБОДНЫЙ ВХОД (Возможно)</>")
		color.Println("<gray>Бот успешно зашел, сообщений о регистрации в чате не обнаружено.</ gray>")
	}
	color.Println("<green>==================================================</>")
}
