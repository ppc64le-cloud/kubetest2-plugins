# kubetest2-plugins

This project contains the [kubetest2](https://github.com/kubernetes-sigs/kubetest2) plugins for deploying the k8s on different ppc64le cloud and run the tests on it.

## kubetest2-powervs

kubetest2-powervs is a deployer created for deploying on [IBM Cloud Power Virtual Server](https://www.ibm.com/in-en/cloud/power-virtual-server) infrastructure

## Development
```shell
$ export TF_DATA=`pwd`/data/data
$ ./bin/kubetest2-tf
```
