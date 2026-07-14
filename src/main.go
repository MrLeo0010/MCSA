package main

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/gookit/color"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		color.Print("<cyan>Введи IP сервера (или IP:порт, exit для выхода): </>")
		input, err := reader.ReadString('\n')
		if err != nil {
			color.Printf("<red>Ошибка ввода: %v</>\n", err)
			continue
		}

		targetServer := strings.TrimSpace(input)
		if targetServer == "" {
			continue
		}

		if targetServer == "exit" || targetServer == "quit" {
			color.Println("<yellow>Выход...</>")
			break
		}

		color.Printf("<gray>Опрос %s...</>\n", targetServer)

		start := time.Now()
		res, err := PingServer(targetServer, 5*time.Second)
		if err != nil {
			color.Printf("<red>[ОШИБКА] %v</>\n\n", err)
			continue
		}
		duration := time.Since(start)

		// Вывод результатов с красивыми тегами gookit/color
		color.Println("\n<green>=== РЕЗУЛЬТАТЫ АНАЛИЗА ===</>")

		// Пинг
		color.Printf("<blue>Пинг:</>        %v\n", duration.Round(time.Millisecond))

		// Ядро и Версия
		color.Printf("<blue>Ядро/Версия:</> <magenta>%s</> (Протокол: %d)\n", res.Version.Name, res.Version.Protocol)

		// Игроки
		color.Printf("<blue>Онлайн:</>      <green>%d</> / <red>%d</>\n", res.Players.Online, res.Players.Max)

		// MOTD
		motd := strings.TrimSpace(res.ParseMOTD())
		color.Printf("<blue>MOTD:</>        %s\n", motd)

		// Список игроков (если отдает)
		if len(res.Players.Sample) > 0 {
			color.Println("<blue>Игроки онлайн (выборка):</>")
			for _, player := range res.Players.Sample {
				color.Printf("  - <yellow>%s</> (<gray>%s</>)\n", player.Name, player.ID)
			}
		} else {
			if res.Players.Online > 0 {
				color.Println("<blue>Игроки онлайн:</> <gray>[список скрыт настройками сервера]</>")
			} else {
				color.Println("<blue>Игроки онлайн:</> <gray>[нет игроков на сервере]</>")
			}
		}

		color.Println("<green>==========================</>\n")
	}
}
