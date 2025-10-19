# Aragomodoro

[![Run Go Tests](https://github.com/aureliomalheiros/Aragomodoro/actions/workflows/tests.yaml/badge.svg)](https://github.com/aureliomalheiros/Aragomodoro/actions/workflows/tests.yaml)

[![Release Aragomodoro](https://github.com/aureliomalheiros/Aragomodoro/actions/workflows/release.yaml/badge.svg)](https://github.com/aureliomalheiros/Aragomodoro/actions/workflows/release.yaml)

![Aragomodoro](assets/img/aragorn.png)

**Aragomodoro** is a Pomodoro timer written in Go, inspired by the honor and discipline of *Aragorn*, son of Arathorn. March through your deep work like a Ranger of the North — 25 minutes at a time.

Available in both **command-line** and **web interface** modes.

## Features

- 🧭 **CLI Mode**: Traditional terminal-based countdown timer
- 🌐 **Web Interface**: Modern browser-based GUI with real-time updates
- ⚙️ Configurable focus and break durations
- 🔁 Multiple Pomodoro cycles support
- 🔊 Optional sound notifications (`.wav`)
- 📱 Responsive web design for desktop and mobile
- 🎯 WebSocket-powered real-time timer updates
- 🏗️ Modular structure with `cobra-cli` and internal packages
- 🧙‍♂️ Inspired by the world of Tolkien (because why not?)

## Structure

```bash
aragomodoro/
├── assets/           
│   └── sounds/
├── cmd/               
│   ├── root.go
│   └── web.go        # 🌐 Web server command
├── internal/          
│   ├── pomodoro/      
│   ├── sound/         
│   └── web/          # 🌐 Web interface
│       ├── handlers.go
│       ├── server.go
│       └── templates/
│           └── index.html
├── main.go
└── README.md
```

### 🔧 Prerequisites

- Go 1.20+
- ALSA development libs (Linux only – for audio):

```bash
sudo apt install libasound2-dev pkg-config
```

## 🛠️ Installation

### Prerequisites

- [Go](https://golang.org/dl/) version **1.20 or higher** installed and configured (`go` available in your terminal).
- Compatible operating system: **Linux**, **macOS**, or **Windows**.

### Install with `go install`

To install Aragomodoro from source, run:

```bash
go install github.com/aureliomalheiros/aragomodoro@latest
```

#### Verify installation 

After installation, check if the command is available:

```bash
aragomodoro --help
```

## 🚀 Getting Started

Aragomodoro offers two ways to boost your productivity:

### 🧭 CLI Mode (Terminal)

Traditional command-line interface for terminal lovers:

```bash
# Basic usage with default settings (25min focus, 5min break)
aragomodoro

# Custom durations
aragomodoro --focus 30 --break 10

# Multiple cycles
aragomodoro --focus 25 --break 5 --repeat 4

# Available flags
aragomodoro --help
```

### 🌐 Web Interface Mode

Modern browser-based interface with real-time updates:

```bash
# Start web server (default port 8080)
aragomodoro web

# Custom port
aragomodoro web --port 3000

# Then open your browser at:
# http://localhost:8080
```

#### Web Interface Features

- 📱 **Responsive Design**: Works on desktop, tablet, and mobile
- ⚡ **Real-time Updates**: WebSocket-powered live timer
- 🎨 **Modern UI**: Clean, intuitive interface with Aragorn-inspired design
- 🔧 **Easy Configuration**: Set focus/break times and cycles via web form
- ⏸️ **Timer Control**: Start, stop, and monitor progress in real-time

### Command Reference

```bash
Usage:
  aragomodoro [flags]
  aragomodoro [command]

Available Commands:
  web         Start the Aragomodoro web interface
  help        Help about any command

Flags:
  -b, --break int    Break duration in minutes (default 5)
  -c, --continue     Continue the timer during breaks
  -f, --focus int    Focus duration in minutes (default 25)
  -r, --repeat int   Number of Pomodoros before a long break (default 1)
  -h, --help         help for aragomodoro

Web Command:
  aragomodoro web [flags]
  
Web Flags:
  -p, --port int     Port for the web server (default 8080)
```
