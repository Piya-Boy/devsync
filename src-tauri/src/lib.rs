use serde::{Deserialize, Serialize};
use std::{fs, path::PathBuf, process::Command};

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

fn config_path() -> Result<PathBuf, String> {
    let cwd = std::env::current_dir().map_err(|err| format!("failed to get current directory: {err}"))?;
    Ok(cwd.join(".devsync.json"))
}

#[tauri::command]
fn run_devsync(server: String, folders: Vec<String>, dry_run: bool) -> Result<String, String> {
    let cli = cli_path()?;
    let mut command = Command::new(&cli);
    command.arg("push").arg("--server").arg(server).arg("--yes");

    if dry_run {
        command.arg("--dry-run");
    }

    for folder in folders {
        command.arg("--folder").arg(folder);
    }

    let output = command
        .output()
        .map_err(|err| format!("failed to execute {}: {err}", cli.display()))?;

    let mut combined = String::new();
    combined.push_str(&String::from_utf8_lossy(&output.stdout));
    combined.push_str(&String::from_utf8_lossy(&output.stderr));

    if output.status.success() {
        Ok(combined)
    } else {
        Err(if combined.trim().is_empty() {
            format!("devsync exited with status {}", output.status)
        } else {
            combined
        })
    }
}

fn cli_path() -> Result<PathBuf, String> {
    if let Ok(path) = std::env::var("DEVSYNC_CLI") {
        return Ok(PathBuf::from(path));
    }

    let exe_name = if cfg!(windows) { "devsync.exe" } else { "devsync" };
    let cwd = std::env::current_dir().map_err(|err| format!("failed to get current directory: {err}"))?;
    let cwd_cli = cwd.join(exe_name);
    if cwd_cli.exists() {
        return Ok(cwd_cli);
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
        .invoke_handler(tauri::generate_handler![load_config, run_devsync])
        .run(tauri::generate_context!())
        .expect("error while running DevSync");
}
