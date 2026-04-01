# Go Serial Terminal

A lightweight, high-performance serial terminal built with **Go** and the **Fyne** GUI toolkit. Designed for developers working with microcontrollers (Arduino, ESP32, STM32) who need a reliable tool with real-time data visualization.

## ✨ Features

* **Real-time Monitoring:** Smooth serial data reception and transmission.
* **HEX Mode:** Toggle between ASCII and HEX views with multiple formatting options (`XX`, `0xXX`, `\xXX`).
* **Customizable UI:** * Adjustable **Font Size** for the receiver log (saved in settings).
    * Multiple text colors for TX and RX.
    * Bold text toggles for better readability.
* **Advanced Logging:**
    * Optional **Timestamps** (HH:MM:SS.mmm) for every incoming packet.
    * Line numbering for easy data tracking.
* **Hardware Control:** * Visual TX/RX blinkers.
    * Configurable Baud Rate, Data Bits, Parity, and Stop Bits.
    * Local Echo and Line Render (`[CR][LF]`) toggles.
* **Auto-Configuration:** Automatically remembers your last used Port and settings via `settings.json`.

## 🚀 Getting Started

### Prerequisites
* **Go** (1.20 or later recommended)
* **Fyne Dependencies:** (Internal CGO requirements for Linux/Windows)

### Installation
1. Clone the repository:
   ```bash
   git clone [https://github.com/gufim/TERMINAL-GO.git]
   cd TERMINAL-GO