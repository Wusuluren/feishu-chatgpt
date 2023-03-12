package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	maxTokens   = 2000
	temperature = 0.7
	engine      = "gpt-3.5-turbo"
)

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatGPTChoiceItem    `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}
type ChatGPTChoiceItem struct {
	Message      Messages `json:"message"`
	Index        int      `json:"index"`
	FinishReason string   `json:"finish_reason"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string     `json:"model"`
	Messages         []Messages `json:"messages"`
	MaxTokens        int        `json:"max_tokens"`
	Temperature      float32    `json:"temperature"`
	TopP             int        `json:"top_p"`
	FrequencyPenalty int        `json:"frequency_penalty"`
	PresencePenalty  int        `json:"presence_penalty"`
}

//func (gpt ChatGPT) Completions(msg []Messages) (resp Messages, err error) {
//	requestBody := ChatGPTRequestBody{
//		Model:            engine,
//		Messages:         msg,
//		MaxTokens:        maxTokens,
//		Temperature:      temperature,
//		TopP:             1,
//		FrequencyPenalty: 0,
//		PresencePenalty:  0,
//	}
//	gptResponseBody := &ChatGPTResponseBody{}
//	err = gpt.sendRequestWithBodyType(gpt.ApiUrl+"/v1/chat/completions", "POST",
//		jsonBody,
//		requestBody, gptResponseBody)
//
//	if err == nil && len(gptResponseBody.Choices) > 0 {
//		resp = gptResponseBody.Choices[0].Message
//	} else {
//		resp = Messages{}
//		err = errors.New("openai 请求失败")
//	}
//	return resp, err
//}

// curl -XPOST https://ai117.com/ -d '{"msg":"11","token":"","style":"0"}' -H "content-type: application/json"
// {"code":200,"msg":"success","data":["11","\n\nI'm sorry"]}
func (gpt ChatGPT) Completions(msg []Messages) (resp Messages, err error) {
	type MyReq struct {
		Msg   string `json:"msg"`
		Token string `json:"token"`
		Style string `json:"style"`
	}
	type MyResp struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Data []string `json:"data"`
	}

	contentReq := MyReq{
		Msg:   msg[len(msg)-1].Content,
		Token: "",
		Style: "0",
	}
	bReq, _ := json.Marshal(contentReq)
	rsp, err := http.Post("https://ai117.com/", "application/json", bytes.NewReader(bReq))
	bResp, _ := ioutil.ReadAll(rsp.Body)
	var contentResp MyResp
	json.Unmarshal(bResp, &contentResp)
	fmt.Printf("req:%s, resp:%s, err:%v", string(bReq), string(bResp), err)
	if err == nil && len(contentResp.Data) > 1 {
		resp.Role = msg[len(msg)-1].Role
		resp.Content = contentResp.Data[1]
	} else {
		resp = Messages{}
		err = errors.New("openai 请求失败")
	}
	return resp, err
}
