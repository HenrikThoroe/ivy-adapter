// Package conf provides functionality to configure the bahviour
// of the different sub packages and commands.
//
// The configuration is loaded from the file ivyconf.yaml in the
// current working directory (”./”), the home directory (”$HOME”) or the Ivy
// directory (”$IVY_PATH”).
//
// Make sure to call Load() before using any other function in this package.
package conf

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

// A ServerConfig contains the configuration for a server.
// A configuration for a server has to have the following structure in
// a configuration file (yaml, json, etc.):
//
//	<name>:
//		host: <string>
//		port: <int>
//		secure: <bool>
type ServerConfig struct {
	Host   string // The host of the server.
	Port   int    // The port of the server.
	Secure bool   // Whether the server uses HTTPS / WSS.
}

// Configurations for the different servers and storage options.
var (
	evc         ServerConfig // Configuration of the server which provides the engine version control.
	gameServer  ServerConfig // Configuration of the server which provides the game server.
	gameManager ServerConfig // Configuration of the server which provides the game manager.
	test        ServerConfig // Configuration of the server which provides the test server.
	engineStore string       // The path to the directory where the engines are stored.
)

// Load loads the configuration from the file ivyconf.yaml in the
// current working directory (”./”), the home directory (”$HOME”) or the Ivy
// directory (”$IVY_PATH”).
//
// This function should be called before any other function in this package
// to ensure that the configuration is loaded.
func Load(path string) {
	if path == "" {
		viper.AddConfigPath(".")
		viper.AddConfigPath("$IVY_PATH/Ivy")
		viper.AddConfigPath("$HOME/Ivy")
		viper.SetConfigType("yaml")
		viper.SetConfigName("ivyconf")
	} else {
		viper.SetConfigFile(path)
	}

	viper.ReadInConfig()

	initServerConfig(&evc, "evc", "localhost", 4500, false)
	initServerConfig(&gameServer, "game", "localhost", 4502, false)
	initServerConfig(&gameManager, "game-manager", "localhost", 4501, false)
	initServerConfig(&test, "test", "localhost", 4504, false)

	engineStore = viper.GetString("engine-store")

	engineStore, _ = filepath.Abs(engineStore)

	if engineStore == "" {
		home := os.Getenv("HOME")
		ivy := os.Getenv("IVY_PATH")
		base := home + "/Ivy"

		if ivy != "" {
			base = ivy + "/Ivy"
		}

		engineStore = base + "/engines"
	}
}

// GetTestServerConfig returns the configuration of the server which provides the test server.
func GetTestServerConfig() *ServerConfig {
	return &test
}

// GetEngineStore returns the path to the directory where the engines are stored.
func GetEngineStore() string {
	return engineStore
}

// GetEVCConfig returns the configuration of the server which provides the engine version control.
func GetEVCConfig() *ServerConfig {
	return &evc
}

// GetGameServerConfig returns the configuration of the server which provides the game server.
func GetGameServerConfig() *ServerConfig {
	return &gameServer
}

// GetGameManagerConfig returns the configuration of the server which provides the game manager.
func GetGameManagerConfig() *ServerConfig {
	return &gameManager
}

// GetURL returns the URL of the server.
// The URL is constructed from the host, port and whether the server uses HTTPS / WSS.
// The returned URL has the format:
//
// (http(s?)|ws(s?))://<host>:<port>
func (sc *ServerConfig) GetURL() string {
	protocol := "http"
	port := ""

	if sc == &gameServer || sc == &test {
		protocol = "ws"
	}

	if sc.Secure {
		protocol += "s"
	}

	if sc.Port > 0 {
		port = ":" + strconv.Itoa(sc.Port)
	}

	return protocol + "://" + sc.Host + port
}

// initServerConfig initializes the configuration of a server.
// The configuration is loaded using viper identified by the prefix.
func initServerConfig(sc *ServerConfig, prefix string, defHost string, defPort int, defSec bool) {
	viper.SetDefault(prefix+".host", defHost)
	viper.SetDefault(prefix+".port", defPort)
	viper.SetDefault(prefix+".secure", defSec)

	sc.Host = viper.GetString(prefix + ".host")
	sc.Port = viper.GetInt(prefix + ".port")
	sc.Secure = viper.GetBool(prefix + ".secure")
}
