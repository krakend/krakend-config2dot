# krakend-config2dot

Translate your KrakenD config file into a dot graph

## Installation

```
$ go install github.com/krakendio/krakend-config2dot/v2/cmd/krakend-config2dot@latest
```

If you have your `$GOPATH/bin` in your `$PATH`, this is how you can create the `.dot` representation of your config file:

```
$ krakend-config2dot -c /path/to/your/config/file.json
```

## Graph generation

Just pipe the output into `dot`.

```
$ krakend-config2dot -c /path/to/your/config/file.json | dot -Tpng -o config.png
```
