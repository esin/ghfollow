package main

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {

	gitHubToken := os.Getenv("GITHUB_TOKEN")
	if gitHubToken == "" {
		log.Panicln("Env variable GITHUB_TOKEN not found")
	}

	gitHubUserName := os.Getenv("GITHUB_USERNAME")
	if gitHubUserName == "" {
		log.Panicln("Env variable GITHUB_USERNAME not found")
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

	allFollowing := make(map[string]bool)

	followingsPerPage := 100

	// Count whom I'm following
	for i := 0; i < (*currentUser.Following/followingsPerPage)+1; i++ {

		following, _, err := ghClient.Users.ListFollowing(ctx, "", &github.ListOptions{PerPage: followingsPerPage, Page: i})
		if err != nil {
			log.Println(err)
		}

		for _, fl := range following {
			allFollowing[*fl.Login] = true
		}
	}

	// Count, who following me
	for i := 0; i < (*currentUser.Followers/followingsPerPage)+1; i++ {

		followers, _, err := ghClient.Users.ListFollowers(ctx, "", &github.ListOptions{PerPage: followingsPerPage, Page: i})
		if err != nil {
			log.Println(err)
		}

		for _, fl := range followers {
			if _, ok := allFollowing[*fl.Login]; !ok {
				log.Println("Gonna follow:", *fl.Login)
				_, err := ghClient.Users.Follow(ctx, *fl.Login)

				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}
