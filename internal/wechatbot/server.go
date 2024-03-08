package wechatbot

import (
	"context"
	"fmt"
	"os"
	"wechatbot/internal/pkg/code"
	genericoptions "wechatbot/internal/pkg/options"
	genericapiserver "wechatbot/internal/pkg/server"
	"wechatbot/internal/pkg/utils/barkutil"
	"wechatbot/internal/wechatbot/config"
	"wechatbot/internal/wechatbot/schedule"
	"wechatbot/internal/wechatbot/service"
	srv "wechatbot/internal/wechatbot/service"
	"wechatbot/internal/wechatbot/store"
	"wechatbot/internal/wechatbot/store/postgres"
	"wechatbot/internal/wechatbot/store/redis"
	"wechatbot/pkg/log"
	"wechatbot/pkg/shutdown"
	"wechatbot/pkg/shutdown/shutdownmanagers/posixsignal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type apiServer struct {
	gs               *shutdown.GracefulShutdown
	redisOptions     *genericoptions.RedisOptions
	gRPCAPIServer    *grpcAPIServer
	genericAPIServer *genericapiserver.GenericAPIServer
}

type preparedAPIServer struct {
	*apiServer
}

// ExtraConfig defines extra configuration for the iam-apiserver.
type ExtraConfig struct {
	Addr            string
	MaxMsgSize      int
	ServerCert      genericoptions.GeneratableKeyCert
	PostgresOptions *genericoptions.PostgresOptions
	// MySQLOptions    *genericoptions.MySQLOptions
	RedisOptions *genericoptions.RedisOptions
	// etcdOptions      *genericoptions.EtcdOptions
}

func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	extraConfig, err := buildExtraConfig(cfg)

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	extraServer, err := extraConfig.complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		gs:               gs,
		redisOptions:     cfg.RedisOptions,
		genericAPIServer: genericServer,
		gRPCAPIServer:    extraServer,
	}

	return server, nil
}

var (
	leaseCtx    context.Context
	leaseCancel context.CancelFunc
)

func (s *apiServer) PrepareRun() preparedAPIServer {
	// s.initRedisStore()
	leaseCtx, leaseCancel = context.WithCancel(context.Background())
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		log.Info("gs close start")
		if leaseCancel != nil {
			leaseCancel()
			log.Info("lease canceled")
		}
		pgStore, _ := postgres.GetPostgreStore(nil)
		if pgStore != nil {
			_ = pgStore.Close()
		}

		s.gRPCAPIServer.Close()
		s.genericAPIServer.Close()
		log.Info("gs close finished")
		return nil
	}))

	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {
	// 项目数据初始化
	storeIns, _ := postgres.GetPostgreStore(nil)
	// 数据库脚本
	if config.GCfg.PostgresOptions.FlywayPath != "" {
		if config.GCfg.GenericServerRunOptions.Runtime == code.RuntimeK8S {
			if config.GCfg.PostgresOptions.FlywayPath != "" {
				flyways, errFlyway := postgres.Migrate(storeIns.DB().Debug(), config.GCfg.PostgresOptions.FlywayPath)
				msg := config.GCfg.PostgresOptions.FlywayPath
				if errFlyway != nil {
					msg += fmt.Sprintf(" failed:%s", errFlyway.Error())
				} else {
					msg += fmt.Sprintf(" succeed %v file(s)", len(flyways))
				}
				barkutil.SendMsg("post-install-flyway", msg)
				os.Exit(0)
			}
		} else {
			postgres.Migrate(storeIns.DB().Debug(), config.GCfg.PostgresOptions.FlywayPath)
		}
	}
	cacheIns, _ := redis.GetRedisCache(nil)
	storeIns.SetCache(cacheIns)
	service := srv.NewService(storeIns)
	err := service.Init()
	if err != nil {
		return err
	}
	// 启动定时初始化
	go func() {
		scheduler, _ := schedule.GetScheduler()
		scheduler.Start()
	}()
	// 启动grpc
	if s.gRPCAPIServer != nil {
		go s.gRPCAPIServer.Run()
	}
	// 初始化api server
	s.genericAPIServer.InitGenericAPIServer()
	// 配置路由表
	initRouter(s.genericAPIServer.Engine)
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}
	// 启动wechat机器人
	go func() {
		errBot := service.WechatBot().Start()
		if errBot != nil {
			log.Panic(errBot.Error())
		}
	}()
	// 启动http(s) server
	return s.genericAPIServer.Run()
}

type completedExtraConfig struct {
	*ExtraConfig
}

// Complete fills in any fields not set that are required to have valid data and can be derived from other fields.
func (c *ExtraConfig) complete() *completedExtraConfig {
	// if c.Addr == "" {
	// 	c.Addr = "127.0.0.1:8081"
	// }

	return &completedExtraConfig{c}
}

// New create a grpcAPIServer instance. 启动grpc和数据库服务
func (c *completedExtraConfig) New() (*grpcAPIServer, error) {
	// creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	// if err != nil {
	// 	log.Fatalf("Failed to generate credentials %s", err.Error())
	// }
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize)}
	grpcServer := grpc.NewServer(opts...)

	storeIns, err := postgres.GetPostgreStore(c.PostgresOptions)
	store.SetClient(storeIns)

	if err != nil {
		log.Fatalf("Failed to get db instance: %s", err.Error())
	}

	_, err = redis.GetRedisCache(c.RedisOptions)
	if err != nil {
		log.Fatalf("Failed to get cache instance: %s", err.Error())
	}

	reflection.Register(grpcServer)

	return &grpcAPIServer{grpcServer, c.Addr}, nil
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	genericConfig = genericapiserver.NewConfig()
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	// if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
	// 	return
	// }
	if cfg.SecureServing != nil && cfg.SecureServing.Enable {
		if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
			return
		}
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	genericConfig.HealthzHandler = service.Healthz

	return
}

// nolint: unparam
func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	return &ExtraConfig{
		Addr:            fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
		MaxMsgSize:      cfg.GRPCOptions.MaxMsgSize,
		ServerCert:      cfg.SecureServing.ServerCert,
		PostgresOptions: cfg.PostgresOptions,
		RedisOptions:    cfg.RedisOptions,
		// etcdOptions:      cfg.EtcdOptions,
	}, nil
}
