/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package wch_otd_api

import (
	"testing"
	"time"

	"github.com/dscottboggs/attest"
)

func TestKeyForDay(t *testing.T) {
	var (
		test = attest.New(t)
		day  = time.Date(2023, 5, 6, 0, 0, 0, 0, time.UTC)
	)
	test.Equals("otd-5-6", keyForDay(&day))
}

func TestSetForDay(t *testing.T) {
	var (
		test    = attest.New(t)
		day     = time.Date(2023, 5, 6, 0, 0, 0, 0, time.UTC)
		example = []OurResponse{
			{Title: "test",
				Content:  "test content",
				MoreInfo: "",
				Excerpt:  "test content",
				Author: struct {
					Name  string "json:\"name\""
					Url   string "json:\"url\""
					Email string "json:\"email\""
				}{Name: "test author"},
				Url:   "https://stories.workingclasshistory.com/article/10646/alejandro-finisterre-born",
				Media: nil}}
		c      = test.EatError(connectToCache()).(*cache)
		result = test.EatError(c.setForDay(example, &day)).(string)
	)
	test.Equals(
		`[{"title":"test","content":"test content","more_info":"","excerpt":"test content","author":{"name":"test author","url":"","email":""},"url":"https://stories.workingclasshistory.com/article/10646/alejandro-finisterre-born","media":null}]`+"\n",
		result,
	)
}

func TestGetDay(t *testing.T) {
	var (
		test        = attest.New(t)
		exampleData = "asoifoinsfoindlfksfoiwne"
		day         = time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC)
		c           = test.EatError(connectToCache()).(*cache)
	)
	test.Handle(c.client.Set(c.ctx, "otd-2-29", exampleData, 30*time.Second).Err())
	test.Equals(exampleData, test.EatError(c.getDay(&day)).(string))
}
