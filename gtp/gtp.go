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
    Index        int         `json:"index"`
    Message      RoleContent `json:"message"`
    FinishReason string      `json:"finish_reason"`
}

type RoleContent struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
    Model            string        `json:"model"`
    Prompt           string        `json:"prompt,omitempty"`
    MaxTokens        int           `json:"max_tokens"`
    Temperature      float32       `json:"temperature"`
    TopP             int           `json:"top_p"`
    FrequencyPenalty int           `json:"frequency_penalty"`
    PresencePenalty  float32       `json:"presence_penalty"`
    Stop             []string      `json:"stop"`
    Messages         []RoleContent `json:"messages,omitempty"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
    requestBody := ChatGPTRequestBody{
        Model:            "gpt-3.5-turbo",
        Temperature:      0.8, // 每次返回的答案的相似度0-1（0：每次都一样，1：每次都不一样）
        MaxTokens:        4000,
        TopP:             1,
        FrequencyPenalty: 0,
        PresencePenalty:  0.6,
        Stop:             []string{" Human:", " AI:"},
    }
    if config.LoadConfig().Model != "" {
        requestBody.Model = config.LoadConfig().Model
    }
    if requestBody.Model == "text-davinci-003" {
        requestBody.Prompt = msg
    } else {
        requestBody.Messages = []RoleContent{
            {"user", msg},
        }
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
    req, err := http.NewRequest("POST", base+"v1/chat/completions", bytes.NewBuffer(requestData))
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

    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
        log.Printf("read response body error : %s, %v", string(body), err)
        return "", err
    }

    if response.StatusCode != http.StatusOK {
        log.Printf("request gtp error : %s %s", response.Status, string(body))
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
            if v["text"] != nil {
                reply = v["text"].(string)
            } else if v["message"] != nil {
                v2 := v["message"]
                msg := v2.(map[string]interface{})
                reply = msg["content"].(string)
            }
            break
        }
    }
    //log.Printf("gpt response text: %s \n", reply)
    return reply, nil
}
