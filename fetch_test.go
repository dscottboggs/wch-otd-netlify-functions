package wch_otd_netlify_functions

import (
	"testing"
	"time"

	"github.com/dscottboggs/attest"
)

func TestMakeUrl(t *testing.T) {
	test := attest.New(t)
	date := time.Date(2023, 4, 20, 0, 0, 0, 0, time.Local)
	test.Equals(
		makeUrl(date).String(),
		"https://api.baserow.io/api/database/rows/table/33215/?filter__field_177139__equal=4&filter__field_177140__equal=20&size=100&user_field_names=true",
	)
}
