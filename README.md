# FuncGO
 
## Description

This is experimental repository, implementing on demand functions based on unix rootless container.

## Running containers 

1. In order to prep environment run `make build`, this will require `root` as this is required in order to manage network interfaces.
1. run `./bin/controller`

## State

1. [*] Container from scratch
1. [*] Container networking with IP address provisioning from manager 
1. [*] Function process creation on demand with isolated container
1. [*] Function pooling with downscale
1. [*] Basic function handler implementation

## Next steps

1. API GW implementation for Function invocation
1. Function routing - ability to manage more functions per host than 1
1. Function deployment over several hosts
1. Create HA deployment 
1. API for function management 
