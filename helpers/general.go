package helpers

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

// Generate Data-Key from Event-Key
// Ex: event:foo -> data:foo
func GetDataKey(eventKey string) (dataKey string, err error) {
	sliceKey := strings.Split(eventKey, ":")
	if len(sliceKey) != 2 {
		return dataKey, fmt.Errorf("event key not valid : %s len %d", eventKey, len(sliceKey))
	}
	return fmt.Sprintf("data:%s", sliceKey[1]), err
}

// Generate channel name use db selected
// This channel listening key expired
// And triger event by key

func GetZapLoggerSetup() zap.Config {
	cfg := zap.NewProductionConfig()
	if path := os.Getenv("LOG_PATH"); path != "" {
		cfg.OutputPaths = []string{
			"stdout",
			path,
		}
	}
	return cfg
}
