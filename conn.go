package conn

import (
	"errors"

	goconfluence "github.com/virtomize/confluence-go-api"
)

// Client is a confluence-go-api client.
type Client struct {
	*goconfluence.API
}

// Content is goconfluence.Content
type Content = goconfluence.Content

// New returns a new Client.
func New(endpoint string, username string, password string) (*Client, error) {
	api, err := goconfluence.NewAPI(endpoint, username, password)
	client := &Client{api}
	return client, err
}

// CreateSubPageContent :
func (client *Client) CreateSubPageContent(parentPageID string, data *goconfluence.Content) (*goconfluence.Content, error) {
	ancestors := []goconfluence.Ancestor{
		{
			ID: parentPageID,
		},
	}
	data.Ancestors = append(data.Ancestors, ancestors...)
	return client.CreateContent(data)
}

// CreateSubPageContentWithLatest :
func (client *Client) CreateSubPageContentWithLatest(parentPageID string, with func(data *goconfluence.Content) *goconfluence.Content) (*goconfluence.Content, error) {
	content, err := client.GetLatestChildPageContent(parentPageID)
	if err != nil {
		return nil, err
	}
	data := with(content)
	return client.CreateSubPageContent(parentPageID, data)
}

// CreateSubPageContentWith :
func (client *Client) CreateSubPageContentWith(parentPageID string, with func() *goconfluence.Content) (*goconfluence.Content, error) {
	return client.CreateSubPageContent(parentPageID, with())
}

// GetLatestChildPageContent :
func (client *Client) GetLatestChildPageContent(parentPageID string) (*goconfluence.Content, error) {
	content, err := client.GetChildPageContentWith(parentPageID, func(i int, _ goconfluence.Results, list []goconfluence.Results) bool {
		if i == len(list)-1 {
			return true
		} else {
			return false
		}
	})
	if err != nil {
		return nil, err
	}
	return content, nil
}

// GetChildPageContentByID :
func (client *Client) GetChildPageContentByID(parentPageID string, id string) (*goconfluence.Content, error) {
	content, err := client.GetChildPageContentWith(parentPageID, func(_ int, results goconfluence.Results, _ []goconfluence.Results) bool {
		if results.ID == id {
			return true
		} else {
			return false
		}
	})
	if err != nil {
		return nil, err
	}
	return content, nil
}

// GetChildPageContentWith :
func (client *Client) GetChildPageContentWith(parentPageID string, with func(index int, results goconfluence.Results, list []goconfluence.Results) bool) (*goconfluence.Content, error) {
	res, err := client.GetChildPages(parentPageID)
	if err != nil {
		return nil, err
	}
	for i, r := range res.Results {
		if with(i, r, res.Results) {
			content, err := client.GetContentByID(r.ID, goconfluence.ContentQuery{
				Expand: []string{"body.storage", "space"},
			})
			if err != nil {
				return nil, err
			}
			return content, nil
		}
	}
	return nil, errors.New("cannot find a page")
}
