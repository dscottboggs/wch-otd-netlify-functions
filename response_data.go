package wch_otd_netlify_functions

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/mborgerson/gotruncatehtml/truncatehtml"
)

type DbResponse struct {
	Count   uint            `json:"count"`
	Results []DbResponseRow `json:"results"`
}
type DbResponseRow struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Year         string `json:"year"`
	Description  string `json:"description"`
	MoreInfo     string `json:"more_info"`
	Media        string `json:"media"`
	MediaCredit  string `json:"media_credit"`
	MediaCaption string `json:"media_caption"`
	AuthorName   string `json:"author_name"`
	AuthorURL    string `json:"author_url"`
	AuthorEmail  string `json:"author_email"`
}

func (r *DbResponseRow) UrlEncodedTitle() string {
	return url.QueryEscape(strings.ToLower(strings.ReplaceAll(r.Title, " ", "-")))
}

func (r *DbResponseRow) ArticleUrl() (*url.URL, error) {
	return url.Parse(fmt.Sprintf("https://stories.workingclasshistory.com/article/%v/%s", r.ID, r.UrlEncodedTitle()))
}

// Return the first 500 visible characters of the HTML content, preserving tags.
//
// Return an error if invalid HTML is encountered
func (r *DbResponseRow) Excerpt() (string, error) {
	summary, err := truncatehtml.TruncateHtml([]byte(r.Description), 500, "...")
	if err != nil {
		return "", err
	}
	return string(summary), err
}

func (r *DbResponseRow) Transform() (*OurResponse, error) {
	var (
		excerpt, err = r.Excerpt()
		media        *MediaInfo
	)
	if err != nil {
		return nil, err
	}
	articleUrl, err := r.ArticleUrl()
	if err != nil {
		return nil, err
	}
	if len(r.Media) == 0 {
		media = nil
	} else {
		media = &MediaInfo{
			Url:     r.Media,
			Credit:  r.MediaCredit,
			Caption: r.MediaCaption,
		}
	}
	return &OurResponse{
		Title:    r.Title,
		Content:  r.Description,
		MoreInfo: r.MoreInfo,
		Media:    media,
		Excerpt:  excerpt,
		Author: struct {
			Name  string "json:\"name\""
			Url   string "json:\"url\""
			Email string "json:\"email\""
		}{
			Name:  r.AuthorName,
			Url:   r.AuthorURL,
			Email: r.AuthorEmail,
		},
		Url: articleUrl.String(),
	}, nil
}

type MediaInfo struct {
	Url     string `json:"url"`
	Credit  string `json:"credit"`
	Caption string `json:"caption"`
}

type OurResponse struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	MoreInfo string `json:"more_info"`
	Excerpt  string `json:"excerpt"`
	Author   struct {
		Name  string `json:"name"`
		Url   string `json:"url"`
		Email string `json:"email"`
	} `json:"author"`
	Url   string     `json:"url"`
	Media *MediaInfo `json:"media"`
}

func (r DbResponse) Transform() ([]OurResponse, error) {
	result := make([]OurResponse, 0, r.Count)
	for _, row := range r.Results {
		transformed, err := row.Transform()
		if err != nil {
			return result, err
		}
		result = append(result, *transformed)
	}
	return result, nil
}
