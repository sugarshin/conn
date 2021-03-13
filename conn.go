package conn

import (
	goconfluence "github.com/virtomize/confluence-go-api"
)

// Client is a confluence-go-api client.
type Client struct {
	api *goconfluence.API
}

// New returns a new API.
func New(endpoint string, username string, password string) (*Client, error) {
	api, err := goconfluence.NewAPI(endpoint, username, password)
	client := &Client{
		api,
	}
	return client, err
}

// CreateSubPage :
func (client *Client) CreateSubPage(parentPageID string, data *goconfluence.Content) (*goconfluence.Content, error) {
	ancestors := []goconfluence.Ancestor{
		{
			ID: parentPageID,
		},
	}
	data.Ancestors = append(data.Ancestors, ancestors...)
	return client.api.CreateContent(data)
}

// CreateSubPageWithLatest :
func (client *Client) CreateSubPageWithLatest(parentPageID string, with func(data *goconfluence.Content) *goconfluence.Content) (*goconfluence.Content, error) {
	res, err := client.api.GetChildPages(parentPageID)
	if err != nil {
		return nil, err
	}

	latest := res.Results[len(res.Results)-1]
	content, err := client.api.GetContentByID(latest.ID, goconfluence.ContentQuery{
		Expand: []string{"body.storage", "space"},
	})
	if err != nil {
		return nil, err
	}
	data := with(content)
	return client.CreateSubPage(parentPageID, data)
}
