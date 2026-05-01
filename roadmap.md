# DevSync — Full Cross-Platform Roadmap
## 🎯 Goal

Build a **cross-platform deploy tool** with:

* Go CLI (core engine)
* Tauri GUI (desktop app)
* Support **Windows / Linux / macOS**
* Support **SSH (primary)** + **SMB (Windows fallback)**
* Self-update system using **GitHub Releases**

---

# 🧠 System Contract (DO NOT BREAK)

* GUI = thin layer (no business logic)
* CLI = single source of truth
* Default transport = SSH
* SMB = fallback only (Windows legacy)
* Must work offline (local → server)
* Update system must be separate from deploy logic

---

# ✅ PHASE 1 — CLI CORE

## Tasks

* Implement `devsync push`
* Load `.devsync.json`
* Resolve server + folders

---

# ✅ PHASE 2 — TRANSPORT ABSTRACTION

## Tasks

* Define Transport interface
* Create:

  * SSHTransport
  * SMBTransport

---

# ✅ PHASE 3 — SSH (PRIMARY)

## Tasks

* Use rsync via exec
* Flags: `-avz --delete`
* Stream logs

---

# ✅ PHASE 4 — SMB (WINDOWS FALLBACK)

## Tasks

* Run `net use`
* Run `robocopy /MIR`
* Handle credentials

---

# 🧱 PHASE 5 — AUTO TRANSPORT DETECTION

## Tasks

* Check port 22
* If available → SSH
* Else → SMB

---

# 🧱 PHASE 6 — OS ADAPTER

## Tasks

* Detect OS via `runtime.GOOS`
* Normalize paths
* Execute commands safely

---

# 🧱 PHASE 7 — SYNC ENGINE

## Tasks

* Loop folders
* Build remote path
* Call transport.Sync

---

# 🧱 PHASE 8 — SAFETY FEATURES

## Tasks

* Dry-run mode
* Confirmation prompt

---

# 🧱 PHASE 9 — BUILD MULTI-OS BINARIES

```bash id="b1"
GOOS=windows GOARCH=amd64 go build -o devsync.exe
GOOS=linux GOARCH=amd64 go build -o devsync
GOOS=darwin GOARCH=amd64 go build -o devsync
```

---

# 🧱 PHASE 10 — TAURI GUI SETUP

## Tasks

* Create Tauri + React app
* Minimal UI

---

# 🧱 PHASE 11 — GUI DATA LAYER

## Tasks

* Load `.devsync.json`
* Render:

  * server dropdown
  * folder checkboxes

---

# 🧱 PHASE 12 — GUI → CLI BRIDGE

## Tasks

* Execute CLI from GUI
* Pass arguments dynamically

---

# 🧱 PHASE 13 — REALTIME LOG

## Tasks

* Stream CLI output
* Show logs in UI

---

# 🧱 PHASE 14 — BUNDLE CLI

## Tasks

* Place binary in:

```id="b2"
src-tauri/bin/
```

---

# 🧱 PHASE 15 — BUILD DESKTOP APPS

```bash id="b3"
npm run tauri build
```

Outputs:

* Windows → `.exe / .msi`
* macOS → `.app`
* Linux → `.AppImage`

---

# 🧱 PHASE 16 — UX POLISH

## Tasks

* Disable button while running
* Loading state
* Clean layout

---

# 🧱 PHASE 17 — PROD SAFETY

## Tasks

* Detect production server
* Show warning dialog

---

# 🧱 PHASE 18 — UPDATE SYSTEM (GitHub Releases)

## 🎯 Goal

Use GitHub Releases as backend for updates

---

## Tasks

### 1. Version System

```go id="b4"
const Version = "0.1.0"
```

```bash id="b5"
devsync version
```

---

### 2. Fetch Latest Release

```text id="b6"
https://api.github.com/repos/{owner}/{repo}/releases/latest
```

---

### 3. Parse Response

```json id="b7"
{
  "tag_name": "v0.2.0",
  "assets": [
    {
      "name": "devsync-windows.exe",
      "browser_download_url": "https://github.com/..."
    }
  ]
}
```

---

### 4. OS-Based Asset Selection

```go id="b8"
switch runtime.GOOS {
case "windows": return "devsync-windows.exe"
case "linux": return "devsync-linux"
case "darwin": return "devsync-macos"
}
```

---

### 5. CLI Self Update

Flow:

```id="b9"
download → temp file → replace binary → restart
```

---

### 6. GUI Update (Tauri)

```json id="b10"
"updater": {
  "active": true,
  "endpoints": [
    "https://api.github.com/repos/yourname/devsync/releases/latest"
  ]
}
```

---

### 7. UI Button

```id="b11"
[ 🔄 Check Update ]
```

---

# 🧱 PHASE 19 — RELEASE WORKFLOW

## Tasks

```text id="b12"
1. Build binaries (Windows/Linux/macOS)
2. git tag vX.X.X
3. push tag
4. create GitHub release
5. upload assets
```

---

# 🧱 PHASE 20 — AUTOMATION (OPTIONAL)

## GitHub Actions

* Auto build
* Auto upload release assets

---

# 🧪 TEST MATRIX

| Case             | Expected      |
| ---------------- | ------------- |
| SSH Linux        | success       |
| SSH macOS        | success       |
| Windows SSH      | success       |
| Windows SMB      | success       |
| No SSH           | fallback SMB  |
| Invalid config   | error         |
| No folder        | prompt        |
| Update available | notify        |
| Update fail      | safe fallback |

---

# 🚀 FINAL UX

```text id="b13"
Open app → Select server → Select folders → Deploy → Confirm → Done
```

---

# 🔥 FINAL DIRECTIVE FOR AI

* Do NOT redesign architecture
* Do NOT add extra features
* Complete phases sequentially
* Keep UI minimal
* Separate update system from deploy logic
* Prefer reliability over optimization

---

# 🧠 META

This system is:

→ A cross-platform deploy tool

* GitHub-powered update system
  = A production-ready developer product
