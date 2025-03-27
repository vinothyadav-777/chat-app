package utils

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/vinothyadav-777/chat-app/constants"
	configs "github.com/vinothyadav-777/chat-app/utils/go-config-client"
)

type ClientConfig struct {
	configs.Client
	env map[string]string
}

// for test mode initialisation, put a sync once
var o sync.Once

// Client is the instance of the config client to be used by the application
var clientConfig *ClientConfig

// InitTestModeConfigs is used to initialize the configs from local repo
func InitConfigs(directory string, configNames ...string) error {
	var err error
	var c configs.Client
	o.Do(func() {
		c, err = configs.New(configs.Options{
			Provider: configs.FileBased,
			Params: map[string]interface{}{
				constants.ConfigDirectoryKey: directory,
				constants.ConfigNamesKey:     configNames,
				constants.ConfigTypeKey:      "json",
				constants.SecretNamesKey:     []string{},
				constants.SecretDirectoryKey: directory,
				constants.SecretTypeKey:      "json",
			},
		})
		if err == nil {
			clientConfig = getClient(c)
		}

	})
	return err
}

// Get is used to get the instance of the client
func GetClient() *ClientConfig {
	return clientConfig
}

// GetStringWithEnv is used to get the config by filling the variables from environment variables
// for example, say a config value is ${XYZ}/abc, and the value of environment variable XYZ is ABC,
// then this function will return XYZ/abc.
func (c *ClientConfig) GetStringWithEnv(config, key string) (string, error) {
	// first fetch the config value
	s, err := c.GetString(config, key)
	// if error no pointing moving ahead
	if err != nil {
		return s, err
	}
	// now time to look for and replace with all the environment variables
	for k, v := range c.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s, nil
}

// GetStringWithEnvD is used to get the config with default value by filling the variables from environment variables
// for example, say a config value is ${XYZ}/abc, and the value of environment variable XYZ is ABC,
// then this function will return XYZ/abc.
func (c *ClientConfig) GetStringWithEnvD(config, key, defaultValue string) string {
	// first fetch the config value
	s, err := c.GetString(config, key)
	// if error no pointing moving ahead
	if err != nil {
		return defaultValue
	}
	// now time to look for and replace with all the environment variables
	for k, v := range c.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s
}

func (c *ClientConfig) GetInterfaceEnvD(key interface{}, defaultValue string) string {
	s, ok := key.(string)
	if !ok {
		return defaultValue
	}
	for k, v := range c.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s
}

func getEnvironment() map[string]string {
	env := os.Environ()
	result := make(map[string]string)
	for _, e := range env {
		s := strings.Split(e, "=")
		if len(s) >= 2 {
			result[s[0]] = strings.Join(s[1:], "=")
		}
	}
	return result
}

func getClient(c configs.Client) *ClientConfig {
	return &ClientConfig{Client: c, env: getEnvironment()}
}
