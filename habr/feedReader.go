package habr

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedReader interface {
	GetBestFeed(allowedTags []string) []FeedItem
}

type FeedItem struct {
	LinkToImage string
	Message     string
	ID          string
}

func stripTag(textWithTags string, tagNameToRemove string) string {
	r1 := regexp.MustCompile(`(?s)(<` + tagNameToRemove + `[^>]*>)`)
	strippedText := r1.ReplaceAllString(textWithTags, " ")

	r2 := regexp.MustCompile(`(?s)(</` + tagNameToRemove + `[^>]*>)`)
	strippedText = r2.ReplaceAllString(strippedText, " ")

	return strippedText
}

func getAllTags(textWithTags string) []string {
	var tags []string
	r1 := regexp.MustCompile(`(?s)<([^\s>/]*)`)
	match := r1.FindAllStringSubmatch(textWithTags, -1)
	for _, val := range match {
		if val[1] != "" {
			tags = append(tags, val[1])
		}
	}
	return tags
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

func isTagAllowed(tag string, allowedTags []string) bool {
	for _, allowedTag := range allowedTags {
		if tag == allowedTag {
			return true
		}
	}
	return false
}

func stripTags(textWithTags string, allowedTags []string) string {
	allTags := getAllTags(textWithTags)
	strippedString := textWithTags
	for _, tag := range allTags {
		if isTagAllowed(tag, allowedTags) {
			continue
		}
		strippedString = stripTag(strippedString, tag)
	}
	return strippedString
}

type HabrReader struct{}

func NewHabrReader() FeedReader {
	return &HabrReader{}
}

func (HabrReader) GetBestFeed(allowedTags []string) []FeedItem {
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
		linkToImage := getFirstImageLink(item.Description)
		msg := "<a href=\"" + item.Link + "\">" + item.Title + "</a>\n"
		msg += stripTags(item.Description, allowedTags)
		postID := getPostID(item.Link)
		response = append(response, FeedItem{
			LinkToImage: linkToImage,
			Message:     msg,
			ID:          postID})
	}
	return response
}
