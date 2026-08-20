package main

import (
	"encoding/json"
	"fmt"
	"minecraft_server_analyser/src"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
)

func backupMCC() {
	os.Rename("MinecraftClient.exe", "MinecraftClient.exe.bak")
}

func restoreMCC() {
	os.Remove("MinecraftClient.exe")
	os.Rename("MinecraftClient.exe.bak", "MinecraftClient.exe")
}

func removeBackup() {
	os.Remove("MinecraftClient.exe.bak")
}

func testMCC(mccPath string) (bool, error) {
	cmd := exec.Command(mccPath)

	if err := cmd.Start(); err != nil {
		return false, err
	}

	_ = cmd.Process.Kill()
	_ = cmd.Wait()

	return true, nil
}

func updateMCC() {
	latestURL := findLatestVersionURL()
	if latestURL == "" {
		color.Println("<red>Не удалось найти URL последней версии</>")
		return
	}

	executablePath, err := os.Executable()
	if err != nil {
		color.Printf("<red>Ошибка: %s\n</>", err.Error())
		return
	}
	os.Chdir(filepath.Base(executablePath))
	os.Chdir(".\\MCC")
	color.Println("<green>Обновление:</>")
	fmt.Println("	Создание резервной копии")
	backupMCC()
	fmt.Print("	Загрузка новой версии\n    ")
	src.DownloadFile(latestURL, "MinecraftClient.exe", false, "\t")
	fmt.Println("	Тестирование новой версии")
	testResult, err := testMCC(".\\MinecraftClient.exe")
	if !testResult {
		color.Printf("<red>	Тестирование новой версии не прошло с ошибкой: %s\n</>", err)
		color.Println("<red>	Восстанавление резервной копии</>")
		restoreMCC()
		return
	}
	fmt.Println("	Удаление резервной копии")
	removeBackup()
	color.Println("<green>Обновление завершено</>")
}

func findLatestVersionURL() string {
	release_url := fmt.Sprintf("%s/releases/latest", src.BuildGitHubAPIURL(MCCOwner, MCCRepo))
	resp, err := http.Get(release_url)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer resp.Body.Close()

	release := new(Release)

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	for _, asset := range release.Assets {
		if strings.HasPrefix(asset.Name, "MinecraftClient-") &&
			strings.HasSuffix(asset.Name, "-win-x64.exe") {

			color.Println("<green>Ссылка найдена</>")
			return asset.DownloadURL
		}
	}
	return ""
}
