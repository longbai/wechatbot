package gtp

import (
    "bytes"
    "encoding/json"
    "io"
    "log"
    "net/http"

    "github.com/longbai/wechatbot/config"
)

const BASEURL = "https://api.openai.com/"

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
    ID      string                   `json:"id"`
    Object  string                   `json:"object"`
    Created int                      `json:"created"`
    Model   string                   `json:"model"`
    Choices []map[string]interface{} `json:"choices"`
    Usage   map[string]interface{}   `json:"usage"`
}

type ChoiceItem struct {
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
    Model            string   `json:"model"`
    Prompt           string   `json:"prompt"`
    MaxTokens        int      `json:"max_tokens"`
    Temperature      float32  `json:"temperature"`
    TopP             int      `json:"top_p"`
    FrequencyPenalty int      `json:"frequency_penalty"`
    PresencePenalty  float32  `json:"presence_penalty"`
    Stop             []string `json:"stop"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
    requestBody := ChatGPTRequestBody{
        Model:            "text-davinci-003",
        Prompt:           msg,
        Temperature:      0.9, // 每次返回的答案的相似度0-1（0：每次都一样，1：每次都不一样）
        MaxTokens:        4000,
        TopP:             1,
        FrequencyPenalty: 0,
        PresencePenalty:  0.6,
        Stop:             []string{" Human:", " AI:"},
    }
    requestData, err := json.Marshal(requestBody)

    if err != nil {
        return "", err
    }
    log.Printf("request gtp json string : %v", string(requestData))
    base := BASEURL
    if config.LoadConfig().Proxy != "" {
        base = config.LoadConfig().Proxy
    }
    req, err := http.NewRequest("POST", base+"v1/completions", bytes.NewBuffer(requestData))
    if err != nil {
        return "", err
    }

    apiKey := config.LoadConfig().ApiKey
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)
    client := &http.Client{}
    response, err := client.Do(req)
    if err != nil {
        return "", err
    }
    if response.StatusCode != http.StatusOK {
        log.Printf("request gtp error : %v", response.Status)
        return "", err
    }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
        log.Printf("read response body error : %s, %v", string(body), err)
        return "", err
    }

    gptResponseBody := &ChatGPTResponseBody{}
    log.Println(string(body))
    err = json.Unmarshal(body, gptResponseBody)
    if err != nil {
        return "", err
    }
    var reply string
    if len(gptResponseBody.Choices) > 0 {
        for _, v := range gptResponseBody.Choices {
            reply = v["text"].(string)
            break
        }
    }
    log.Printf("gpt response text: %s \n", reply)
    return reply, nil
}
