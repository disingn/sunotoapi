package serve

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
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

func IsJWTExpired() (Jwt string, err error) {
	if Jwt == "" {
		Jwt, err = GetJwtToken()
		return
	}
	parts := strings.Split(Jwt, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid JWT format. Expected format: header.payload.signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Print(err)
		return "", err
	}
	var claims Claims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		log.Print(err)
		return "", err
	}
	expTime := time.Unix(claims.Exp, 0)
	if time.Now().After(expTime) {
		Jwt, err = GetJwtToken()
		return
	}
	return Jwt, nil
}
