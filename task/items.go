package task

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/wfunc/go/xlog"
	"golang.org/x/net/html"
)

var (
	CodeAPI = true
	CodeURL = ""
)

func BootstrapConfig() {
	CodeURL = os.Getenv("CodeURL")
}

const (
	canceled = "context canceled"
	deadline = "context deadline exceeded"
)

// ExtractLinksWithPrefix parses the HTML content, extracts <a> tag href attributes
// that contain any of the keys, and updates the result map with the specified prefix.
func ExtractLinksWithPrefix(htmlContent string, result map[string]string, prefix string, keys []string) error {
	// Parse the HTML document
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return err // Return parsing error
	}

	// Recursive function to traverse the HTML node tree
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		// Process only <a> tags
		if node.Type == html.ElementNode && node.Data == "a" {
			href := getAttributeValue(node, "href")
			if href != "" {
				for _, key := range keys {
					if strings.Contains(href, key) {
						result[href] = prefix
						break // Exit the loop after matching a key
					}
				}
			}
		}
		// Recursively process child nodes
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	// Start traversing the HTML nodes
	traverse(doc)
	return nil
}

// getAttributeValue retrieves the value of the specified attribute from an HTML node.
func getAttributeValue(node *html.Node, attrName string) string {
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

// ParseShopHTML parses the provided outerHTML and extracts relevant shop data into the shopMap.
// The function looks for specific patterns in the HTML content to extract keys and values.
func ParseShopHTML(outerHTML string, shopMap map[string]string) error {
	// Parse the HTML document
	doc, err := html.Parse(strings.NewReader(outerHTML))
	if err != nil {
		return err // Return parsing error
	}

	var lastNodeData string // Stores the last node's data for context-sensitive parsing
	var currentKey string   // The key being constructed

	// Recursive function to traverse the HTML node tree
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		// Context-sensitive logic based on the last node's data
		if lastNodeData == "img" {
			currentKey = "🌹" + node.Data
			if strings.Contains(node.Data, ":") && strings.Contains(node.Data, "级") && !strings.Contains(node.Data, "等级:") {
				currentKey = "🪴" + node.Data
			}
		} else if strings.Contains(node.Data, ":") && strings.Contains(node.Data, "级") && !strings.Contains(node.Data, "等级:") {
			currentKey = "🪴" + node.Data
		} else if strings.Contains(node.Data, "单价:") {
			parts := strings.Split(node.Data, "\n\t\t\t")
			if len(parts) > 1 {
				currentKey += "：" + strings.Split(parts[1], "单价:")[1]
			}
		} else if strings.Contains(node.Data, ":") && strings.Contains(node.Data, "开心币") {
			lines := strings.Split(node.Data, "\n")
			if len(lines) > 1 {
				currentKey = lines[1]
			}
		} else if strings.Contains(node.Data, ":") && strings.Contains(node.Data, "哇币") {
			lines := strings.Split(node.Data, "\n")
			if len(lines) > 1 {
				currentKey += "：" + lines[1]
			} else {
				currentKey = node.Data
			}
		}

		// Update lastNodeData for the next iteration
		lastNodeData = node.Data

		// Process <a> tags to extract href attributes
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					switch {
					case strings.Contains(attr.Val, "seedsInfo.do?seedsId="):
						shopMap[currentKey] = strings.Split(attr.Val, "=")[1]
					case strings.Contains(attr.Val, "paySeedsInfo.do?paySeedsId="):
						shopMap[currentKey] = strings.Split(strings.Split(attr.Val, "=")[1], "&")[0]
					case strings.Contains(attr.Val, "buyFemale.do?femaleid="):
						shopMap[currentKey] = strings.Split(attr.Val, "=")[1]
					}
					return
				}
			}
		}

		// Recursively traverse child nodes
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	// Start traversing the HTML nodes
	traverse(doc)
	return nil
}

func (t *BaseTask) clickNext(ctx context.Context, selector string) bool {
	var exists bool
	if len(selector) < 1 {
		selector = `//a[contains(text(), ">>下页")]`
	}
	script := fmt.Sprintf(`document.evaluate('%s', document, null, XPathResult.BOOLEAN_TYPE, null).booleanValue`, selector)
	fmt.Println(script)
	var err = chromedp.Evaluate(script, &exists).Do(ctx)
	if err != nil {
		if t.Verbose {
			xlog.Infof("FarmTask(%v) extractTimes failed with err %v", t.Username, err)
		}
		return false
	}
	return exists
}
