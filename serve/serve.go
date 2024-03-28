package serve

import (
	"bytes"
	"encoding/json"
	"errors"
	"fksunoapi/cfg"
	"fksunoapi/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	SessionExp int64
	Session    string
)

func GetSession() string {
	_url := "https://clerk.suno.ai/v1/client?_clerk_js_version=4.70.5"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, _url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "__client="+cfg.Config.App.Client)
	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return ""
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Print("Error")
		return ""
	}
	body, _ := io.ReadAll(res.Body)
	var data models.GetSessionData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Print(err)
		return ""
	}
	SessionExp = data.Response.Sessions[0].ExpireAt
	return data.Response.Sessions[0].Id
}

func GetJwtToken() (string, error) {
	if time.Now().After(time.Unix(SessionExp, 0)) {
		Session = GetSession()
	}
	_url := fmt.Sprintf("https://clerk.suno.ai/v1/client/sessions/%s/tokens?_clerk_js_version=4.70.5", Session)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, _url, nil)

	if err != nil {
		log.Print(err)
		return "", err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "__client="+cfg.Config.App.Client)

	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		//log.Print(string(body))
		return "", fmt.Errorf(string(body))
	}
	var data models.GetTokenData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Print(err)
		return "", err
	}
	//有效时间 1 分钟
	if len(data.Jwt) == 0 {
		log.Print("GetJwtToken: ", data.Jwt)
		return "", err
	}
	return data.Jwt, nil
}

func sendRequest(url, method string, data []byte) ([]byte, error) {
	jwt, err := IsJWTExpired()
	if err != nil {
		log.Println("Error getting JWT: ", err)
		return nil, err
	}

	client := &http.Client{}
	var req *http.Request
	if data != nil {
		req, err = http.NewRequest(method, url, bytes.NewReader(data))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		log.Print(err)
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Authorization", "Bearer "+jwt)

	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	return body, nil
}

func V2Generate(d map[string]interface{}) ([]byte, error) {
	_url := "https://studio-api.suno.ai/api/generate/v2/"
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Fatalf("Error marshalling request data: %v", err)
		return nil, err
	}
	body, err := sendRequest(_url, "POST", jsonData)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func V2GetFeedTask(ids string) ([]byte, error) {
	ids = url.QueryEscape(ids)
	_url := "https://studio-api.suno.ai/api/feed/?ids=" + ids
	body, err := sendRequest(_url, "GET", nil)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GenerateLyrics(d map[string]interface{}) ([]byte, error) {
	_url := "https://studio-api.suno.ai/api/generate/lyrics/"
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Fatalf("Error marshalling request data: %v", err)
		return nil, err
	}
	body, err := sendRequest(_url, "POST", jsonData)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetLyricsTask(ids string) ([]byte, error) {
	_url := "https://studio-api.suno.ai/api/generate/lyrics/" + ids
	body, err := sendRequest(_url, "GET", nil)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func SunoChat(c models.OpenaiCompletionsData) (interface{}, error) {
	lastUserContent := getLastUserContent(c)
	d := map[string]interface{}{
		"mv":                     c.Model,
		"gpt_description_prompt": lastUserContent,
		"prompt":                 "",
		"make_instrumental":      false,
	}
	body, err := V2Generate(d)
	var v2GenerateData models.GenerateData
	if err = json.Unmarshal(body, &v2GenerateData); err != nil {
		log.Print(err)
		return nil, err
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
			return nil, errors.New("timeout exceeded")
		case <-tick:
			body, err = V2GetFeedTask(ids)
			if err != nil {
				return nil, err
			}

			var v2GetFeedData []map[string]interface{}
			if err = json.Unmarshal(body, &v2GetFeedData); err != nil {
				log.Print(err)
				return nil, err
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
					"model":   c.Model,
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
