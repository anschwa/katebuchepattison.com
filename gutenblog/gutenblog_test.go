package gutenblog

import (
	"encoding/json"
	"testing"
	"time"
)

func newMonth(year int, month time.Month) Date {
	return NewDate(year, month, 1)
}

func prettyJSON(args ...interface{}) ([]string, error) {
	results := make([]string, 0, len(args))

	for _, x := range args {
		b, err := json.MarshalIndent(x, "", "\t")
		if err != nil {
			return nil, err
		}

		results = append(results, string(b))
	}

	return results, nil
}

func Test_getArchive(t *testing.T) {
	// Test for out-of-order posts as well as duplicate publish dates
	posts := []blogPost{
		{Title: "H", Date: NewDate(2008, 1, 1)},
		{Title: "A", Date: NewDate(2006, 1, 1)},
		{Title: "B", Date: NewDate(2006, 1, 1)},
		{Title: "C", Date: NewDate(2006, 1, 12)},
		{Title: "D", Date: NewDate(2006, 2, 1)},
		{Title: "G", Date: NewDate(2007, 11, 22)},
		{Title: "F", Date: NewDate(2007, 5, 6)},
		{Title: "E", Date: NewDate(2007, 1, 1)},
		{Title: "J", Date: NewDate(2008, 10, 11)},
		{Title: "I", Date: NewDate(2008, 10, 10)},
	}

	want := blogArchive{
		{Date: newMonth(2006, time.January), Posts: []blogPost{
			{Title: "A", Date: NewDate(2006, 1, 1)},
			{Title: "B", Date: NewDate(2006, 1, 1)},
			{Title: "C", Date: NewDate(2006, 1, 12)}}},
		{Date: newMonth(2006, time.February), Posts: []blogPost{
			{Title: "D", Date: NewDate(2006, 2, 1)}}},
		{Date: newMonth(2007, time.January), Posts: []blogPost{
			{Title: "E", Date: NewDate(2007, 1, 1)}}},
		{Date: newMonth(2007, time.May), Posts: []blogPost{
			{Title: "F", Date: NewDate(2007, 5, 6)}}},
		{Date: newMonth(2007, time.November), Posts: []blogPost{
			{Title: "G", Date: NewDate(2007, 11, 22)}}},
		{Date: newMonth(2008, time.January), Posts: []blogPost{
			{Title: "H", Date: NewDate(2008, 1, 1)}}},
		{Date: newMonth(2008, time.October), Posts: []blogPost{
			{Title: "I", Date: NewDate(2008, 10, 10)},
			{Title: "J", Date: NewDate(2008, 10, 11)}}},
	}

	archive := makeArchive(posts)
	if w, g := len(want), len(archive); w != g {
		t.Errorf("wrong number of months: want=%d; got=%d", w, g)
	}

	for i, month := range archive {
		if w, g := len(want[i].Posts), len(month.Posts); w != g {
			t.Errorf("wrong number of posts for %s: want=%d; got=%d", month.Title(), w, g)
		}

		if w, g := want[i].Title(), month.Title(); w != g {
			t.Errorf("wrong title for month: want=%q; got=%q", w, g)
		}

		for j, post := range month.Posts {
			if w, g := want[i].Posts[j], post; w != g {
				t.Errorf("wrong post: want=%q; got=%q", w.Title, g.Title)
			}
		}
	}

	// Pretty print results
	res, err := prettyJSON(want, archive)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("want: %+v \n\n got: %+v", res[0], res[1])
}
