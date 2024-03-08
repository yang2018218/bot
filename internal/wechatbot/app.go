package wechatbot

import (
	"strconv"
	"strings"
	"wechatbot/internal/pkg/code"
	"wechatbot/internal/wechatbot/config"
	"wechatbot/internal/wechatbot/options"
	"wechatbot/pkg/app"
	"wechatbot/pkg/log"

	"github.com/spf13/viper"
)

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("wechatbot",
		basename,
		app.WithOptions(opts),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)
	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		// 处理环境变量的配置
		serverMode := viper.GetString("ENV")
		if serverMode != "" {
			opts.GenericServerRunOptions.Mode = serverMode
		}
		pgdsn := viper.GetString("PG_DSN")
		if pgdsn != "" {
			opts.PostgresOptions.DSN = pgdsn
		}
		flywaypath := viper.GetString("FLYWAY_PATH")
		if flywaypath != "" {
			opts.PostgresOptions.FlywayPath = flywaypath
		}
		httpBindAddress := viper.GetString("INSECURE_BIND_ADDRESS")
		if httpBindAddress != "" {
			opts.InsecureServing.BindAddress = httpBindAddress
		}
		httpsBindAddress := viper.GetString("SECURE_BIND_ADDRESS")
		if httpBindAddress != "" {
			opts.SecureServing.BindAddress = httpsBindAddress
		}
		grpcBindAddress := viper.GetString("GRPC_BIND_ADDRESS")
		if grpcBindAddress != "" {
			opts.InsecureServing.BindAddress = grpcBindAddress
		}
		logOutputPaths := viper.GetString("LOG_OUTPUT_PATHS")
		if logOutputPaths != "" {
			opts.Log.OutputPaths = strings.Split(logOutputPaths, ",")
		}
		logErrorOutputPaths := viper.GetString("LOG_ERROR_OUTPUT_PATHS")
		if logErrorOutputPaths != "" {
			opts.Log.ErrorOutputPaths = strings.Split(logErrorOutputPaths, ",")
		}
		redisAddr := viper.GetString("REDIS_ADDR")
		if redisAddr != "" {
			opts.RedisOptions.Addr = redisAddr
		}
		redisPassword := viper.GetString("REDIS_PASSWORD")
		if redisPassword != "" {
			opts.RedisOptions.Password = redisPassword
		}
		redisDB := viper.GetString("REDIS_DB")
		if redisDB != "" {
			redisDBInt, errInt := strconv.Atoi(redisDB)
			if errInt == nil {
				opts.RedisOptions.DB = &redisDBInt
			}
		}
		utilsIP2LocationDBPath := viper.GetString("UTILS_IP2LOCATION_DB_PATH")
		if utilsIP2LocationDBPath != "" {
			opts.UtilsOptions.IP2LocationDBPath = utilsIP2LocationDBPath
		}
		if opts.GenericServerRunOptions.Mode != code.ServerModelDebug {
			opts.Log.Development = false
			opts.Log.EnableColor = false
			opts.Log.Level = "info"
		}
		appStoragePath := viper.GetString("APP_STORAGE")
		if appStoragePath != "" {
			opts.AppOptions.Storage = appStoragePath
		}
		appSetu := viper.GetString("APP_SETU")
		if appSetu != "" {
			opts.AppOptions.Setu = appSetu
		}

		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}
		config.GCfg = *cfg
		return Run(cfg)
	}
}
