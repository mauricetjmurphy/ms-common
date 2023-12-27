package utils

import (
	"encoding/json"
	"os"
	"strings"
)

func IsLambdaEnv() bool {
	return len(os.Getenv("LAMBDA_TASK_ROOT")) > 0
}

func ParseArgs() (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if len(os.Args) > 1 {
		// Accept raw json from argument
		// {\"id\": 1, \"name\": \"arg\"}
		if len(os.Args[1:]) == 1 {
			var raw = os.Args[1:][0]
			err := json.Unmarshal([]byte(raw), &params)
			return params, err
		}
		// Else, Accept key = value argument
		// id = 2 name = arg
		for _, a := range os.Args[1:] {
			kv := strings.Split(a, "=")
			if len(kv) == 2 {
				k := kv[0]
				v := kv[1]
				params[k] = v
			}
		}
	}
	return params, nil
}
