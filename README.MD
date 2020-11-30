
## How to

Working with directory
```shell script
$ go install ./... && export PATH=$(go env GOPATH)/bin:$PATH
$ deployer gen -w ./example -c ./example/deployer.hcl -o /tmp/job.nomad
$ deployer diff -w ./example -c ./example/deployer.hcl
$ deployer apply -w ./example -c ./example/deployer.hcl
```

```shell script
$ deployer git gen -c example/deployer.hcl -u https://github.com/abc/xyz.git
$ deployer git diff -c example/deployer.hcl -u https://github.com/abc/xyz.git
$ deployer git deploy -c example/deployer.hcl -u https://github.com/abc/xyz.git
```


## TODO
- [x] Git repo support. Clone and checkout git ref according to input.   
- [ ] Support configuration versioning to support hooking.
- [ ] Vault integration
