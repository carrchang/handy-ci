package config

import (
  "github.com/spf13/viper"

  "github.com/carrchang/handy-ci/util"
)

var HandyCiConfig *Config

type Config struct {
  Workspaces []Workspace `yaml:"workspaces"`
}

type Workspace struct {
  Name   string  `yaml:"name"`
  Root   string  `yaml:"root"`
  Groups []Group `yaml:"groups"`
}

type Group struct {
  Name         string       `yaml:"name"`
  Repositories []Repository `yaml:"repositories"`
}

type Repository struct {
  Name    string      `yaml:"name"`
  Remotes []GitRemote `yaml:"remotes"`
  Npms    []Npm       `yaml:"npms"`
}

type GitRemote struct {
  Name string `yaml:"name"`
  URL  string `yaml:"url"`
}

type Npm struct {
  Name string `yaml:"name"`
  Path string `yaml:"path"`
}

func Initialize() {
  err := viper.Unmarshal(&HandyCiConfig)
  if err != nil {
    util.Printf("Unable to decode into config struct, %v", err)
  }
}
