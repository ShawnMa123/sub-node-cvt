package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

var (
	oauthConf *oauth2.Config
	// 在生产环境中，您应该为每个授权请求生成一个唯一的、随机的状态字符串，
	// 并将其存储在会话中进行验证，以防止 CSRF 攻击。
	// 为简化示例，我们使用一个静态字符串。
	oauthStateString = "pseudo-random"
)

// InitOAuth 使用环境变量初始化 OAuth 配置
func InitOAuth() {
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"gist"},
		Endpoint:     githuboauth.Endpoint,
	}
}

// HandleGitHubLogin 将用户重定向到 GitHub 授权页面
func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleGitHubCallback 处理从 GitHub 返回的回调请求
func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// 将访问令牌存储在 HttpOnly、Secure 的 cookie 中
	http.SetCookie(w, &http.Cookie{
		Name:     "github_token",
		Value:    token.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // 如果是 HTTPS 连接，则设置 Secure 标志
		MaxAge:   3600,         // 1 小时
		SameSite: http.SameSiteLaxMode,
	})

	// 重定向回主页
	http.Redirect(w, r, "/", http.StatusFound)
}

// HandleCreateGist 使用用户的 token 创建一个新的私密 Gist
func HandleCreateGist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("github_token")
	if err != nil {
		http.Error(w, "Unauthorized: Please login with GitHub first.", http.StatusUnauthorized)
		return
	}
	accessToken := cookie.Value

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestData struct {
		Description string `json:"description"`
		Content     string `json:"content"`
		Filename    string `json:"filename"`
	}
	if err := json.Unmarshal(body, &requestData); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	isPublic := false
	gistToCreate := &github.Gist{
		Description: &requestData.Description,
		Public:      &isPublic,
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(requestData.Filename): {
				Content: &requestData.Content,
			},
		},
	}

	createdGist, _, err := client.Gists.Create(ctx, gistToCreate)
	if err != nil {
		log.Printf("Failed to create gist: %v", err)
		http.Error(w, "Failed to create Gist on GitHub", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created a new secret gist for a user at %s\n", *createdGist.HTMLURL)

	var rawURL string
	for _, file := range createdGist.Files {
		if file.RawURL != nil && *file.RawURL != "" {
			rawURL = *file.RawURL
			break
		}
	}

	if rawURL == "" {
		http.Error(w, "Could not find raw_url in the created Gist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"raw_url": rawURL})
}

// HandleUserInfo 使用用户的 token 获取其 GitHub 信息
func HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("github_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	accessToken := cookie.Value

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
