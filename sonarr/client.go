package sonarr

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"gopkg.in/resty.v1"
)

var (
	apiRgx = regexp.MustCompile(`[a-z0-9]{32}`)
)

func NewClient(c Config) (*Client, error) {
	if c.Hostname == "" {
		return nil, fmt.Errorf("hostname is empty")
	}

	if match := apiRgx.MatchString(c.APIKey); !match {
		return nil, fmt.Errorf("api key is invalid format: %s", c.APIKey)
	}

	baseURL := createApiURL(c)

	r := resty.New()
	r.SetHostURL(baseURL)
	r.SetHeader("Accept", "application/json")
	r.SetQueryParam("apikey", c.APIKey)
	if c.Username != "" && c.Password != "" {
		r.SetBasicAuth(c.Username, c.Password)
	}

	client := Client{
		apiKey:     c.APIKey,
		maxResults: c.MaxResults,
		username:   c.Username,
		password:   c.Password,
		baseURL:    baseURL,
		client:     r,
	}
	return &client, nil
}

func createApiURL(c Config) string {
	c.Hostname = strings.TrimPrefix(c.Hostname, "http://")
	c.Hostname = strings.TrimPrefix(c.Hostname, "https://")
	c.URLBase = strings.TrimPrefix(c.URLBase, "/")

	u := url.URL{}
	if c.SSL {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	if c.Port == 80 {
		u.Host = c.Hostname
	} else {
		u.Host = fmt.Sprintf("%s:%d", c.Hostname, c.Port)
	}

	if c.URLBase != "" {
		u.Path = fmt.Sprintf("%s/api/v3", c.URLBase)
	} else {
		u.Path = "/api/v3"
	}

	fmt.Println("The URL for Sonarr is", u.String())

	return u.String()
}

type Client struct {
	apiKey     string
	username   string
	password   string
	baseURL    string
	maxResults int
	client     *resty.Client
}

func (c *Client) SearchTVShows(term string) ([]TVShow, error) {
	resp, err := c.client.R().SetQueryParam("term", term).SetResult([]TVShow{}).Get("series/lookup")
	if err != nil {
		return nil, err
	}

	tvShows := *resp.Result().(*[]TVShow)
	if len(tvShows) > c.maxResults {
		tvShows = tvShows[:c.maxResults]
	}
	return tvShows, nil
}

func (c *Client) GetFolders() ([]Folder, error) {
	resp, err := c.client.R().SetResult([]Folder{}).Get("rootfolder")
	if err != nil {
		return nil, err
	}

	folders := *resp.Result().(*[]Folder)
	return folders, nil
}

func (c *Client) GetProfile(prfl string) ([]Profile, error) {

	resp, err := c.client.R().SetResult([]Profile{}).Get(prfl)
	if err != nil {
		return nil, err
	}
	profile := *resp.Result().(*[]Profile)

	return profile, nil

}

func (c *Client) AddTVShow(m TVShow, languageProfile int, qualityProfile int, path string) (tvShow TVShow, err error) {

	request := AddTVShowRequest{
		Title:             m.Title,
		TitleSlug:         m.TitleSlug,
		Images:            m.Images,
		QualityProfileID:  qualityProfile,
		LanguageProfileID: languageProfile,
		TVDBID:            m.TVDBID,
		RootFolderPath:    path,
		Monitored:         true,
		Year:              m.Year,
		Seasons:           m.Seasons,
		AddOptions:        AddTVShowOptions{SearchForMissingEpisodes: true},
	}

	resp, err := c.client.R().SetBody(request).SetResult(TVShow{}).Post("series")
	if err != nil {
		fmt.Println(err)
		return
	}

	tvShow = *resp.Result().(*TVShow)
	return
}
