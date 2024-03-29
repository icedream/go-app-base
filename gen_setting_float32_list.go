// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package appbase

import "fmt"

func (setting *Setting) Float32List(defaultValue []float32) (target *[]float32) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue []float32
	target = &realFinalValue

	setting.kingpinFlag.Default(fmt.Sprint(defaultValue)).Float32List()
	setting.readIn = func() (err error) {
		err = setting.application.viper.UnmarshalKey(setting.viperPath, target)
		return
	}
	return
}
