package mexcapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/google/go-querystring/query"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	key, secret string
	subaccount  string
	client      *http.Client
	window      int
}

func New(key, secret, subaccount string) *Client {
	hc := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &Client{
		key:        key,
		secret:     secret,
		subaccount: subaccount,
		client:     hc,
		window:     5000,
	}
}

// in milliseconds, default is 5000
func (c *Client) SetRecvWindow(recvWindow int) {
	c.window = recvWindow
}

func (c *Client) do(product, method, path string, data interface{}, sign bool, stream bool) (response []byte, err error) {
	var ENDPOINT string
	switch product {
	case "spot":
		ENDPOINT = "https://www.mexc.com"
	}
	values, err := query.Values(data)
	if err != nil {
		return nil, err
	}
	payload := values.Encode()
	if sign {
		payload = fmt.Sprintf("%s&timestamp=%v&recvWindow=%d", payload, time.Now().UnixNano()/(1000*1000), c.window)
		mac := hmac.New(sha256.New, []byte(c.secret))
		_, err = mac.Write([]byte(payload))
		if err != nil {
			return nil, err
		}
		payload = fmt.Sprintf("%s&signature=%s", payload, hex.EncodeToString(mac.Sum(nil)))
	}
	var req *http.Request
	if method == http.MethodGet {
		var url string
		if payload == "" {
			url = fmt.Sprintf("%s/%s", ENDPOINT, path)
		} else {
			url = fmt.Sprintf("%s/%s?%s", ENDPOINT, path, payload)
		}
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, fmt.Sprintf("%s/%s", ENDPOINT, path), strings.NewReader(payload))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if sign || stream {
		req.Header.Add("X-MBX-APIKEY", c.key)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d: %v", resp.StatusCode, string(response))
	}
	return response, err
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
