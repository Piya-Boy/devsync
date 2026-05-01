# DevSync — Cross-Platform Roadmap (AI-Ready)

## 🎯 Goal

Build a **cross-platform deploy tool** with:

* Go CLI (core engine)
* Tauri GUI (desktop app)
* Support **Windows / Linux / macOS**
* Support **SSH (primary)** + **SMB (Windows fallback)**

---

# 🧠 System Contract (DO NOT BREAK)

* GUI = thin layer (no business logic)
* CLI = single source of truth
* Default transport = SSH
* SMB = fallback only (Windows legacy)
* Must work offline (local → server)

---

# 🧱 PHASE 1 — CLI CORE (CROSS-PLATFORM READY)

## Tasks

* Implement `devsync push`
* Load `.devsync.json`
* Resolve server + folders

## Acceptance

* Works on at least one OS (dev machine)

## AI Prompt

```id="p1"
Create a Go CLI using cobra.
Implement command: devsync push [folders...]
Load config from .devsync.json.
Resolve server and folders.
Print resolved values.
```

---

# 🧱 PHASE 2 — TRANSPORT ABSTRACTION

## Tasks

* Define Transport interface
* Create:

  * SSHTransport
  * SMBTransport (stub)

## Acceptance

* Code compiles
* Transport selected by config

## AI Prompt

```id="p2"
Design a Transport interface with method Sync(local, remote string).
Create SSHTransport and SMBTransport structs.
Add factory function to select transport by server.type.
```

---

# 🧱 PHASE 3 — SSH (PRIMARY, ALL OS)

## Tasks

* Use rsync via exec
* Support flags: -avz --delete
* Handle stdout/stderr

## Acceptance

* Sync works on Linux/macOS
* Works on Windows if rsync available

## AI Prompt

```id="p3"
Implement SSHTransport.Sync using rsync.
Command: rsync -avz --delete local/ user@host:/remote/
Stream stdout/stderr to console.
Return errors properly.
```

---

# 🧱 PHASE 4 — SMB (WINDOWS FALLBACK)

## Tasks

* Run `net use`
* Run `robocopy /MIR`
* Handle credentials

## Acceptance

* Works on Windows without SSH

## AI Prompt

```id="p4"
Implement SMBTransport.Sync.
Run:
1. net use \\host /user:user password
2. robocopy local \\host\remote /MIR
Ensure output is visible.
```

---

# 🧱 PHASE 5 — AUTO DETECT TRANSPORT

## Tasks

* Check port 22
* Fallback to SMB if unavailable

## Acceptance

* Correct transport auto-selected

## AI Prompt

```id="p5"
Implement transport auto detection.
Try TCP connect to host:22 (timeout 2s).
If success → SSH
Else → SMB
```

---

# 🧱 PHASE 6 — OS ADAPTER (CRITICAL)

## Tasks

* Detect runtime OS
* Normalize paths
* Execute commands safely

## Acceptance

* Same code runs on Win/Linux/macOS

## AI Prompt

```id="p6"
Implement OS adapter layer.
Detect OS using runtime.GOOS.
Normalize paths:
- Windows: backslash
- Unix: slash
Provide helper for executing commands cross-platform.
```

---

# 🧱 PHASE 7 — SYNC ENGINE

## Tasks

* Loop folders
* Build remote path (basePath + folder.remote)
* Call transport.Sync

## Acceptance

* Multiple folders sync correctly

## AI Prompt

```id="p7"
Implement sync engine.
For each folder:
- Build remote path
- Call selected transport
Print logs for each step.
```

---

# 🧱 PHASE 8 — SAFETY FEATURES

## Tasks

* Dry-run
* Confirmation prompt

## Acceptance

* No accidental overwrite

## AI Prompt

```id="p8"
Add dry-run flag.
rsync: add -n
robocopy: add /L
Add confirmation prompt before execution.
```

---

# 🧱 PHASE 9 — BUILD MULTI-OS BINARIES

## Tasks

* Cross-compile

## Commands

```id="p9"
GOOS=windows GOARCH=amd64 go build -o devsync.exe
GOOS=linux GOARCH=amd64 go build -o devsync
GOOS=darwin GOARCH=amd64 go build -o devsync
```

## Acceptance

* CLI runs on all OS

---

# 🧱 PHASE 10 — TAURI GUI SETUP

## Tasks

* Create Tauri + React app
* Clean UI

## Acceptance

* App launches on dev OS

## AI Prompt

```id="p10"
Create a Tauri app with React + TypeScript.
Render a minimal UI with title "DevSync".
```

---

# 🧱 PHASE 11 — GUI DATA LAYER

## Tasks

* Load `.devsync.json`
* Render:

  * server dropdown
  * folder checkboxes

## Acceptance

* UI reflects real config

## AI Prompt

```id="p11"
Read .devsync.json in frontend.
Render dropdown for servers and checkbox list for folders.
```

---

# 🧱 PHASE 12 — GUI → CLI BRIDGE

## Tasks

* Use Tauri shell plugin
* Execute bundled CLI

## Acceptance

* Button triggers CLI

## AI Prompt

```id="p12"
Use Tauri shell plugin to run devsync binary.
Pass selected folders as arguments.
Capture stdout.
```

---

# 🧱 PHASE 13 — REALTIME LOG

## Tasks

* Stream stdout to UI

## Acceptance

* Logs update live

## AI Prompt

```id="p13"
Stream CLI stdout in Tauri.
Append logs to UI in real time.
Auto-scroll log panel.
```

---

# 🧱 PHASE 14 — BUNDLE CLI WITH GUI

## Tasks

* Place binary in:

```id="p14"
src-tauri/bin/
```

## Acceptance

* GUI runs without external CLI install

---

# 🧱 PHASE 15 — BUILD DESKTOP APPS

## Commands

```id="p15"
npm run tauri build
```

## Outputs

* Windows → .exe / .msi
* macOS → .app
* Linux → AppImage

---

# 🧱 PHASE 16 — UX POLISH

## Tasks

* Disable button while running
* Loading state
* Clean layout

---

# 🧱 PHASE 17 — PROD SAFETY

## Tasks

* Detect prod server
* Show warning dialog

---

# 🧪 TEST MATRIX

| Case           | Expected     |
| -------------- | ------------ |
| SSH Linux      | success      |
| SSH macOS      | success      |
| Windows SSH    | success      |
| Windows SMB    | success      |
| No SSH         | fallback SMB |
| Invalid config | error        |
| No folder      | prompt       |

---

# 🚀 FINAL UX

```id="p16"
Open app → Select server → Select folders → Deploy → Confirm → Done
```

---

# 🔥 FINAL DIRECTIVE FOR AI

* Do NOT add features
* Do NOT redesign architecture
* Complete phases sequentially
* Keep UI minimal
* Prefer reliability over optimization

---

# 🧠 META

This system is:
→ A cross-platform deploy tool
→ With adaptive transport
→ Optimized for real developer workflow
