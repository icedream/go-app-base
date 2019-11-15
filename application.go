package appbase

import (
	"fmt"
	"os"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type ApplicationDescription struct {
	// Vendor is an optional string containing the name of the developing team or author.
	// It is used in generated file paths.
	VendorID string

	// ID is a identifying name for this application. It is used in generated file paths.
	// Example: "vendor-my-example-app"
	ID string

	// Name is the application's full actual name.
	// Example: "My Example App"
	Name string

	// ShortDescription is a text explaining the purpose of this application in a single sentence.
	ShortDescription string

	// AllowConfigurationInCurrentDirectory indicates whether configuration files should be scanned for in the current directory.
	// These files will have priority over other scanned paths but not for environment variables.
	AllowConfigurationInCurrentDirectory bool
}

type Application struct {
	*ApplicationDescription

	viper    *viper.Viper
	kingpin  *kingpin.Application
	xdg      *xdg.XDG
	settings *Settings
}

func New(description ApplicationDescription) *Application {
	app := &Application{
		ApplicationDescription: &description,
	}
	app.Initialize()
	return app
}

// EnvPrefix derives a prefix for environment variables for Viper from the application ID.
func (app *Application) EnvPrefix() string {
	return strcase.ToScreamingSnake(app.ID)
}

// Initialize configures Viper and Kingpin respectively.
// This must be run before accessing the configuration and CLI objects.
func (app *Application) Initialize() {
	app.xdg = xdg.New(app.VendorID, app.ID)
	app.kingpin = kingpin.New(os.Args[0], app.ShortDescription)

	envPrefix := app.EnvPrefix()

	v := viper.New()
	v.SetConfigName("config")

	// Current directory has top priority (containers)
	if app.AllowConfigurationInCurrentDirectory {
		v.AddConfigPath(".")
	}

	// XDG specification
	// Read all configs in alphabetical order from config directories
	configDirs := append([]string{
		app.xdg.ConfigHome(),
	}, app.xdg.ConfigDirs()...)
	for _, dir := range configDirs {
		viper.AddConfigPath(dir)
	}

	// System-wide config
	v.AddConfigPath(fmt.Sprintf("/etc/%s", strcase.ToKebab(app.ID)))

	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	v.ReadInConfig()

	app.viper = v

	// Settings
	app.settings = &Settings{
		application:        app,
		registeredSettings: []*Setting{},
	}
}

func (app *Application) Configuration() *viper.Viper {
	return app.viper
}

func (app *Application) CLI() *kingpin.Application {
	return app.kingpin
}

func (app *Application) Settings() *Settings {
	return app.settings
}
