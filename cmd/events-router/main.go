package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.mpi-internal.com/Yapo/events-router/pkg/infrastructure"
	"github.mpi-internal.com/Yapo/events-router/pkg/interfaces/handlers"
	"github.mpi-internal.com/Yapo/events-router/pkg/interfaces/loggers"
	"github.mpi-internal.com/Yapo/events-router/pkg/interfaces/repository"
	"github.mpi-internal.com/Yapo/events-router/pkg/usecases"
)

func main() {
	var shutdownSequence = infrastructure.NewShutdownSequence()
	var conf infrastructure.Config
	fmt.Printf("Etag:%d\n", conf.CacheConf.InitEtag())
	shutdownSequence.Listen()
	infrastructure.LoadFromEnv(&conf)
	if jconf, err := json.MarshalIndent(conf, "", "    "); err == nil {
		fmt.Printf("Config: \n%s\n", jconf)
	} else {
		fmt.Printf("Config: \n%+v\n", conf)
	}

	fmt.Printf("Setting up Prometheus\n")
	prometheus := infrastructure.MakePrometheusExporter(
		conf.PrometheusConf.Port,
		conf.PrometheusConf.Enabled,
	)

	fmt.Printf("Setting up logger\n")
	logger, err := infrastructure.MakeYapoLogger(&conf.LoggerConf,
		prometheus.NewEventsCollector(
			"events_router_service_events_total",
			"events tracker counter for events-router service",
		),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	shutdownSequence.Push(prometheus)

	logger.Info("Initializing resources")

	// HealthHandler
	var healthHandler handlers.HealthHandler

	remoteConfig, errRconf := infrastructure.NewRconf(
		conf.EtcdConf.Host,
		conf.EtcdConf.Router,
		conf.EtcdConf.Prefix,
		logger,
	)

	if errRconf != nil {
		logger.Error("Error loading remote conf")
		os.Exit(1)
	} else {
		logger.Info("Remote Conf loaded")
	}

	kafkaProducer, err := infrastructure.NewKafkaProducer(
		conf.KafkaProducerConf.Host,
		conf.KafkaProducerConf.Port,
		conf.KafkaProducerConf.Acks,
		conf.KafkaProducerConf.CompressionType,
		conf.KafkaProducerConf.Retries,
		conf.KafkaProducerConf.LingerMS,
		conf.KafkaProducerConf.RequestTimeoutMS,
		conf.KafkaProducerConf.EnableIdempotence,
	)
	if err != nil {
		logger.Error("Error starting kafka producer: %+v", err)
		os.Exit(1)
	}
	shutdownSequence.Push(kafkaProducer)

	kafkaConsumer, err := infrastructure.NewKafkaConsumer(
		conf.KafkaConsumerConf.Host,
		conf.KafkaConsumerConf.Port,
		conf.KafkaConsumerConf.GroupID,
		conf.KafkaConsumerConf.OffsetReset,
		conf.KafkaConsumerConf.RebalanceEnable,
		conf.KafkaConsumerConf.ChannelEnable,
		conf.KafkaConsumerConf.PartitionEOF,
		conf.KafkaConsumerConf.TimeOut,
		conf.KafkaConsumerConf.Topics,
		logger,
	)
	if err != nil {
		logger.Error("Error starting kafka consumer: %+v", err)
		os.Exit(1)
	}
	shutdownSequence.Push(kafkaConsumer)

	interactor := &usecases.DispatchInteractor{
		Producer: repository.MakeProducer(kafkaProducer),
		Router:   repository.MakeRouter(remoteConfig),
		Logger:   loggers.MakeDispatchInteractorlogger(logger),
	}

	dispatcherHandler := handlers.NewDispatchEventHandler(
		kafkaConsumer,
		interactor,
		loggers.MakeDispatchEventHandlerlogger(logger),
	)

	// Setting up router
	maker := infrastructure.RouterMaker{
		Logger:        logger,
		Cors:          conf.CorsConf,
		WrapperFuncs:  []infrastructure.WrapperFunc{prometheus.TrackHandlerFunc},
		WithProfiling: conf.ServiceConf.Profiling,
		Routes: infrastructure.Routes{
			{
				// This is the base path, all routes will start with this prefix
				Prefix: "/api/v{version:[1-9][0-9]*}",
				Groups: []infrastructure.Route{
					{
						Name:    "Check service health",
						Method:  "GET",
						Pattern: "/healthcheck",
						Handler: &healthHandler,
					},
				},
			},
		},
	}

	router := maker.NewRouter()
	server := infrastructure.NewHTTPServer(
		fmt.Sprintf("%s:%d", conf.Runtime.Host, conf.Runtime.Port),
		router,
		logger,
	)
	logger.Info("Starting request serving")
	go server.ListenAndServe()
	go kafkaConsumer.Listen()
	go dispatcherHandler.Consume()
	shutdownSequence.Wait()
	logger.Info("Server exited normally")

}
