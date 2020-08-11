package deployer

import (
	"github.com/spf13/pflag"
	"path/filepath"
	"sigs.k8s.io/kubetest2/pkg/types"
)

// Name is the name of the deployer
const Name = "powervs"

// New implements deployer.New for kind
func New(opts types.Options) (types.Deployer, *pflag.FlagSet) {
	// create a deployer object and set fields that are not flag controlled
	d := &deployer{
		commonOptions: opts,
		logsDir:       filepath.Join(opts.ArtifactsDir(), "logs"),
	}
	// register flags and return
	return d, bindFlags(d)
}

func bindFlags(d *deployer) *pflag.FlagSet {
	flags := pflag.NewFlagSet(Name, pflag.ContinueOnError)
	return flags
}

// assert that New implements types.NewDeployer
var _ types.NewDeployer = New

type deployer struct {
	// generic parts
	commonOptions types.Options
	logsDir        string // dir to export logs to
}

func (d *deployer) Up() error {
	panic("implement me")
}

func (d *deployer) Down() error {
	panic("implement me")
}

func (d *deployer) IsUp() (up bool, err error) {
	panic("implement me")
}

func (d *deployer) DumpClusterLogs() error {
	panic("implement me")
}

func (d *deployer) Build() error {
	panic("implement me")
}

// assert that deployer implements types.Deployer
var _ types.Deployer = &deployer{}