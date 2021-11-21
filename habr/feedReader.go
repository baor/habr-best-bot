package habr

import (
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
	req.Header.Set("User-Agent", "curl/7.64.0")
	req.Header.Set("Accept", "*/*")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln("Error requesting habr RSS:", err.Error())
	}
	defer resp.Body.Close()

	log.Printf("Parse RSS feed")
	feed, _ := fp.Parse(resp.Body)
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
