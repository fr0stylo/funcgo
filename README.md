# FuncGO
 
## Description

This is experimental repository, implementing on demand functions based on unix rootless container. Fully implemented in GO, works only on unix based systems as it is using unix syscalls and unix namespaces

## Running containers 

1. In order to prep environment run `make build`, this will require `root` as this is required in order to manage network interfaces.
1. run `./bin/controller` 

## Unet
 Unet is app from https://github.com/LK4D4/unc/blob/master/unet/main.go which is used to create and maintains links between container namespaces. This applications requires elevated access to run, that's why make uses `sudo`

## State

- [x] Container from scratch
- [x] Container networking with IP address provisioning from manager 
- [x] Function process creation on demand with isolated container
- [x] Function pooling with downscale
- [x] Basic function handler implementation

## Next steps

- [ ] API GW implementation for Function invocation
- [ ] Function routing - ability to manage more functions per host than 1
- [ ] Function deployment over several hosts
- [ ] Create HA deployment 
- [ ] API for function management 
