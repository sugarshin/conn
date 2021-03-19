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

// Ancestor is goconfluence.Ancestor
type Ancestor = goconfluence.Ancestor

// Version is goconfluence.Version
type Version = goconfluence.Version

// Body is goconfluence.Body
type Body = goconfluence.Body

// Storage is goconfluence.Storage
type Storage = goconfluence.Storage

// Results is goconfluence.Results
type Results = goconfluence.Results

// New returns a new Client.
func New(endpoint string, username string, password string) (*Client, error) {
	api, err := goconfluence.NewAPI(endpoint, username, password)
	client := &Client{api}
	return client, err
}

// CreateChildPageContent :
func (client *Client) CreateChildPageContent(parentPageID string, data *Content) (*Content, error) {
	ancestors := []goconfluence.Ancestor{
		{
			ID: parentPageID,
		},
	}
	data.Ancestors = append(data.Ancestors, ancestors...)
	return client.CreateContent(data)
}

// CreateChildPageContentWithLatest :
func (client *Client) CreateChildPageContentWithLatest(parentPageID string, with func(data *Content) *Content) (*Content, error) {
	content, err := client.GetLatestChildPageContent(parentPageID)
	if err != nil {
		return nil, err
	}
	data := with(content)
	return client.CreateChildPageContent(parentPageID, data)
}

// CreateChildPageContentWith :
func (client *Client) CreateChildPageContentWith(parentPageID string, with func() *Content) (*Content, error) {
	return client.CreateChildPageContent(parentPageID, with())
}

// GetLatestChildPageContent :
func (client *Client) GetLatestChildPageContent(parentPageID string) (*Content, error) {
	content, err := client.GetChildPageContentWith(parentPageID, func(i int, _ Results, list []Results) bool {
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
func (client *Client) GetChildPageContentByID(parentPageID string, id string) (*Content, error) {
	content, err := client.GetChildPageContentWith(parentPageID, func(_ int, results Results, _ []Results) bool {
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
func (client *Client) GetChildPageContentWith(parentPageID string, with func(index int, results Results, list []Results) bool) (*Content, error) {
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
