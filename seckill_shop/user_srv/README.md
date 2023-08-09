# Zhiliao_user_srv Service

This is the Zhiliao_user_srv service

Generated with

```
micro new zhiliao_user_srv --namespace=zhiliao.user --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: zhiliao.user.srv.zhiliao_user_srv
- Type: srv
- Alias: zhiliao_user_srv

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
./zhiliao_user_srv-srv
```

Build a docker image
```
make docker
```