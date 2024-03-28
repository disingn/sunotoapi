package serve

import (
	"encoding/base64"
	"encoding/json"
	"fksunoapi/models"
	"fmt"
	"log"
	"strings"
	"time"
)

type Claims struct {
	Exp int64 `json:"exp"`
}

var Jwt string

func getLastUserContent(data models.OpenaiCompletionsData) string {
	var lastUserContent string
	for _, message := range data.Messages {
		if message.Role == "user" {
			lastUserContent = message.Content
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
