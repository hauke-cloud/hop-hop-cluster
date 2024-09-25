package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	General   General `mapstructure:"general"`
	Logger    Logger  `mapstructure:"logger"`
	Cluster   Cluster `mapstructure:"cluster"`
	TLS       TLS     `mapstructure:"tls"`
	ClientTLS *tls.Config
}

type General struct {
	Mode       string `mapstructure:"mode" validate:"required,oneof=production development"`
	ListenPort int    `mapstructure:"listen_port" default:"8080"`
	NodeName   string `mapstructure:"node_name" validate:"required"`
	NodeIP     string `mapstructure:"node_ip" validate:"required"`
}

type Cluster struct {
	ID       string   `mapstructure:"id" validate:"required"`
	Interval int      `mapstructure:"interval" default:"1"`
	Retries  int      `mapstructure:"retries" default:"5"`
	Members  []Member `mapstructure:"members" validate:"required"`
}

type Member struct {
	Name string `mapstructure:"name" validate:"required"`
	IP   string `mapstructure:"ip" validate:"required"`
}

type TLS struct {
	CAFile   string `mapstructure:"ca_file" default:"/etc/hophop/pki/ca.crt"`
	CertFile string `mapstructure:"cert_file" default:"/etc/hophop/pki/server.crt"`
	KeyFile  string `mapstructure:"key_file" default:"/etc/hophop/pki/server.key"`
}

type Logger struct {
	Encoding          string `mapstructure:"encoding" validate:"required,oneof=json console"`
	Level             string `mapstructure:"level" validate:"required,oneof=debug info warn error dpanic panic fatal"`
	Development       bool   `mapstructure:"development" validate:"required"`
	DisableCaller     bool   `mapstructure:"disable_caller" validate:"omitempty"`
	DisableStacktrace bool   `mapstructure:"disable_stacktrace" validate:"omitempty"`
}

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigType("yaml")
	v.AddConfigPath("config")
	v.SetConfigName(filename)
	v.SetEnvPrefix("HOP_HOP")
	v.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		return Config{}, err
	}

	// Set default values
	setDefaults(&c)

	// Validate the configuration
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return Config{}, fmt.Errorf("validation error: field '%s' failed on the '%s' tag with value '%v'", err.Field(), err.Tag(), err.Value())
		}
		return Config{}, err
	}

	return c, nil
}

// SetDefaults sets the default values for the config based of the "default"-tag defined in the struct
// Parameters:
// - v: The struct to set the defaults in
// Returns:
// - None
func setDefaults(v interface{}) {
	val := reflect.ValueOf(v).Elem()
	setDefaultsRecursively(val)
}

// setDefaultsRecursively sets the default values for the struct based of the "default"-tag defined in the struct
// Parameters:
// - val: The reflect.Value of the struct to set the defaults in
// Returns:
// - None
func setDefaultsRecursively(val reflect.Value) {
	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)
		defaultValue := fieldType.Tag.Get("default")

		if field.Kind() == reflect.Struct {
			setDefaultsRecursively(field)
		} else if defaultValue != "" && field.CanSet() && isEmptyValue(field) {
			switch field.Kind() {
			case reflect.String:
				field.SetString(defaultValue)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if intVal, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
					field.SetInt(intVal)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if uintVal, err := strconv.ParseUint(defaultValue, 10, 64); err == nil {
					field.SetUint(uintVal)
				}
			case reflect.Float32, reflect.Float64:
				if floatVal, err := strconv.ParseFloat(defaultValue, 64); err == nil {
					field.SetFloat(floatVal)
				}
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(defaultValue); err == nil {
					field.SetBool(boolVal)
				}
			}
		}
	}
}

// isEmptyValue checks if the value is empty
// Parameters:
// - v: The reflect.Value to check
// Returns:
// - bool: True if the value is empty,
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}
	return false
}

// LoadClientCertificates loads the client certificates
// Parameters:
// - config: The configuration
// Returns:
// - Config: The configuration
func LoadClientCertificates(config Config) (Config, error) {
	// Client authentication
	// Load the server certificate and key
	serverCert, err := tls.LoadX509KeyPair(config.TLS.CertFile, config.TLS.KeyFile)
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}

	// Load the CA certificate (for client verification)
	caCertPool := x509.NewCertPool()
	caCert, err := os.ReadFile(config.TLS.CAFile)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure the TLS settings
	config.ClientTLS = &tls.Config{
		Certificates: []tls.Certificate{serverCert},  // Server certificate
		ClientCAs:    caCertPool,                     // CA certificates
		RootCAs:      caCertPool,                     // CA certificates
		ClientAuth:   tls.RequireAndVerifyClientCert, // Require and verify client certificate
	}

	return config, nil
}
