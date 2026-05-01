package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const ConfigFile = ".devsync.json"

type Server struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Port     int    `json:"port"`
	Password string `json:"password,omitempty"`
	KeyPath  string `json:"key_path,omitempty"`
}

type Folder struct {
	Name   string `json:"name"`
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

func (f Folder) ID() string {
	return f.Local + "::" + f.Remote
}

type Config struct {
	Servers []Server `json:"servers"`
	Folders []Folder `json:"folders"`
}

// Load reads .devsync.json from the current working directory.
func Load() (*Config, error) {
	path, err := findConfigFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func findConfigFile() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get working directory: %w", err)
	}
	path := filepath.Join(cwd, ConfigFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("%s not found in %s", ConfigFile, cwd)
	}
	return path, nil
}

func (c *Config) validate() error {
	if len(c.Servers) == 0 {
		return fmt.Errorf("config must have at least one server")
	}
	for i, s := range c.Servers {
		if s.Name == "" {
			return fmt.Errorf("server[%d]: name is required", i)
		}
		if s.Host == "" {
			return fmt.Errorf("server %q: host is required", s.Name)
		}
		if s.User == "" {
			return fmt.Errorf("server %q: user is required", s.Name)
		}
		if s.Port == 0 {
			c.Servers[i].Port = 22
		}
	}
	if len(c.Folders) == 0 {
		return fmt.Errorf("config must have at least one folder")
	}
	for i, f := range c.Folders {
		if f.Local == "" {
			return fmt.Errorf("folder[%d]: local path is required", i)
		}
		if f.Remote == "" {
			return fmt.Errorf("folder[%d]: remote path is required", i)
		}
	}
	return nil
}

// FindServer returns the server with the given name, or the first server if name is empty.
func (c *Config) FindServer(name string) (*Server, error) {
	if name == "" {
		return &c.Servers[0], nil
	}
	for i := range c.Servers {
		if c.Servers[i].Name == name {
			return &c.Servers[i], nil
		}
	}
	return nil, fmt.Errorf("server %q not found in config", name)
}

// SelectFolders returns folders matching the given selectors. Empty selectors return all folders.
// A selector may be a folder name, local path, remote path, or the stable local::remote folder ID.
func (c *Config) SelectFolders(selectors []string) ([]Folder, error) {
	if len(selectors) == 0 {
		return c.Folders, nil
	}

	selected := make([]Folder, 0, len(selectors))
	missing := make([]string, 0)
	used := make(map[string]bool, len(selectors))

	for _, selector := range selectors {
		found := false
		for _, folder := range c.Folders {
			if folderMatchesSelector(folder, selector) {
				if !used[folder.ID()] {
					selected = append(selected, folder)
					used[folder.ID()] = true
				}
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, selector)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("folder selector(s) not found in config: %v", missing)
	}
	return selected, nil
}

func folderMatchesSelector(folder Folder, selector string) bool {
	return selector == folder.Name ||
		selector == folder.Local ||
		selector == folder.Remote ||
		selector == folder.ID()
}
