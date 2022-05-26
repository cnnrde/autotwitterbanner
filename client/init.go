package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/dghubble/oauth1"
)

func CreateHttpClient() *http.Client {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))

	return config.Client(oauth1.NoContext, token)
}

type UpdateAccountBannerParams struct {
	Banner     string
	Width      string
	Height     string
	OffsetLeft string
	OffsetTop  string
}

func UpdateAccountBanner(params UpdateAccountBannerParams) {
	baseUrl := "https://api.twitter.com/1.1/account/update_profile_banner.json"
	form := url.Values{}
	form.Set("url", "https://twitter.com/")
	form.Set("banner", params.Banner)
	form.Set("width", params.Width)
	form.Set("height", params.Height)
	form.Set("offset_left", params.OffsetLeft)
	form.Set("offset_top", params.OffsetTop)

	client := CreateHttpClient()
	res, err := client.PostForm(baseUrl, form)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
