use serde::{Deserialize, Serialize};
use std::{fs, path::PathBuf};

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

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![load_config])
        .run(tauri::generate_context!())
        .expect("error while running DevSync");
}
