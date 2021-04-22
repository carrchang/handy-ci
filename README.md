## Handy CI

Handy CI is a tool for managing and building multi-repository source code on developer host.

### Configuration Concepts

* Workspace
  * Name
  * Root Path
  * Source Group
    * Name
    * Repository
      * Name
      * Remote
        * Name
        * URL
      * Cmd
        * Name
        * Paths

### Usage

```

Handy CI is a tool for managing and building multi-repository source code on developer host.

Usage:

  handy-ci COMMAND [OPTIONS]

Commands:
  exec        Execute any command
  git         Execute Git command

Options:
  -W, --workspace string    Execute command in workspace
  -G, --group string        Execute command in group
  -R, --repository string   Execute command in repository
  -C, --continue            Skip failed command and continue
      --skip string         Skip execution in comma-delimited list of repositories
      --config string       Config file (default is $HOME/.handy-ci/config.yaml)

Options can be in front of, behind, or on both sides of the command.

Original options of any command can be as additional options, and be in behind of the command.

Use "handy-ci COMMAND --help" for more information about a command.

```

### Example Configuration

```
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
      npms:
      - path:
    - name: soupe
      remotes:
      - name: origin
        url: git@gitlab.com:keepnative/soupe.git
      cmds:
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
    - name: data-flow
      remotes:
        - name: origin
          url: git@gitlab.com:keepnative/spring-cloud/data-flow.git
        - name: spring-cloud
          url: git@github.com:spring-cloud/spring-cloud-dataflow.git    
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

#### Fetch fetch all changes from remote git repository in group `spring-cloud`

```
handy-ci git fetch --all -W keepnative -G spring-cloud
```

#### `-W` option can be ignoreg if group `spring-cloud` is unique in all workspaces

```
handy-ci git fetch --all -G spring-cloud
```

#### `-G` option also can be ignore if repository `deployer-kubernetes` is unique in all workspaces

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

### Build and Install the Binaries from Source

#### Prerequisite Tools

* [Git](https://git-scm.com/)
* [Go (at least Go 1.11)](https://golang.org/dl/)

#### Install from GitHub

```
git clone https://github.com/carrchang/handy-ci.git
cd handy-ci
go install
```
