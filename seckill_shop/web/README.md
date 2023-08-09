# Zhiliao_web Service

This is the Zhiliao_web service

Generated with

```
micro new zhiliao_web --namespace=zhiliao.web --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: zhiliao.web.web.zhiliao_web
- Type: web
- Alias: zhiliao_web

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
./zhiliao_web-web
```

Build a docker image
```
make docker
```