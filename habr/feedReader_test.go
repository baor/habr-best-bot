package habr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripTag_BrTag_NoTags(t *testing.T) {
	input := `notags`
	outputExpected := input
	outputActual := stripTags(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestStripTag_BrTag(t *testing.T) {
	input := `<br>text<br>`
	outputExpected := " text "
	outputActual := stripTags(input)
	assert.Equal(t, outputExpected, outputActual)
}
func TestStripTag_BrTag2(t *testing.T) {
	input := `<br/>text<br/>`
	outputExpected := " text "
	outputActual := stripTags(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestStripTag_ATag(t *testing.T) {
	input := `<b>text</b>`
	outputExpected := " text"
	outputActual := stripTags(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestStripTag_NestedBTagWithArgsAndText(t *testing.T) {
	input := `<b>text1 <a arg="aaa">text2</a> text3</b>`
	outputExpected := ` text1 <a arg="aaa">text2</a> text3`
	outputActual := stripTags(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestStripTag_ATagOneByOne(t *testing.T) {
	input := `<b>text1</b> text2 <b>text3</b>`
	outputExpected := ` text1 text2  text3`
	outputActual := stripTags(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestStripTag_ImgTag(t *testing.T) {
	input := `<img>text`
	outputExpected := " text"
	outputActual := stripTags(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestStripTag_ImgTagWithArgs(t *testing.T) {
	input := `<img src="" size="">text`
	outputExpected := " text"
	outputActual := stripTags(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestStripTags_AllowedTags(t *testing.T) {
	input := `<a>text1</a> text2 <z>text3</z>`
	outputExpected := "<a>text1</a> text2  text3"
	outputActual := stripTags(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestGetFirstImageLink_NoLink(t *testing.T) {
	input := `notags`
	outputExpected := ""
	outputActual := getFirstImageLink(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestGetFirstImageLink_WithArgs(t *testing.T) {
	input := `<img src="https://site/pic.png" alt="image" width="300">`
	outputExpected := "https://site/pic.png"
	outputActual := getFirstImageLink(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestGetFirstImageLink_FirstLink(t *testing.T) {
	input := `<img src="https://site/pic1.png"> <img src="https://site/pic2.png">`
	outputExpected := "https://site/pic1.png"
	outputActual := getFirstImageLink(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestGetPostID(t *testing.T) {
	input := `https://habr.com/post/413925/?utm_source=habrahabr&utm_medium=rss&utm_campaign=413925`
	outputExpected := "413925"
	outputActual := getPostID(input)
	assert.Equal(t, outputExpected, outputActual)
}

func TestStripTags_UrlDecode(t *testing.T) {
	input := `<p>ого <a href="https://audio-v-text.silero.ai/a%3BRedis%3BRocksDB"> text </a> Читать дальше &rarr; ого —&nbsp;ого</p>`
	outputActual := stripTags(input)
	outputExpected := " ого <a href=\"https://audio-v-text.silero.ai/a%3BRedis%3BRocksDB\"> text </a> Читать дальше → ого —\u00a0ого"
	assert.Equal(t, outputExpected, outputActual)
}

// func TestManual(t *testing.T) {
// 	c := NewHabrReader()
// 	feeds := c.GetBestFeed()
// 	assert.Greater(t, len(feeds), 1)
// }
