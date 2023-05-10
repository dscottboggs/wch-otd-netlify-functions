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

func TestSetForDay(t *testing.T) {
	var (
		test        = attest.New(t)
		day         = time.Date(2023, 5, 6, 0, 0, 0, 0, time.UTC)
		example     = make(DaysData, 0, 1)
		exampleText = `[{"title":"Alejandro Finisterre born","content":"On 6 May 1919, Alexandre Campos Ramírez (more commonly known as Alejandro Finisterre), anarchist poet and inventor of the Spanish version of table football or \"foosball\" was born. His political idealism seemed to be realised with the outbreak of the revolution in 1936 which defeated the military rising of general Franco. However his house was bombed by nationalist forces and he was severely injured. Seeing many injured children who were unable to play football with their friends, he built a table out of pine and attached players to steel bars — inventing the modern version of table football commonly played in Spain. After the war he had to flee the Franco regime to France, and was imprisoned for four years in Morocco before heading to the Americas. More info in this fascinating short biography of him here: https://libcom.org/history/finisterre-alejandro-1919-2007","more_info":"","excerpt":"On 6 May 1919, Alexandre Campos Ramírez (more commonly known as Alejandro Finisterre), anarchist poet and inventor of the Spanish version of table football or \"foosball\" was born. His political idealism seemed to be realised with the outbreak of the revolution in 1936 which defeated the military rising of general...","author":{"name":"Working Class History","url":"https://workingclasshistory.com","email":""},"url":"https://stories.workingclasshistory.com/article/10646/alejandro-finisterre-born","media":null}]
`
		reader = strings.NewReader(exampleText)
	)
	test.Handle(json.NewDecoder(reader).Decode(&example))
	var (
		c      = test.EatError(connectToCache()).(*cache)
		result = test.EatError(c.setForDay(example, &day)).(DaysData)[0]
	)
	test.Equals(example[0].Author.Email, result.Author.Email)
	test.Equals(example[0].Author.Name, result.Author.Name)
	test.Equals(example[0].Author.Url, result.Author.Url)
	test.Equals(example[0].Content, result.Content)
	test.Equals(example[0].Title, result.Title)
	test.Equals(example[0].Excerpt, result.Excerpt)
	test.Equals(example[0].Url, result.Url)
	test.Equals(example[0].MoreInfo, result.MoreInfo)
	if result.Media != nil {
		test.Fatalf("media was not nil, it was: %#+v", result.Media)
	}
}

func TestGetDay(t *testing.T) {
	var (
		test        = attest.New(t)
		day         = time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC)
		c           = test.EatError(connectToCache()).(*cache)
		exampleText = `[{"title":"Alejandro Finisterre born","content":"On 6 May 1919, Alexandre Campos Ramírez (more commonly known as Alejandro Finisterre), anarchist poet and inventor of the Spanish version of table football or \"foosball\" was born. His political idealism seemed to be realised with the outbreak of the revolution in 1936 which defeated the military rising of general Franco. However his house was bombed by nationalist forces and he was severely injured. Seeing many injured children who were unable to play football with their friends, he built a table out of pine and attached players to steel bars — inventing the modern version of table football commonly played in Spain. After the war he had to flee the Franco regime to France, and was imprisoned for four years in Morocco before heading to the Americas. More info in this fascinating short biography of him here: https://libcom.org/history/finisterre-alejandro-1919-2007","more_info":"","excerpt":"On 6 May 1919, Alexandre Campos Ramírez (more commonly known as Alejandro Finisterre), anarchist poet and inventor of the Spanish version of table football or \"foosball\" was born. His political idealism seemed to be realised with the outbreak of the revolution in 1936 which defeated the military rising of general...","author":{"name":"Working Class History","url":"https://workingclasshistory.com","email":""},"url":"https://stories.workingclasshistory.com/article/10646/alejandro-finisterre-born","media":null}]
`
		reader  = strings.NewReader(exampleText)
		example = make(DaysData, 0, 1)
	)
	test.Handle(json.NewDecoder(reader).Decode(&example))
	test.Handle(c.client.Set(c.ctx, "otd-2-29", exampleText, 30*time.Second).Err())
	result := test.EatError(c.getDay(&day)).(DaysData)[0]
	test.Equals(example[0].Author.Email, result.Author.Email)
	test.Equals(example[0].Author.Name, result.Author.Name)
	test.Equals(example[0].Author.Url, result.Author.Url)
	test.Equals(example[0].Content, result.Content)
	test.Equals(example[0].Title, result.Title)
	test.Equals(example[0].Excerpt, result.Excerpt)
	test.Equals(example[0].Url, result.Url)
	test.Equals(example[0].MoreInfo, result.MoreInfo)
	if result.Media != nil {
		test.Fatalf("media was not nil, it was: %#+v", result.Media)
	}
}
