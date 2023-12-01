package config

import (
	"github.com/spf13/viper"

	"invoice-service/utils/helper"

	"log"
	"os"
)

var Config AppConfig

type AppConfig struct {
	Port                       int      `json:"port" yaml:"port"`
	AppName                    string   `json:"appName" yaml:"appName"`
	AppEnv                     string   `json:"appEnv" yaml:"appEnv"`
	AppDebug                   bool     `json:"appDebug" yaml:"appDebug"`
	SignatureKey               string   `json:"signatureKey" yaml:"signatureKey"`
	StaticKey                  string   `json:"staticKey" yaml:"staticKey"`
	Database                   Database `json:"database" yaml:"database"`
	SentryDsn                  string   `json:"sentryDsn" yaml:"sentryDsn"`
	SentrySampleRate           float64  `json:"sentrySampleRate" yaml:"sentrySampleRate"`
	SentryEnableTracing        bool     `json:"SentryEnableTracing" yaml:"SentryEnableTracing"`
	RateLimiterMaxRequest      float64  `json:"rateLimiterMaxRequest" yaml:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond      int      `json:"rateLimiterTimeSecond" yaml:"rateLimiterTimeSecond"`
	GCSType                    string   `json:"gcsType" yaml:"gcsType"`
	GCSProjectID               string   `json:"gcsProjectID" yaml:"gcsProjectID"`
	GCSPrivateKeyID            string   `json:"gcsPrivateKeyID" yaml:"gcsPrivateKeyID"`
	GCSPrivateKey              string   `json:"gcsPrivateKey" yaml:"gcsPrivateKey"`
	GCSClientEmail             string   `json:"gcsClientEmail" yaml:"gcsClientEmail"`
	GCSClientID                string   `json:"gcsClientID" yaml:"gcsClientID"`
	GCSAuthURI                 string   `json:"gcsAuthURI" yaml:"gcsAuthURI"`
	GCSTokenURI                string   `json:"gcsTokenURI" yaml:"gcsTokenURI"`
	GCSAuthProviderX509CertURL string   `json:"gcsAuthProviderX509CertURL" yaml:"gcsAuthProviderX509CertURL"`
	GCSClientX509CertURL       string   `json:"gcsClientX509CertURL" yaml:"gcsClientX509CertURL"`
	GCSBucketName              string   `json:"gcsBucketName" yaml:"gcsBucketName"`
	GCSSignedURLTimeInMinutes  uint     `json:"gcsSignedURLTimeInMinutes" yaml:"gcsSignedURLTimeInMinutes"`
}

type Database struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Name     string `json:"name" yaml:"name"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Timeout  int    `json:"timeout" yaml:"timeout"`
}

func Init() {
	err := helper.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		log.Printf("failed load cold config from file: %s", viper.ConfigFileUsed())
		err = helper.BindFromConsul(&Config, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			panic(err)
		}
	}
}
