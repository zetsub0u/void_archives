// helper to map cobra pflags to viper configs

package cmd

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type viperPFlagBinding struct {
	configName string
	flagValue  pflag.Value
}

type viperPFlagHelper struct {
	bindings []viperPFlagBinding
}

func (vch *viperPFlagHelper) BindPFlag(configName string, flag *pflag.Flag) (err error) {
	err = viper.BindPFlag(configName, flag)
	if err == nil {
		vch.bindings = append(vch.bindings, viperPFlagBinding{configName, flag.Value})
	}
	return
}

func (vch *viperPFlagHelper) setPFlagsFromViper() {
	for _, v := range vch.bindings {
		v.flagValue.Set(viper.GetString(v.configName))
	}
}
