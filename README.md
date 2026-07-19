### Language:
- [Русский](README-RUS.md)
- English (Current)

# MCSA (Minecraft Server Analyser)

**MCSA** is a modular Go-based tool designed for scanning, filtering, and analyzing Minecraft game servers. The project automates the entire verification process—from mass filtering open ports to detailed configuration analysis of specific servers.

> ⚠️ **Note:** The project is under active development. The documentation may slightly lag behind the current state of the code, but the overall concept and structure always remain accurate.

---

## System Modules

### 1. MCSA Minimal (Mass Filter)
A high-speed tool for processing large lists of IP addresses and ports (e.g., outputs from Masscan or Advanced Port Scanner).
* Operates in multi-threaded mode (Worker Pool).
* Pings ports using the Minecraft protocol and filters out third-party services (SSH, HTTP, RDP).
* Saves the list of valid game servers to a separate file.

### 2. MCSA Base (Basic Ping)
A CLI utility for quick, single scans of a specific server.
* Retrieves information via Server List Ping.
* Displays ping, protocol version, current MOTD, online player count, and a sample of online players.

### 3. MCSA Extended (Deep Analysis)
A module for detailed reconnaissance of server login parameters.
* Integrates with an isolated *Minecraft Console Client (MCC)* instance.
* Analyzes system messages and chat in real-time, highlighting discovered triggers.
* Automatically determines the authentication type: **Premium** (requires a Microsoft session), **Cracked** (requires registration/login), or **Open Access** (free entry).

---

## Auxiliary Tools

The project includes a companion Python script for data preprocessing. It cleans the "dirty" output from Advanced Port Scanner, extracts port numbers using regular expressions, and formats them into `IP:PORT` format for subsequent loading into `MCSA Minimal`.

---

## Building the Project

Compilation of all modules for the required platform is carried out using ready-made automation scripts:

* **Windows:** `build_win.bat`
* **Linux:** coming soon...
