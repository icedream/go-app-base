// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package appbase

import "fmt"

func (setting *Setting) String(defaultValue string) (target *string) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue string
	target = &realFinalValue

	setting.kingpinFlag.Default(fmt.Sprint(defaultValue)).String()
	setting.readIn = func() (err error) {
		err = setting.application.viper.UnmarshalKey(setting.viperPath, target)
		return
	}
	return
}
