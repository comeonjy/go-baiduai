package body

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/comeonjy/go-baiduai/lib"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Interface .
type Interface interface{}

// Reply .
type Reply struct {
	BaseReply
	Interface
}

// BaseReply .
type BaseReply struct {
	ErrorCode int64  `json:"error_code"`
	LogId     int64  `json:"log_id"`
	ErrorMsg  string `json:"error_msg"`
}

// Body .
type Body struct {
	Token *lib.AccessToken
}

// New .
func New(appKey string, appSecret string, store lib.Storage) *Body {
	return &Body{
		Token: lib.NewToken(appKey, appSecret, store),
	}
}

// PostJson .
func (b *Body) PostJson(url string, v []byte, res *Reply) error {
	if err := b.Token.SetAccessToken(); err != nil {
		return errors.WithStack(err)
	}
	resp, err := http.Post(url+"?access_token="+b.Token.AccessToken, "application/json", bytes.NewReader(v))
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := json.Unmarshal(body, res); err != nil {
		return errors.WithStack(err)
	}

	if res.ErrorCode != 0 {
		return errors.New(fmt.Sprintf("%d:%s", res.ErrorCode, res.ErrorMsg))
	}
	return nil
}

// PostForm .
func (b *Body) PostForm(url string, v url.Values, res interface{}) error {
	if err := b.Token.SetAccessToken(); err != nil {
		return errors.WithStack(err)
	}

	v.Add("access_token", b.Token.AccessToken)
	resp, err := http.PostForm(url, v)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := json.Unmarshal(body, res); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *Body) checkErr(baseReply BaseReply, err error) error {
	if err != nil {
		return err
	}
	if baseReply.ErrorCode != 0 {
		return errors.New(fmt.Sprintf("%d:%s", baseReply.ErrorCode, baseReply.ErrorMsg))
	}
	return nil
}
