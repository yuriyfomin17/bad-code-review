package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/yuriyfomin17/bad-code-review/common"
)

type Config struct {
	Port                       string
	NumOfWorkers               int
	BatchNumOrdersIdsToProcess int
	HttpClientTimeoutSeconds   int
}

func Read() (Config, error) {
	var config Config

	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if exists {
		config.Port = serverPort
	}

	numOfWorkers, err := getIntEnvByKey("NUM_OF_WORKERS")
	switch {
	case errors.Is(err, common.ErrEnvVarNotSet):
		config.NumOfWorkers = 5
	case err != nil:
		return Config{}, fmt.Errorf("NUM_OF_WORKERS env variable is not a number")
	default:
		config.NumOfWorkers = numOfWorkers
	}

	batchNumOrdersIdsToProcess, err := getIntEnvByKey("BATCH_NUM_ORDERS_IDS_TO_PROCESS")

	switch {
	case errors.Is(err, common.ErrEnvVarNotSet):
		config.BatchNumOrdersIdsToProcess = 10
	case err != nil:
		return Config{}, fmt.Errorf("BATCH_NUM_ORDERS_IDS_TO_PROCESS env variable is not a number")
	default:
		config.BatchNumOrdersIdsToProcess = batchNumOrdersIdsToProcess
	}

	httpClientTimeoutSeconds, err := getIntEnvByKey("HTTP_CLIENT_TIMEOUT_SECONDS")

	switch {
	case errors.Is(err, common.ErrEnvVarNotSet):
		config.HttpClientTimeoutSeconds = 5
	case err != nil:
		return Config{}, fmt.Errorf("HTTP_CLIENT_TIMEOUT_SECONDS env variable is not a number")
	case httpClientTimeoutSeconds < 1:
		return Config{}, fmt.Errorf("HTTP_CLIENT_TIMEOUT_SECONDS env variable must be greater than 0")
	default:
		config.HttpClientTimeoutSeconds = httpClientTimeoutSeconds
	}

	return config, nil
}

func getIntEnvByKey(key string) (int, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return 0, common.ErrEnvVarNotSet
	}
	numEnv, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s env variable is not a number", key)
	}
	return numEnv, nil
}
