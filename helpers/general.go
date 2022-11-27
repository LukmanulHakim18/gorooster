package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
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

type RedisSetup struct {
	Host     string
	SelectDB int
	Password string
}

// Get Default falue from redis setup
// If .env not found
func GetRedisSetup() RedisSetup {
	godotenv.Load()
	rs := RedisSetup{
		Host:     "localhost:6379",
		SelectDB: 14,
		Password: "",
	}

	if host := os.Getenv("REDIS_SERVER_IP"); host != "" {
		rs.Host = host
	}
	if db := os.Getenv("REDIS_SELECT_DB"); db != "" {
		rs.SelectDB, _ = strconv.Atoi(db)
	}
	if password := os.Getenv("REDIS_SERVER_PASSWORD"); password != "" {
		rs.Password = password
	}
	return rs
}

// Generate channel name use db selected
// This channel listening key expired
// And triger event by key
func (rs RedisSetup) GenerateKeyEventChannel() string {
	return fmt.Sprintf("__keyevent@%d__:expired", rs.SelectDB)
}

func GetZapLoggerSetup() zap.Config {
	jsonFile, err := os.Open("logger/config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	rawJSON := []byte(byteValue)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	return cfg
}
