use serde::{Deserialize, Serialize};
use std::{
    fs,
    io::{BufRead, BufReader, Read},
    path::PathBuf,
    process::{Command, Stdio},
    sync::{Arc, Mutex},
    thread,
};
use tauri::{Emitter, Manager};

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
struct DevSyncConfig {
    servers: Vec<Server>,
    folders: Vec<Folder>,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
struct Server {
    name: String,
    host: String,
    user: String,
    #[serde(default = "default_ssh_port")]
    port: u16,
    #[serde(default)]
    password: String,
    #[serde(default, rename = "key_path")]
    key_path: String,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
struct Folder {
    #[serde(default)]
    name: String,
    local: String,
    remote: String,
}

fn default_ssh_port() -> u16 {
    22
}

#[tauri::command]
fn load_config() -> Result<DevSyncConfig, String> {
    let path = config_path()?;
    let raw = fs::read_to_string(&path)
        .map_err(|err| format!("failed to read {}: {err}", path.display()))?;
    serde_json::from_str(&raw).map_err(|err| format!("invalid JSON in {}: {err}", path.display()))
}

#[tauri::command]
fn check_update(app: tauri::AppHandle) -> Result<String, String> {
    let cli = cli_path(&app)?;
    let output = Command::new(&cli)
        .arg("update")
        .arg("--check")
        .output()
        .map_err(|err| format!("failed to execute {}: {err}", cli.display()))?;

    let mut combined = String::new();
    combined.push_str(&String::from_utf8_lossy(&output.stdout));
    combined.push_str(&String::from_utf8_lossy(&output.stderr));

    if output.status.success() {
        Ok(combined)
    } else {
        Err(if combined.trim().is_empty() {
            format!("devsync update check exited with status {}", output.status)
        } else {
            combined
        })
    }
}

fn config_path() -> Result<PathBuf, String> {
    let cwd = std::env::current_dir().map_err(|err| format!("failed to get current directory: {err}"))?;
    Ok(cwd.join(".devsync.json"))
}

#[tauri::command]
fn run_devsync(
    app: tauri::AppHandle,
    server: String,
    folders: Vec<String>,
    dry_run: bool,
) -> Result<String, String> {
    let cli = cli_path(&app)?;
    let mut command = Command::new(&cli);
    command.arg("push").arg("--server").arg(server).arg("--yes");

    if dry_run {
        command.arg("--dry-run");
    }

    for folder in folders {
        command.arg("--folder").arg(folder);
    }

    let mut child = command
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .map_err(|err| format!("failed to execute {}: {err}", cli.display()))?;

    let collected = Arc::new(Mutex::new(String::new()));
    let stdout = child
        .stdout
        .take()
        .ok_or_else(|| "failed to capture devsync stdout".to_string())?;
    let stderr = child
        .stderr
        .take()
        .ok_or_else(|| "failed to capture devsync stderr".to_string())?;

    let stdout_handle = stream_reader(app.clone(), stdout, "", collected.clone());
    let stderr_handle = stream_reader(app.clone(), stderr, "", collected.clone());

    let status = child
        .wait()
        .map_err(|err| format!("failed to wait for devsync process: {err}"))?;

    stdout_handle
        .join()
        .map_err(|_| "failed to join stdout reader thread".to_string())?;
    stderr_handle
        .join()
        .map_err(|_| "failed to join stderr reader thread".to_string())?;

    let combined = collected
        .lock()
        .map_err(|_| "failed to collect devsync output".to_string())?
        .clone();

    if status.success() {
        Ok(combined)
    } else {
        Err(if combined.trim().is_empty() {
            format!("devsync exited with status {status}")
        } else {
            combined
        })
    }
}

fn stream_reader<R: Read + Send + 'static>(
    app: tauri::AppHandle,
    reader: R,
    prefix: &'static str,
    collected: Arc<Mutex<String>>,
) -> thread::JoinHandle<()> {
    thread::spawn(move || {
        for line in BufReader::new(reader).lines() {
            match line {
                Ok(line) => {
                    let payload = format!("{prefix}{line}");
                    if let Ok(mut buffer) = collected.lock() {
                        buffer.push_str(&payload);
                        buffer.push('\n');
                    }
                    let _ = app.emit("devsync-log", payload);
                }
                Err(err) => {
                    let payload = format!("failed to read devsync output: {err}");
                    if let Ok(mut buffer) = collected.lock() {
                        buffer.push_str(&payload);
                        buffer.push('\n');
                    }
                    let _ = app.emit("devsync-log", payload);
                    break;
                }
            }
        }
    })
}

fn cli_path(app: &tauri::AppHandle) -> Result<PathBuf, String> {
    if let Ok(path) = std::env::var("DEVSYNC_CLI") {
        return Ok(PathBuf::from(path));
    }

    let exe_name = if cfg!(windows) { "devsync.exe" } else { "devsync" };
    let cwd = std::env::current_dir().map_err(|err| format!("failed to get current directory: {err}"))?;
    let cwd_cli = cwd.join(exe_name);
    if cwd_cli.exists() {
        return Ok(cwd_cli);
    }

    let resource_cli = app
        .path()
        .resource_dir()
        .map_err(|err| format!("failed to resolve resource directory: {err}"))?
        .join(exe_name);
    if resource_cli.exists() {
        return Ok(resource_cli);
    }

    let current_exe = std::env::current_exe().map_err(|err| format!("failed to get current executable: {err}"))?;
    let exe_dir = current_exe
        .parent()
        .ok_or_else(|| "failed to locate current executable directory".to_string())?;
    Ok(exe_dir.join(exe_name))
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![
            check_update,
            load_config,
            run_devsync
        ])
        .run(tauri::generate_context!())
        .expect("error while running DevSync");
}
