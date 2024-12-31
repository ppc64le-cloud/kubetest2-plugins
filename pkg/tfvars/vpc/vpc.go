package vpc

type TFVars struct {
	VPCName       string `json:"vpc_name"`
	SubnetName    string `json:"vpc_subnet_name"`
	Apikey        string `json:"vpc_api_key,omitempty"`
	SSHKey        string `json:"vpc_ssh_key"`
	Region        string `json:"vpc_region"`
	Zone          string `json:"vpc_zone"`
	ResourceGroup string `json:"vpc_resource_group"`
	NodeImageName string `json:"node_image"`
	NodeProfile   string `json:"node_profile"`
}
