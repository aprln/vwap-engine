package env

import (
	"os"
	"strconv"
	"strings"
)

const delimiter = "|"

func LoadEnvString(key, defVal string) string {
	val, found := os.LookupEnv(key)
	if !found {
		return defVal
	}

	return val
}

func LoadEnvStringSlice(key string, defVal []string) []string {
	val, found := os.LookupEnv(key)
	if !found {
		return defVal
	}

	return strings.Split(val, delimiter)
}

func MustLoadEnvPositiveInt(key string, defVal int) int {
	val, found := os.LookupEnv(key)
	if !found {
		return defVal
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic("invalid int value: " + val)
	}

	if intVal <= 0 {
		panic("invalid positive int value: " + val)
	}

	return intVal
}
