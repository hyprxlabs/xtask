# xtask

A cross platform task runner.

```yaml
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

```bash
xtask build
xtask ssh
```

