package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/skip2/go-qrcode"
)

const (
	// 微信开放平台API地址
	apiBaseURL = "https://api.weixin.qq.com"
	// 获取二维码的API路径
	qrcodePath = "/sns/oauth2/access_token"
	// 获取用户信息的API路径
	userInfoPath = "/sns/userinfo"
	// 二维码存储路径
	qrcodeDir = "static/qrcodes"
)

// Client 微信客户端
type Client struct {
	appID      string
	appSecret  string
	httpClient *http.Client
}

// NewClient 创建新的微信客户端
func NewClient(appID, appSecret string) *Client {
	// 确保二维码存储目录存在
	if err := os.MkdirAll(qrcodeDir, 0755); err != nil {
		panic(fmt.Sprintf("创建二维码存储目录失败: %v", err))
	}

	return &Client{
		appID:     appID,
		appSecret: appSecret,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// QRCodeResponse 二维码响应
type QRCodeResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}

// UserInfo 用户信息
type UserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
}

// GetQRCode 获取微信登录二维码
func (c *Client) GetQRCode() (string, error) {
	// 生成唯一的文件名
	filename := fmt.Sprintf("wx_login_%d.png", time.Now().UnixNano())
	filepath := filepath.Join(qrcodeDir, filename)

	// 生成授权URL
	authURL := fmt.Sprintf("https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=STATE#wechat_redirect",
		c.appID,
		"http://your-domain.com/api/users/wx-callback", // 需要替换为实际的回调地址
	)

	// 生成二维码
	err := qrcode.WriteFile(authURL, qrcode.Medium, 256, filepath)
	if err != nil {
		return "", fmt.Errorf("生成二维码失败: %v", err)
	}

	// 返回二维码URL
	return fmt.Sprintf("/static/qrcodes/%s", filename), nil
}

// CheckLogin 检查微信登录状态
func (c *Client) CheckLogin(code string) (*QRCodeResponse, error) {
	url := fmt.Sprintf("%s%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		apiBaseURL, qrcodePath, c.appID, c.appSecret, code)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求微信API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result QRCodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// GetUserInfo 获取用户信息
func (c *Client) GetUserInfo(accessToken, openID string) (*UserInfo, error) {
	url := fmt.Sprintf("%s%s?access_token=%s&openid=%s",
		apiBaseURL, userInfoPath, accessToken, openID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求微信API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result UserInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}
