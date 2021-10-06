package conman

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Config struct {
	config        map[string]interface{}
	filename      string
	mu            sync.RWMutex
	duration      time.Duration
	quitWatchChan chan struct{}
}

// returns new config
func NewConfig() *Config {
	return &Config{
		quitWatchChan: make(chan struct{}),
		duration:      -1,
	}
}

// sets filename
func (c *Config) SetFilename(fileName string) {
	c.filename = fileName
}

// sets config according to the byte array data provided
func (c *Config) Set(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return json.Unmarshal(data, &c.config)
}

// sets config from file
func (c *Config) SetFromFile() error {
	data, err := ioutil.ReadFile(c.filename)
	if err != nil {
		return err
	}
	return c.Set(data)
}

// returns current file name
func (c *Config) GetFileName() string {
	return c.filename
}

// write configuration to file
func (c *Config) WriteConfigToFile(data []byte) error {
	return ioutil.WriteFile(c.filename, data, fs.FileMode(os.O_RDWR))
}

// returns current watch interval
func (c *Config) GetCurrentWatchInterval() time.Duration {
	return c.duration
}

// enable watching file changes in regular intervals using go routine
func (c *Config) WatchFileChanges(duration time.Duration) {
	c.duration = duration
	ticker := time.NewTicker(duration)

	go func() {
		for {
			select {
			case <-ticker.C:
				c.SetFromFile()
			case <-c.quitWatchChan:
				return
			}
		}
	}()
}

// disable watching to file changes
func (c *Config) UnWatchFileChanges() {
	c.duration = -1
	c.quitWatchChan <- struct{}{}
}

// returns the interface of the requested config
func (c *Config) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if config, ok := c.config[key]; ok {
		return config
	}
	return nil
}

// returns all the config
func (c *Config) GetAll() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}
