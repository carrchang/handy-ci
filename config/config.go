package config

import (
  "github.com/spf13/viper"

  "github.com/carrchang/handy-ci/util"
)

var HandyCiConfig *Config

type Config struct {
  ScriptDefinitions []ScriptDefinition `yaml:"scriptDefinitions"`
  Workspaces        []Workspace        `yaml:"workspaces"`
}

type ScriptDefinition struct {
  Name        string `yaml:"name"`
  DefaultArgs string `yaml:"defaultArgs"`
}

type Workspace struct {
  Name   string  `yaml:"name"`
  Path   string  `yaml:"path"`
  Groups []Group `yaml:"groups"`
}

type Group struct {
  Name              string       `yaml:"name"`
  NameIgnoredInPath bool         `yaml:"nameIgnoredInPath"`
  Path              string       `yaml:"path"`
  Repositories      []Repository `yaml:"repositories"`
}

type Repository struct {
  Name              string      `yaml:"name"`
  NameIgnoredInPath bool        `yaml:"nameIgnoredInPath"`
  Path              string      `yaml:"path"`
  Remotes           []GitRemote `yaml:"remotes"`
  Scripts           []Script    `yaml:"scripts"`
}

type GitRemote struct {
  Name string `yaml:"name"`
  URL  string `yaml:"url"`
}

type Script struct {
  Name    string   `yaml:"name"`
  Default bool     `yaml:"default"`
  Paths   []string `yaml:"paths"`
}

func Initialize() {
  err := viper.Unmarshal(&HandyCiConfig)
  if err != nil {
    util.Printf("Unable to decode into config struct, %v", err)
  }
}
