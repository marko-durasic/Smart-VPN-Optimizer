package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"vpn_speed/speedtest"

	"github.com/go-ping/ping"
)

type ServerSpeed struct {
	server string
	speed  float64
}

type Servers struct {
	Servers []ServerInfo
	Date    time.Time
}

type ServerInfo struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
}

type ServerPingResult struct {
	Server ServerInfo
	AvgRtt time.Duration
}

func main() {
	// modes
	// 	-mode speed_server
	// - mode ping_city
	// - mode ping_server
	// - mode speed_test

	mode := flag.String("mode", "speed_server", "The mode in which you want to run the script.")

	countryPtr := flag.String("country", "us", "The country to which you want to connect your VPN.")
	speedPtr := flag.Float64("speed", 250, "The desired Mbit/s download speed.")
	maxAttempts := flag.Int("attempts", 3, "The maximum number of attempts to connect to a server with the desired speed.")
	flag.Parse()

	// validate input
	if *mode != "speed_server" && *mode != "ping_city" && *mode != "ping_server" && *mode != "speed_test" {
		fmt.Println("Invalid mode. Exiting...")
		return
	}

	if *mode == "ping_server" {
		// make sure you are disconnected from VPN
		if isConnected() {
			disconnectVPN()
		}
		fmt.Println("Getting servers for country...")
		servers := getServersForCountry(*countryPtr)
		var pingResults []ServerPingResult

		fmt.Println("Pinging servers...")
		for _, server := range servers.Servers {
			avgRtt, err := pingServer(server.Hostname)
			if err != nil {
				fmt.Printf("Error pinging server %s: %s\n", server.Hostname, err)
				continue
			}
			pingResults = append(pingResults, ServerPingResult{
				Server: server,
				AvgRtt: avgRtt,
			})
		}

		// Sort pingResults based on AvgRtt (you'd define ServerPingResult and the sort logic)
		var fastestServer ServerPingResult
		for _, pingResult := range pingResults {
			if fastestServer.Server.Hostname == "" || pingResult.AvgRtt < fastestServer.AvgRtt {
				fastestServer = pingResult
			}
		}

		fmt.Printf("Fastest server: %s, with average ping time of %s\n", fastestServer.Server.Hostname, fastestServer.AvgRtt)

		serverName := strings.Split(fastestServer.Server.Hostname, ".")[0]

		// connect to fastest server
		success := connectToVPN(serverName)

		if !success {
			fmt.Println("Failed to connect to fastest server. Exiting...")
			return
		}
		// show speed of this server
		currentSpeedData := speedtest.FetchAndRunSpeedTest()
		fmt.Println("Extracting current speed...")
		currentSpeed := speedtest.ExtractDownloadSpeed(currentSpeedData)
		fmt.Printf("Current speed: %.2f Mbit/s\n", currentSpeed)
	} else if *mode == "ping_city" {
	} else if *mode == "speed_server" {
		var speeds []ServerSpeed
		var maxRecordedSpeed ServerSpeed

		currentAttempt := 1
		for currentAttempt <= *maxAttempts {
			fmt.Printf("Attempt %d/%d:\n", currentAttempt, *maxAttempts)

			if !isConnectedTo(*countryPtr) {
				fmt.Printf("Not connected to %s\n", *countryPtr)
				if !connectToVPN(*countryPtr) {
					fmt.Println("Debbuging")
					fmt.Println("Failed to connect to VPN. Exiting...")
					currentAttempt++
					continue
				}

				fmt.Printf("Connected to %s\n", *countryPtr)
			}
			server := getCurrentServer()
			if server == "" {
				server = "N/A"
			}
			fmt.Printf("Connected to server: %s\n", server)
			currentSpeedData := speedtest.FetchAndRunSpeedTest()
			fmt.Println("Extracting current speed...")
			currentSpeed := speedtest.ExtractDownloadSpeed(currentSpeedData)
			fmt.Printf("Current speed: %.2f Mbit/s\n", currentSpeed)
			speeds = append(speeds, ServerSpeed{server, currentSpeed})

			if currentSpeed >= *speedPtr {
				fmt.Printf("Connected to a server with desired speed: %.2f Mbit/s\n", currentSpeed)
				return
			} else {
				fmt.Printf("Current server (%s) speed %.2f Mbit/s is below desired speed of %.2f . Reconnecting...\n", server, currentSpeed, *speedPtr)
				if currentSpeed > maxRecordedSpeed.speed {
					maxRecordedSpeed = ServerSpeed{server, currentSpeed}
				}

				// if not last attempt, disconnect and reconnect
				if currentAttempt < *maxAttempts {
					disconnectVPN()
					connectToVPN(*countryPtr)
				}
			}

			currentAttempt++

			// add space between attempts
			fmt.Println()
		}

		fmt.Printf("Max attempts reached without achieving desired speed.\n")
		// if current server is not the fastest recorded server, connect to the fastest recorded server
		if getCurrentServer() != maxRecordedSpeed.server {
			fmt.Printf("Connecting to the fastest recorded server (%s), with maximum recorded speed of (%.2f) Mbit/s...\n", maxRecordedSpeed.server, maxRecordedSpeed.speed)
			result := exec.Command("nordvpn", "c", maxRecordedSpeed.server).Run()
			if result != nil {
				fmt.Println("Error connecting to the fastest recorded server. Exiting...")
				return
			}
			fmt.Println("Connected to the fastest recorded server.")
		} else {
			fmt.Println("Fastest recorded server is already connected.")
		}
	} else if *mode == "speed_test" {
		fmt.Println("Running speed test...")
		currentSpeedData := speedtest.FetchAndRunSpeedTest()
		fmt.Println("Extracting current speed...")
		currentSpeed := speedtest.ExtractDownloadSpeed(currentSpeedData)
		fmt.Printf("Current speed: %.2f Mbit/s\n", currentSpeed)
	}
}

// cache 24h
var servers Servers

func getServersForCountry(country string) Servers {
	// check if cache is not empty and if not older than 24h
	if servers.Servers != nil && time.Now().Sub(servers.Date) < 24*time.Hour {
		fmt.Println("Using cached server info.")
		return servers
	}
	resp, err := http.Get("https://api.nordvpn.com/v1/servers?limit=9999999")
	if err != nil {
		fmt.Println("Error fetching server info:", err)
		return Servers{}
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var servers []ServerInfo
	err = json.Unmarshal(data, &servers)
	if err != nil {
		fmt.Println("Error parsing server info:", err)
		return Servers{}
	}

	var countryServers []ServerInfo
	for _, server := range servers {
		// to lower
		serverCountry := strings.ToLower(server.Name)
		if strings.Contains(serverCountry, country) {
			countryServers = append(countryServers, server)
		}
	}

	return Servers{
		Servers: countryServers,
		Date:    time.Now(),
	}
}

func pingServer(hostname string) (time.Duration, error) {
	pinger, err := ping.NewPinger(hostname)
	if err != nil {
		return 0, err
	}

	pinger.Count = 1 // send only one ping request
	pinger.Timeout = 3 * time.Second

	err = pinger.Run()
	if err != nil {
		return 0, err
	}

	return pinger.Statistics().AvgRtt, nil
}

func getCurrentServer() string {
	out, err := exec.Command("nordvpn", "status").Output()
	if err != nil {
		return ""
	}

	re := regexp.MustCompile(`Hostname: (\S+)\.nordvpn\.com`)
	matches := re.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func extractHostname(status string) string {
	re := regexp.MustCompile(`Hostname: (\S+)\.nordvpn\.com`)
	matches := re.FindStringSubmatch(status)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func isConnectedTo(country string) bool {
	fmt.Println("Checking current VPN status...")
	output, _ := exec.Command("nordvpn", "status").CombinedOutput()
	hostname := extractHostname(string(output))
	if strings.Contains(strings.ToLower(string(output)), "connected") && !strings.Contains(strings.ToLower(string(output)), "disconnected") && strings.Contains(strings.ToLower(string(hostname)), strings.ToLower(country)) {
		fmt.Println("Connected to the correct country.")
		return true
	}
	fmt.Println("Not connected to the correct country.")
	return false
}

func isConnected() bool {
	fmt.Println("Checking current VPN status...")
	output, _ := exec.Command("nordvpn", "status").CombinedOutput()
	statusLine := regexp.MustCompile(`(?i)Status: (\w+)`)
	matches := statusLine.FindStringSubmatch(string(output))
	if len(matches) > 1 && strings.ToLower(matches[1]) == "connected" {
		fmt.Println("Connected to VPN.")
		return true
	}
	fmt.Println("Not connected to VPN.")
	return false
}

func connectToVPN(argument string) bool {
	fmt.Printf("Connecting to %s via NordVPN...\n", argument)
	exec.Command("nordvpn", "c", argument).Run()

	timeout := time.After(60 * time.Second) // max timeout after 30 seconds
	tick := time.Tick(2 * time.Second)      // check every 2 seconds

	fmt.Println("Waiting for VPN connection to establish...")
	for {
		select {
		case <-timeout:
			fmt.Println("Timed out waiting for VPN connection.")
			return false
		case <-tick:
			output, err := exec.Command("nordvpn", "status").Output()
			if err == nil && strings.Contains(string(output), "Connected") {
				fmt.Println("VPN connection established.")
				return true
			}
		}
	}
}

func disconnectVPN() {
	fmt.Println("Disconnecting from current VPN server...")
	exec.Command("nordvpn", "d").Run()

	// wait till disconnected not using timeout
	for isConnected() {
		time.Sleep(2 * time.Second)
	}
	fmt.Println("Disconnected from VPN.")
}
