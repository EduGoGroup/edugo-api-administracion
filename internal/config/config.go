package config

import (
	"fmt"
	"time"
)

type Config struct {
	Environment string         `mapstructure:"environment"`
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	Logging     LoggingConfig  `mapstructure:"logging"`
	Auth        AuthConfig     `mapstructure:"auth"`
	Redis       RedisConfig    `mapstructure:"redis"`
	Defaults    DefaultsConfig `mapstructure:"defaults"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	UseMockRepositories bool           `mapstructure:"use_mock_repositories"` // Toggle para usar repositorios mock (true) o reales (false)
	Postgres            PostgresConfig `mapstructure:"postgres"`
	MongoDB             MongoDBConfig  `mapstructure:"mongodb"`
}

type PostgresConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	Database       string `mapstructure:"database"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	MaxConnections int    `mapstructure:"max_connections"`
	SSLMode        string `mapstructure:"ssl_mode"`
}

type MongoDBConfig struct {
	URI      string        `mapstructure:"uri"`
	Database string        `mapstructure:"database"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// AuthConfig contiene toda la configuración de autenticación centralizada
type AuthConfig struct {
	JWT              JWTConfig              `mapstructure:"jwt"`
	Password         PasswordConfig         `mapstructure:"password"`
	RateLimit        RateLimitConfig        `mapstructure:"rate_limit"`
	InternalServices InternalServicesConfig `mapstructure:"internal_services"`
	Cache            AuthCacheConfig        `mapstructure:"cache"`
}

// JWTConfig configuración de tokens JWT
type JWTConfig struct {
	Secret               string        `mapstructure:"secret"`                 // ENV: AUTH_JWT_SECRET
	Issuer               string        `mapstructure:"issuer"`                 // ENV: AUTH_JWT_ISSUER - debe ser "edugo-central"
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`  // ENV: AUTH_JWT_ACCESS_TOKEN_DURATION
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"` // ENV: AUTH_JWT_REFRESH_TOKEN_DURATION
	Algorithm            string        `mapstructure:"algorithm"`              // HS256 por defecto
}

// PasswordConfig configuración de validación de passwords
type PasswordConfig struct {
	MinLength        int  `mapstructure:"min_length"`
	RequireUppercase bool `mapstructure:"require_uppercase"`
	RequireLowercase bool `mapstructure:"require_lowercase"`
	RequireNumber    bool `mapstructure:"require_number"`
	RequireSpecial   bool `mapstructure:"require_special"`
	BcryptCost       int  `mapstructure:"bcrypt_cost"`
}

// RateLimitConfig configuración de rate limiting
type RateLimitConfig struct {
	Login            LoginRateLimitConfig    `mapstructure:"login"`
	InternalServices ServiceRateLimitConfig  `mapstructure:"internal_services"`
	ExternalClients  ServiceRateLimitConfig  `mapstructure:"external_clients"`
}

// LoginRateLimitConfig rate limiting para intentos de login
type LoginRateLimitConfig struct {
	MaxAttempts   int           `mapstructure:"max_attempts"`   // ENV: AUTH_RATE_LIMIT_LOGIN_ATTEMPTS
	Window        time.Duration `mapstructure:"window"`         // ENV: AUTH_RATE_LIMIT_LOGIN_WINDOW
	BlockDuration time.Duration `mapstructure:"block_duration"` // ENV: AUTH_RATE_LIMIT_LOGIN_BLOCK
}

// ServiceRateLimitConfig rate limiting para servicios
type ServiceRateLimitConfig struct {
	MaxRequests int           `mapstructure:"max_requests"`
	Window      time.Duration `mapstructure:"window"`
}

// InternalServicesConfig configuración de servicios internos autorizados
type InternalServicesConfig struct {
	APIKeys  string `mapstructure:"api_keys"`  // ENV: AUTH_INTERNAL_SERVICES_API_KEYS formato: "servicio:key,servicio:key"
	IPRanges string `mapstructure:"ip_ranges"` // ENV: AUTH_INTERNAL_SERVICES_IP_RANGES formato CIDR
}

// AuthCacheConfig configuración de cache para autenticación
type AuthCacheConfig struct {
	TokenValidation CacheItemConfig `mapstructure:"token_validation"`
	UserInfo        CacheItemConfig `mapstructure:"user_info"`
}

// CacheItemConfig configuración de un tipo de cache
type CacheItemConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	TTL     time.Duration `mapstructure:"ttl"`
	MaxSize int           `mapstructure:"max_size"`
}

// RedisConfig configuración de Redis para cache
type RedisConfig struct {
	Host     string `mapstructure:"host"`     // ENV: REDIS_HOST
	Port     int    `mapstructure:"port"`     // ENV: REDIS_PORT
	Password string `mapstructure:"password"` // ENV: REDIS_PASSWORD
	DB       int    `mapstructure:"db"`       // ENV: REDIS_DB
}

func (c *PostgresConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// GetRedisAddr retorna la dirección de Redis en formato host:port
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// DefaultsConfig contiene todas las configuraciones de valores por defecto
type DefaultsConfig struct {
	School SchoolDefaults `mapstructure:"school"`
}

// SchoolDefaults contiene los valores por defecto para escuelas
type SchoolDefaults struct {
	Country          string `mapstructure:"country"`           // ENV: DEFAULT_SCHOOL_COUNTRY
	SubscriptionTier string `mapstructure:"subscription_tier"` // ENV: DEFAULT_SCHOOL_SUBSCRIPTION_TIER
	MaxTeachers      int    `mapstructure:"max_teachers"`      // ENV: DEFAULT_SCHOOL_MAX_TEACHERS
	MaxStudents      int    `mapstructure:"max_students"`      // ENV: DEFAULT_SCHOOL_MAX_STUDENTS
}
