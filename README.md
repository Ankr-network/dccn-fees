# Dcmgr Service

This is the Dcmgr service

Generated with

```
micro new github.com/Ankr-network/refactor/app_dccn_dcmgr --namespace=network.ankr --alias=dcmgr --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: network.ankr.srv.v1
- Type: srv
- Alias: v1

## Dependencies

Micro services depend on service discovery. The default is consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./dcmgr-srv
```

Build a docker image
```
make docker
```

## Test
### Postman
    $http://192.168.0.102:8080/rpc
```
{
    "service": "go.micro.srv.v1.dcmgr",
	"method": "DcMgr.Create",
	"request": {
	    "name": "dc0",
	    "id": 2,
	    "status": 1
	}
}
```

### curl
```
curl -d 'service=go.micro.srv.v1.dcmgr' \
	 -d 'method=DcMgr.Create' \
	 -d 'request={"id": 1, "name": "dc01", "status": 1}' \
	 http://localhost:8080/rpc
	 ```
