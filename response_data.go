/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * This file contains the data structures for deserializing database results and
 * serializing API responses, and their associated methods.
 */

package wch_otd_api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/mborgerson/gotruncatehtml/truncatehtml"
)

// The response format from the Baserow API
type DbResponse struct {
	Count   uint            `json:"count"`
	Results []DbResponseRow `json:"results"`
}

// The relevant data from a "row" in the database, as returned by the Baserow
// API. Please note that many fields are omitted, and we may have a use in the
// future for those omitted fields.
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

// Transforms the title into the format used as a part of the URL to the full
// article.
func (r *DbResponseRow) UrlEncodedTitle() string {
	return url.QueryEscape(strings.ToLower(strings.ReplaceAll(r.Title, " ", "-")))
}

// Return the fully constructed URL which points to the full article in the
// Stories app.
func (r *DbResponseRow) ArticleUrl() (*url.URL, error) {
	return url.Parse(fmt.Sprintf("https://stories.workingclasshistory.com/article/%v/%s", r.ID, r.UrlEncodedTitle()))
}

// Return the first 500 visible characters of the HTML content, preserving tags.
//
// Return an error if invalid HTML is encountered (that is, if a closing tag
// was encountered which had not been open, like "<p></p></p>", see test)
func (r *DbResponseRow) Excerpt() (string, error) {
	summary, err := truncatehtml.TruncateHtml([]byte(r.Description), 500, "...")
	if err != nil {
		return "", err
	}
	return string(summary), err
}

// Return the API-response formatted data about this row in the database.
//
// Returns an error if there is a problem constructing the URL or excerpt.
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

// Metadata for media attached to this event.
type MediaInfo struct {
	Url     string `json:"url"`
	Credit  string `json:"credit"`
	Caption string `json:"caption"`
}

// The format of a single historical event as output by the API
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

// Transform a full response from the database into a list of rows in our
// chosen format.
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
