package conman

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"
)

type Config struct {
	config        map[string]interface{}
	filename      string
	mu            sync.RWMutex
	quitWatchChan chan struct{}
}

func NewConfig() *Config {
	return &Config{
		quitWatchChan: make(chan struct{}),
	}
}
func (c *Config) SetFilename(fileName string) {
	c.filename = fileName
}
func (c *Config) SetFromFile() error {
	data, err := ioutil.ReadFile(c.filename)
	if err != nil {
		return err
	}
	return c.Set(data)
}
func (c *Config) WatchFileChanges(duration time.Duration) {
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
func (c *Config) UnWatchFileChanges() {
	c.quitWatchChan <- struct{}{}
}
func (c *Config) Set(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return json.Unmarshal(data, &c.config)
}
func (c *Config) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if config, ok := c.config[key]; ok {
		return config, nil
	}
	return nil, nil
}
func (c *Config) GetAll() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}
