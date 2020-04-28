package habr

import (
	"log"
	"regexp"

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
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://habr.com/rss/best/")

	log.Printf("Pull RSS feed. Description %s, Published: %s", feed.Description, feed.Published)
	var response []FeedItem
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
