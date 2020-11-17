protodep
=======

[![Circle CI](https://circleci.com/gh/stormcat24/protodep.svg?style=shield&circle-token=f53432c65ac4fd4bd4b8b778892690e4032ea141)](https://circleci.com/gh/stormcat24/protodep)
[![Language](https://img.shields.io/badge/language-go-brightgreen.svg?style=flat)](https://golang.org/)
[![issues](https://img.shields.io/github/issues/stormcat24/protodep.svg?style=flat)](https://github.com/stormcat24/protodep/issues?state=open)
[![License: MIT](https://img.shields.io/badge/license-MIT-orange.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/stormcat24/protodep?status.png)](https://godoc.org/github.com/stormcat24/protodep)

Dependency tool for Protocol Buffers IDL file (.proto) vendoring tool.


## Motivation

In building Microservices architecture, gRPC with Protocol Buffers is effective. When using gRPC, your application will depend on many remote services.

If you manage proto files in a git repository, what will you do? Most remote services are managed by git and they will be versioned. We need to control which dependency service version that application uses.


## Install

### go get

```bash
$ go get github.com/stormcat24/protodep
```

### from binary

Support as follows:

* protodep_darwin_amd64.tar.gz
* protodep_linux_386.tar.gz
* protodep_linux_amd64.tar.gz
* protodep_linux_arm.tar.gz
* protodep_linux_arm64.tar.gz

```bash
$ wget https://github.com/stormcat24/protodep/releases/download/0.0.8/protodep_darwin_amd64.tar.gz
$ tar -xf protodep_darwin_amd64.tar.gz
$ mv protodep /usr/local/bin/
```

## Usage

### protodep.toml

Proto dependency management is defined in `protodep.toml`.

```Ruby
proto_outdir = "./proto"

[[dependencies]]
  target = "github.com/stormcat24/protodep/protobuf"
  branch = "master"

[[dependencies]]
  target = "github.com/grpc-ecosystem/grpc-gateway/examples/examplepb"
  revision = "v1.2.2"
  path = "grpc-gateway/examplepb"

[[dependencies]]
  target = "github.com/kubernetes/helm/_proto/hapi"
  branch = "master"
  path = "helm/hapi"
  ignores = ["./release", "./rudder", "./services", "./version"]
```

### protodep up

In same directory, execute this command.

```bash
$ protodep up
```

If succeeded, `protodep.lock` is generated.

### protodep up -f (force update)

Even if protodep.lock exists, you can force update dependenies.

```bash
$ protodep up -f
```

### [Attention] Changes from 0.1.0

From protodep 0.1.0 supports ssh-agent, and this is the default.
In other words, in order to operate protodep without options as before, it is necessary to set with ssh-add.

As the follows:

```bash
$ ssh-add ~/.ssh/id_rsa
$ protodep up
```

### Getting via HTTPS

If you want to get it via HTTPS, do as follows.

```bash
$ protodep up --use-https
```

And also, if Basic authentication is required, do as follows.
If you have 2FA enabled, specify the Personal Access Token as the password. 

```bash
$ protodep up --use-https \
    --basic-auth-username=your-github-username \
    --basic-auth-password=your-github-password
```

License
===
See [LICENSE](LICENSE).

Copyright © stromcat24. All Rights Reserved.
