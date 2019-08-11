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
      * Npm
        * Path

### Usage

```

Handy CI is a tool for managing and building multi-repository source code on developer host.

Usage:
  handy-ci [Options] COMMAND

Commands:
  exec        Execute any command
  git         Execute Git command
  mvn         Execute Apache Maven command
  npm         Execute npm command

Options:
      --config string       Config file (default is $HOME/.handy-ci/config.yaml)
  -C, --continue            Skip failed command and continue
  -G, --group string        Execute command in group
  -R, --repository string   Execute command in repository
  -W, --workspace string    Execute command in workspace

You can use original options of "git", "mvn", "npm" or any command line tools as additional options.

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
      npms:
      - path: soupe-ida/soupe-ida-ui/src/main/node
      - path: soupe-modern-ui/src/main/node
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
handy-ci mvn clean install -R deployer-kubernetes 
```

#### Use `-C` option can skip previous execution error and continue to next execution

```
handy-ci npm outdated -C
``` 

### Build and Install the Binaries from Source

#### Prerequisite Tools

* [Git](https://git-scm.com/)
* [Go (at least Go 1.11)](https://golang.org/dl/)

#### Setup GO

Set $GOPATH environment variable:
```
$ export GOPATH=$(go env GOPATH)
```

To learn more about the GOPATH environment variable, see 'go help gopath'.

To use a custom workspace location, [set the GOPATH environment variable](https://golang.org/wiki/SettingGOPATH).

Add the workspace's bin subdirectory to your PATH:
```
$ export PATH=$PATH:$GOPATH/bin
```

#### Install from GitHub

```bash
go get github.com/carrchang/handy-ci
cdt  $GOPATH/src/github.com/carrchang/handy-ci
go install
```
