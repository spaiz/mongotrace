# Mongotrace

Tool written in GO for printing currently executing MongoDb queries from multiple servers.

Tested on: mongodb v4.0.11

# How it looks

![Output of raw and formatted query log](docs/images/example1.png)

![Output of raw and formatted query log](docs/images/example2.png)

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
