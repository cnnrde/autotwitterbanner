package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	client2 "github.com/connordennison/autotwitterbanner/client"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

type Credentials struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

// getTwitterClient returns a Twitter client authenticated with the provided consumer key and secret
func getTwitterClient(creds *Credentials) (*twitter.Client, error) {
	// pass consumerkey and secret
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// pass token and secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return client, nil
}

func loadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	loadDotEnv()
	fmt.Println("AutoTwitterBanner")
	creds := Credentials{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
		AccessToken:    os.Getenv("ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("ACCESS_SECRET"),
	}

	// fmt.Printf("%+v\n", creds)

	client, err := getTwitterClient(&creds)
	if err != nil {
		fmt.Println("Error getting Twitter client: ", err)
		os.Exit(1)
	}

	user, _, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		fmt.Println("Error getting user: ", err)
		os.Exit(1)
	}

	// fmt.Printf("%+v\n", client)
	// print number of followers
	followers, _, err := client.Followers.List(&twitter.FollowerListParams{ScreenName: "cnnrde", Count: 1})
	if err != nil {
		fmt.Println("Error getting followers: ", err)
		os.Exit(1)
	}
	// print followers
	// fmt.Printf("%+v\n", followers)
	fmt.Printf("Fetched %v follower\n", len(followers.Users))
	for _, follower := range followers.Users {
		fmt.Printf("%v\n", follower.Name)
	}
	file, err := os.Open("template.svg")
	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}
	defer file.Close()
	// read file in chunks of 4 bytes
	template := ""
	b := make([]byte, 4)
	for {
		readTotal, err := file.Read(b)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		template += string(b[:readTotal])
		// template = string(b[:readTotal])
	}
	// replace text in template
	template = strings.Replace(template, "[mrf tag]", "@"+followers.Users[0].ScreenName, -1)
	template = strings.Replace(template, "[pfp source]", followers.Users[0].ProfileImageURLHttps, -1)
	// get username of logged in user
	template = strings.Replace(template, "[tag]", "@"+user.ScreenName, -1)
	template = strings.Replace(template, "[username]", user.Name, -1)
	// save to new file
	file, err = os.Create("banner.svg")
	if err != nil {
		fmt.Println("Error creating file: ", err)
		os.Exit(1)
	}
	defer file.Close()
	file.Write([]byte(template))

	// this refers to a node package installed globally because i'm lazy af
	out, err := exec.Command("svgexport", "banner.svg", "banner.png").Output()
	if err != nil {
		fmt.Println(string(out))
		fmt.Println("Error creating banner png: ", err)
		os.Exit(1)
	}

	// convert banner.png to base64
	banner, err := os.Open("banner.png")
	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}
	defer banner.Close()
	// read file in chunks of 4 bytes
	b = make([]byte, 4)
	image := ""
	for {
		readTotal, err := banner.Read(b)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		image += string(b[:readTotal])
		// bannerBase64 = string(b[:readTotal])
	}
	bannerBase64 := base64.StdEncoding.EncodeToString([]byte(image))

	// update user
	client2.UpdateAccountBanner(client2.UpdateAccountBannerParams{Banner: bannerBase64, Width: "1500", Height: "500", OffsetLeft: "0", OffsetTop: "0"})
}
