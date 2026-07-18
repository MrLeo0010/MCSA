package main

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"time"

	// Импортируем твой общий пакет протокола
	"minecraft_server_analyser/src"

	"github.com/gookit/color"
)

// Конфиг сканера
const (
	MaxWorkers = 500             // Сколько горутин фигачат одновременно
	Timeout    = 3 * time.Second // Тайм-аут на один сервер (чтобы не зависать)
)

func main() {
	color.Println("<green>==================================================</>")
	color.Println("<green>       MCSA MINIMAL: МАСС-ФИЛЬТР MINECRAFT        </>")
	color.Println("<green>==================================================</>")

	inputFile := "targets.txt"
	outputFile := "live_servers.txt"

	// 1. Читаем цели из файла
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
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Если порт не указан, автоматом ставим 25565
		if !strings.Contains(line, ":") {
			line = line + ":25565"
		}
		targets = append(targets, line)
	}

	totalTargets := len(targets)
	color.Printf("<cyan>[i]</> Загружено целей для проверки: <yellow>%d</>\n", totalTargets)
	color.Printf("<gray>[i]</> Запуск пула из %d воркеров...\n\n", MaxWorkers)

	// 2. Каналы для распределения задач
	jobs := make(chan string, totalTargets)
	results := make(chan string, totalTargets)

	var wg sync.WaitGroup

	// 3. Запускаем воркеры
	for w := 1; w <= MaxWorkers; w++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// 4. Закидываем задачи в канал
	for _, target := range targets {
		jobs <- target
	}
	close(jobs) // Больше задач не будет

	// 5. Горутина для ожидания завершения воркеров и закрытия канала результатов
	go func() {
		wg.Wait()
		close(results)
	}()

	// 6. Сбор результатов и запись в файл «на лету»
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
		// Выводим в консоль красивый лог
		color.Printf("<green>[МАТЧ]</> Найдена жизнь: <lightGreen>%s</>\n", res)
		// Пишем в файл
		_, _ = writer.WriteString(res + "\n")
		_ = writer.Flush() // Сбрасываем в файл сразу, чтобы не потерять при прерывании
	}

	color.Println("\n<green>==================================================</>")
	color.Printf("<cyan>[КОНЕЦ]</> Сканирование завершено!\n")
	color.Printf("Проверено: <yellow>%d</> | Валидных серверов майна: <green>%d</>\n", totalTargets, liveCount)
	color.Printf("Результаты сохранены в: <magenta>%s</>\n", outputFile)
	color.Println("<green>==================================================</>")
}

// Воркер, выполняющий Server List Ping
func worker(jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for target := range jobs {
		// Делаем полноценный запрос статуса сервера Minecraft
		res, err := src.PingServer(target, Timeout)
		if err != nil {
			// Если порт закрыт, таймаут или там не майн — просто игнорим
			continue
		}

		// Если получили вменяемый ответ (проверяем версию или описание)
		if res.Version.Name != "" || res.ParseMOTD() != "" {
			results <- target
		}
	}
}
