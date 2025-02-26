package config

type Config struct {
	MinioEndpoint     string
	BucketName        string
	MinioRootUser     string
	MinioRootPassword string
	MinioUseSSL       bool
}
