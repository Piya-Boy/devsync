import { useEffect, useMemo, useState, type FormEvent } from 'react';
import { invoke } from '@tauri-apps/api/core';

type Server = {
  name: string;
  host: string;
  user: string;
  port: number;
  password?: string;
  keyPath?: string;
};

type Folder = {
  name?: string;
  local: string;
  remote: string;
};

type DevSyncConfig = {
  servers: Server[];
  folders: Folder[];
};

function App() {
  const [config, setConfig] = useState<DevSyncConfig | null>(null);
  const [selectedServer, setSelectedServer] = useState('');
  const [selectedFolders, setSelectedFolders] = useState<Set<string>>(() => new Set());
  const [error, setError] = useState('');
  const [deployOutput, setDeployOutput] = useState('');
  const [isRunning, setIsRunning] = useState(false);
  const [dryRun, setDryRun] = useState(true);

  useEffect(() => {
    let cancelled = false;

    invoke<DevSyncConfig>('load_config')
      .then((loadedConfig) => {
        if (cancelled) {
          return;
        }

        setConfig(loadedConfig);
        setSelectedServer(loadedConfig.servers[0]?.name ?? '');
        setSelectedFolders(new Set(loadedConfig.folders.map(folderId)));
        setError('');
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(String(err));
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  const selectedServerDetails = useMemo(
    () => config?.servers.find((server) => server.name === selectedServer),
    [config?.servers, selectedServer],
  );

  const selectedFolderIds = useMemo(() => Array.from(selectedFolders), [selectedFolders]);

  const toggleFolder = (id: string) => {
    setSelectedFolders((current) => {
      const next = new Set(current);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  const handleDeploy = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    setDeployOutput('');
    setIsRunning(true);

    try {
      const output = await invoke<string>('run_devsync', {
        server: selectedServer,
        folders: selectedFolderIds,
        dryRun,
      });
      setDeployOutput(output || 'devsync completed with no output.');
    } catch (err) {
      setError(String(err));
    } finally {
      setIsRunning(false);
    }
  };

  return (
    <main className="app-shell">
      <section className="hero">
        <p className="eyebrow">DevSync</p>
        <h1>Cross-platform deploys from one simple desktop app.</h1>
        <p className="lede">
          Select a configured server and folders, then deploy through the Go CLI engine.
        </p>
      </section>
      <section className="panel">
        <h2>Deploy</h2>
        {error ? <p className="error">{error}</p> : null}
        {!config && !error ? <p className="muted">Loading .devsync.json...</p> : null}
        {config ? (
          <form className="deploy-form" onSubmit={handleDeploy}>
            <label className="field">
              <span>Server</span>
              <select value={selectedServer} onChange={(event) => setSelectedServer(event.target.value)}>
                {config.servers.map((server) => (
                  <option key={server.name} value={server.name}>
                    {server.name}
                  </option>
                ))}
              </select>
            </label>

            {selectedServerDetails ? (
              <p className="muted">
                {selectedServerDetails.user}@{selectedServerDetails.host}:{selectedServerDetails.port}
              </p>
            ) : null}

            <fieldset className="folder-list">
              <legend>Folders</legend>
              {config.folders.map((folder) => {
                const id = folderId(folder);
                return (
                  <label className="folder-item" key={id}>
                    <input
                      type="checkbox"
                      checked={selectedFolders.has(id)}
                      onChange={() => toggleFolder(id)}
                    />
                    <span>
                      <strong>{folder.name || folder.local}</strong>
                      <small>
                        {folder.local} → {folder.remote}
                      </small>
                    </span>
                  </label>
                );
              })}
            </fieldset>
            <label className="check-row">
              <input type="checkbox" checked={dryRun} onChange={(event) => setDryRun(event.target.checked)} />
              <span>Dry run</span>
            </label>
            <button type="submit" disabled={isRunning || selectedFolderIds.length === 0 || !selectedServer}>
              {isRunning ? 'Deploying...' : 'Deploy'}
            </button>
            {deployOutput ? <pre className="log-output">{deployOutput}</pre> : null}
          </form>
        ) : null}
      </section>
    </main>
  );
}

function folderId(folder: Folder) {
  return `${folder.local}::${folder.remote}`;
}

export default App;
