package appbase

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

//go:generate genny -in setting_func.gotpl -out gen_setting_bool.go gen "KingpinType=Bool ActualType=bool"
//go:generate genny -in setting_func.gotpl -out gen_setting_bool_list.go gen "KingpinType=BoolList ActualType=[]bool"
//go:generate genny -in setting_func.gotpl -out gen_setting_duration.go gen "KingpinType=Duration ActualType=time.Duration"
//go:generate genny -in setting_func.gotpl -out gen_setting_duration_list.go gen "KingpinType=DurationList ActualType=[]time.Duration"
//go:generate genny -in setting_func.gotpl -out gen_setting_float32.go gen "KingpinType=Float32 ActualType=float32"
//go:generate genny -in setting_func.gotpl -out gen_setting_float32_list.go gen "KingpinType=Float32List ActualType=[]float32"
//go:generate genny -in setting_func.gotpl -out gen_setting_float64.go gen "KingpinType=Float64 ActualType=float64"
//go:generate genny -in setting_func.gotpl -out gen_setting_float64_list.go gen "KingpinType=Float64List ActualType=[]float64"
//go:generate genny -in setting_func.gotpl -out gen_setting_int.go gen "KingpinType=Int ActualType=int"
//go:generate genny -in setting_func.gotpl -out gen_setting_ints.go gen "KingpinType=Ints ActualType=[]int"
//go:generate genny -in setting_func.gotpl -out gen_setting_int8.go gen "KingpinType=Int8 ActualType=int8"
//go:generate genny -in setting_func.gotpl -out gen_setting_int8_list.go gen "KingpinType=Int8List ActualType=[]int8"
//go:generate genny -in setting_func.gotpl -out gen_setting_int16.go gen "KingpinType=Int16 ActualType=int16"
//go:generate genny -in setting_func.gotpl -out gen_setting_int16_list.go gen "KingpinType=Int16List ActualType=[]int16"
//go:generate genny -in setting_func.gotpl -out gen_setting_int32.go gen "KingpinType=Int32 ActualType=int32"
//go:generate genny -in setting_func.gotpl -out gen_setting_int32_list.go gen "KingpinType=Int32List ActualType=[]int32"
//go:generate genny -in setting_func.gotpl -out gen_setting_int64.go gen "KingpinType=Int64 ActualType=int64"
//go:generate genny -in setting_func.gotpl -out gen_setting_int64_list.go gen "KingpinType=Int64List ActualType=[]int64"
//go:generate genny -in setting_func.gotpl -out gen_setting_string.go gen "KingpinType=String ActualType=string"
//go:generate genny -in setting_func.gotpl -out gen_setting_strings.go gen "KingpinType=Strings ActualType=[]string"
//go:generate genny -in setting_func.gotpl -out gen_setting_uint.go gen "KingpinType=Uint ActualType=uint"
//go:generate genny -in setting_func.gotpl -out gen_setting_uints.go gen "KingpinType=Uints ActualType=[]uint"
//go:generate genny -in setting_func.gotpl -out gen_setting_uint8.go gen "KingpinType=Uint8 ActualType=uint8"
//go:generate genny -in setting_func.gotpl -out gen_setting_uint8_list.go gen "KingpinType=Uint8List ActualType=[]uint8"
//go:generate genny -in setting_func.gotpl -out gen_setting_uint16.go gen "KingpinType=Uint16 ActualType=uint16"
//go:generate genny -in setting_func.gotpl -out gen_setting_uint16_list.go gen "KingpinType=Uint16List ActualType=[]uint16"
//go:generate genny -in setting_func.gotpl -out gen_setting_uint32.go gen "KingpinType=Uint32 ActualType=uint32"
//go:generate genny -in setting_func.gotpl -out gen_setting_uint32_list.go gen "KingpinType=Uint32List ActualType=[]uint32"
//go:generate genny -in setting_func.gotpl -out gen_setting_uint64.go gen "KingpinType=Uint64 ActualType=uint64"
//go:generate genny -in setting_func.gotpl -out gen_setting_uint64_list.go gen "KingpinType=Uint64List ActualType=[]uint64"

type Settings struct {
	application        *Application
	registeredSettings []*Setting
}

func (settings *Settings) Create(key string, description string) *Setting {
	setting := &Setting{
		application: settings.application,
		viperPath:   key,
		kingpinFlag: settings.application.kingpin.Flag(key, description),
	}

	setting.kingpinFlag = setting.kingpinFlag.
		Action(settings.application.viperAdapter). // Read into viper
		NoEnvar()                                  // Viper already reads values from environment

	settings.registeredSettings = append(settings.registeredSettings, setting)

	return setting
}

func (settings *Settings) ReadIn() (command string, errors []error) {
	command = kingpin.MustParse(settings.application.kingpin.Parse(os.Args[1:]))

	errors = []error{}
	for _, setting := range settings.registeredSettings {
		err := setting.ReadIn()
		if err != nil {
			errors = append(errors, &SettingsReadInError{
				error: fmt.Errorf("%s: %s", setting.viperPath, err.Error()),
				Path:  setting.viperPath,
				Inner: err,
			})
		}
	}
	return
}

type SettingsReadInError struct {
	error
	Path  string
	Inner error
}

type Setting struct {
	application *Application
	viperPath   string
	kingpinFlag *kingpin.Clause

	readIn func() error
}

func (app *Application) viperAdapter(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
	var name string
	switch {
	case element.OneOf.Flag != nil:
		name = element.OneOf.Flag.Model().Name
	case element.OneOf.Arg != nil:
		name = element.OneOf.Arg.Model().Name
	default:
		return nil
	}
	app.viper.Set(name, *element.Value)
	return nil
}

func (setting *Setting) Hidden() *Setting {
	setting.kingpinFlag = setting.kingpinFlag.Hidden()
	return setting
}

func (setting *Setting) HintAction(action kingpin.HintAction) *Setting {
	setting.kingpinFlag = setting.kingpinFlag.HintAction(action)
	return setting
}

func (setting *Setting) PlaceHolder(placeholder string) *Setting {
	setting.kingpinFlag = setting.kingpinFlag.PlaceHolder(placeholder)
	return setting
}

func (setting *Setting) Short(name rune) *Setting {
	setting.kingpinFlag = setting.kingpinFlag.Short(name)
	return setting
}

/*func (setting *Setting) Required() *Setting {
	setting.kingpinFlag.Required()
	return setting
}*/

func (setting *Setting) ReadIn() (err error) {
	// Sanity check - postReadIn must be set to convert interface{} back to expected type
	if setting.readIn == nil {
		panic(fmt.Sprintf("no setting type was set for %s", setting.viperPath))
	}

	err = setting.readIn()
	return
}

func (setting *Setting) Counter(defaultValue int) (target *int) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue int
	target = &realFinalValue

	setting.kingpinFlag.Default(fmt.Sprint(defaultValue)).Counter()
	setting.readIn = func() (err error) {
		err = setting.application.viper.UnmarshalKey(setting.viperPath, target)
		return
	}
	return
}

func (setting *Setting) Enum(defaultValue string, options ...string) (target *string) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue string
	target = &realFinalValue

	setting.kingpinFlag.Default(defaultValue).Enum(options...)
	setting.readIn = func() (err error) {
		var rawValue string
		err = setting.application.viper.UnmarshalKey(setting.viperPath, &rawValue)
		if err != nil {
			return
		}

		for _, option := range options {
			if strings.EqualFold(option, rawValue) {
				*target = rawValue
				return
			}
		}

		err = fmt.Errorf("invalid value %q, allowed values are: %q", defaultValue, options)
		return
	}
	return
}

func (setting *Setting) Enums(defaultValue []string, options ...string) (target *[]string) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue []string
	target = &realFinalValue

	setting.kingpinFlag.Default(defaultValue...).Enums(options...)
	setting.readIn = func() (err error) {
		var rawValues []string
		err = setting.application.viper.UnmarshalKey(setting.viperPath, &rawValues)
		if err != nil {
			return
		}

	rawValueLoop:
		for _, rawValue := range rawValues {
			for _, option := range options {
				if strings.EqualFold(option, rawValue) {
					continue rawValueLoop
				}
			}
			err = fmt.Errorf("invalid value %q, allowed values are: %q", defaultValue, options)
			return
		}

		*target = rawValues
		return
	}
	return
}

func (setting *Setting) IP(defaultValue net.IP) (target *net.IP) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue.String())

	var realFinalValue net.IP
	target = &realFinalValue

	setting.kingpinFlag.Default(fmt.Sprint(defaultValue)).String()
	setting.readIn = func() (err error) {
		var rawValue string
		if err = setting.application.viper.UnmarshalKey(setting.viperPath, &rawValue); err != nil {
			return
		}

		*target = net.ParseIP(rawValue)
		return
	}
	return
}

func (setting *Setting) IPList(defaultValue []net.IP) (target *[]net.IP) {
	defaultValues := make([]string, len(defaultValue))
	for i, value := range defaultValue {
		defaultValues[i] = value.String()
	}
	setting.application.viper.SetDefault(setting.viperPath, defaultValues)

	var realFinalValue []net.IP
	target = &realFinalValue

	setting.kingpinFlag.String() // TODO - default
	setting.readIn = func() (err error) {
		var rawValues []string
		if err = setting.application.viper.UnmarshalKey(setting.viperPath, &rawValues); err != nil {
			return
		}

		*target = make([]net.IP, len(rawValues))
		for i, rawValue := range rawValues {
			(*target)[i] = net.ParseIP(rawValue)
		}
		return
	}
	return
}

func (setting *Setting) Time(format string, defaultValue time.Time) (target *time.Time) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue time.Time
	target = &realFinalValue

	setting.kingpinFlag.Default(defaultValue.Format(format)).Time(format)
	setting.readIn = func() (err error) {
		err = setting.application.viper.UnmarshalKey(setting.viperPath, target)
		return
	}
	return
}

func (setting *Setting) TimeList(format string, defaultValue time.Time) (target *[]time.Time) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue)

	var realFinalValue []time.Time
	target = &realFinalValue

	setting.kingpinFlag.TimeList(format) // TODO - default
	setting.readIn = func() (err error) {
		err = setting.application.viper.UnmarshalKey(setting.viperPath, target)
		return
	}
	return
}

func (setting *Setting) URL(defaultValue *url.URL, mustBeAbsolute bool) (target **url.URL) {
	setting.application.viper.SetDefault(setting.viperPath, defaultValue.String())

	var realFinalValue *url.URL
	target = &realFinalValue

	setting.kingpinFlag.Default(fmt.Sprint(defaultValue)).String()
	setting.readIn = func() (err error) {
		var rawValue string
		if err = setting.application.viper.UnmarshalKey(setting.viperPath, &rawValue); err != nil {
			return
		}

		u, err := url.Parse(rawValue)
		if err != nil {
			return
		}

		*target = u
		return
	}
	return
}

// TODO - URLList
