package main

import (
	"context"
	"html/template"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	tw "github.com/g8rswimmer/go-twitter/v2"
	"github.com/goodsign/monday"
	"github.com/joho/godotenv"
	"github.com/x-gkm/sbtufood/menu"
	"github.com/x-gkm/sbtufood/twitter"
)

var tweetTemplate = template.Must(template.New("tweet").Parse(
	`Tarih: {{ .Date }}
Menü:
{{- range .Items }}
- {{ . }}
{{- end }}
`))

type tweetData struct {
	Date  string
	Items []string
}

func renderTemplate(menu menu.Menu) (string, error) {
	items := make([]string, 0, len(menu.Items))

	for _, item := range menu.Items {
		items = append(items, capitalize(item))
	}

	data := tweetData{
		Date:  monday.Format(menu.Date, monday.DefaultFormatTrTRFull, monday.LocaleTrTR),
		Items: items,
	}

	buf := &strings.Builder{}
	err := tweetTemplate.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	client := twitter.New(ctx, twitter.Auth{
		BearerToken:    os.Getenv("TWITTER_BEARER_TOKEN"),
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_KEY_SECRET"),
		AccessKey:      os.Getenv("TWITTER_ACCESS_KEY"),
		AccessSecret:   os.Getenv("TWITTER_ACCESS_KEY_SECRET"),
	})

	today, err := menu.Today()
	if err != nil {
		log.Panicln(err)
	}

	text, err := renderTemplate(today)
	if err != nil {
		log.Panicln(err)
	}

	resp, err := client.CreateTweet(ctx, tw.CreateTweetRequest{
		Text: text,
	})
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("%#v\n", resp)
}
