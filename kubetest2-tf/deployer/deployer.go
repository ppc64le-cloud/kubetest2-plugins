package deployer

import (
	"encoding/json"
	goflag "flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"text/template"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/pflag"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/kubetest2/pkg/types"

	"github.com/ppc64le-cloud/kubetest2-plugins/kubetest2-tf/deployer/options"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/ansible"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/build"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/providers"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/providers/common"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/providers/powervs"
	"github.com/ppc64le-cloud/kubetest2-plugins/pkg/terraform"
)

const (
	Name              = "tf"
	inventoryTemplate = `[masters]
{{range .Masters}}{{.}}
{{end}}
[workers]
{{range .Workers}}{{.}}
{{end}}
`
)

var GitTag string

type AnsibleInventory struct {
	Masters []string
	Workers []string
}

// Add additional Linux package dependencies here, used by checkDependencies()
var dependencies = []string{"terraform", "ansible"}

func (i *AnsibleInventory) addMachine(mtype string, value string) {
	v := reflect.ValueOf(i).Elem().FieldByName(mtype)
	if v.IsValid() {
		v.Set(reflect.Append(v, reflect.ValueOf(value)))
	}
}

type deployer struct {
	commonOptions types.Options
	BuildOptions  *options.BuildOptions
	logsDir       string
	doInit        sync.Once
	tmpDir        string
	provider      providers.Provider

	RepoRoot              string            `desc:"The path to the root of the local kubernetes repo. Necessary to call certain scripts. Defaults to the current directory. If operating in legacy mode, this should be set to the local kubernetes/kubernetes repo."`
	IgnoreClusterDir      bool              `desc:"Ignore the cluster folder if exists"`
	AutoApprove           bool              `desc:"Terraform Auto Approve"`
	RetryOnTfFailure      int               `desc:"Retry on Terraform Apply Failure"`
	BreakKubetestOnUpfail bool              `desc:"Breaks kubetest2 when up fails"`
	Playbook              string            `desc:"name of ansible playbook to be run"`
	ExtraVars             map[string]string `desc:"Passes extra-vars to ansible playbook, enter a string of key=value pairs"`
	SetKubeConfig         bool              `desc:"Flag to set kubeconfig"`
}

func (d *deployer) Version() string {
	return GitTag
}

func (d *deployer) init() error {
	var err error
	d.doInit.Do(func() { err = d.initialize() })
	return err
}

func (d *deployer) initialize() error {
	fmt.Println("Check if package dependencies are installed in the environment")
	if d.commonOptions.ShouldBuild() {
		if err := d.verifyBuildFlags(); err != nil {
			return fmt.Errorf("init failed to check build flags: %s", err)
		}
	}
	if err := d.checkDependencies(); err != nil {
		return err
	}
	d.provider = powervs.PowerVSProvider
	common.CommonProvider.Initialize()
	d.tmpDir = common.CommonProvider.ClusterName
	if _, err := os.Stat(d.tmpDir); os.IsNotExist(err) {
		err := os.Mkdir(d.tmpDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create dir: %s", d.tmpDir)
		}
	} else if !d.IgnoreClusterDir {
		return fmt.Errorf("directory named %s already exist, please choose a different cluster-name", d.tmpDir)
	}
	return nil
}

var _ types.Deployer = &deployer{}

func New(opts types.Options) (types.Deployer, *pflag.FlagSet) {
	d := &deployer{
		commonOptions: opts,
		logsDir:       filepath.Join(opts.RunDir(), "logs"),
		BuildOptions: &options.BuildOptions{
			CommonBuildOptions: &build.Options{
				Builder:         &build.NoopBuilder{},
				Stager:          &build.NoopStager{},
				Strategy:        "make",
				TargetBuildArch: "linux/ppc64le",
			},
		},
		RetryOnTfFailure: 1,
		Playbook:         "install-k8s.yml",
	}
	flagSet, err := gpflag.Parse(d)
	if err != nil {
		klog.Fatalf("couldn't parse flagset for deployer struct: %s", err)
	}
	klog.InitFlags(nil)
	flagSet.AddGoFlagSet(goflag.CommandLine)
	fs := bindFlags(d)
	flagSet.AddFlagSet(fs)
	return d, flagSet
}

func bindFlags(d *deployer) *pflag.FlagSet {
	flags := pflag.NewFlagSet(Name, pflag.ContinueOnError)
	common.CommonProvider.BindFlags(flags)
	powervs.PowerVSProvider.BindFlags(flags)

	return flags
}

func (d *deployer) Up() error {
	if err := d.init(); err != nil {
		return fmt.Errorf("up failed to init: %s", err)
	}

	err := common.CommonProvider.DumpConfig(d.tmpDir)
	if err != nil {
		return fmt.Errorf("failed to dump common flags: %s", d.tmpDir)
	}

	err = d.provider.DumpConfig(d.tmpDir)
	if err != nil {
		return fmt.Errorf("failed to dumpconfig to: %s and err: %+v", d.tmpDir, err)
	}

	for i := 0; i <= d.RetryOnTfFailure; i++ {
		path, err := terraform.Apply(d.tmpDir, "powervs", d.AutoApprove)
		op, oerr := terraform.Output(d.tmpDir, "powervs")
		if err != nil {
			if i == d.RetryOnTfFailure {
				fmt.Printf("terraform.Output: %s\nterraform.Output error: %v\n", op, oerr)
				if !d.BreakKubetestOnUpfail {
					return fmt.Errorf("terraform Apply failed. Error: %v", err)
				}
				klog.Infof("Terraform Apply failed. Look into it and delete the resources")
				klog.Infof("terraform.Apply error: %v", err)
				os.Exit(1)
			}
			continue
		} else {
			fmt.Printf("terraform.Output: %s\nterraform.Output error: %v\n", op, oerr)
			fmt.Printf("Terraform State at: %s\n", path)
			break
		}
	}
	inventory := AnsibleInventory{}
	for _, machineType := range []string{"Masters", "Workers"} {
		var tmp []interface{}
		op, err := terraform.Output(d.tmpDir, "powervs", "-json", strings.ToLower(machineType))

		if err != nil {
			return fmt.Errorf("terraform.Output failed: %v", err)
		}
		klog.Infof("%s: %s", strings.ToLower(machineType), op)
		err = json.Unmarshal([]byte(op), &tmp)
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %v", err)
		}
		for index := range tmp {
			inventory.addMachine(machineType, tmp[index].(string))
		}
	}
	klog.Infof("inventory: %v", inventory)
	t := template.New("Ansible inventory file")

	t, err = t.Parse(inventoryTemplate)
	if err != nil {
		return fmt.Errorf("template parse failed: %v", err)
	}

	inventoryFile, err := os.Create(filepath.Join(d.tmpDir, "hosts"))
	if err != nil {
		log.Println("create file: ", err)
		return fmt.Errorf("failed to create inventory file: %v", err)
	}

	err = t.Execute(inventoryFile, inventory)
	if err != nil {
		return fmt.Errorf("template execute failed: %v", err)
	}

	common.CommonProvider.ExtraCerts = strings.Join(inventory.Masters, ",")

	commonJSON, err := json.Marshal(common.CommonProvider)
	if err != nil {
		return fmt.Errorf("failed to marshal provider into JSON: %v", err)
	}
	klog.Infof("commonJSON: %v", string(commonJSON))
	//Unmarshalling commonJSON into map to add extra-vars
	final := map[string]interface{}{}
	json.Unmarshal([]byte(commonJSON), &final)
	//Iterating through extra-vars and adding them to map
	for k := range d.ExtraVars {
		final[k] = d.ExtraVars[k]
	}
	//Marshalling back the map to JSON
	finalJSON, err := json.Marshal(final)
	if err != nil {
		return fmt.Errorf("failed to marshal provider into JSON: %v", err)
	}
	klog.Infof("finalJSON with extra vars: %v", string(finalJSON))

	exitcode, err := ansible.Playbook(d.tmpDir, filepath.Join(d.tmpDir, "hosts"), string(finalJSON), d.Playbook)
	if err != nil {
		return fmt.Errorf("failed to run ansible playbook: %v\n with exit code: %d", err, exitcode)
	}

	if d.SetKubeConfig {
		if err = setKubeconfig(inventory.Masters[0]); err != nil {
			return fmt.Errorf("failed to setKubeconfig: %v", err)
		}
		fmt.Printf("KUBECONFIG set to: %s\n", os.Getenv("KUBECONFIG"))
	}
	return nil
}

// setKubeconfig overrides the server IP addresses in the kubeconfig and set the KUBECONFIG environment
func setKubeconfig(host string) error {
	_, err := os.Stat(common.CommonProvider.KubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to locate the kubeconfig file: %v", err)
	}

	config, err := clientcmd.LoadFromFile(common.CommonProvider.KubeconfigPath)
	if err != nil {
		fmt.Printf("failed to load the kuneconfig file")
	}
	for i := range config.Clusters {
		surl, err := url.Parse(config.Clusters[i].Server)
		if err != nil {
			return fmt.Errorf("failed while Parsing the URL: %s", config.Clusters[i].Server)
		}
		_, port, err := net.SplitHostPort(surl.Host)
		if err != nil {
			return fmt.Errorf("errored while SplitHostPort")
		}
		surl.Host = net.JoinHostPort(host, port)
		config.Clusters[i].Server = surl.String()
	}
	clientcmd.WriteToFile(*config, common.CommonProvider.KubeconfigPath)
	kubecfgAbsPath, err := filepath.Abs(common.CommonProvider.KubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to create absolute path for the kubeconfig file: %v", err)
	}
	if err = os.Setenv("KUBECONFIG", kubecfgAbsPath); err != nil {
		return fmt.Errorf("failed to set the KUBECONFIG environment variable")
	}
	return nil
}

func (d *deployer) Down() error {
	if err := d.init(); err != nil {
		return fmt.Errorf("down failed to init: %s", err)
	}
	err := terraform.Destroy(d.tmpDir, "powervs", d.AutoApprove)
	if err != nil {
		if common.CommonProvider.IgnoreDestroy {
			klog.Infof("terraform.Destroy failed: %v", err)
		} else {
			return fmt.Errorf("terraform.Destroy failed: %v", err)
		}
	}
	return nil
}

func (d *deployer) IsUp() (up bool, err error) {
	panic("implement me")
}

func (d *deployer) DumpClusterLogs() error {
	panic("implement me")
}

// checkDependencies determines if the required packages are installed before
// the test execution begins, providing a fail-fast route for exit if the packages are not found.
func (d *deployer) checkDependencies() error {
	for _, dependency := range dependencies {
		if _, err := exec.LookPath(dependency); err != nil {
			return fmt.Errorf("failed to find %s in the test environment: %s", dependency, err)
		}
	}
	return nil
}
