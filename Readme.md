# Mongotrace

Tool written in GO for printing currently executing MongoDb queries from multiple servers.

Tested on: mongodb v4.0.11

# Build

```shell
make build
make install
```

# Usage examples

```shell
mongotrace enable --config="./confs/config.local.json"
mongotrace disable --config="./confs/config.local.json"
mongotrace tail --config="./confs/config.local.json"
mongotrace tail --config="./confs/config.local.json" --debug
mongotrace tail --config="./confs/config.local.json" --raw
```
