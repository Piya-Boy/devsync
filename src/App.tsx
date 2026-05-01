import { useEffect, useMemo, useState, type FormEvent } from 'react';
import { invoke } from '@tauri-apps/api/core';
import { listen } from '@tauri-apps/api/event';

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

  useEffect(() => {
    let unlisten: (() => void) | undefined;

    listen<string>('devsync-log', (event) => {
      setDeployOutput((current) => `${current}${event.payload}\n`);
    }).then((cleanup) => {
      unlisten = cleanup;
    });

    return () => {
      unlisten?.();
    };
  }, []);

  const selectedServerDetails = useMemo(
    () => config?.servers.find((server) => server.name === selectedServer),
    [config?.servers, selectedServer],
  );
  const isProduction = selectedServerDetails ? isProductionServer(selectedServerDetails) : false;

  const selectedFolderIds = useMemo(() => Array.from(selectedFolders), [selectedFolders]);
  const deployStatus = isRunning
    ? 'Deploy is running. Logs will appear below.'
    : `${selectedFolderIds.length} folder${selectedFolderIds.length === 1 ? '' : 's'} selected.`;

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

    if (!dryRun && isProduction) {
      const confirmed = window.confirm(
        `You are deploying to production server "${selectedServer}". Continue?`,
      );
      if (!confirmed) {
        return;
      }
    }

    setIsRunning(true);

    try {
      const output = await invoke<string>('run_devsync', {
        server: selectedServer,
        folders: selectedFolderIds,
        dryRun,
      });
      if (!output) {
        setDeployOutput((current) => current || 'devsync completed with no output.');
      }
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
            {isProduction ? (
              <p className="warning">Production server detected. Non-dry-run deploys require confirmation.</p>
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
            <div className="action-row">
              <button type="submit" disabled={isRunning || selectedFolderIds.length === 0 || !selectedServer}>
                {isRunning ? 'Deploying...' : 'Deploy'}
              </button>
              <p className="status-text" aria-live="polite">
                {deployStatus}
              </p>
            </div>
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

function isProductionServer(server: Server) {
  const value = `${server.name} ${server.host}`.toLowerCase();
  return value.includes('prod') || value.includes('production');
}

export default App;
