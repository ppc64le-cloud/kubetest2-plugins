package providers

import "github.com/spf13/pflag"

type Provider interface {
	BindFlags(*pflag.FlagSet)
	DumpConfig(string) error
	Initialize() error
}
