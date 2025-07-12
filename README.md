# Aragomodoro

[![Run Go Tests](https://github.com/aureliomalheiros/Aragomodoro/actions/workflows/tests.yaml/badge.svg)](https://github.com/aureliomalheiros/Aragomodoro/actions/workflows/tests.yaml)

[![Release Aragomodoro](https://github.com/aureliomalheiros/Aragomodoro/actions/workflows/release.yaml/badge.svg)](https://github.com/aureliomalheiros/Aragomodoro/actions/workflows/release.yaml)

![Aragomodoro](assets/img/aragorn.png)

*“Focus is not given to men. It is earned.”*  
— **Aragomodoro**, son of Go, heir to the sacred cycles.

**Aragomodoro** is a command-line Pomodoro timer written in Go, inspired by the honor and discipline of *Aragorn*, son of Arathorn. March through your deep work like a Ranger of the North — 25 minutes at a time.

## Features

- Configurable focus and break durations
- Terminal-based countdown timer
- Optional sound notification (`.wav`)
- Modular structure with `cobra-cli` and internal packages
- Inspired by the world of Tolkien (because why not?)

## Structure

```bash
aragomodoro/
├── assets/           
│   └── sounds/
├── cmd/               
│   └── root.go
├── internal/          
│   ├── pomodoro/      
│   └── sound/         
├── main.go
└── README.md
```

## 🚀 Getting Started

```bash
Aragomodoro is a playful take on the Pomodoro technique, inspired by the spirit of Aragorn from The Lord of the Rings.

Usage:
  aragomodoro [flags]

Flags:
  -b, --break int   Break duration in minutes (default 5)
  -f, --focus int   Focus duration in minutes (default 25)
  -h, --help        help for aragomodoro
```

### 🔧 Prerequisites

- Go 1.20+
- ALSA development libs (Linux only – for audio):

```bash
  sudo apt install libasound2-dev
```
