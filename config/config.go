package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// Singleton object
var (
	Config *Configuration
)

// Configuration is the overall config that brightgauge log plugin will use.
type Configuration struct {
	ClientId       string `json:"clientid"`
	ClientSecret   string `json:"client_secret"`
	AppName        string `json:"app_name"`
	RedirectURI    string `json:"redirect_uri"`
	Scopes         string `json:"scopes"`
	SpotifyAPIURL  string `json:"spotify_api_url"`
	SpotifyAuthURL string `json:"spotify_auth_url"`
	AppURL         string `json:"app_url"`
}

// ServiceConfigFile holds info for where config file is stored
type ServiceConfigFile struct {
	Path string `json:"path"`
}

// UnmarshalJSON custom json parser
func (d *ServiceConfigFile) UnmarshalJSON(b []byte) error {
	var aux = struct {
		Path string `json:"path"`
	}{}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}

	d.Path, err = filepath.Abs(aux.Path)
	return err
}

// Load reads and loads configuration
func Load(configFilePath string) error {
	if Config != nil {
		return nil
	}
	contents, err := ioutil.ReadFile(filepath.Clean(configFilePath))
	if err != nil {
		return err
	}

	return json.NewDecoder(bytes.NewBuffer(contents)).Decode(&Config)
}
