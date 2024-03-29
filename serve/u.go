package serve

import (
	"strings"
)

type Claims struct {
	Exp int64 `json:"exp"`
}

var Jwt string

func getLastUserContent(data map[string]interface{}) string {
	var lastUserContent string
	messages, ok := data["messages"].([]interface{})
	if !ok {
		return lastUserContent
	}

	for i := len(messages) - 1; i >= 0; i-- {
		message, ok := messages[i].(map[string]interface{})
		if !ok {
			continue
		}

		role, ok := message["role"].(string)
		if !ok {
			continue
		}

		if role == "user" {
			lastUserContent, _ = message["content"].(string)
			break
		}
	}

	return lastUserContent
}

func ParseToken(authorizationHeader string) string {
	if authorizationHeader == "" {
		return ""
	}
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}

	return ""
}
