# VPN Speed Application TODO List

This file contains a list of tasks and ideas for improving the VPN speed application.

## Tasks

1. Add an option to update the nordvpn package. First, check if there is a new version available by running `nordvpn status`. If there is, run the command `sh <(curl -sSf https://downloads.nordcdn.com/apps/linux/install.sh)` to update the package.

2. Improve the terminal outputs by making them prettier. This can be done by adding colors, formatting, ASCII art, etc.

## Modes

- "speed server": Test the speed of different servers and choose the fastest one.
- "ping city": Test the ping latency to different cities and choose the closest one.
- "ping server": Test the ping latency to different servers and choose the best one so far.
- "speed test": Perform a speed test without connecting to any specific server.
- "lowest load": Choose the server with the lowest load.
- "historical": Choose the server with the best historical speed data.
- "ultimate": Combine all modes and historical data to intelligently choose the best server.

## Optimizations

- Try concurrency and better algorithms to improve performance.
- Add an option for no-log to enhance privacy.
- Make the maximum number of connection attempts dynamic.

## Additional Features

- Run the application as a daemon.
- Create a web app to interact with the application.
- Store speed data in a long-term database for smarter server selection.

## Dependencies

- Install the `github.com/go-ping/ping` package using the command `go get -u github.com/go-ping/ping`.
- Install the nordvpn package.

## Country and City Codes

- Create a command to list country names and their codes.
- Create a command to list city names and their codes.
- Create a command to list server names and their codes.
- Alternatively, create your own synonyms for countries and map them to the real ones.

## Additional Tasks

- Add a link to the instructions for installing the nordvpn tool.
- Write tests for different scenarios, such as connecting to different servers and checking the results.

Note: The current mapping for countries is limited. For example, only "us" is recognized as the code for the United States. Consider improving the mapping to accept alternative names like "usa" or "united states".
