package powervs

type TFVars struct {
	ResourceGroup string  `json:"powervs_resource_group"`
	DNSName       string  `json:"powervs_dns"`
	DNSZone       string  `json:"powervs_dns_zone"`
	Apikey        string  `json:"powervs_api_key,omitempty"`
	Region        string  `json:"powervs_region"`
	Zone          string  `json:"powervs_zone"`
	ServiceID     string  `json:"powervs_service_id"`
	NetworkName   string  `json:"powervs_network_name"`
	ImageName     string  `json:"powervs_image_name"`
	Memory        float64 `json:"powervs_memory"`
	Processors    float64 `json:"powervs_processors"`
	SSHKey        string  `json:"powervs_ssh_key"`
}
