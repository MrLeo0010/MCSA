package main

import (
	"bufio"
	"minecraft_server_analyser/src"
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

		// === ЛОГ ВЫПОЛНЕНИЯ ===
		color.Printf("\n<gray>Опрос %s...</>\n", targetServer)
		start := time.Now()
		res, err := src.PingServer(targetServer, 5*time.Second)
		if err != nil {
			color.Printf("<red>[ОШИБКА] Не удалось опросить сервер: %v</>\n\n", err)
			continue
		}
		duration := time.Since(start)

		// === ЧИСТЫЙ КРАСИВЫЙ ФИНАЛЬНЫЙ ВЫВОД ===
		color.Println("\n<green>==================================================</>")
		color.Println("<green>               РЕЗУЛЬТАТЫ АНАЛИЗА                 </>")
		color.Println("<green>==================================================</>")

		// Основные метрики
		color.Printf("<blue>Пинг:</>        %v\n", duration.Round(time.Millisecond))
		color.Printf("<blue>Ядро/Версия:</> <magenta>%s</> (Протокол: %d)\n", res.Version.Name, res.Version.Protocol)

		motd := strings.TrimSpace(res.ParseMOTD())
		color.Printf("<blue>MOTD:</>        %s\n", motd)

		color.Printf("<blue>Онлайн:</>      <green>%d</> / <red>%d</>\n", res.Players.Online, res.Players.Max)

		// Игроки
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

		color.Println() // Пустая строка для читаемости

		// Вердикт по типу сервера (Лицензия/Пиратка)
		if res.IsPirate {
			color.Printf("<red>❌ Режим: ПИРАТСКИЙ (online-mode = false)</>\n")
			color.Printf("   Основание: <yellow>%s</>\n", res.PirateReason)
		} else {
			color.Printf("<green>✅ Режим: ЛИЦЕНЗИОННЫЙ (online-mode = true)</>\n")
			color.Printf("   Основание: %s\n", res.PirateReason)
		}

		color.Println("<green>==================================================</>\n")
	}
}
