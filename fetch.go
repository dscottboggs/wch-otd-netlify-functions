/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * This file contains the high-level functions for querying the API which are
 * called directly by the netlify function scripts. They wrap up the whole
 * process of querying the database and transforming the results.
 */

package wch_otd_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var BASEROW_API_KEY string

func init() {
	for _, arg := range os.Args {
		if arg[:5] == "-test" {
			fmt.Println("test mode, not checking for API key")
			return
		}
	}
	var found bool
	BASEROW_API_KEY, found = os.LookupEnv("REACT_APP_BASEROW_TOKEN")
	if !found {
		log.Fatalln("REACT_APP_BASEROW_TOKEN environment variable not found")
	}
}
func makeUrl(date time.Time) *url.URL {
	query := url.Values{}
	query.Set("filter__field_177140__equal", strconv.Itoa(date.Day()))
	query.Set("filter__field_177139__equal", strconv.Itoa(int(date.Month())))
	query.Set("size", "100")
	query.Set("user_field_names", "true")
	it, _ := url.Parse("https://api.baserow.io/api/database/rows/table/33215/")
	it.RawQuery = query.Encode()
	return it
}

func FetchForDate(date time.Time, tz *time.Location) ([]OurResponse, error) {
	if tz == nil {
		tz = time.UTC
	}
	date = date.In(tz)
	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Token %s", BASEROW_API_KEY))
	headers.Set("Accept", "application/json")
	client := http.Client{
		// Netlify enforces a 10-second max run time, so...
		Timeout: 9 * time.Second,
	}
	request := http.Request{
		Method: "GET",
		URL:    makeUrl(date),
		Header: headers,
	}
	result, err := client.Do(&request)
	if err != nil {
		return nil, err
	}
	if result.StatusCode == 200 {
		var dbResponse DbResponse
		err = json.NewDecoder(result.Body).Decode(&dbResponse)
		if err != nil {
			return nil, err
		}
		return dbResponse.Transform()
	} else {
		return nil, fmt.Errorf("database request failed: %s", result.Status)
	}
}

func FetchToday(tz *time.Location) ([]OurResponse, error) {
	return FetchForDate(time.Now(), tz)
}
