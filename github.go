package main

import (
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	token *oauth2.Token
}

// Token satisfies oauth2.TokenSource interface
func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

// GithubClient holds the GitHub client and the owner from config
type GithubClient struct {
	client *github.Client
	owner  string
}

// NewGithubClient returns a GitHub client
func NewGithubClient(token, owner string) GithubClient {
	client := github.NewClient(oauth2.NewClient(oauth2.NoContext, &tokenSource{
		&oauth2.Token{
			AccessToken: token,
			TokenType:   "token",
		},
	}))
	return GithubClient{client: client, owner: owner}
}

// GetKeys uses the GitHub API to get the SSH keys of a given user
func (c *GithubClient) GetKeys(user github.User) ([]github.Key, error) {
	keys, _, err := c.client.Users.ListKeys(*user.Login, nil)
	return keys, err
}

// GetTeamMembers uses the GitHub API to get the members of a given team
func (c *GithubClient) GetTeamMembers(name string) ([]github.User, error) {
	var team github.Team
	teams, _, err := c.client.Organizations.ListTeams(c.owner, nil)
	if err != nil {
		panic(err)
	}
	for _, t := range teams {
		if strings.EqualFold(*t.Name, name) {
			team = t
			break
		}
	}
	users, _, err := c.client.Organizations.ListTeamMembers(*team.ID, nil)
	return users, err
}

// GetTeamKeys uses the GitHub API to get the SSH keys of each member of a team
func (c *GithubClient) GetTeamKeys(users []github.User) []github.Key {
	ch := make(chan []github.Key)
	keys := []github.Key{}
	remaining := len(users)

	for _, user := range users {
		go func(user github.User) {
			k, err := c.GetKeys(user)
			if err != nil {
				panic(err)
			}
			ch <- k
		}(user)
	}

	for {
		select {
		case res := <-ch:
			keys = append(keys, res...)
			remaining--
			if remaining <= 0 {
				return keys
			}
		}
	}

	return keys
}
