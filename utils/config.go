package utils

import (
	"fmt"
	"os"
	"reflect"
	"time"
	//"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment            string        `mapstructure:"ENVIRONMENT"`
	HTTPServerAddress      string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration   time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmpSignupTokenDuration time.Duration `mapstructure:"EMP_SIGNUP_TOKEN_DURATION"`
	EmailSenderName        string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress     string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword    string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

var ServerConfig Config

func LoadConfig() error {

	ATD, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	RTD, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	ESTD, _ := time.ParseDuration(os.Getenv("EMP_SIGNUP_TOKEN_DURATION"))

	ServerConfig = Config{
		Environment:            os.Getenv("ENVIRONMENT"),
		HTTPServerAddress:      os.Getenv("HTTP_SERVER_ADDRESS"),
		TokenSymmetricKey:      os.Getenv("TOKEN_SYMMETRIC_KEY"),
		AccessTokenDuration:    ATD,
		RefreshTokenDuration:   RTD,
		EmpSignupTokenDuration: ESTD,
		EmailSenderName:        os.Getenv("EMAIL_SENDER_NAME"),
		EmailSenderAddress:     os.Getenv("EMAIL_SENDER_ADDRESS"),
		EmailSenderPassword:    os.Getenv("EMAIL_SENDER_PASSWORD"),
	}
	st := reflect.TypeOf(ServerConfig)

	// Iterate over the fields of the struct
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		value := reflect.ValueOf(ServerConfig).Field(i)
		fmt.Printf("%s: %v\n", field.Name, value.Interface())
	}
	return nil
}

/*
// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&ServerConfig)
	return
}
*/
