package config

type AwsS3 struct {
	Bucket           string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Region           string `mapstructure:"region" json:"region" yaml:"region"`
	Endpoint         string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	S3ForcePathStyle bool   `mapstructure:"s3-force-path-style" json:"s3ForcePathStyle" yaml:"s3-force-path-style"`
	DisableSSL       bool   `mapstructure:"disable-ssl" json:"disableSSL" yaml:"disable-ssl"`
	SecretID         string `mapstructure:"secret-id" json:"secretID" yaml:"secret-id"`
	SecretKey        string `mapstructure:"secret-key" json:"secretKey" yaml:"secret-key"`
	BaseURL          string `mapstructure:"base-url" json:"baseURL" yaml:"base-url"`
	PathPrefix       string `mapstructure:"path-prefix" json:"pathPrefix" yaml:"path-prefix"`
}
