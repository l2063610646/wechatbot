package gtp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/869413421/wechatbot/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const BASEURL = "https://api.openai.com/v1/"

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
	Text         string            `json:"text"`
	Index        int               `json:"index"`
	Logprobs     int               `json:"logprobs"`
	FinishReason string            `json:"finish_reason"`
	Message      ChoiceItemMessage `json:"message"`
}

type ChoiceItemMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}

type ChatGPTRequestBody struct {
	Messages    []ChatGPTMessage `json:"messages"`
	Temperature float64          `json:"temperature"`
	MaxTokens   int              `json:"max_tokens"`
	TopP        int              `json:"top_p"`
	Model       string           `json:"model"`
}

type ChatGPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func Completions(msg string) (string, error) {
	if config.LoadConfig().UseGPT {
		//使用gptTurbo的模型
		return ChatGPTCompletions(msg)
	} else {
		//使用text-davinci-003的模型
		return OpenAICompletions(msg)
	}
}

// OpenAICompletions gtp文本模型回复
// curl https://api.openai.com/v1/completions
// -H "Content-Type: application/json"
// -H "Authorization: Bearer your chatGPT key"
// -d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func OpenAICompletions(msg string) (string, error) {
	requestBody := OpenAIRequestBody{
		Model:            "text-davinci-003",
		Prompt:           msg,
		MaxTokens:        1024,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL+"completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	if config.LoadConfig().UseProxy {
		log.Printf("已开启代理 : %v", config.LoadConfig().ProxyUrl)
		proxy, _ := url.Parse(config.LoadConfig().ProxyUrl)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("gtp api status code not equals 200,code is %d", response.StatusCode))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
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
		reply = gptResponseBody.Choices[0].Text
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}

func ChatGPTCompletions(msg string) (string, error) {
	var message []ChatGPTMessage
	err := json.Unmarshal([]byte(msg), &message)
	requestBody := ChatGPTRequestBody{
		Model:       "gpt-3.5-turbo",
		Messages:    message,
		MaxTokens:   1024,
		Temperature: 0.7,
		TopP:        1,
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL+"chat/completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	if config.LoadConfig().UseProxy {
		log.Printf("已开启代理 : %v", config.LoadConfig().ProxyUrl)
		proxy, _ := url.Parse(config.LoadConfig().ProxyUrl)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("gtp api status code not equals 200,code is %d", response.StatusCode))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
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
		reply = gptResponseBody.Choices[0].Message.Content
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
