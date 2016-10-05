# tor 
This is tor proxy controller package.

## what is Tor?
* [Wikipedia - Tor](https://en.wikipedia.org/wiki/Tor)

## install tor_multi
```
$ go install github.com/goshinobi/tor_multi/cmd/tor_multi
```

## install tor
### OS X
```
$ brew install tor
```
## run tor_multi
```
$ go run ${TOR_MULTI_REPO_DIR}/cmd/tor_multi/main.go
```

or

```
$ tor_multi
```

## API
### list proxy
```
/list
```

### add proxy
```
/add
```

### kill all proxy
```
/killAll
```
