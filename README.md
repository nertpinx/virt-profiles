# virt-profiles

This is an implementation in golang of the [virt profiles/virtuned concept](https://github.com/nertpinx/virt-manager/pull/1)

[![Documentation](https://godoc.org/github.com/fromanirh/virt-profiles/pkg/profiler?status.svg)](http://godoc.org/github.com/fromanirh/virt-profiles/pkg/profiler)

## license
Apache v2

## content
The components of this project are:

```
.
├── cmd/virtprofilesd             - Exporting functions as REST APIs
├── cmd/tools/                    - Command line tools, see README.md here
|            └── virtprofilectl   - Example client/debug tool for virtprofilesd
├── collection                    - Collection of the actual profiles
├── pkg/profiler                  - Package for applying profiles
└── pkg/catalogue                 - Utility package to access collection of data
```

This packages reuses CRDs and types from KubeVirt and Kubernetes project
