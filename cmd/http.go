package cmd

import (
	"encoding/base64"
	"fmt"
	"invoice-service/common/gcs"
	constant "invoice-service/constant/error"

	controllerRegistry "invoice-service/controllers"
	repositoryRegistry "invoice-service/repositories"
	routeRegistry "invoice-service/routes"
	serviceRegistry "invoice-service/services"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"invoice-service/common/sentry"
	"invoice-service/config"
	"invoice-service/middlewares"
	"invoice-service/utils/response"

	"net/http"
	"time"
)

var restCmd = &cobra.Command{
	Use:   "serve",
	Short: "Command to start http server",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load() //nolint:errcheck
		config.Init()
		db, err := config.NewDatabaseConnection()
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}
		time.Local = loc

		// Sentry for error tracking
		sentry := sentry.NewSentry(
			sentry.WithDsn(config.Config.SentryDsn),
			sentry.WithDebug(config.Config.AppDebug),
			sentry.WithEnv(config.Config.AppEnv),
			sentry.WithSampleRate(config.Config.SentrySampleRate),
			sentry.WithEnableTracing(config.Config.SentryEnableTracing),
		)

		gcs := initGCS()

		repository := repositoryRegistry.NewRepositoryRegistry(db, sentry)
		service := serviceRegistry.NewServiceRegistry(repository, sentry, gcs)
		controller := controllerRegistry.NewControllerRegistry(service, sentry)

		router := gin.Default()
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  constant.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})

		lmt := tollbooth.NewLimiter(
			config.Config.RateLimiterMaxRequest,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			},
		)
		router.Use(middlewares.RateLimiter(lmt))
		group := router.Group("/api/v1")
		route := routeRegistry.NewRouteRegistry(controller, group)
		route.Serve()

		port := fmt.Sprintf(":%d", config.Config.Port)
		err = router.Run(port)
		if err != nil {
			panic(err)
		}
	},
}

func initGCS() *gcs.GCSPackage {
	decodeGCSPrivateKey, err := base64.StdEncoding.DecodeString(config.Config.GCSPrivateKey)
	if err != nil {
		panic(err)
	}

	var stringGCSPrivateKey = string(decodeGCSPrivateKey)
	gcsServiceAccount := gcs.ServiceAccountKeyJSON{
		Type:                    config.Config.GCSType,
		ProjectId:               config.Config.GCSProjectID,
		PrivateKeyId:            config.Config.GCSPrivateKeyID,
		PrivateKey:              stringGCSPrivateKey,
		ClientEmail:             config.Config.GCSClientEmail,
		ClientId:                config.Config.GCSClientID,
		AuthUri:                 config.Config.GCSAuthURI,
		TokenUri:                config.Config.GCSTokenURI,
		AuthProviderX509CertUrl: config.Config.GCSAuthProviderX509CertURL,
		ClientX509CertUrl:       config.Config.GCSClientX509CertURL,
	}
	gcsClient := gcs.NewGCSClient(
		gcs.WithServiceAccountKeyJSON(gcsServiceAccount),
		gcs.WithSignedURLTimeInMinutes(config.Config.GCSSignedURLTimeInMinutes),
		gcs.WithBucketName(config.Config.GCSBucketName))

	return gcsClient
}

func Run() {
	err := restCmd.Execute()
	if err != nil {
		panic(err)
	}
}
