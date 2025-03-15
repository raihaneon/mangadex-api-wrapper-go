package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// Base URL for MangaDex API
const (
	BaseURL = "https://api.mangadex.org"
)

// MangaDex API response structures
type MangaResponse struct {
	Result   string      `json:"result"`
	Response string      `json:"response"`
	Data     []MangaData `json:"data"`
	Limit    int         `json:"limit"`
	Offset   int         `json:"offset"`
	Total    int         `json:"total"`
}

type MangaData struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Attributes    MangaAttributes `json:"attributes"`
	Relationships []Relationship  `json:"relationships"`
}

type MangaAttributes struct {
	Title                  map[string]string   `json:"title"`
	AltTitles              []map[string]string `json:"altTitles"`
	Description            map[string]string   `json:"description"`
	IsLocked               bool                `json:"isLocked"`
	Links                  map[string]string   `json:"links"`
	OriginalLanguage       string              `json:"originalLanguage"`
	LastVolume             string              `json:"lastVolume"`
	LastChapter            string              `json:"lastChapter"`
	PublicationDemographic string              `json:"publicationDemographic"`
	Status                 string              `json:"status"`
	Year                   int                 `json:"year"`
	ContentRating          string              `json:"contentRating"`
	Tags                   []TagData           `json:"tags"`
	CreatedAt              string              `json:"createdAt"`
	UpdatedAt              string              `json:"updatedAt"`
}

type TagData struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Attributes TagAttributes `json:"attributes"`
}

type TagAttributes struct {
	Name        map[string]string `json:"name"`
	Description map[string]string `json:"description"`
	Group       string            `json:"group"`
}

type Relationship struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type ChapterResponse struct {
	Result   string        `json:"result"`
	Response string        `json:"response"`
	Data     []ChapterData `json:"data"`
	Limit    int           `json:"limit"`
	Offset   int           `json:"offset"`
	Total    int           `json:"total"`
}

type ChapterData struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Attributes    ChapterAttributes `json:"attributes"`
	Relationships []Relationship    `json:"relationships"`
}

type ChapterAttributes struct {
	Title              string `json:"title"`
	Volume             string `json:"volume"`
	Chapter            string `json:"chapter"`
	Pages              int    `json:"pages"`
	TranslatedLanguage string `json:"translatedLanguage"`
	ExternalURL        string `json:"externalUrl"`
	PublishAt          string `json:"publishAt"`
	ReadableAt         string `json:"readableAt"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type ChapterPages struct {
	Result  string  `json:"result"`
	BaseURL string  `json:"baseUrl"`
	Chapter Chapter `json:"chapter"`
}

type Chapter struct {
	Hash      string   `json:"hash"`
	Data      []string `json:"data"`
	DataSaver []string `json:"dataSaver"`
}

// API client struct
type MangaDexClient struct {
	httpClient *http.Client
	browser    *rod.Browser
}

// Initialize a new MangaDex client
func NewMangaDexClient() *MangaDexClient {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	// Launch browser headless mode for Rod
	path := launcher.New().Set("--no-sandbox").Headless(true).MustLaunch()
	browser := rod.New().ControlURL(path).MustConnect()

	return &MangaDexClient{
		httpClient: client,
		browser:    browser,
	}
}

// Close the Rod browser
func (mdc *MangaDexClient) Close() {
	mdc.browser.MustClose()
}

// Search manga by title
func (mdc *MangaDexClient) SearchManga(title string, limit int) (*MangaResponse, error) {
	url := fmt.Sprintf("%s/manga?title=%s&limit=%d&includes[]=cover_art&includes[]=author",
		BaseURL, title, limit)

	resp, err := mdc.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mangaResp MangaResponse
	if err := json.NewDecoder(resp.Body).Decode(&mangaResp); err != nil {
		return nil, err
	}

	return &mangaResp, nil
}

// Get manga by ID
func (mdc *MangaDexClient) GetMangaByID(id string) (*MangaResponse, error) {
	url := fmt.Sprintf("%s/manga/%s?includes[]=cover_art&includes[]=author", BaseURL, id)

	resp, err := mdc.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mangaResp MangaResponse
	if err := json.NewDecoder(resp.Body).Decode(&mangaResp); err != nil {
		return nil, err
	}

	return &mangaResp, nil
}

// Get chapters for a manga
func (mdc *MangaDexClient) GetChapters(mangaID string, translatedLanguage string, limit int, offset int) (*ChapterResponse, error) {
	url := fmt.Sprintf("%s/manga/%s/feed?translatedLanguage[]=%s&limit=%d&offset=%d&order[volume]=desc&order[chapter]=desc",
		BaseURL, mangaID, translatedLanguage, limit, offset)

	resp, err := mdc.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var chapterResp ChapterResponse
	if err := json.NewDecoder(resp.Body).Decode(&chapterResp); err != nil {
		return nil, err
	}

	return &chapterResp, nil
}

// Get pages for a chapter
func (mdc *MangaDexClient) GetChapterPages(chapterID string) (*ChapterPages, error) {
	url := fmt.Sprintf("%s/at-home/server/%s", BaseURL, chapterID)

	resp, err := mdc.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pagesResp ChapterPages
	if err := json.NewDecoder(resp.Body).Decode(&pagesResp); err != nil {
		return nil, err
	}

	return &pagesResp, nil
}

// Download a manga image using Rod (for cases when direct API access fails)
func (mdc *MangaDexClient) DownloadMangaImage(url, outputPath string) error {
	page := mdc.browser.MustPage(url)
	defer page.Close()

	page.MustWaitStable()

	// Wait for the image to load
	img := page.MustElement("img")
	img.MustWaitVisible()

	// Take a screenshot of the image
	data := img.MustScreenshot()

	// Save the image
	return os.WriteFile(outputPath, data, 0644)
}

// Main function
func main() {
	r := gin.Default()

	// Create MangaDex client
	client := NewMangaDexClient()
	defer client.Close()

	// API routes
	api := r.Group("/api")
	{
		// Search manga
		api.GET("/manga/search", func(c *gin.Context) {
			title := c.Query("title")
			limitStr := c.DefaultQuery("limit", "10")
			limit, _ := strconv.Atoi(limitStr)

			results, err := client.SearchManga(title, limit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, results)
		})

		// Get manga by ID
		api.GET("/manga/:id", func(c *gin.Context) {
			id := c.Param("id")

			manga, err := client.GetMangaByID(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, manga)
		})

		// Get chapters for a manga
		api.GET("/manga/:id/chapters", func(c *gin.Context) {
			id := c.Param("id")
			lang := c.DefaultQuery("lang", "en")
			limitStr := c.DefaultQuery("limit", "30")
			offsetStr := c.DefaultQuery("offset", "0")

			limit, _ := strconv.Atoi(limitStr)
			offset, _ := strconv.Atoi(offsetStr)

			chapters, err := client.GetChapters(id, lang, limit, offset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, chapters)
		})

		// Get pages for a chapter
		api.GET("/chapter/:id/pages", func(c *gin.Context) {
			id := c.Param("id")

			pages, err := client.GetChapterPages(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, pages)
		})

		// Download a chapter page
		api.GET("/chapter/:id/download/:page", func(c *gin.Context) {
			id := c.Param("id")
			pageNum := c.Param("page")
			quality := c.DefaultQuery("quality", "data")

			// Get chapter pages info first
			pages, err := client.GetChapterPages(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var pageURL string
			pageNumInt, _ := strconv.Atoi(pageNum)

			if pageNumInt < 0 || pageNumInt >= len(pages.Chapter.Data) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
				return
			}

			if quality == "data-saver" && len(pages.Chapter.DataSaver) > pageNumInt {
				pageURL = fmt.Sprintf("%s/data-saver/%s/%s",
					pages.BaseURL, pages.Chapter.Hash, pages.Chapter.DataSaver[pageNumInt])
			} else {
				pageURL = fmt.Sprintf("%s/data/%s/%s",
					pages.BaseURL, pages.Chapter.Hash, pages.Chapter.Data[pageNumInt])
			}

			// Download the image
			resp, err := http.Get(pageURL)
			if err != nil {
				// Fallback to Rod if direct download fails
				tempFile := fmt.Sprintf("temp_%s_%s.jpg", id, pageNum)
				err = client.DownloadMangaImage(pageURL, tempFile)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download page"})
					return
				}

				// Serve the file
				c.File(tempFile)
				// Clean up
				defer os.Remove(tempFile)
				return
			}
			defer resp.Body.Close()

			// Set content type
			contentType := resp.Header.Get("Content-Type")
			if contentType == "" {
				contentType = "image/jpeg"
			}
			c.Header("Content-Type", contentType)

			// Copy the image data to response
			io.Copy(c.Writer, resp.Body)
		})
	}

	// Start server
	r.Run(":8080")
}
