/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package wch_otd_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	ctx    context.Context
	client *redis.Client
}

func connectToCache() (*cache, error) {
	var (
		loc, locFound         = os.LookupEnv("API_CACHE_LOCATION")
		pass                  = os.Getenv("API_CACHE_PASSWORD")
		username              = os.Getenv("API_CACHE_USERNAME")
		dbAsText, dbFound     = os.LookupEnv("API_CACHE_DB")
		db                int = 0
		err               error
	)
	if !locFound {
		loc = "localhost:6379"
	}
	if dbFound {
		db, err = strconv.Atoi(dbAsText)
		if err != nil {
			return nil, err
		}
	}
	return &cache{
		ctx:    context.Background(),
		client: redis.NewClient(&redis.Options{Addr: loc, Password: pass, DB: db, Username: username}),
	}, nil
}

func keyForDay(date *time.Time) string {
	_, month, day := date.Date()
	return fmt.Sprintf("otd-%d-%v", month, day)
}

// Store the given day's database result data in the cache for 48 hours in the cache for 48 hours
func (c *cache) setForDay(data DaysData, day *time.Time) (DaysData, error) {
	var (
		key   = keyForDay(day)
		buf   = new(bytes.Buffer)
		err   = json.NewEncoder(buf).Encode(data)
		value string
	)
	if err != nil {
		return data, err
	}
	value = buf.String()
	return data, c.client.Set(c.ctx, key, value, 48*time.Hour).Err()
}

// Return the given day's data from the redis cache, if present. Returns
// redis.Nil if not found.
func (c cache) getDay(day *time.Time) (DaysData, error) {
	var (
		data      = make(DaysData, 0)
		key       = keyForDay(day)
		text, err = c.client.Get(c.ctx, key).Result()
		reader    = strings.NewReader(text)
	)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(reader).Decode(&data)
	return data, err
}
