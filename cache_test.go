/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package wch_otd_api

import (
	"encoding/json"
	"strings"
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

const exampleText = `[{"title": "Ludlow Massacre","content": "<p>On 20 April 1914, the Ludlow massacre took place when US troops opened fire with machine guns on a camp of striking miners and their families in Ludlow, Colorado.</p> <p>12,000 miners had gone out on strike the previous September against the Rockefeller family-owned Colorado Fuel and Iron Corporation (CF&amp;I) following the killing of an activist of the United Mine Workers of America (UMWA). They then demanded better safety at work, and to be paid in money, instead of company scrip (tokens which could only be redeemed in the company store).</p> <p>The Rockefellers evicted the striking miners and their families from their homes, and so they set up \"tent cities\" to live in collectively, which miners' wives helped run. Company thugs harassed strikers, and occasionally drove by camps riddling them with machine-gun fire, killing and injuring workers and their children.</p> <p>Eventually, the national guard was ordered to evict all the strike encampments, and the morning of April 20 they attacked the largest camp in Ludlow. They opened fire with machine guns on the tents of the workers and their families, who then returned fire. The main organiser of the camp, Louis Tikas, went to visit the officer in charge of the national guard to arrange a truce. But he was beaten to the ground then shot repeatedly in the back, killing him. That night, troops entered the camp and set fire to it, killing 11 children and two women, in addition to 13 other people who were killed in the fighting. The youngest victim was Elvira Valdez, aged just 3 months.</p> <p>Protests against the massacre broke out across the country, but the workers at CF&amp;I were defeated, and many of them were subsequently sacked and replaced with non-union miners. Over the course of the strike 66 people were killed, but no guardsmen or company thugs were prosecuted.</p>","more_info": "","excerpt": "<p>On 20 April 1914, the Ludlow massacre took place when US troops opened fire with machine guns on a camp of striking miners and their families in Ludlow, Colorado.</p> <p>12,000 miners had gone out on strike the previous September against the Rockefeller family-owned Colorado Fuel and Iron Corporation (CF...</p>","author": "Working Class History","url": "https://stories.workingclasshistory.com/article/9243/ludlow-massacre","media":{"url": "https://workingclasshistory.com/wp-content/uploads/2023/02/04.20-x-Ludlow_striker_family_in_front_of_tent.jpg","credit": "Denver Library/Wikimedia Commons","caption": "Striking miners wives and children in the strikers' tent camp, 1914"}}]`

func TestSetForDay(t *testing.T) {
	var (
		test    = attest.New(t)
		day     = time.Date(2023, 5, 6, 0, 0, 0, 0, time.UTC)
		example = make(DaysData, 0, 1)
		reader  = strings.NewReader(exampleText)
	)
	test.Handle(json.NewDecoder(reader).Decode(&example))
	var (
		c      = test.EatError(connectToCache()).(*cache)
		result = test.EatError(c.setForDay(example, &day)).(DaysData)[0]
	)
	test.Equals(example[0].Author, result.Author)
	test.Equals(example[0].Content, result.Content)
	test.Equals(example[0].Title, result.Title)
	test.Equals(example[0].Excerpt, result.Excerpt)
	test.Equals(example[0].Url, result.Url)
	test.Equals(example[0].MoreInfo, result.MoreInfo)
	test.Equals(example[0].Media.Url, result.Media.Url)
	test.Equals(example[0].Media.Credit, result.Media.Credit)
	test.Equals(example[0].Media.Caption, result.Media.Caption)
	cached := test.EatError(c.client.Get(c.ctx, "otd-5-6").Result()).(string)
	var cacheResponse DaysData
	test.Handle(json.NewDecoder(strings.NewReader(cached)).Decode(&cacheResponse))
	result = cacheResponse[0]
	test.Equals(example[0].Author, result.Author)
	test.Equals(example[0].Content, result.Content)
	test.Equals(example[0].Title, result.Title)
	test.Equals(example[0].Excerpt, result.Excerpt)
	test.Equals(example[0].Url, result.Url)
	test.Equals(example[0].MoreInfo, result.MoreInfo)
	test.Equals(example[0].Media.Url, result.Media.Url)
	test.Equals(example[0].Media.Credit, result.Media.Credit)
	test.Equals(example[0].Media.Caption, result.Media.Caption)
}

func TestGetDay(t *testing.T) {
	var (
		test    = attest.New(t)
		day     = time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC)
		c       = test.EatError(connectToCache()).(*cache)
		reader  = strings.NewReader(exampleText)
		example = make(DaysData, 0, 1)
		result  OurResponse
	)
	test.Handle(json.NewDecoder(reader).Decode(&example))
	test.Handle(c.client.Set(c.ctx, "otd-2-29", exampleText, 30*time.Second).Err())
	result = test.EatError(c.getDay(&day)).(DaysData)[0]
	test.Equals(example[0].Author, result.Author)
	test.Equals(example[0].Content, result.Content)
	test.Equals(example[0].Title, result.Title)
	test.Equals(example[0].Excerpt, result.Excerpt)
	test.Equals(example[0].Url, result.Url)
	test.Equals(example[0].MoreInfo, result.MoreInfo)
	test.Equals(example[0].Media.Url, result.Media.Url)
	test.Equals(example[0].Media.Credit, result.Media.Credit)
	test.Equals(example[0].Media.Caption, result.Media.Caption)
}
