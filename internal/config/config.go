package config

import (
	"MusicService/internal/model"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"PORT"`
	} `mapstructure:"SERVER"`
	Database struct {
		Host     string `mapstructure:"HOST"`
		Port     string `mapstructure:"PORT"`
		User     string `mapstructure:"USER"`
		Password string `mapstructure:"PASSWORD"`
		Name     string `mapstructure:"NAME"`
	} `mapstructure:"DATABASE"`
	MinIO struct {
		Endpoint   string `mapstructure:"ENDPOINT"`
		AccessKey  string `mapstructure:"ACCESS_KEY"`
		SecretKey  string `mapstructure:"SECRET_KEY"`
		BucketName string `mapstructure:"BUCKET_NAME"`
		UseSSL     bool   `mapstructure:"USE_SSL"`
	} `mapstructure:"MINIO"`
	JWT struct {
		SecretKey string `mapstructure:"SECRET_KEY"`
	} `mapstructure:"JWT"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Track{},
		&model.Playlist{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate models: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}
