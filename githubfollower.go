package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mmcdole/gofeed"
	"golang.org/x/oauth2"
)

func isFollowing(following []*github.User, user string) bool {

	if len(following) <= 0 {
		return false
	}
	for _, u := range following {

		if *u.HTMLURL == user {
			return true
		}
	}
	return false
}

func main() {

	gitHubToken := os.Getenv("GITHUB_TOKEN")
	if gitHubToken == "" {
		log.Panicln("Env variable GITHUB_TOKEN not found")
	}
	gitHubRss := os.Getenv("GITHUB_RSS")
	if gitHubRss == "" {
		log.Panicln("Env variable GITHUB_RSS not found")
	}
	gitHubUserName := os.Getenv("GITHUB_USERNAME")
	if gitHubUserName == "" {
		log.Panicln("Env variable GITHUB_USERNAME not found")
	}

	client := http.Client{}
	request, err := http.NewRequest("GET", gitHubRss, nil)

	if err != nil {
		log.Panicln(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Panicln(err)
	}

	feed, err := gofeed.NewParser().Parse(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	ghClient := github.NewClient(tc)

	currentUser, _, err := ghClient.Users.Get(ctx, gitHubUserName)
	if err != nil {
		log.Panicln(err)
	}

	allFollowing := []*github.User{}

	followingsPerPage := 100

	for i := 0; i < (*currentUser.Following/followingsPerPage)+1; i++ {

		following, _, err := ghClient.Users.ListFollowing(ctx, "", &github.ListOptions{PerPage: followingsPerPage, Page: i})
		if err != nil {
			log.Println(err)
		}

		for _, fl := range following {
			allFollowing = append(allFollowing, fl)
		}
	}

	log.Printf("Following count: %d", *currentUser.Following)

	followCount := 0
	for _, f := range feed.Items {
		if strings.Contains(f.GUID, "FollowEvent") {
			if !isFollowing(allFollowing, f.Link) {
				s := strings.Split(f.Link, "github.com/")
				log.Println("Gonna follow:", s[1])
				_, err := ghClient.Users.Follow(ctx, s[1])

				if err != nil {
					log.Println(err)
					continue
				}
				followCount++
			}
		}
	}
	if followCount > 0 {
		log.Println("Follow count:", followCount)
	}

	if followCount == 0 {
		log.Println("No new followings")
	}
}
