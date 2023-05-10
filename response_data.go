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
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mborgerson/gotruncatehtml/truncatehtml"
)

// How many characters the excerpt should be at most
var EXCERPT_MAX_CHARS int

func init() {
	if maxChars, found := os.LookupEnv("API_EXCERPT_MAX_CHARS"); found {
		var err error
		EXCERPT_MAX_CHARS, err = strconv.Atoi(maxChars)
		if err != nil {
			log.Panicf("invalid value for $API_EXERPT_MAX_CHARS \"%s\": %v", maxChars, err)
		}
	} else if EXCERPT_MAX_CHARS < 0 {
		log.Panicf("invalid value for $API_EXERPT_MAX_CHARS \"%s\": length cannot be less than zero", maxChars)
	} else {
		// Default
		EXCERPT_MAX_CHARS = 250
	}
}

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
	summary, err := truncatehtml.TruncateHtml([]byte(r.Description), EXCERPT_MAX_CHARS, "...")
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
		Author:   r.AuthorName,
		Url:      articleUrl.String(),
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
	Title    string     `json:"title"`
	Content  string     `json:"content"`
	MoreInfo string     `json:"more_info"`
	Excerpt  string     `json:"excerpt"`
	Author   string     `json:"author"`
	Url      string     `json:"url"`
	Media    *MediaInfo `json:"media"`
}

// Transform a full response from the database into a list of rows in our
// chosen format.
func (r DbResponse) Transform() (DaysData, error) {
	result := make(DaysData, 0, r.Count)
	for _, row := range r.Results {
		transformed, err := row.Transform()
		if err != nil {
			return result, err
		}
		result = append(result, *transformed)
	}
	return result, nil
}

// The data for a single day
type DaysData []OurResponse

// Return an appropriate HTTP response from this dataset. If there is an error
// while serializ, error)ing the data, an Internal Server Error (status 500) response
// will be returned instead
func (data DaysData) MakeResponse() *events.APIGatewayProxyResponse {
	var (
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(data)
	)
	if err != nil {
		return InternalServerError("encoding JSON", err)
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       buf.String(),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}
