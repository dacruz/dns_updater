# DNS Updater
Small tool to handle my home IP changes by updating poiuytre.nl A record on GoDaddy

## Testing
```
$ go test ./...
```

## Coverage report
```
$ go test -coverprofile=cover.out  ./... && go tool cover -html=cover.out
```

## How to run it

### Requisite
To run the application, you first need to decrypt secrets.yaml file. Or write your own and replace the one you've found here...

I'm using the following for encyption:
* age - https://github.com/FiloSottile/age
* SOPS - https://github.com/mozilla/sops

#### Decrypt:
```
$ sops -d --output-type yaml secrets.enc.yaml > secrets.yaml
```

#### Encrypt:
```
$ sops -e --input-type yaml --age ageEXAMPLE secrets.yaml > secrets.enc.yaml 
```

## Deploying to your k8s cluster
```
$ kubectl apply -f secrets.yaml && kubectl apply -f config.yaml 
```

```
$ kubectl apply -f cronjob.yaml
```