## Handy CI

Handy CI is a tool for managing and building multi-repository source code on developer host.

### Configuration Concepts

```go
package config

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
```

### Usage

```

Handy CI is a tool for managing and building multi-repository source code on developer host.

Usage:

  handy-ci COMMAND [OPTIONS]

Commands:
  exec        Execute any command
  git         Execute Git command

Options:
  -W, --workspace string      Execute command in workspace
  -G, --group string          Execute command in group
  -R, --repositories string   Execute command in comma-delimited list of repositories
  -C, --continue              Skip failed command and continue
      --dry-run               Only print the command and execution path
      --skip string           Skip execution in comma-delimited list of repositories
  -F, --from string           Execute command from repository to end
      --config string         Config file (default is /Users/carrchang/.handy-ci/config.yaml)

Options can be in front of, behind, or on both sides of the command.

Original options of any command can be as additional options, and be in behind of the command.

Use "handy-ci COMMAND --help" for more information about a command.

```

### Example Configuration

```
scriptDefinitions:
  - name: mvn
    defaultArgs: clean install -nsu
  - name: npm
    defaultArgs: outdated
workspaces:
  - name: home
    root: /Users
    groups:
      - name: carrchang
        repositories:
          - name: .handy-ci
            remotes:
              - name: origin
              - url: git@github.com:carrchang/handy-ci-config.git
  - name: carrchang-go
    root: /coding/go/src/github.com
    groups:
      - name: carrchang
        repositories:
          - name: handy-ci
            remotes:
              - name: origin
              - url: git@github.com:carrchang/handy-ci.git
  - name: keepnative
    root: /coding/keepnative
    groups:
      - name: next
        repositories:
          - name: java
            remotes:
              - name: origin
                url: git@gitlab.com:keepnative/java.git
          - name: soupe-ui-components
            remotes:
              - name: origin
                url: git@gitlab.com:keepnative/soupe-ui-components.git
            scripts:
              - name: npm
          - name: soupe
            remotes:
              - name: origin
                url: git@gitlab.com:keepnative/soupe.git
            scripts:
              - name: mvn
                default: true
              - name: npm
                paths:
                  - soupe-ida/soupe-ida-ui/src/main/node
                  - soupe-modern-ui/src/main/node
      - name: spring-cloud
        repositories:
          - name: deployer-kubernetes
            remotes:
              - name: origin
                url: git@gitlab.com:keepnative/spring-cloud/deployer-kubernetes.git
              - name: spring-cloud
                url: git@github.com:spring-cloud/spring-cloud-deployer-kubernetes.git
            scripts:
              - name: mvn
          - name: data-flow
            remotes:
              - name: origin
                url: git@gitlab.com:keepnative/spring-cloud/data-flow.git
              - name: spring-cloud
                url: git@github.com:spring-cloud/spring-cloud-dataflow.git
            scripts:
              - name: mvn
```

### Examples

#### Get git repository status in all workspace

```
handy-ci git status
``` 

#### Get git repository status in workspace `keepnative`

```
handy-ci git status -W keepnative
```

#### Fetch all changes from remote git repository in group `spring-cloud`

```
handy-ci git fetch --all -W keepnative -G spring-cloud
```

#### `-W` option not needed if group `spring-cloud` is unique in all workspaces

```
handy-ci git fetch --all -G spring-cloud
```

#### `-G` option also  not needed if repository `deployer-kubernetes` is unique in all workspaces

```
handy-ci exec mvn clean install -R deployer-kubernetes 
```

#### Use `-C` option can skip previous execution error and continue to next execution

```
handy-ci exec npm outdated -C
```

#### Use `--skip` option can skip execution in repository `deployer-kubernetes`

```
handy-ci exec mvn clean install -G spring-cloud --skip deployer-kubernetes
```

#### Execute default script, first script will be executed when default not specified

```
handy-ci exec
```

### Build and Install the Binaries from Source

#### Prerequisite Tools

* [Git](https://git-scm.com/)
* [Go (at least Go 1.18)](https://golang.org/dl/)

#### Install from GitHub

```
git clone https://github.com/carrchang/handy-ci.git
cd handy-ci
go install
```
