package main

import (
	"context"
	"log"
	"net/http"
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
			// log.Println("*u.HTMLURL", *u.HTMLURL)
			return true
		}
	}
	return false
}

func main() {
	client := http.Client{}
	request, err := http.NewRequest("GET", "https://github.com/esin.private.atom?token=AAARBB3Y4R2GCVJMTMHFUIV5TW3ZU", nil)

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
		&oauth2.Token{AccessToken: "dcbd38866b051566559e13f577f94502cd00272b"},
	)
	tc := oauth2.NewClient(ctx, ts)

	ghClient := github.NewClient(tc)

	currentUser, _, err := ghClient.Users.Get(ctx, "esin")
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
			// log.Println("Tadam: ", f.Link)
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
