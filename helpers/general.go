package helpers

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

// Generate Data-Key from Event-Key
// Ex: client:event:foo -> client:data:foo
func GetDataKey(eventKey string) (dataKey string, err error) {
	sliceKey := strings.Split(eventKey, ":")
	if len(sliceKey) != 3 {
		return dataKey, fmt.Errorf("event key not valid : %s len %d", eventKey, len(sliceKey))
	}
	return fmt.Sprintf("%s:data:%s", sliceKey[0], sliceKey[2]), err
}

// Setup zap.config and locaition file
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

func GenerateKeyEvent(clientName, key string) string {
	keyEvent := fmt.Sprintf("%s:event:%s", clientName, key)
	return keyEvent
}

func GenerateKeyData(clientName, key string) string {
	keyData := fmt.Sprintf("%s:data:%s", clientName, key)
	return keyData
}

func ValidatorClinetNameAndKey(str string) bool {
	return !strings.Contains(str, ":") && str != ""

}
