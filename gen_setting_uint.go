// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package appbase

import "fmt"

func (setting *Setting) Uint(defaultValue uint) (target *uint) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue uint
	target = &realFinalValue

	setting.kingpinFlag.Default(fmt.Sprint(defaultValue)).Uint()
	setting.readIn = func() (err error) {
		err = setting.application.viper.UnmarshalKey(setting.viperPath, target)
		return
	}
	return
}