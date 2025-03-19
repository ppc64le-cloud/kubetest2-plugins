🚨 **This repository is now archived and no longer maintained.** 🚨  

This means:  
- No further updates, bug fixes, or feature additions will be made.  
- Issues and pull requests will not be reviewed.  
- The repository remains available for historical reference.  

For alternative solutions or continued development, please check:  
👉 https://github.com/kubernetes-sigs/provider-ibmcloud-test-infra/tree/main/kubetest2-tf

# kubetest2-plugins

This project contains the [kubetest2](https://github.com/kubernetes-sigs/kubetest2) plugins for deploying k8s on different ppc64le cloud and run the tests on it. This plugin predominantly uses terraform for infrastructure provisioning and ansible for setting up k8s on the deployed infrastructure.

## kubetest2-powervs

kubetest2-powervs is a deployer created for deploying on [IBM Cloud Power Virtual Server](https://www.ibm.com/in-en/cloud/power-virtual-server) infrastructure.

### Plugin Installation

##### Using make:
The plugin can be installed by executing the following command from repository root:
```
# make install-deployer-tf
```


##### Manual Installation:
The kubetest2-plugin uses the `powervs` and `k8s-ansible` resources as embedded files.

As [`k8s-ansible`](https://github.com/ppc64le-cloud/k8s-ansible) is a submodule, it requires initialisation and update to clone the repository to the data/k8s-ansible path before the binary is built.


```
# git submodule update --init
```
Install the kubetest2-tf plugin using the following command:
```
# go install ./...
```
To update the repository with the latest commits in submodules:
```
# git submodule update --remote
```
