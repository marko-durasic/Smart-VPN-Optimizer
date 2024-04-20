# VPN Connection Optimizer

Optimize your VPN connection with various strategies to achieve the best speed possible. Initially built as a wrapper for the NordVPN Linux tool, this Go-powered tool offers multiple modes to intelligently select the fastest server based on your requirements. While it currently supports NordVPN, the architecture allows for potential expansion to other VPN providers in the future. Additional tools may be added as well.

The motivation for creating this tool stemmed from noticing that VPN providers do not always connect you to the fastest server available. As I often watch international sports streams, it was crucial to ensure the best possible speed. Therefore, I developed this tool to consistently find the best server.

## Features

- **Speed Server Mode**: Directly connect and check the speed.
  
- (TODO) **Ping City Mode**: Check latency to different cities within a country.
  
- **Ping Server Mode**: Granularly check latency to each server in a country.
  
- **Speed Test Mode**: Evaluate the speed of your current connection without changing servers.
  
- (TODO) **Lowest Load Mode**: Connect to the server with the lowest load for potentially better speeds.
  
- (TODO) **Historical Mode**: Rely on historical data to connect to servers that have previously provided good speeds.
  
- (TODO) **ULTIMATE MODE**: A combination of all modes, providing the most thorough evaluation for the best connection.

## Installation

1. Ensure you have `Go` and `NordVPN` installed on your system (Linux and macOS are supported).
2. Clone this repository.
3. Run `go build` in the directory to compile the executable. `go build -o vpn_optimizier ./cmd` is the command to compile the executable.

## Usage

Run the executable with the following command:
```sh
./vpn_optimizer -mode=<MODE> -country=<COUNTRY_NAME> -speed=<DESIRED_SPEED> -attempts=<MAX_ATTEMPTS>
```

Replace `<MODE>`, `<COUNTRY_NAME>`, `<DESIRED_SPEED>` in Mbit/s, and `<MAX_ATTEMPTS>` with your desired values.

For example:
```sh
./vpn_optimizer -mode=speed_server -country=serbia -speed=190 -attempts=2
```

## Modes

### Speed Server Mode

This is the traditional method of connecting and then checking the speed. It's a direct method but may not always yield the best results.

```sh
./vpn_optimizer -mode=speed_server -country=<COUNTRY_NAME> -speed=<DESIRED_SPEED> -attempts=<MAX_ATTEMPTS>
```
**Note:** Country is optional in this mode, as it will connect to the US as the default.
Speed is optional in this mode, as it will default to 25 Mbps.
Attempts is optional in this mode, as it will default to 3.

### Ping City Mode (TODO)

This checks the latency to different cities in a country. Lower latency can sometimes correlate to faster speeds, but not always. It's a quicker way than connecting to each server but may not always be the most accurate in predicting speeds. 

```sh
./vpn_optimizer -mode=ping_city -country=<COUNTRY_NAME> -speed=<DESIRED_SPEED> -attempts=<MAX_ATTEMPTS>
```

### Ping Server Mode

This mode checks the latency to each individual server within a country. It can sometimes yield much better results. If you need a quick solution, this method may be the best choice. The tool retrieves the list of all servers using the NordVPN API and caches it for improved performance.


```sh
./vpn_optimizer -mode=ping_server -country=<COUNTRY_NAME>
```

### Speed Test Mode

This mode just tests the current connection's speed, without trying to change servers. Useful as a standalone utility or for diagnostic purposes.

```sh
./vpn_optimizer -mode=speed_test
```
### ULTIMATE MODE (TODO)

This mode combines data and strategies from all the other modes. It first checks historical data, then if there's no conclusive decision, it checks server loads, then ping times, and then finally does direct speed tests. It's the most thorough but also the most time-consuming.


```sh
./vpn_optimizer -mode=ultimate -country=serbia -speed=190 -attempts=2
```

## Contributing

Contributions are welcome! Please fork this repository and create a pull request with your changes.

## License

This project is licensed under the MIT License.