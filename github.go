package main

import (
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	MAX_PAGE_SIZE = 100
)

//- naive oauth setup
type AccessToken struct {
	token *oauth2.Token
}

func (a AccessToken) Token() (*oauth2.Token, error) {
	return a.token, nil
}

func newAccessToken(token string) AccessToken {
	t := oauth2.Token{AccessToken: token}
	return AccessToken{token: &t}
}

//- models & namespacing
type GithubClient struct {
	client github.Client
	owner  string
}

func (c *GithubClient) GetKeys(user github.User) ([]*github.Key, error) {
	keys, _, err := c.client.Users.ListKeys(*user.Login, nil)
	return keys, err
}

func (c *GithubClient) GetTeamMembersByID(teamID int) ([]*github.User, error) {
	var opt = &github.OrganizationListTeamMembersOptions{
		ListOptions: github.ListOptions{
			PerPage: MAX_PAGE_SIZE,
		},
	}

	var allMembers []*github.User
	for {
		members, resp, err := c.client.Organizations.ListTeamMembers(teamID, opt)
		if err != nil {
			return nil, err
		}

		allMembers = append(allMembers, members...)
		if resp.LastPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allMembers, nil
}

func (c *GithubClient) GetTeamMembers(name string) ([]*github.User, error) {
	var team *github.Team
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

	return c.GetTeamMembersByID(*team.ID)
}

func (c *GithubClient) GetTeamKeys(users []*github.User) []*github.Key {
	ch := make(chan []*github.Key)
	keys := []*github.Key{}
	remaining := len(users)

	for _, user := range users {
		go func(user github.User) {
			k, err := c.GetKeys(user)
			if err != nil {
				panic(err)
			}
			ch <- k
		}(*user)
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

func NewGithubClient(token, owner string) GithubClient {
	c := oauth2.NewClient(oauth2.NoContext, newAccessToken(token))
	return GithubClient{client: *github.NewClient(c), owner: owner}
}
