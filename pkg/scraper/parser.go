package scraper

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type htmlParser struct {
	config *Config 
}

func NewParser(cfg *Config) Parser {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	
	return &htmlParser{
		config: cfg,
	}
}

func (p *htmlParser) Parse(html []byte) (*ScrapedData, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err 
	}
	
	data := &ScrapedData{
		OGTags:   make(map[string]string),
		Headings: make(map[string][]string),
	}
	
	p.extractMetadata(doc, data)
	
	p.extractContent(doc, data)
	
	p.extractStructure(doc, data)	

	if p.config.IncludeRawHTML { 
		data.RawHTML = string(html)
	}
	
	return data, nil
}

func (p *htmlParser) extractMetadata(doc *goquery.Document, data *ScrapedData) {
	
	data.Title = doc.Find("title").First().Text()
	data.Title = strings.TrimSpace(data.Title)
	data.MetaDescription = p.getMetaContent(doc, "name", "description")
	
	data.MetaKeywords = p.getMetaContent(doc, "name", "keywords")
	
	
	canonical, exists := doc.Find("link[rel='canonical']").Attr("href")
	if exists {
		data.CanonicalURL = strings.TrimSpace(canonical)
	}
	
	
	doc.Find("meta[property^='og:']").Each(func(i int, s *goquery.Selection) {
		property, propExists := s.Attr("property")
		content, contentExists := s.Attr("content")
		
		if propExists && contentExists {
			data.OGTags[property] = content
		}
	})
}
//text content
func (p *htmlParser) extractContent(doc *goquery.Document, data *ScrapedData) {
	for i := 1; i <= 6; i++ {
		selector := "h" + string(rune('0'+i))
		headings := []string{}
		
		doc.Find(selector).Each(func(index int, s *goquery.Selection) {
			text := s.Text()
			text = strings.TrimSpace(text)
			if text != "" {
				headings = append(headings, text)
			}
		})
		
		if len(headings) > 0 {
			data.Headings[selector] = headings
		}
	}
	//paragraphs
	data.Paragraphs = []string{}
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		text = strings.TrimSpace(text)
		if text != "" {
			data.Paragraphs = append(data.Paragraphs, text)
		}
	})
	
	
	bodyText := doc.Find("body").Text()
	data.TextContent = p.cleanText(bodyText)
	
	if p.config.MaxTextLength > 0 && len(data.TextContent) > p.config.MaxTextLength {
		data.TextContent = data.TextContent[:p.config.MaxTextLength]
	} // 
}

// extractStructure extracts links and images from the HTML document
func (p *htmlParser) extractStructure(doc *goquery.Document, data *ScrapedData) {
	// Extract all links
	// <a href="https://example.com">Link Text</a>
	data.Links = []Link{}
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return // Skip if no href
		}
		
		// Get link text
		text := s.Text()
		text = strings.TrimSpace(text)
		
		// Get rel attribute (for nofollow, noopener, etc.)
		rel, _ := s.Attr("rel")
		
		link := Link{
			URL:  href,
			Text: text,
			Rel:  rel,
		}
		
		data.Links = append(data.Links, link)
	})
	
	// Extract all images
	// <img src="image.jpg" alt="Description">
	data.Images = []Image{}
	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists || src == "" {
			return // Skip if no src
		}
		
		// Get alt text
		alt, _ := s.Attr("alt")
		
		image := Image{
			Src: src,
			Alt: alt,
		}
		
		data.Images = append(data.Images, image)
	})
}

// getMetaContent extracts content from meta tags
// attrName: "name" or "property"
// attrValue: "description", "keywords", etc.
func (p *htmlParser) getMetaContent(doc *goquery.Document, attrName, attrValue string) string {
	// Build selector: meta[name="description"] or meta[property="og:title"]
	selector := "meta[" + attrName + "='" + attrValue + "']"
	
	content, exists := doc.Find(selector).First().Attr("content")
	if exists {
		return strings.TrimSpace(content)
	}
	
	return ""
}

// cleanText cleans and normalizes text content
func (p *htmlParser) cleanText(text string) string {
	// Trim leading/trailing whitespace
	text = strings.TrimSpace(text)
	
	// Replace multiple whitespace with single space
	// This includes newlines, tabs, etc.
	text = p.normalizeWhitespace(text)
	
	return text
}

func (p *htmlParser) normalizeWhitespace(text string) string {
	fields := strings.Fields(text)
		return strings.Join(fields, " ")
}