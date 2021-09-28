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
	duration      time.Duration
	quitWatchChan chan struct{}
}

func NewConfig() *Config {
	return &Config{
		quitWatchChan: make(chan struct{}),
		duration:      -1,
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
func (c *Config) GetFileName() string {
	return c.filename
}
func (c *Config) GetCurrentWatchInterval() time.Duration {
	return c.duration
}
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
func (c *Config) UnWatchFileChanges() {
	c.duration = -1
	c.quitWatchChan <- struct{}{}
}
func (c *Config) Set(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return json.Unmarshal(data, &c.config)
}
func (c *Config) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if config, ok := c.config[key]; ok {
		return config
	}
	return nil
}
func (c *Config) GetAll() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}
