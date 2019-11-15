// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package appbase

import "fmt"

func (setting *Setting) Int64List(defaultValue []int64) (target *[]int64) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue []int64
	target = &realFinalValue

	setting.kingpinFlag.Default(fmt.Sprint(defaultValue)).Int64List()
	setting.readIn = func() (err error) {
		err = setting.application.viper.UnmarshalKey(setting.viperPath, target)
		return
	}
	return
}
