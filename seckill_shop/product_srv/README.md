# Zhiliao_product_srv Service

This is the Zhiliao_product_srv service

Generated with

```
micro new zhiliao_product_srv --namespace=zhiliao.product --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: zhiliao.product.srv.zhiliao_product_srv
- Type: srv
- Alias: zhiliao_product_srv

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
./zhiliao_product_srv-srv
```

Build a docker image
```
make docker
```