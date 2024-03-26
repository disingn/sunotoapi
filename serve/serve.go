package serve

import (
	"bytes"
	"encoding/json"
	"fksunoapi/cfg"
	"fksunoapi/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	SessionExp int64
	Session    string
)

func GetSession() string {
	url := "https://clerk.suno.ai/v1/client?_clerk_js_version=4.70.5"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
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

func GetJwtToken() string {
	if time.Now().After(time.Unix(SessionExp, 0)) {
		Session = GetSession()
	}
	url := fmt.Sprintf("https://clerk.suno.ai/v1/client/sessions/%s/tokens?_clerk_js_version=4.70.5", Session)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Print(err)
		return ""
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "__client="+cfg.Config.App.Client)

	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return ""
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		log.Print("err body: ", string(body))
		return ""
	}
	var data models.GetTokenData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Print(err)
		return ""
	}
	//有效时间 1 分钟
	if len(data.Jwt) == 0 {
		log.Print("GetJwtToken: ", data.Jwt)
		return ""
	}
	return data.Jwt
}

func V2Generate(d models.GenerateCreateData) *models.GenerateData {
	jwt := IsJWTExpired(Jwt)
	url := "https://studio-api.suno.ai/api/generate/v2/"
	method := "POST"
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Fatalf("Error marshalling request data: %v", err)
		return nil
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonData))
	if err != nil {
		log.Print(err)
		return nil
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Authorization", "Bearer "+jwt)
	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		log.Print("create err body: ", string(body))
		return nil
	}
	var data models.GenerateData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Print(err)
		return nil
	}
	return &data
}

func V2GetFeedJop(ids string) []byte {
	jwt := IsJWTExpired(Jwt)
	ids = url.QueryEscape(ids)
	_url := "https://studio-api.suno.ai/api/feed/?ids=" + ids
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, _url, nil)
	if err != nil {
		log.Print(err)
		return nil
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Authorization", "Bearer "+jwt)
	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		log.Print("get err body: ", string(body))
		return nil
	}
	return body
}
