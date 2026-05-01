package menu

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/goodsign/monday"
)

var NotFoundError = errors.New("couldn't find today's menu")
var NoServiceError = errors.New("food service is not available during weekends")

var locationTr *time.Location

func init() {
	var err error
	locationTr, err = time.LoadLocation("Europe/Istanbul")
	if err != nil {
		panic(err)
	}
}

type Menu struct {
	Date  time.Time
	Items []string
}

func extract(reader io.Reader) ([]Menu, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var menus []Menu
	for _, s := range doc.Find(".yemekhane-primary").EachIter() {
		title := s.Find(".yemekhane-title").Text()
		date, err := monday.ParseInLocation("Monday 02 January 2006 15:04", title, locationTr, monday.LocaleTrTR)
		if err != nil {
			return nil, err
		}

		var items []string
		for _, s := range s.Find("li").EachIter() {
			raw := s.Contents().Last().Text()
			for item := range strings.SplitSeq(raw, "/") {
				item = strings.TrimSpace(item)
				item = strings.ToLowerSpecial(unicode.TurkishCase, item)
				if item != "" && item != "..." {
					items = append(items, item)
				}
			}
		}

		if items != nil {
			menus = append(menus, Menu{date, items})
		}
	}

	return menus, nil
}

func Month(ctx context.Context) ([]Menu, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://sivas.edu.tr/yemek-listesi", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return extract(resp.Body)
}

func dateEqual(a, b time.Time) bool {
	ya, ma, da := a.Date()
	yb, mb, db := b.Date()
	return ya == yb && ma == mb && da == db
}

func Today(ctx context.Context) (Menu, error) {
	now := time.Now().In(locationTr)

	switch now.Weekday() {
	case time.Saturday, time.Sunday:
		return Menu{}, NoServiceError
	}

	month, err := Month(ctx)
	if err != nil {
		return Menu{}, err
	}

	for _, day := range month {
		if dateEqual(day.Date, now) {
			return day, nil
		}
	}

	return Menu{}, NotFoundError
}
