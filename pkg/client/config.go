package client

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	ConfigFile    = "whisper.conf"
	ServerAddress = "127.0.0.1:6667"
)

type Config struct {
	UID     uint64 `yaml:"uid"`
	Address string `yaml:"address"`
}

func (c *Client) ReadConfig() error {

	path := filepath.Join(c.dir, ConfigFile)

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return c.createNewConfig(path)
	} else if err != nil {
		return fmt.Errorf("failed to get status of config file: %v", err)
	}

	if info.IsDir() {
		return fmt.Errorf("config file is a directory, not file: %s", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}

	d := make([]byte, 1024)
	n, err := f.Read(d)
	if err != nil {
		return fmt.Errorf("faild to read from config file: %v", err)
	}
	d = d[:n]

	config := new(Config)
	if err := yaml.Unmarshal(d, config); err != nil {
		return fmt.Errorf("failed to unmashal config from file: %v", err)
	}
	c.config = config

	return nil
}

func (c *Client) createNewConfig(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer f.Close()

	if err := f.Chmod(os.FileMode(0600)); err != nil {
		return fmt.Errorf("failed to set permissions config file: %v", err)
	}

	config := &Config{
		UID:     0,
		Address: ServerAddress,
	}

	d, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to mashal config: %v", err)
	}

	_, err = f.Write(d)
	if err != nil {
		return fmt.Errorf("failed to write config to file: %v", err)
	}

	c.config = config

	return nil
}
