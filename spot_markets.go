package mexcapi

import (
	"bytes"
	"net/http"
)

type TickerInfoResponse struct {
	Code int `json:"code"`
	Data []struct {
		Symbol     string `json:"symbol"`
		Volume     string `json:"volume"`
		High       string `json:"high"`
		Low        string `json:"low"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Open       string `json:"open"`
		Last       string `json:"last"`
		Time       int64  `json:"time"`
		ChangeRate string `json:"change_rate"`
	} `json:"data"`
}

func (c *Client) GetSpotTickerInfo(symbol string) (*TickerInfoResponse, error) {
	var buffer bytes.Buffer
	buffer.WriteString("open/api/v2/market/ticker?symbol=")
	buffer.WriteString(symbol)
	res, err := c.do("spot", http.MethodGet, buffer.String(), nil, false, false)
	if err != nil {
		return nil, err
	}
	result := new(TickerInfoResponse)
	err = json.Unmarshal(res, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
