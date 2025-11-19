package app

const (
	Dev     Environment = "dev"
	Stage   Environment = "stage"
	Acc     Environment = "acc"
	Sandbox Environment = "sandbox"
	Prod    Environment = "prod"
)

type Environment string

type Configuration struct {
	Environment Environment
	LogLevel    string
	HTTPPort    string
}
