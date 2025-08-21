# xtask

A cross platform task runner.

## Sample Yaml

Schema for xtaskfile is available at [jsonschema/xtask.schema.json](jsonschema/xtask.schema.json).

It is also availble using <https://raw.githubusercontent.com/hyprxlabs/xtask/refs/heads/main/jsonschema/xtask.schema.json>.

```yaml

config:
  # sets the default shell to use for tasks
  shell: "bash"
  # enables command substitution in env variables
  # e.g $(az keyvault secret show --name mysecret --vault-name myvault --query value -o tsv)
  # environment expansion is enabled by default e.g. ${HOME}
  command-substitution: true
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
