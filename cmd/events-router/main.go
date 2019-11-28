package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.mpi-internal.com/Yapo/events-router/pkg/infrastructure"
	"github.mpi-internal.com/Yapo/events-router/pkg/interfaces/handlers"

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
			"events-router_service_events_total",
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
	// To handle http connections you can use an httpHandler
	HTTPHandler := infrastructure.NewHTTPHandler(logger)

	// CLONE-RCONF REMOVE START
	// Initialize remote conf example
	lastUpdate, errRconf := infrastructure.NewRconf(
		conf.EtcdConf.Host,
		conf.EtcdConf.LastUpdate,
		conf.EtcdConf.Prefix,
		logger,
	)

	if errRconf != nil {
		logger.Error("Error loading remote conf")
	} else {
		logger.Info("Remote Conf Updated at %s", lastUpdate.Content.Node.Value)
	}
	// CLONE-RCONF REMOVE END

	useBrowserCache := handlers.Cache{
		MaxAge:  conf.CacheConf.MaxAge,
		Etag:    conf.CacheConf.Etag,
		Enabled: conf.CacheConf.Enabled,
	}
	// Setting up router
	maker := infrastructure.RouterMaker{
		Logger:        logger,
		Cors:          conf.CorsConf,
		Cache:         useBrowserCache,
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
	shutdownSequence.Push(server)
	logger.Info("Starting request serving")
	go server.ListenAndServe()
	shutdownSequence.Wait()
	logger.Info("Server exited normally")

}
