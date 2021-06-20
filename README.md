# Super go mod proxy

Yet another go modules proxy

# About

This project aim at creating a very configurable go module proxy that can work in many environments.

This project can work as: 

1. A proxy for the standard goproxy online service
2. A caching proxy for gomodules
3. A single entrypoint for all go modules.    
4. A configurable proxy (rewrite module name, force VCS, add authentication, ...)

## Getting Started

Binaries for multiple platforms are available as download in the releases page. 
Start the binary with a config file then, enjoy ! 

## Protocol defintion

This go module follows the goproxy protocol described here: https://golang.org/ref/mod#goproxy-protocol

## Configuration

The configuration file is a Json file. It can be given to the executable with the `-config` flag. 

### General 

General options for the proxy:

* `defaultRelayProxy` : set the default upstream go module proxy

### Phases

To allow deep configuration and override of the proxy behaviour, phases represents the main
steps done to collect a gomodule.

* `Receive` : when the request is received
* `Prefetch`: just before fetching the content
* `Fetch`: when downloading / getting the content

In each phase one or more plugins can be used, given they do something for that phase.

### Plugins

Note: this is not a real plugin system as you have to declare/register the "plugin" in the source code so that it can be used.

Using the term plugin is a way to say that the phases execution and result can change based on the configuration. 

#### Default

A Default plugin that implement all phases and redirect all requests to upstream goproxy. 

#### Private plugin

A plugin to specify a private plugin. It changes the remote url from default go proxy to 
a git url to your repository. Authentication can be defined in the plugin config

#### Rewrite plugin 

A plugin to rename / rewrite the module being targeted to another one.
This will impact the fetched repository as only the replacement module name will be kept.

#### VCS plugin

A plugin to arbitrary define a custom VCS for a given module or group of module. 
This might be usefull in private environnement where the domain used is not where the modules are
hosted. For example, module name would be something like `myCompany.com/myGroup/module` and the repository is located at 
`http://git.private.local/myGroup/module`. This is then transparent for the go client. 

 