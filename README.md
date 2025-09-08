# xtask

A cross platform task runner.

## Sample Yaml

Schema for xtaskfile is available at [assets/schema/xtaskfile.schema.json](jsonschema/xtaskfile.schema.json).

It is also availble using <https://raw.githubusercontent.com/hyprxlabs/xtask/refs/heads/main/assets/schema/xtaskfile.schema.json>.

```yaml

config:
  # sets the default shell to use for tasks
  shell: "bash"
  # enables command substitution in env variables
  # e.g $(az keyvault secret show --name mysecret --vault-name myvault --query value -o tsv)
  # environment expansion is enabled by default e.g. ${HOME}
  substitution: true
  # loads before everything else
  env:
    CUSTOM_VAR: "Hello World"
    RELATIVE_VAR: "${XTASK_DIR}/relative"

  # prepend paths to the PATH environment variable.
  # these paths can be absolute or relative and may
  # contain platform specific paths. e.g. windows, linux, darwin
  prepend-paths:
    - windows: "C:\\Program Files\\Git\\usr\\bin"
    - ./my/bin

# loads after config.env
dotenv:
  - ${RELATIVE_VAR}/.env

# loads after config.env and dotenv
# used to override dotenv variables
env:
  ONE: "one"
  TWO: "${TWO:-two}"
  SECRET: "$(az keyvault secret show --name mysecret --vault-name myvault --query value -o tsv)"

tasks:
  build:
    run: |
       go build -o bin/xtask ./cmd/xtask
       echo "${ONE} ${TWO}"
    uses: bash
    
  ssh:
    run: |
        uptime
        echo "Hello from ${HOSTNAME}!"
    uses: ssh://user@host

  echo: 'echo "Hello from ${CUSTOM_VAR} and ${SECRET}"'
```

## Sample Usage

```bash
xtask build
xtask ssh
xtask ls # lists all task
```

## CLI

```bash
xtask [command] [options]
```

`run` will execute a task defined in the xtaskfile. You can pass additional
arguments to the task using `--` followed by the arguments or using
any flag that starts with `-`.

```bash
xtask run [options] [...task] [additional args]
xtask run my-task
xtask run my-task --arg1 value1 --arg2 value2
xtask run my-task -- /arg1 value1 /arg2 value2
```

`ls` will list all tasks defined in the xtaskfile.

```bash
xtask ls [options]
xtask ls --match bu*"
```

`exec` will execute a command using the current environment variables and PATH
from the xtaskfile.

```bash
xtask exec [options] [command] [args...]
xtask exec bash -c 'echo "Hello from ${CUSTOM_VAR}"'
```

## xtaskfile YAML Format

## Config Section

```yaml
config:
  shell: "bash" # default shell to use for tasks
  substitution: true # enable command substitution in env variables
  env: # environment variables to set before for all tasks and other dotenv files
    CUSTOM_VAR: "Hello World"
  prepend-paths: # paths to prepend to the PATH environment variable
    - windows: "C:\\Program Files\\Git\\usr\\bin" # only windows
    - linux: "/usr/local/git/bin" # only linux
```

## Dotenv Section

Specify a list of dotenv files to load.  These files are loaded after the
`config.env` section and before the `env` section of the xtaskfile.

Optional files are supported by prefixing the file with a `?` at the
end of the path.

```yaml
dotenv:
  - ./path/to/.env # relative to the xtaskfile directory
  - ${XDG_CONFIG_HOME}/.env? # using environment variable
```

## Env Section

Specify environment variables tofor all running tasks.

```yaml
env:
  ONE: "one"
  TWO: "${TWO:-two}" # use existing value or default to "two"
  SECRET: "$(az keyvault secret show --name mysecret --vault-name myvault --query value -o tsv)" # command substitution
```

## Hosts Section

```yaml
hosts:
  - path/to/import
  - host: "host1"
    user: "user1"
    identity: "id_rsa"
    password: "password"
    groups:
      - "group1"
      - "group2"
    meta:
      key: "value"
    os:
      platform: "linux"
      version: "ubuntu-20.04"
      arch: "amd64"
```

```yaml
hosts:
  vm1:
    host: "192.168.1.3"
    user: "test"
    identity: "~/.ssh/id_rsa"
  vm2:
    host: "192.168.1.4"
    user: "test"
    identity: "~/.ssh/id_rsa"
```

## Tasks Section

The tasks section defines the tasks that can be run.  Each task has a name and
a set of fields that define how the task is run.

```yaml
tasks:
  build:
    run: |
       go build -o bin/xtask ./cmd/xtask
       echo "${ONE} ${TWO}"
    uses: bash
```

### Task Fields

```yaml
tasks:
  my-task:
    name: "My Task" # optional name for the task
    desc: "My Task Description" # optional description for the task
    needs: ["other-task"] # list of tasks that must be run before this task
    run: | # the command or script to run for the task
      echo "Hello from my-task"
    uses: bash # the type of task to run (see Builtin Task Types below)
    timeout: "1h30m" # optional timeout for the task
    cwd: "./path/to/dir" # optional working directory for the task
    # optional list of hosts for ssh/scp tasks. this may be the name of the host or 
    # the name of a group.  If it is a group, then all members of the group will be used.
    # this is only valid for ssh/scp tasks.
    hosts: ["host1", "host2"] 
    env: # environment variables to set for the task
      TASK_VAR: "task value"
    dotenv: # dotenv files to load for the task
      - ./path/to/.env
```

#### Sample SCP Task

files are in a list of source:destination pairs.

```yaml
tasks:
  copy-files:
    uses: scp://user@host
    with:
      files:
        - "./source/file1.txt:/opt/dest/file1.txt"
```

Using multiple hosts.

```yaml
hosts:
  vm1:
    host: "192.168.1.3"
    user: "test"
    identity: "~/.ssh/id_rsa"
  vm2:
    host: "192.168.1.4"
    user: "test"
    identity: "~/.ssh/id_rsa"

tasks:
  copy-files:
    uses: scp
    hosts: ["vm1", "vm2"]
    with:
      files:
        - "./source/file1.txt:/opt/dest/file1.txt"
        - "./source/file2.txt:/opt/dest/file2.txt"
```

#### Sample SSH Task

```yaml
tasks:
  remote-cmd:
    uses: ssh://user@host
    run: |
      echo "Hello from remote host"
      uptime
```

## xhostfile YAML Format

All password are mapped to environment variables. So you must never store
the actuall password in the xtaskfile or xhostfile.  Instead use environment
variables to store the passwords using something like command substitution.

If the password is password: MY_PASSWORD, it will look for the environment variable
MY_PASSWORD and use that value as the password.

If the identity value is set, the password is used to decrypt the identity file if it is
encrypted.  Otherwise, the password is used to authenticate the user.

### default section

The default section is used to define default values for all hosts.

```yaml

default:
  user: "default_user"
  identity: "~/.ssh/id_rsa"
  password: "MY_PASSWORD"
  port: 22
  groups: ["group1", "group2"]
  meta:
    key: "value"
    key2: 10
    key2:
     - nested: "value"
  os:
    platform: "linux"
    version: "20.04"
    arch: "amd64"
    family: "debian"
    codename: "focal"
    variant: "ubuntu"
```

### defaults section

If you need multiple defaults, you can use the `defaults` section.

```yaml
defaults:
  key1:
    user: "default_user"
    identity: "~/.ssh/id_rsa"
    password: "MY_PASSWORD"
    port: 22
    groups: ["group1", "group2"]
  key2:
    user: "other_user"
    identity: "~/.ssh/id_rsa_other"
    password: "OTHER_PASSWORD"
    port: 2222
    groups: ["group3", "group4"]
```

### Hosts section

```yaml
hosts:
   key1:
     host: "10.0.0.1"
     user: "user1"
     identity: "~/.ssh/id_rsa"
     password: "MY_PASSWORD"
   key2:
     host: "10.0.0.2"
     user: "user2"
     identity: "~/.ssh/id_rsa_other"
     password: "OTHER_PASSWORD"
```

### Host Fields

```yaml
hosts:
  host1:
    host: "10.0.0.1" # cname or ip address of the host
    user: "user1" # username to use for ssh
    identity: "~/.ssh/id_rsa" # path to the private key file
    password: "MY_PASSWORD" # password or passphrase for the private key file
    port: 22 # port to use for ssh
    groups: # groups that the host belongs to
      - "group1"
      - "group2"
    meta: # arbitrary metadata for the host
      key: "value"
      key2: 10
    os:
      platform: "linux" # platform of the host
      version: "20.04" # version of the host
      arch: "amd64" # architecture of the host
      family: "debian" # family of the host
      codename: "focal" # codename of the host
      variant: "ubuntu" # variant of the host
```

## Builtin Task Types

- **ssh** or `ssh://user@host` - Execute the task on a remote host using SSH.
- **scp** or `scp://user@host` - Copy files to a remote host using SCP.
- **shell** - The shell task uses a given shell or script interpreter to execute the task. Supported
  shells are:
  - `bash`
  - `sh`
  - `powershell`
  - `pwsh`
  - `python`
  - `ruby`
  - `deno`
  - `node`
  - `bun`
- **tmpl** - A template task that renders a template file using Go templates and environment
  variable substitution. Either Environment substitution and go templates may be disabled
- **import** - Loads a single task using the file:// or relative path syntax so long as the file
  ends with .xtask.yaml or .xtask.yml. This allows you to split your tasks into multiple files and
  reuse them.

## Environment Variables

Process environment variables are not modified directly. Instead a workflow gets a copy of the
process environment variables and then any changes are applied to that copy.  This allows
xtask to run to possibly run parallel tasks in the future without affecting each other and
simulates what CI/CD/pipelines typically when running tasks as a separate process.

Tasks may also modify the environment variables for just that task.  A task an modify
subsquent tasks by writing file provided by the `XTASK_ENV` environment variable or
prepend paths writing to the file provided by the `XTASK_PATH` environment variable.

For dotenv files and any `env` section of the xtaskfile, the environment variables
are loaded in the order they are defined.  This allows you to override environment
variables by defining them later or using concatenation using previously defined
environment variables within the same section or file.

Environment variables are loaded in the following order.

- Process environment variables.
- `config:env` section of the current xtaskfile.
- dotenv files found in the `XTASK_CONFIG_HOME` directory.
- dotenv files found in the `XTASK_ETC_DIR` directory.
- dotenv files found in the `XTASK_DIR` directory.
- dotnet files specified in the `dotenv` section of the xtaskfile.
- `env` section of the xtaskfile.
- `--env-file` or `-E` specified on the command line.
- `--env` or `-e` specified on the command line.
- dotenv files specified in the `task:<name>:dotenv` section for the current
  task being executed.
- `task:<name>:env` section of the xtaskfile for the current task being executed.

For dotenv files, the autoload directories are from the following locations.

- `$XTASK_CONFIG_HOME` - Default is `$XDG_CONFIG_HOME/xtask` or `$HOME/.config/xtask`.
- `$XTASK_ETC_DIR` - Default is `./.xtask/etc` if the XTASK_ETC_DIR environment variable is not set.
- `$XTASK_DIR` - The directory where the xtaskfile is located.

The .env file is consider to be the shared environment file and is loaded first. If the current
`$XTASK_CONTEXT` variable is empty, then `default` is used. If the a file with the context
name exists, then it is loaded after the .env file. For example, if the current context is 'production',
then the following .env.production file is loaded after the .env file for each of the autoload directories.

The environment variables are expanded using the `${VAR}`, `$VAR` and `${VAR:-default}` syntax.

Command substitution which uses the `$(...)` syntax, has limited support. Command
substitution only supports a single command and is not recursive. It is primarily
intended to be used to load secrets from a secret manager like Azure Key Vault,
AWS Secrets Manager, HashiCorp Vault, sops, using command line apps.

### SSH and Environment Variables

When using the `ssh://user@host` syntax for the `uses` field of a task, only the **environment variables defined
on the task will be passed to the remote host**.  The remote must have the environment variables enabled in
the sshd_config file using the `AcceptEnv` directive.  For example, to accept all environment variables,

```text
AcceptEnv XTASK_*
```

### XTASK Environment Variables

- `XTASK_DIR` - The directory where the xtaskfile is located.
- `XTASK_FILE` - The xtaskfile name. Default is `xtaskfile`.
- `XTASK_CONTEXT` - The current context name. Default is 'default'.
- `XTASK_SHELL` - The default shell to use for tasks. Default is 'powershell' on Windows and 'bash' on other platforms.
- `XTASK_DOT_DIR` - The base directory for xtask configuration. Default is `./.xtask`.
- `XTASK_ETC_DIR` - The etc directory. Default is `./.xtask/etc`.
- `XTASK_APPS_DIRS` - The apps directories. Default is `./.xtask/apps`.
- `XTASK_SCRIPTS_DIR` - The scripts directory. Default is `./.xtask/scripts`.
- `XTASK_BIN_DIR` - The bin directory. Default is `./.xtask/bin`. This is used to load scripts or binaries.
- `XTASK_CONFIG_HOME` - The config home directory. Default is `$XDG_CONFIG_HOME/xtask` or `$HOME/.config/xtask`.
- `XTASK_DATA_HOME` - The data home directory. Default is `$XDG_DATA_HOME/xtask` or `$HOME/.local/share/xtask`.
- `XTASK_CACHE_HOME` - The cache home directory. Default is `$XDG_CACHE_HOME/xtask` or `$HOME/.cache/xtask`.
- `XTASK_STATE_HOME` - The state home directory. Default is `$XDG_STATE_HOME/xtask` or `$HOME/.local/state/xtask`.
- `XTASK_RUNTIME_DIR` - The runtime directory. Default is `$XDG_RUNTIME_DIR/xtask` or `/run/user/$UID/xtask`.

## Config

The `config` section of the xtaskfile is used to configure the behavior of xtask.

- `shell` - The default shell to use for tasks. Default is 'powershell' on Windows and 'bash' on other platforms.
- `substitution` - Enable command substitution in environment variables. Default is true. This allows you to use
  simple command substitution using `$(...)` syntax in environment variables. Only a single command is supported.
- `env` - Environment variables to set before for all tasks and other dotenv files.
- `prepend-paths` - Paths to prepend to the PATH environment variable. These paths can be absolute or relative and may
  contain platform specific paths. e.g. windows, linux, darwin.

## Prepend Paths

The `config.prepend-paths` section allows you to specify paths to be prepended to the PATH environment variable
To set os platform specific paths, you can use the following keys:

```yaml
config:
  prepend-paths:
    - windows: "C:\\Program Files\\Git\\usr\\bin" # only windows
    - linux: "/usr/local/git/bin" # only linux
    - darwin: "/usr/local/git/bin" # only darwin (macOS)
    - ./my/bin # all platforms
```

Environment variables in the paths are also expanded. So if you need to use common windows 
environment variables like `%ProgramFiles%`, you can do so as shown below:

```yaml
config:
  prepend-paths:
    - windows: "%ProgramFiles%\\Git\\usr\\bin"
```

You can also dynamically set prepend paths when using multiple tasks by using the `$XTASK_PATH` environment
variable.  The XTASK_PATH represents a path to a file that contains a list of paths that should be prepended to the PATH
environment variable.

```yaml

tasks:
  setup:
    run: |
      echo "%ProgramFiles%\\Git\\usr\\bin" > $XTASK_PATH
      echo "./my/bin" >> $XTASK_PATH
    uses: bash

  build:
    run: |
      git --version
      echo $PATH
    uses: bash
```

### Auto Prepend PATH

The following directories are automatically prepended to the PATH environment variable in the given order.

1. `$HOME/.local/bin` on non-windows platforms when SUDO_USER is set.
2. `./node_modules/.bin` if it exists.
3. `./bin` if it exists.
4. `${XTASK_BIN_DIR}` if it exists.
5. Any directories in the `config.prepend-paths` section of the xtaskfile.

### Context Switching

The context is a way to switch between different cotexts when calling tasks. The context may be set
using the --context or -c flag on the command line or by setting the `XTASK_CONTEXT` environment variable.

When autoloading dotenv files, the context is used to load context specific dotenv files.  For example, if the current context is 'production', then the following .env.production file is loaded after the .env file for each of the autoload directories.

The context can be used by tasks to branch logic.  

```yaml

tasks:
  build:
    run: |
      dotnet build -c "${XTASK_CONTEXT}"
    uses: bash
```

```bash
xtask build -c "Release"
xtask run -c "Release" deploy
```

For runing lifecyle tasks like `build`, `deploy`, `test`, `install` etc, you can define use
the context to branch tasks

```yaml

tasks:
  build:before:
    run: |
      echo "Running pre-build steps"

  build:
    run:  |
      # execute large complex build logic here

  build:ci:before:
    run: |
      echo "Running pre-build steps for ci/cd pipelines"

  build:ci:
    run: |
      # execute large altnerative build logic here for ci/cd pipelines
```

```bash
xtask build -c "ci"
xtask lc -c "ci" build
```
