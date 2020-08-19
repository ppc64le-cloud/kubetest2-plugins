package main

import (
	"sigs.k8s.io/kubetest2/pkg/app"

	"github.com/ppc64le-cloud/kubetest2-plugins/kubetest2-tf/deployer"
)

func main() {
	app.Main(deployer.Name, deployer.New)
}
