// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package appbase

import (
	"fmt"
	"time"
)

func (setting *Setting) DurationList(defaultValue []time.Duration) (target *[]time.Duration) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue []time.Duration
	target = &realFinalValue

	setting.kingpinFlag.Default(fmt.Sprint(defaultValue)).DurationList()
	setting.readIn = func() (err error) {
		err = setting.application.viper.UnmarshalKey(setting.viperPath, target)
		return
	}
	return
}
