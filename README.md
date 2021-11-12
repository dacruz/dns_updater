# dns_updater
Small tool to handle my home IP changes by updating poiuytre.nl A record on GoDaddy

## Testing
```
$ go test ./...
```

## Coverage report
```
$ go test -coverprofile=cover.out  ./... && go tool cover -html=cover.out
```

```
$ go tool cover -html=cover.out
```