package habr

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

type FeedReader interface {
	GetBestFeed() []FeedItem
}

type FeedItem struct {
	Message string
	ID      string
}

func getPostID(link string) string {
	r1 := regexp.MustCompile(`(?s)/post/(\d*)/`)
	match := r1.FindAllStringSubmatch(link, -1)
	if len(match) == 0 {
		return ""
	}
	if len(match[0]) < 2 {
		return ""
	}
	return match[0][1]
}

func stripTags(textWithTags string) (string, error) {
	tokenizer := html.NewTokenizer(strings.NewReader(textWithTags))

	var b bytes.Buffer
	isInAToken := false
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			if tokenizer.Err() != io.EOF {
				log.Println("Error html tokenizer token", tokenizer.Err().Error())
				return "", tokenizer.Err()
			}
			res := b.String()
			log.Println(res)
			return res, nil
		case html.TextToken:
			if isInAToken {
				_, _ = b.Write(tokenizer.Raw())
				continue
			}
			_, _ = b.Write(tokenizer.Text())
		case html.StartTagToken, html.SelfClosingTagToken:
			tn, _ := tokenizer.TagName()
			switch string(tn) {
			case "img":
				_, _ = b.WriteString(" ")
				continue
			case "a":
				isInAToken = true
				_, _ = b.Write(tokenizer.Raw())
				continue
			default:
				_, _ = b.WriteString(" ")
			}

		case html.EndTagToken:
			tn, _ := tokenizer.TagName()
			switch string(tn) {
			case "a":
				_, _ = b.Write(tokenizer.Raw())
				isInAToken = false
				continue
			}
		}
	}
}

type HabrReader struct{}

func NewHabrReader() FeedReader {
	return &HabrReader{}
}

func processItem(item *gofeed.Item) (FeedItem, error) {
	body, err := stripTags(item.Description)
	if err != nil {
		return FeedItem{}, err
	}
	msg := "<a href=\"" + item.Link + "\">" + item.Title + "</a>\n" + body
	postID := getPostID(item.Link)

	return FeedItem{
		Message: msg,
		ID:      postID}, nil
}

func (HabrReader) GetBestFeed() []FeedItem {
	var processedFeed []FeedItem

	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	fp := gofeed.NewParser()

	log.Printf("Pull RSS feed")
	req, err := http.NewRequest("GET", "https://habr.com/ru/rss/best/", nil)
	if err != nil {
		log.Fatalln("Error creating request object to habr RSS")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "habr_web_home=ARTICLES_LIST_TOP; hl=ru; fl=ru")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln("Error requesting habr RSS. ", err.Error())
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatalln("Error decoding gzip body. ", err.Error())
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	log.Printf("Parse RSS feed")
	feed, err := fp.Parse(reader)
	if err != nil {
		body, _ := ioutil.ReadAll(reader)
		log.Printf("Feed response: %s", body)
		log.Fatalln("Error parsing RSS feed:", err.Error())
	}
	log.Printf("RSS feed is pulled. Published: %s, Number of items: %d", feed.Published, len(feed.Items))

	for _, item := range feed.Items {
		processedItem, err := processItem(item)
		if err != nil {
			log.Println("Error processing feed item:", err.Error())
			continue
		}
		processedFeed = append(processedFeed, processedItem)
	}
	return processedFeed
}
