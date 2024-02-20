package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

func showFollowCount(followCount int) {
	if followCount > 0 {
		log.Println("Follow count:", followCount)
		return
	}

	log.Println("No new followings")
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
	allFollowing_q := make(map[string]bool)

	followingsPerPage := 100

	log.Printf("Following count: %d", *currentUser.Following)
	log.Println("Finding all your followings...")

	for i := 0; i < (*currentUser.Following/followingsPerPage)+1; i++ {

		following, _, err := ghClient.Users.ListFollowing(ctx, "", &github.ListOptions{PerPage: followingsPerPage, Page: i})
		if err != nil {
			log.Println(err)
		}

		for _, fl := range following {
			allFollowing = append(allFollowing, fl)
			allFollowing_q[*fl.Login] = true
			//TODO: Need to rewrite
		}
		time.Sleep(300 * time.Millisecond)
	}

	log.Println("Finding all your followings...done")
	followCount := 0

	// Followback mode
	// Finds follower which are not followed by you and tries to follow them
	if os.Getenv("FOLLOWBACK") == "1" {
		log.Println("Finding all your followers and try to follow then")
		for i := 0; i < (*currentUser.Followers/followingsPerPage)+1; i++ {

			followers, _, err := ghClient.Users.ListFollowers(ctx, "", &github.ListOptions{PerPage: followingsPerPage, Page: i})
			if err != nil {
				log.Println(err)
			}

			for _, fl := range followers {
				if _, ok := allFollowing_q[*fl.Login]; !ok {
					log.Println("Gonna follow:", *fl.Login)

					_, err := ghClient.Users.Follow(ctx, *fl.Login)
					if err != nil {
						log.Println(err)
						continue
					}
					followCount++
					time.Sleep(3 * time.Second)
				}
			}
		}
		showFollowCount(followCount)
		os.Exit(0)
	}

	log.Println("Downloading RSS feed...")

	client := http.Client{}
	request, err := http.NewRequest("GET", gitHubRss, nil)

	if err != nil {
		log.Panicln(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Downloading RSS feed...done")

	feed, err := gofeed.NewParser().Parse(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Try to follow new users")

	for _, f := range feed.Items {
		if strings.Contains(f.GUID, "FollowEvent") {
			if !isFollowing(allFollowing, f.Link) {
				s := strings.Split(f.Link, "github.com/")
				log.Println("Gonna follow:", s[1])
				_, err := ghClient.Users.Follow(ctx, s[1])

				time.Sleep(500 * time.Millisecond)
				if err != nil {
					log.Println(err)
					continue
				}
				followCount++
			}
		}
	}
	showFollowCount(followCount)

}
