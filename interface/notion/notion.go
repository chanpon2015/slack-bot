package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// 9a60ebe973e541dc8cb1ec291d19eb9c
const (
	// EndPointPages is Notion
	EndPointPages     = "https://api.notion.com/v1/pages"
	EndPointDatabases = "https://api.notion.com/v1/databases/%s"
	EndPointBlocks    = "https://api.notion.com/v1/blocks/%s"
)

// Notion is
type Notion interface {
	GetDatabases(context.Context, string) error
	CreatePage(context.Context, string) error
}

type notionImpl struct {
	token  string
	client *http.Client
}

type Record struct {
	Parent     Parent   `json:"parent"`
	Properties Property `json:"properties"`
}

type Parent struct {
	DatabaseID string `json:"database_id"`
}

type Property struct {
	Name Name `json:"Name"`
}

type Name struct {
	Title []Title `json:"title"`
}

type Title struct {
	Text Text `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

// NewNotion is
func NewNotion(token string, client *http.Client) Notion {
	return &notionImpl{
		token:  token,
		client: client,
	}
}

func (n *notionImpl) Do(req *http.Request) (*http.Response, error) {
	n.baseRequest(req)
	return n.client.Do(req)
}

func (n *notionImpl) baseRequest(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+n.token)
	req.Header.Add("Notion-Version", "2021-05-13")
	req.Header.Add("Content-Type", "application/json")
}

func (n *notionImpl) GetDatabases(ctx context.Context, databaseID string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(EndPointDatabases, databaseID),
		nil,
	)
	if err != nil {
		return err
	}
	resp, err := n.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("http status %s", resp.Status)
	}
	aaa, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(aaa))
	return nil
}

func (n *notionImpl) CreatePage(ctx context.Context, databaseID string) error {
	data := Record{
		Parent: Parent{
			DatabaseID: databaseID,
		},
		Properties: Property{
			Name: Name{
				Title: []Title{
					{
						Text: Text{
							Content: "test",
						},
					},
				},
			},
		},
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println(string(payload))
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		EndPointPages,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}
	resp, err := n.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		aaa, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(aaa))
		return fmt.Errorf("http status %s", resp.Status)
	}
	return nil
}
