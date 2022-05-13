package service

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/api"
	"github.com/ONSdigital/dp-cantabular-metadata-extractor-api/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the API
type Service struct {
	Config      *config.Config
	Server      HTTPServer
	Router      *mux.Router
	Api         *api.CantabularMetadataExtractorAPI
	ServiceList *ExternalServiceList
	HealthCheck HealthChecker
	datasetAPI  *dataset.Client
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {

	log.Info(ctx, "running service")

	log.Info(ctx, "using service configuration", log.Data{"config": cfg})

	// Get HTTP Server
	r := mux.NewRouter()

	s := serviceList.GetHTTPServer(cfg.BindAddr, r)

	// Get health client for dataset-api
	datasetAPIClient := serviceList.GetHealthClient("dataset-api", cfg.DatasetAPIURL)

	// Get health client for cantabular-api-ext - TODO: reinstate when find out from SCC what endpoint to check
	// cantabularExtClient := serviceList.GetHealthClient("cantabular-api-ext", cfg.CantabularExtURL)
	d := dataset.NewAPIClient(cfg.DatasetAPIURL)

	// Setup the API
	a := api.Setup(ctx, r, d)

	// Get HealthCheck
	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}

	if err := registerCheckers(ctx, hc, datasetAPIClient); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)

	// Start healthcheck
	hc.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{
		Config:      cfg,
		Router:      r,
		Api:         a,
		HealthCheck: hc,
		ServiceList: serviceList,
		Server:      s,
		datasetAPI:  d,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.ServiceList.HealthCheck {
			svc.HealthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		// TODO: Close other dependencies, in the expected order
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context, hc HealthChecker, datasetAPIClient *health.Client) (err error) {
	hasErrors := false

	if err = hc.AddCheck("dataset-api", datasetAPIClient.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for dataset-api", err)
	}

	// TODO: reinstate when find out from SCC what endpoint to check
	// if err = hc.AddCheck("cantabular-api-ext", cantabularExtClient.Checker); err != nil {
	// 	hasErrors = true
	// 	log.Error(ctx, "error adding check for cantabular-api-ext", err)
	// }

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}
