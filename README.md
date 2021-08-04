# pulumi production application

This repo is a Pulumi multi language component that creates an example "Production application".

The intent is to show how you can define a Pulumi component using Pulumi's multi language capabilities, and then consume that component using Pulumi's automation API

## SDK

The production app component can be consumed in all of Pulumi's supported languages. You can see examples in [examples](examples).

These SDK components are designed to be used with the Pulumi CLI, driving a familiar infrastructure as code workflow

```
Updating (dev)

View Live: https://app.pulumi.com/jaxxstorm/productionapp-node-example/dev/updates/29

     Type                                 Name                            Status
 +   pulumi:pulumi:Stack                  productionapp-node-example-dev  created
 +   ├─ productionapp:index:Deployment    example                         created
 +   └─ kubernetes:core/v1:Namespace      example                         created
 +      ├─ kubernetes:core/v1:Service     example                         created
 +      └─ kubernetes:apps/v1:Deployment  example                         created
```

## CLI

A custom CLI can also be built using Pulumi's Automation API. The automation API allows you to build custom tooling and workflows.

The CLI is based on amazing work by [@komalali](https://github.com/komalali)

The CLI lives in the [cli](cli) directory and is written using Go.

### Building the CLI

Simply use `go build` to build the CLI for your operating system:

```
go build -o prodapp main.go
```

### Using the CLI

Run the CLI tool and provide a name for your deployment. If a name is not specified, one will be generated for you:

```
./prodapp --help
usage: productionapp [<flags>] <command> [<args> ...]
```

```
./prodapp deploy --image="nginx:latest" --name "cli-example"

⣷ Current step: Running update...

  Updates in progress                                │  Updates completed                                 │
  ───────────────────                                │  ─────────────────                                 │
                                                     │  ✓ kubernetes:apps/v1:Deployment                   │
                                                     │  ✓ kubernetes:core/v1:Namespace                    │
                                                     │  ✓ kubernetes:core/v1:Service                      │
                                                     │  ✓ productionapp:index:Deployment                  │
                                                     │  ✓ pulumi:pulumi:Stack                             │
                                                     │                                                    │
```                                                     


## Web Platform

A web platform built using Python and flask is also available. Again, most of the UI and framework was built by [@komalali](https://github.com/komalali)

It lives in the [web](web) directory. 

### Installing

From withing the web directory..

Create a venv:

```
python -m venv venv
```

Install all the dependencies:

```
venv/bin/pip3 install -r requirements.txt
```

Run the web application:

```
venv/bin/python3 __main__.py
 * Serving Flask app '__main__' (lazy loading)
 * Environment: production
   WARNING: This is a development server. Do not use it in a production deployment.
   Use a production WSGI server instead.
 * Debug mode: off
 * Running on http://127.0.0.1:5000/ (Press CTRL+C to quit)
```

You should now be able to browse to http://localhost:5000 and see your web application platform and deploy to your Kubernetes cluster.

## Building

Building the SDK is done using [go-task](https://github.com/go-task/task)


Install it:

```
brew install go-task/tap/go-task
```

Build the provider and the SDKs
```
task build build:provider generate:sdks build:sdks install:nodejs install:python
```

Add the provider to your `PATH`:

```
export PATH=$PATH:$(pwd)/bin
```
