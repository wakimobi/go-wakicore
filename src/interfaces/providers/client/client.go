package client

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/wakimobi/go-wakicore/src/domain/products"
	"github.com/wakimobi/go-wakicore/src/utils/http_utils"
)

type Client struct {
	Product *products.Product
}

func (c *Client) Postback() {
	urlPostback := c.Product.UrlPostback
	headers := map[string]string{
		"ContentType": "application/json",
	}

	queryParameters := url.Values{}
	queryParameters.Add("a", "a")

	var response map[string]interface{}

	response, err := http_utils.MakeHTTPRequest(urlPostback, "GET", headers, queryParameters, nil, response)
	if err != nil {
		panic(err)
	}
	fmt.Printf("response: %+v", response)
}

func (c *Client) NotifSub() {
	urlNotifSub := c.Product.UrlNotifSub
	headers := map[string]string{
		"ContentType": "application/json",
	}

	// the query parameters to pass
	queryParameters := url.Values{}

	// the body to pass
	body := bytes.NewBufferString(`{"name":"test"}`)

	// the type to unmarshal the response into
	var response map[string]interface{}

	// call the function
	response, err := http_utils.MakeHTTPRequest(urlNotifSub, "POST", headers, queryParameters, body, response)
	if err != nil {
		panic(err)
	}

	// do something with the response
	fmt.Printf("response: %+v", response)
}
