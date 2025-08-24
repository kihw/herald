package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Riot     RiotConfig     `mapstructure:"riot"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	Environment  string        `mapstructure:"environment"`
	Debug        bool          `mapstructure:"debug"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
	Driver   string `mapstructure:"driver"` // sqlite or postgres
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type RiotConfig struct {
	APIKey             string        `mapstructure:"api_key"`
	RateLimitPerSecond int           `mapstructure:"rate_limit_per_second"`
	RateLimitPerMinute int           `mapstructure:"rate_limit_per_minute"`
	BaseURL            string        `mapstructure:"base_url"`
	Timeout            time.Duration `mapstructure:"timeout"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    string `mapstructure:"port"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set default values
	setDefaults()

	// Read environment variables
	viper.AutomaticEnv()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: No config file found, using defaults and environment variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Override with environment variables for critical settings
	overrideWithEnv(&config)

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.idle_timeout", "30s")
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("server.debug", false)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "herald")
	viper.SetDefault("database.password", "herald_dev")
	viper.SetDefault("database.name", "herald_dev")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.driver", "sqlite")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// JWT defaults
	viper.SetDefault("jwt.secret", "change_me_in_production")
	viper.SetDefault("jwt.expiration", "24h")

	// Riot API defaults
	viper.SetDefault("riot.base_url", "https://na1.api.riotgames.com")
	viper.SetDefault("riot.rate_limit_per_second", 20)
	viper.SetDefault("riot.rate_limit_per_minute", 100)
	viper.SetDefault("riot.timeout", "30s")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	// Metrics defaults
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.port", "9091")
}

func overrideWithEnv(config *Config) {
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}

	if env := os.Getenv("ENV"); env != "" {
		config.Server.Environment = env
	}

	if debug := os.Getenv("DEBUG"); debug != "" {
		if val, err := strconv.ParseBool(debug); err == nil {
			config.Server.Debug = val
		}
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		config.Database.Port = dbPort
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}

	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.Name = dbName
	}

	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		config.Redis.Host = redisHost
	}

	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		config.Redis.Port = redisPort
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWT.Secret = jwtSecret
	}

	if riotAPIKey := os.Getenv("RIOT_API_KEY"); riotAPIKey != "" {
		config.Riot.APIKey = riotAPIKey
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// GetDatabaseDSN returns the database DSN string
func (c *Config) GetDatabaseDSN() string {
	switch c.Database.Driver {
	case "sqlite":
		return "./herald.db"
	case "postgres":
		return "host=" + c.Database.Host +
			" port=" + c.Database.Port +
			" user=" + c.Database.User +
			" password=" + c.Database.Password +
			" dbname=" + c.Database.Name +
			" sslmode=" + c.Database.SSLMode
	default:
		return "./herald.db"
	}
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return c.Redis.Host + ":" + c.Redis.Port
}
