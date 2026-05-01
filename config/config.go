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
