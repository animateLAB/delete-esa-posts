package main

import (
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
	"time"
	"strconv"
	"os"
	"github.com/pkg/errors"
)

type EsaPostsResponse struct {
	Posts []struct {
		Number int `json:"number"`
		Name string `json:"name"`
		FullName string `json:"full_name"`
		Wip bool `json:"wip"`
		BodyMd string `json:"body_md"`
		BodyHTML string `json:"body_html"`
		CreatedAt time.Time `json:"created_at"`
		Message string `json:"message"`
		Kind string `json:"kind"`
		CommentsCount int `json:"comments_count"`
		TasksCount int `json:"tasks_count"`
		DoneTasksCount int `json:"done_tasks_count"`
		URL string `json:"url"`
		UpdatedAt time.Time `json:"updated_at"`
		Tags []interface{} `json:"tags"`
		Category string `json:"category"`
		RevisionNumber int `json:"revision_number"`
		CreatedBy struct {
			Name string `json:"name"`
			ScreenName string `json:"screen_name"`
			Icon string `json:"icon"`
		} `json:"created_by"`
		UpdatedBy struct {
			Name string `json:"name"`
			ScreenName string `json:"screen_name"`
			Icon string `json:"icon"`
		} `json:"updated_by"`
		StargazersCount int `json:"stargazers_count"`
		WatchersCount int `json:"watchers_count"`
		Star bool `json:"star"`
		Watch bool `json:"watch"`
		SharingUrls interface{} `json:"sharing_urls"`
	} `json:"posts"`
	PrevPage interface{} `json:"prev_page"`
	NextPage interface{} `json:"next_page"`
	TotalCount int `json:"total_count"`
	Page int `json:"page"`
	PerPage int `json:"per_page"`
	MaxPerPage int `json:"max_per_page"`
}

type Client struct {
	URL string
	Token string
	HTTPClient *http.Client
}

func NewClient(apiUrl, token string) (*Client, error) {
	client := &Client{
		URL: apiUrl,
		HTTPClient: &http.Client{},
		Token: token,
	}
	return client, nil
}

func NewRequest(requestURL string, method string, token string) (*http.Request, error) {
	parsedUrl, err := url.ParseRequestURI(requestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse url")
	}

	req, err := http.NewRequest(method, parsedUrl.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed create http request")
	}

	req.Header.Add("Authorization", "Bearer " + token)
	return req, nil
}

func GetPosts(client *Client, query string, out *EsaPostsResponse) error {
	req, err := NewRequest(client.URL + query, "GET", client.Token)
	if err != nil {
		return errors.Wrap(err, "failed create request")
	}

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed do http request")
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	return errors.Wrap(decoder.Decode(out), "failed decode")
}

func DeletePosts(client *Client, ids []int) error {
	for _, postNumber := range ids {
		query := "/posts/" + strconv.Itoa(postNumber)

		req, err := NewRequest(client.URL + query, "DELETE", client.Token)
		if err != nil {
			return errors.Wrap(err, "failed create request")
		}

		res, err := client.HTTPClient.Do(req)
		if err != nil {
			return errors.Wrap(err, "failed do http request")
		}

		res.Body.Close()

		time.Sleep(1 * time.Second)
	}
	return nil
}

func MapPostsToIds(postsRes *EsaPostsResponse) ([]int, error) {
	var ids []int
	for _, post := range postsRes.Posts {
		ids = append(ids, post.Number)
	}
	return ids, nil
}

func main() {
	token := os.Getenv("ESA_TOKEN")
	if token == "" {
		fmt.Printf("ESA_TOKENが未設定です\n")
		return
	}

	team := os.Getenv("ESA_TEAM")
	if team == "" {
		fmt.Printf("ESA_TEAMが未設定です\n")
		return
	}

	searchQuery := os.Getenv("ESA_SEARCH_QUERY")
	if searchQuery == "" {
		fmt.Printf("ESA_SEARCH_QUERYが未設定です\n")
		return
	}

	url := "https://api.esa.io/v1/teams/" + team
	query := "/posts?q=" + searchQuery
	fmt.Printf("検索条件 %s の記事を削除します\n", searchQuery)

	client, err := NewClient(url, token)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	var postsRes EsaPostsResponse
	err = GetPosts(client, query, &postsRes)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	ids, _ := MapPostsToIds(&postsRes)
	if len(ids) == 0 {
		fmt.Println("削除対象の記事はありません")
		return
	}

	err = DeletePosts(client, ids)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%d件の記事を削除しました\n", len(ids))
}
