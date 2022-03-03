# krakend-config2dot
Transalte your KrakenD config file into a dot graph

## Installation

```
$ go install github.com/devopsfaith/krakend-config2dot/cmd/krakend-config2dot
```

If you have your `$GOPATH/bin` in your `$PATH`, this is how you can create the `.dot` representation of your config file:

```
$ krakend-config2dot -c /path/to/your/config/file.json
```

## Graph generation

Just pipe the output into `dot`

```
$ krakend-config2dot -c /path/to/your/config/file.json | dot -Tpng -o config.png
```
