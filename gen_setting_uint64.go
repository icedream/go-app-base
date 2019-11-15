// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package appbase

import "fmt"

func (setting *Setting) Uint64(defaultValue uint64) (target *uint64) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue uint64
	target = &realFinalValue

	setting.kingpinFlag.Default(fmt.Sprint(defaultValue)).Uint64()
	setting.readIn = func() (err error) {
		err = setting.application.viper.UnmarshalKey(setting.viperPath, target)
		return
	}
	return
}