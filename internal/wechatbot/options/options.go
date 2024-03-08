package options

import (
	cliflag "wechatbot/pkg/cli/flag"
	"wechatbot/pkg/log"

	genericoptions "wechatbot/internal/pkg/options"
)

type Options struct {
	AppOptions              *AppOptions                            `json:"app"      mapstructure:"app"`
	JwtOptions              *genericoptions.JwtOptions             `json:"jwt"      mapstructure:"jwt"`
	GenericServerRunOptions *genericoptions.ServerRunOptions       `json:"server"   mapstructure:"server"`
	GRPCOptions             *genericoptions.GRPCOptions            `json:"grpc"     mapstructure:"grpc"`
	InsecureServing         *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing           *genericoptions.SecureServingOptions   `json:"secure"   mapstructure:"secure"`
	Log                     *log.Options                           `json:"log"      mapstructure:"log"`
	PostgresOptions         *genericoptions.PostgresOptions        `json:"postgres" mapstructure:"postgres"`
	RedisOptions            *genericoptions.RedisOptions           `json:"redis"    mapstructure:"redis"`
	UtilsOptions            *UtilsOptions                          `json:"utils"    mapstructure:"utils"`
	ChatGptOptions          *ChatGptOptions                        `json:"chatgpt"  mapstructure:"chatgpt"`
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	o := Options{
		AppOptions:              NewAppOptions(),
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		GRPCOptions:             genericoptions.NewGRPCOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		JwtOptions:              genericoptions.NewJwtOptions(),
		RedisOptions:            genericoptions.NewRedisOptions(),
		Log:                     log.NewOptions(),
		PostgresOptions:         genericoptions.NewPostgresOptions(),
		UtilsOptions:            NewUtilsOptions(),
		ChatGptOptions:          NewChatGptOptions(),
	}
	return &o
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.AppOptions.AddFlags(fss.FlagSet("app"))
	o.Log.AddFlags(fss.FlagSet("logs"))
	o.PostgresOptions.AddFlags(fss.FlagSet("postgres"))
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("server"))
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.JwtOptions.AddFlags(fss.FlagSet("jwt"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.ChatGptOptions.AddFlags(fss.FlagSet("chatgpt"))
	return fss
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.Log.Validate()...)
	return errs
}
