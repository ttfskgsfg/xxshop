# Zhiliao_seckill_srv Service

This is the Zhiliao_seckill_srv service

Generated with

```
micro new zhiliao_seckill_srv --namespace=zhiliao.seckill --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: zhiliao.seckill.srv.zhiliao_seckill_srv
- Type: srv
- Alias: zhiliao_seckill_srv

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
./zhiliao_seckill_srv-srv
```

Build a docker image
```
make docker
```