package main

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"time"

	"minecraft_server_analyser/src"

	"github.com/gookit/color"
)

const (
	MaxWorkers = 500
	Timeout    = 5 * time.Second
)

func main() {
	color.Println("<green>==================================================</>")
	color.Println("<green>       MCSA MINIMAL: МАСС-ФИЛЬТР MINECRAFT        </>")
	color.Println("<green>==================================================</>")

	inputFile := "targets.txt"
	outputFile := "live_servers.txt"

	file, err := os.Open(inputFile)
	if err != nil {
		color.Printf("<red>[ОШИБКА]</> Не найден файл %s рядом с бинарником!\n", inputFile)
		color.Println("<gray>Создай его и запиши туда IP:порт (каждый с новой строки).</>")
		return
	}
	defer file.Close()

	var targets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if scanner.Err() != nil {
			color.Printf("<red>[ОШИБКА]</> Ошибка чтения файла: %v\n", scanner.Err())
			continue
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, ":") {
			line = line + ":25565"
		}
		targets = append(targets, line)
	}

	totalTargets := len(targets)
	color.Printf("<cyan>[i]</> Загружено целей для проверки: <yellow>%d</>\n", totalTargets)
	color.Printf("<gray>[i]</> Запуск пула из %d воркеров...\n\n", MaxWorkers)

	jobs := make(chan string, totalTargets)
	results := make(chan string, totalTargets)

	var wg sync.WaitGroup

	for w := 1; w <= MaxWorkers; w++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	for _, target := range targets {
		jobs <- target
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	outFile, err := os.Create(outputFile)
	if err != nil {
		color.Printf("<red>[ОШИБКА]</> Не удалось создать файл результатов: %v\n", err)
		return
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	liveCount := 0

	for res := range results {
		liveCount++
		color.Printf("<green>[МАТЧ]</> Найдена жизнь: <lightGreen>%s</>\n", res)
		_, _ = writer.WriteString(res + "\n")
		_ = writer.Flush()
	}

	color.Println("\n<green>==================================================</>")
	color.Printf("<cyan>[КОНЕЦ]</> Сканирование завершено!\n")
	color.Printf("Проверено: <yellow>%d</> | Валидных серверов майна: <green>%d</>\n", totalTargets, liveCount)
	color.Printf("Результаты сохранены в: <magenta>%s</>\n", outputFile)
	color.Println("<green>==================================================</>")
	color.Print("\n<cyan>Нажмите Enter для выхода</>")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func worker(jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for target := range jobs {
		res, err := src.PingServer(target, Timeout)
		if err != nil {
			continue
		}

		if res.Version.Name != "" || res.ParseMOTD() != "" {
			results <- target
		}
	}
}
