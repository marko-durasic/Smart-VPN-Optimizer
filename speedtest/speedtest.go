package speedtest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"vpn_speed/util"
)

func FetchAndRunSpeedTest() string {
	fmt.Println("Running speed test...")
	return runSpeedtest()
}

func runSpeedtest() string {
	EnsureDependencies()
	cmd := exec.Command("speedtest-cli")
	output, _ := cmd.CombinedOutput()
	return string(output)
}

func ExtractDownloadSpeed(output string) float64 {
	re := regexp.MustCompile(`Download: ([\d.]+) Mbit/s`)
	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		var speed float64
		fmt.Sscan(match[1], &speed)
		return speed
	}
	return 0.0
}

func EnsureDependencies() {
	fmt.Println("Checking dependencies for speed test to be performed...")
	if !util.CommandExists("speedtest-cli") {
		fmt.Println("`speedtest-cli` not found. Attempting to install...")
		installSpeedtestCli()
		if !util.CommandExists("speedtest-cli") {
			fmt.Println("Failed to install `speedtest-cli`. Please install manually and rerun the script.")
			os.Exit(1)
		}
	}
	fmt.Println("Dependencies OK.")
}

func installSpeedtestCli() {
	fmt.Println("Installing speedtest-cli...")
	cmd := exec.Command("sudo", "apt", "install", "speedtest-cli", "-y")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error installing speedtest-cli:", err)
	}
}

const (
	speedTestURL = "https://raw.githubusercontent.com/sivel/speedtest-cli/master/speedtest.py"
)

func FetchAndRunSpeedTestScript() string {
	script := fetchSpeedTestScript()
	return runSpeedtestScript(script)
}

// cache speedtest script
var speedTestScript []byte

func fetchSpeedTestScript() []byte {
	if speedTestScript != nil {
		fmt.Println("Using cached speedtest script.")
		return speedTestScript
	}
	fmt.Println("Fetching speedtest script...")
	resp, err := http.Get(speedTestURL)
	if err != nil {
		fmt.Println("Error fetching speedtest script:", err)
		return nil
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	speedTestScript = data

	fmt.Println("Speedtest script fetched.")
	return data
}

func runSpeedtestScript(script []byte) string {
	fmt.Println("Running speed test...")
	var output []byte
	for i := 0; i < 3; i++ {
		cmd := exec.Command("python", "-")
		cmd.Stdin = bytes.NewReader(script)
		output, _ = cmd.CombinedOutput()
		if strings.Contains(string(output), "Cannot retrieve speedtest configuration") {
			fmt.Println("Failed to retrieve speedtest configuration. Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	fmt.Println("Speed test finished.")
	return string(output)
}
