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
	LinkToImage string
	Message     string
	ID          string
}

func getFirstImageLink(textWithTags string) string {
	r1 := regexp.MustCompile(`(?s)<img\s*src="([^"]*)"`)
	match := r1.FindAllStringSubmatch(textWithTags, -1)
	if len(match) == 0 {
		return ""
	}
	if len(match[0]) < 2 {
		return ""
	}
	imageLink := match[0][1]
	return imageLink
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

func traverse(n *html.Node, b *bytes.Buffer) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "html":
		case "head":
		case "body":
		case "a":
			err := html.Render(b, n)
			if err != nil {
				log.Fatalln("Error rendering the node", err.Error())
			}
			return
		default:
			_, _ = b.WriteString(" ")
		}
	}

	if n.Type == html.TextNode {
		_, _ = b.WriteString(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, b)
	}
}

func stripTags(textWithTags string) string {
	parsedHtmlDescription, err := html.Parse(strings.NewReader(textWithTags))
	if err != nil {
		log.Fatalln("Error parsing body as HTML :", textWithTags)
	}

	var b bytes.Buffer
	traverse(parsedHtmlDescription, &b)
	res := b.String()
	log.Println(res)
	return res
}

type HabrReader struct{}

func NewHabrReader() FeedReader {
	return &HabrReader{}
}

func processItem(item *gofeed.Item) FeedItem {
	linkToImage := getFirstImageLink(item.Description)
	msg := "<a href=\"" + item.Link + "\">" + item.Title + "</a>\n"
	msg += stripTags(item.Description)
	postID := getPostID(item.Link)

	return FeedItem{
		LinkToImage: linkToImage,
		Message:     msg,
		ID:          postID}
}

func (HabrReader) GetBestFeed() []FeedItem {
	var response []FeedItem

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
	log.Printf("RSS feed is pulled. Description %s, Published: %s", feed.Description, feed.Published)

	for _, item := range feed.Items {
		response = append(response, processItem(item))
	}
	return response
}
