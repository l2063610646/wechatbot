package service

import (
	"encoding/json"
	"github.com/869413421/wechatbot/config"
	"github.com/869413421/wechatbot/gtp"
	"github.com/patrickmn/go-cache"
	"strings"
	"time"
	"unicode/utf8"
)

// UserServiceInterface 用户业务接口
type UserServiceInterface interface {
	GetUserSessionContext(userId string, question string) string
	SetUserSessionContext(userId string, question, reply string)
	ClearUserSessionContext(userId string, msg string) bool
}

var _ UserServiceInterface = (*UserService)(nil)

// UserService 用戶业务
type UserService struct {
	// 缓存
	cache *cache.Cache
}

func (s *UserService) GetUserSessionContext(userId string, question string) string {
	if config.LoadConfig().UseGPT {
		sessionContext, _ := s.cache.Get(userId)
		msg := "[]"
		if sessionContext != nil || sessionContext == "" {
			msg = sessionContext.(string)
		}
		var message []gtp.ChatGPTMessage
		_ = json.Unmarshal([]byte(msg), &message)
		message = append(message, gtp.ChatGPTMessage{
			Role:    "user",
			Content: question,
		})
		fin, _ := json.Marshal(message)
		return string(fin)
	} else {
		sessionContext, ok := s.cache.Get(userId)
		if !ok {
			return ""
		}
		return sessionContext.(string) + "Human:" + question + "\nAI:"
	}
}

func (s *UserService) SetUserSessionContext(userId string, question, reply string) {
	if config.LoadConfig().UseGPT {
		sessionContext, _ := s.cache.Get(userId)
		msg := "[]"
		if sessionContext != nil || sessionContext == "" {
			msg = sessionContext.(string)
		}
		if msg == "" {
			msg = "[]"
		}
		var message []gtp.ChatGPTMessage
		_ = json.Unmarshal([]byte(msg), &message)
		message = append(message, gtp.ChatGPTMessage{
			Role:    "user",
			Content: question,
		})
		message = append(message, gtp.ChatGPTMessage{
			Role:    "assistant",
			Content: reply,
		})
		fin, _ := json.Marshal(message)
		s.cache.Set(userId, string(fin), time.Second*config.LoadConfig().SessionTimeout)
	} else {
		sessionContext, _ := s.cache.Get(userId)
		question := sessionContext.(string)
		value := question + reply + "\n"
		s.cache.Set(userId, value, time.Second*config.LoadConfig().SessionTimeout)
	}
}

// ClearUserSessionContext 清空GTP上下文，接收文本中包含`我要问下一个问题`，并且Unicode 字符数量不超过20就清空
func (s *UserService) ClearUserSessionContext(userId string, msg string) bool {
	if strings.Contains(msg, "过") && utf8.RuneCountInString(msg) < 20 {
		s.cache.Delete(userId)
		return true
	}
	return false
}

// NewUserService 创建新的业务层
func NewUserService() UserServiceInterface {
	return &UserService{cache: cache.New(time.Second*config.LoadConfig().SessionTimeout, time.Minute*10)}
}
