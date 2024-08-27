package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/casbin/caswaf/conf"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/google/uuid"
)

type verifyResponse struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Sub    string      `json:"sub"`
	Name   string      `json:"name"`
	Data   bool        `json:"data"`
	Data2  interface{} `json:"data2"`
}

var verifiedSession = make(map[string]time.Time)

func redirectToCaptcha(w http.ResponseWriter, r *http.Request) {
	scheme := getScheme(r)
	callbackUrl := fmt.Sprintf("%s://%s/caswaf-captcha-verify", scheme, r.Host)
	captchaUri := fmt.Sprintf(
		"%s/captcha?client_id=%s&redirect_uri=%s&state=%s",
		getCasdoorEndpoint(),
		conf.GetConfigString("clientId"),
		callbackUrl,
		conf.GetConfigString("casdoorApplication"),
	)
	http.Redirect(w, r, captchaUri, http.StatusFound)
}

func handleCaptchaCallback(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	code := r.URL.Query().Get("code")
	typeStr := r.URL.Query().Get("type")
	secret := r.URL.Query().Get("secret")
	applicationId := r.URL.Query().Get("applicationId")
	if code == "" || typeStr == "" || secret == "" || applicationId == "" {
		redirectToCaptcha(w, r)
	}

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	_ = writer.WriteField("captchaToken", code)
	_ = writer.WriteField("clientSecret", secret)
	_ = writer.WriteField("applicationId", applicationId)
	_ = writer.WriteField("captchaType", typeStr)
	_ = writer.Close()
	verifyURL := casdoorsdk.GetUrl("verify-captcha", nil)
	req, err := http.NewRequest("POST", verifyURL, &b)
	if err != nil {
		redirectToCaptcha(w, r)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		redirectToCaptcha(w, r)
	}
	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		redirectToCaptcha(w, r)
	}
	// parse response
	var vr verifyResponse
	err = json.Unmarshal(body, &vr)
	if err != nil {
		redirectToCaptcha(w, r)
	}
	if vr.Status != "ok" || !vr.Data {
		redirectToCaptcha(w, r)
	}
	// set verified session
	uuidStr := uuid.NewString()
	verifiedSession[uuidStr] = time.Now()
	cookie := &http.Cookie{
		Name:       "casdoor_captcha_token",
		Value:      uuidStr,
		Path:       "/",
		Domain:     host,
		Expires:    time.Now().Add(30 * time.Minute),
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   false,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}

	scheme := getScheme(r)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, scheme+"://"+host, http.StatusFound)
	return
}

func isVerifiedSession(r *http.Request) bool {
	cookie, err := r.Cookie("casdoor_captcha_token")
	if err != nil {
		return false
	}
	token := cookie.Value
	if token == "" {
		return false
	}
	t, ok := verifiedSession[token]
	if ok {
		if time.Now().Sub(t) < 30*time.Minute {
			return true
		}
		delete(verifiedSession, token)
	}
	return false
}
