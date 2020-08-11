package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ea3hsp/go-ikea-server/pkg/tradfri"
	"github.com/go-kit/kit/log"
)

const (
	// Default config definitions
	defProcName    = "go-ikea-server"
	defClientID    = "TEST"
	defGatewayAddr = "192.168.0.15"
	defGatewayPort = "5684"
	defGatewayPsk  = "eJUYhaSMaz6W3qV7"
	// Environment variable names
	envProcName    = "PROCESS_NAME"
	envClientID    = "GATEWAY_CLIENT_ID"
	envGatewayAddr = "GATEWAY_IP_ADDR"
	envGatewayPort = "GATEWAY_UDP_PORT"
	envGatewayPsk  = "GATEWAT_PSK"
)

// config struct definition
type config struct {
	processName string
	clientID    string
	gwAddr      string
	gwPort      string
	gwPsk       string
}

// Run main func
func Run() {
	// parse os args
	cfg := loadConfig()
	// Creates logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	dtlsClient := tradfri.NewClient(fmt.Sprintf("%s:%s", cfg.gwAddr, cfg.gwPort), "Client_identity", cfg.gwPsk, logger)
	r, err := dtlsClient.AuthExchange(cfg.clientID)
	if err != nil {
		logger.Log("[error]", "auth exchanging", "msg", err.Error)
		os.Exit(1)
	}
	// response
	res := r.(map[string]interface{})
	//

	// wait until signal
	sig := WaitSignal()
	// exit banner
	logger.Log("[Info]", "Exit", "signal", sig.String())
}

// WaitSignal catching exit signal
func WaitSignal() os.Signal {
	ch := make(chan os.Signal, 2)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			return sig
		}
	}
}

// env get environment variable or fallback to default one
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// load config parameters
func loadConfig() *config {
	return &config{
		processName: env(envProcName, defProcName),
		clientID:    env(envClientID, defClientID),
		gwAddr:      env(envGatewayAddr, defGatewayAddr),
		gwPort:      env(envGatewayPort, defGatewayPort),
		gwPsk:       env(envGatewayPsk, defGatewayPsk),
	}
}
