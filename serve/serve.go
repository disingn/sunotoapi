package serve

import (
	"bytes"
	"encoding/json"
	"fksunoapi/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ErrCodeRequestFailed   = 1001
	ErrCodeResponseInvalid = 1002
	ErrCodeJsonFailed      = 1003
	ErrCodeTimeout         = 1004
)

var (
	SessionExp int64
	Session    string
)

type ErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func NewErrorResponse(errorCode int, errorMsg string) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
	}
}

func NewErrorResponseWithError(errorCode int, err error) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode: errorCode,
		ErrorMsg:  err.Error(),
	}
}

func GetSession(c string) string {
	_url := "https://clerk.suno.ai/v1/client?_clerk_js_version=4.70.5"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, _url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "__client="+c)
	res, err := client.Do(req)
	if err != nil {
		log.Printf("GetSession failed, error: %v", err)
		return ""
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("GetSession failed, invalid status code: %d", res.StatusCode)
		return ""
	}
	body, _ := io.ReadAll(res.Body)
	var data models.GetSessionData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Printf("GetSession failed, json unmarshal error: %v", err)
		return ""
	}
	SessionExp = data.Response.Sessions[0].ExpireAt
	return data.Response.Sessions[0].Id
}

func GetJwtToken(c string) (string, *ErrorResponse) {
	if time.Now().After(time.Unix(SessionExp/1000, 0)) {
		Session = GetSession(c)
	}
	_url := fmt.Sprintf("https://clerk.suno.ai/v1/client/sessions/%s/tokens?_clerk_js_version=4.70.5", Session)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, _url, nil)

	if err != nil {
		log.Printf("GetJwtToken failed, error: %v", err)
		return "", NewErrorResponse(ErrCodeRequestFailed, "create request failed")
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "__client="+c)

	res, err := client.Do(req)
	if err != nil {
		log.Printf("GetJwtToken failed, error: %v", err)
		return "", NewErrorResponse(ErrCodeRequestFailed, "send request failed")
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		log.Printf("GetJwtToken failed, invalid status code: %d, response: %s", res.StatusCode, string(body))
		return "", NewErrorResponse(ErrCodeResponseInvalid, "invalid response")
	}

	var data models.GetTokenData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Printf("GetJwtToken failed, json unmarshal error: %v", err)
		return "", NewErrorResponse(ErrCodeJsonFailed, "parse response failed")
	}

	if len(data.Jwt) == 0 {
		log.Print("GetJwtToken failed, empty jwt token")
		return "", NewErrorResponse(ErrCodeResponseInvalid, "get empty jwt token")
	}
	return data.Jwt, nil
}

func sendRequest(url, method, c string, data []byte) ([]byte, *ErrorResponse) {
	jwt, errResp := GetJwtToken(c)
	if errResp != nil {
		errMsg := fmt.Sprintf("error getting JWT: %s", errResp.ErrorMsg)
		log.Printf("sendRequest failed, %s", errMsg)
		return nil, NewErrorResponse(errResp.ErrorCode, errMsg)
	}

	client := &http.Client{}
	var req *http.Request
	var err error
	if data != nil {
		req, err = http.NewRequest(method, url, bytes.NewReader(data))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		log.Printf("sendRequest failed, error creating request: %v", err)
		return nil, NewErrorResponseWithError(ErrCodeRequestFailed, err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Authorization", "Bearer "+jwt)

	res, err := client.Do(req)
	if err != nil {
		log.Printf("sendRequest failed, error sending request: %v", err)
		return nil, NewErrorResponseWithError(ErrCodeRequestFailed, err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		log.Printf("sendRequest failed, unexpected status code: %d, response body: %s", res.StatusCode, string(body))
		return body, NewErrorResponse(ErrCodeResponseInvalid, fmt.Sprintf("unexpected status code: %d, response body: %s", res.StatusCode, string(body)))
	}

	return body, nil
}

func V2Generate(d map[string]interface{}, c string) ([]byte, *ErrorResponse) {
	_url := "https://studio-api.suno.ai/api/generate/v2/"
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Printf("V2Generate failed, error marshalling request data: %v", err)
		return nil, NewErrorResponseWithError(ErrCodeJsonFailed, err)
	}
	body, errResp := sendRequest(_url, "POST", c, jsonData)
	if errResp != nil {
		return body, errResp
	}
	return body, nil
}

func V2GetFeedTask(ids, c string) ([]byte, *ErrorResponse) {
	ids = url.QueryEscape(ids)
	_url := "https://studio-api.suno.ai/api/feed/?ids=" + ids
	body, errResp := sendRequest(_url, "GET", c, nil)
	if errResp != nil {
		return body, errResp
	}
	return body, nil
}

func GenerateLyrics(d map[string]interface{}, c string) ([]byte, *ErrorResponse) {
	_url := "https://studio-api.suno.ai/api/generate/lyrics/"
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Printf("GenerateLyrics failed, error marshalling request data: %v", err)
		return nil, NewErrorResponseWithError(ErrCodeJsonFailed, err)
	}
	body, errResp := sendRequest(_url, "POST", c, jsonData)
	if errResp != nil {
		return body, errResp
	}
	return body, nil
}

func GetLyricsTask(ids, c string) ([]byte, *ErrorResponse) {
	_url := "https://studio-api.suno.ai/api/generate/lyrics/" + ids
	body, errResp := sendRequest(_url, "GET", c, nil)
	if errResp != nil {
		return body, errResp
	}
	return body, nil
}

func SunoChat(c map[string]interface{}, ck string) (interface{}, *ErrorResponse) {
	lastUserContent := getLastUserContent(c)
	d := map[string]interface{}{
		"mv":                     c["model"].(string),
		"gpt_description_prompt": lastUserContent,
		"prompt":                 "",
		"make_instrumental":      false,
	}
	body, errResp := V2Generate(d, ck)
	if errResp != nil {
		return nil, errResp
	}

	var v2GenerateData models.GenerateData
	if err := json.Unmarshal(body, &v2GenerateData); err != nil {
		log.Printf("SunoChat failed, error unmarshalling generate data: %v, response body: %s", err, string(body))
		return nil, NewErrorResponse(ErrCodeResponseInvalid, fmt.Sprintf("parse generate data failed, response body: %s", string(body)))
	}

	clipIds := make([]string, len(v2GenerateData.Clips))
	for i, clip := range v2GenerateData.Clips {
		clipIds[i] = clip.Id
	}
	ids := strings.Join(clipIds, ",")

	timeout := time.After(3 * time.Minute)
	tick := time.Tick(5 * time.Second)

	for {
		select {
		case <-timeout:
			return nil, NewErrorResponse(ErrCodeTimeout, "get feed task timeout")
		case <-tick:
			body, errResp = V2GetFeedTask(ids, ck)
			if errResp != nil {
				return nil, errResp
			}

			var v2GetFeedData []map[string]interface{}
			if err := json.Unmarshal(body, &v2GetFeedData); err != nil {
				log.Printf("SunoChat failed, error unmarshalling feed data: %v, response body: %s", err, string(body))
				return nil, NewErrorResponse(ErrCodeResponseInvalid, fmt.Sprintf("parse feed data failed, response body: %s", string(body)))
			}

			allComplete := true
			for _, data := range v2GetFeedData {
				if data["status"] != "complete" {
					allComplete = false
					break
				}
			}

			if allComplete {
				var markdown strings.Builder
				markdown.WriteString(fmt.Sprintf("# %s\n\n", v2GetFeedData[0]["title"]))
				markdown.WriteString(fmt.Sprintf("%s\n\n", v2GetFeedData[0]["metadata"].(map[string]interface{})["prompt"]))
				markdown.WriteString(fmt.Sprintf("## 版本一\n\n"))
				markdown.WriteString(fmt.Sprintf("视频链接：%s\n\n", v2GetFeedData[0]["video_url"]))
				markdown.WriteString(fmt.Sprintf("音频链接：%s\n\n", v2GetFeedData[0]["audio_url"]))
				markdown.WriteString(fmt.Sprintf("图片链接：%s\n\n", v2GetFeedData[0]["image_large_url"]))
				markdown.WriteString(fmt.Sprintf("## 版本二\n\n"))
				markdown.WriteString(fmt.Sprintf("视频链接：%s\n\n", v2GetFeedData[1]["video_url"]))
				markdown.WriteString(fmt.Sprintf("音频链接：%s\n\n", v2GetFeedData[1]["audio_url"]))
				markdown.WriteString(fmt.Sprintf("图片链接：%s\n\n", v2GetFeedData[1]["image_large_url"]))

				response := map[string]interface{}{
					"choices": []map[string]interface{}{
						{
							"finish_reason": "stop",
							"index":         0,
							"message": map[string]string{
								"content": markdown.String(),
								"role":    "assistant",
							},
							"logprobs": nil,
						},
					},
					"created": time.Now().Unix(),
					"id":      "chatcmpl-7QyqpwdfhqwajicIEznoc6Q47XAyW",
					"model":   c["model"].(string),
					"object":  "chat.completion",
					"usage": map[string]int{
						"completion_tokens": 17,
						"prompt_tokens":     57,
						"total_tokens":      74,
					},
				}

				return response, nil
			}
		}
	}
}
