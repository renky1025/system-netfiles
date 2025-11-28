package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"redis"`
	Storage  StorageConfig  `yaml:"storage" mapstructure:"storage"`
	MinIO    MinIOConfig    `yaml:"minio" mapstructure:"minio"`
	JWT      JWTConfig      `yaml:"jwt" mapstructure:"jwt"`
}

type ServerConfig struct {
	Port int    `yaml:"port" mapstructure:"port"`
	Mode string `yaml:"mode" mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"dbname" mapstructure:"dbname"`
	SSLMode  string `yaml:"sslmode" mapstructure:"sslmode"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
}

type StorageConfig struct {
	Type      string `yaml:"type" mapstructure:"type"`
	LocalPath string `yaml:"local_path" mapstructure:"local_path"`
}

type MinIOConfig struct {
	Endpoint        string `yaml:"endpoint" mapstructure:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id" mapstructure:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key" mapstructure:"secret_access_key"`
	BucketName      string `yaml:"bucket_name" mapstructure:"bucket_name"`
	UseSSL          bool   `yaml:"use_ssl" mapstructure:"use_ssl"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret" mapstructure:"secret"`
	ExpireHours int    `yaml:"expire_hours" mapstructure:"expire_hours"`
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("E:/netfilessys/internal/config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Set default JWT secret if not provided
	if AppConfig.JWT.Secret == "" {
		AppConfig.JWT.Secret = "your_secret_key_change_in_production"
		log.Println("WARNING: Using default JWT secret. Please set a secure secret in production!")
	}

	if AppConfig.JWT.ExpireHours == 0 {
		AppConfig.JWT.ExpireHours = 24
	}
}
