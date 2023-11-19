package main

import "os"

type ServerArgs struct {
	Username    string
	Password    string
	Auth        string
	Port        string
	Directory   string
	TLSCertFile string
	TLSKeyFile  string
}

func GetEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}

func NewServerArgs() ServerArgs {
	args := ServerArgs{}

	args.Username = GetEnvOrDefault("SCAR_USERNAME", "admin")
	args.Password = GetEnvOrDefault("SCAR_PASSWORD", "spectro")
	args.Auth = GetEnvOrDefault("SCAR_AUTH", "none")
	args.Port = GetEnvOrDefault("SCAR_PORT", "8080")
	args.Directory = GetEnvOrDefault("SCAR_DIRECTORY", "./public")
	args.TLSCertFile = GetEnvOrDefault("SCAR_TLS_CERT", "")
	args.TLSKeyFile = GetEnvOrDefault("SCAR_TLS_KEY", "")

	return args
}
