# MangaDex API Client

A Go-based REST API client for interacting with the MangaDex API. This service allows you to search for manga, retrieve manga details, fetch chapters, and download manga pages.

## Features

- Search manga by title
- Get detailed manga information
- Retrieve chapters for a specific manga
- Get pages for a specific chapter
- Download manga pages with fallback mechanism
- Support for different image qualities

## Requirements

- Go 1.16+
- [Gin](https://github.com/gin-gonic/gin) web framework
- [Rod](https://github.com/go-rod/rod) headless browser automation library

## Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/mangadex-api-client.git
cd mangadex-api-client
```

2. Install dependencies
```bash
go mod tidy
```

3. Run the application
```bash
go run main.go
```

The server will start on port 8080.

## API Endpoints

### Search Manga
```
GET /api/manga/search
```
Parameters:
- `title`: Manga title to search for
- `limit`: Maximum number of results (default: 10)

Example:
```
GET /api/manga/search?title=One%20Piece&limit=5
```

### Get Manga by ID
```
GET /api/manga/:id
```
Parameters:
- `id`: Manga ID

Example:
```
GET /api/manga/a96676e5-8ae2-425e-b549-7f15dd34a6d8
```

### Get Chapters for a Manga
```
GET /api/manga/:id/chapters
```
Parameters:
- `id`: Manga ID
- `lang`: Translation language (default: "en")
- `limit`: Maximum number of chapters (default: 30)
- `offset`: Pagination offset (default: 0)

Example:
```
GET /api/manga/a96676e5-8ae2-425e-b549-7f15dd34a6d8/chapters?lang=en&limit=20&offset=0
```

### Get Pages for a Chapter
```
GET /api/chapter/:id/pages
```
Parameters:
- `id`: Chapter ID

Example:
```
GET /api/chapter/b58d92d9-f351-4d3c-a8dd-9c7e7f5efd8d/pages
```

### Download a Chapter Page
```
GET /api/chapter/:id/download/:page
```
Parameters:
- `id`: Chapter ID
- `page`: Page number (starting from 0)
- `quality`: Image quality (default: "data", alternatives: "data-saver" for lower quality)

Example:
```
GET /api/chapter/b58d92d9-f351-4d3c-a8dd-9c7e7f5efd8d/download/0?quality=data-saver
```

## Implementation Details

- Uses the official MangaDex API at `https://api.mangadex.org`
- Implements fallback to headless browser (Rod) when direct API image downloads fail
- Organized API responses using Go structs for proper type handling
- Implements pagination for manga searches and chapter listings

## Example Usage

### Search for a manga
```bash
curl "http://localhost:8080/api/manga/search?title=Naruto&limit=3"
```

### Get manga details
```bash
curl "http://localhost:8080/api/manga/a96676e5-8ae2-425e-b549-7f15dd34a6d8"
```

### Get chapters
```bash
curl "http://localhost:8080/api/manga/a96676e5-8ae2-425e-b549-7f15dd34a6d8/chapters?lang=en&limit=5"
```

### Download a page
```bash
curl -o page.jpg "http://localhost:8080/api/chapter/b58d92d9-f351-4d3c-a8dd-9c7e7f5efd8d/download/0"
```

## License

MIT

## Disclaimer

This project is not affiliated with or endorsed by MangaDex. It is created for educational purposes. Please respect MangaDex's terms of service and API usage policies.