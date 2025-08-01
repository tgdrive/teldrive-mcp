package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mcp/internal/api"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

type securitySource struct {
	token string
}

func (s *securitySource) BearerAuth(ctx context.Context, _ api.OperationName) (api.BearerAuth, error) {
	return api.BearerAuth{Token: s.token}, nil
}

func (s *securitySource) ApiKeyAuth(ctx context.Context, operationName api.OperationName) (api.ApiKeyAuth, error) {
	return api.ApiKeyAuth{APIKey: s.token}, nil
}

type file struct {
	ID string `json:"id" jsonschema_description:"The unique identifier for the file."`

	Name string `json:"name" jsonschema_description:"The name of the file or folder."`

	MimeType string `json:"mimeType" jsonschema_description:"The MIME type of the file, e.g., 'application/pdf' or 'image/jpeg'."`

	Size int64 `json:"size" jsonschema_description:"The total size of the file in bytes."`

	UpdatedAt time.Time `json:"updatedAt" jsonschema_description:"The timestamp of when the file was last modified."`
}

type meta struct {
	Count      int `json:"count" jsonschema_description:"The total number of items matching the query."`
	TotalPages int `json:"totalPages" jsonschema_description:"The total number of available pages based on the query limit."`
}

type fileList struct {
	Files []file `json:"files" jsonschema_description:"An array of file objects for the current page."`
	Meta  meta   `json:"meta" jsonschema_description:"Pagination metadata including total items, total pages, and the current page."`
}

func mapfileResponse(in *api.FileList) *fileList {
	res := fileList{}
	for _, f := range in.Items {
		res.Files = append(res.Files, file{ID: f.ID.Value, Name: f.Name,
			MimeType: f.MimeType.Value, Size: f.Size.Value, UpdatedAt: f.UpdatedAt.Value})
	}
	res.Meta.Count = in.Meta.Count
	res.Meta.TotalPages = in.Meta.TotalPages
	return &res
}

func main() {

	baseUrl := getEnv("BASE_URL")
	token := getEnv("AUTH_TOKEN")
	client, err := api.NewClient(baseUrl+"/api", &securitySource{token: token})
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
		return
	}
	mcpServer := server.NewMCPServer(
		"TelDrive MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)
	mcpServer.AddTool(mcp.NewTool("search_files",
		mcp.WithDescription("Search or Filter files and folders"),
		mcp.WithString("query",
			mcp.Description("File name or keyword to search"),
		),
		mcp.WithString("name",
			mcp.Description("Exact File name to find"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results to return"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number to return"),
		),
		mcp.WithArray("category",
			mcp.WithStringEnumItems([]string{"document", "image", "video", "audio", "archive", "other"}),
			mcp.Description("filter by category (e.g., document, image, video, audio, archive, other)"),
		),
		mcp.WithString("searchType", mcp.Description("text for literal matches (default), regex for regular expressions")),
		mcp.WithString("type", mcp.Description("filter by type (file or folder)")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.GetString("query", "")
		name := request.GetString("name", "")
		limit := request.GetInt("limit", 20)
		page := request.GetInt("page", 1)
		category := request.GetStringSlice("category", []string{})
		searchType := request.GetString("searchType", "text")
		fileType := request.GetString("type", "")
		params := api.FilesListParams{Query: api.NewOptString(query),
			Limit:     api.NewOptInt(limit),
			Page:      api.NewOptInt(page),
			Operation: api.NewOptFileQueryOperation(api.FileQueryOperationFind),
		}

		if len(category) > 0 {
			for _, cat := range category {
				params.Category = append(params.Category, api.Category(cat))
			}
		}
		if fileType != "" {
			params.Type = api.NewOptFileQueryType(api.FileQueryType(fileType))
		}
		if name != "" {
			params.Name = api.NewOptString(name)
		}
		if searchType != "" {
			params.SearchType = api.NewOptFileQuerySearchType(api.FileQuerySearchType(searchType))
		}
		result, err := client.FilesList(ctx, params)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(mapfileResponse(result))
		return mcp.NewToolResultText(string(data)), nil

	})
	mcpServer.AddTool(mcp.NewTool("list_files",
		mcp.WithDescription("List files in a folder"),
		mcp.WithOutputSchema[fileList](),
		mcp.WithString("folder_id",
			mcp.Description("ID of the folder to list files from"),
			mcp.Required(),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results to return"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number to return"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		folderID := request.GetString("folder_id", "")
		limit := request.GetInt("limit", 50)
		page := request.GetInt("page", 1)
		params := api.FilesListParams{
			ParentId: api.NewOptString(folderID),
			Limit:    api.NewOptInt(limit),
			Page:     api.NewOptInt(page),
		}
		result, err := client.FilesList(ctx, params)
		if err != nil {
			return nil, err
		}
		data, _ := json.Marshal(mapfileResponse(result))
		return mcp.NewToolResultText(string(data)), nil

	})

	mcpServer.AddTool(mcp.NewTool("create_folder",
		mcp.WithDescription("Create a new folder"),
		mcp.WithString("path",
			mcp.Description("Path to create the folder at, e.g., /folder1/folder2 , default is root (/)"),
			mcp.Required(),
		),
		mcp.WithString("name",
			mcp.Description("Name of the folder to create"),
			mcp.Required(),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := request.GetString("path", "")
		name := request.GetString("name", "")
		_, err := client.FilesCreate(ctx, &api.File{
			Name: name,
			Path: api.NewOptString(path),
			Type: api.FileTypeFolder,
		})
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText("folder created"), nil

	})

	mcpServer.AddTool(mcp.NewTool("read_file",
		mcp.WithDescription("Read a file from Teldrive using fileID"),
		mcp.WithString("file_id",
			mcp.Description("The ID of the file to read"),
			mcp.Required(),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileID := request.GetString("file_id", "")
		if fileID == "" {
			return mcp.NewToolResultError("file_id is required"), nil
		}
		res, err := readFile(ctx, client, token, fileID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		switch res.Category {
		case Text:
			return mcp.NewToolResultText(res.Content), nil
		case Image:
			return mcp.NewToolResultImage("", res.Content, res.MIMEType), nil
		case Audio:
			return mcp.NewToolResultAudio("", res.Content, res.MIMEType), nil
		}
		return mcp.NewToolResultError("unsupported content type"), nil

	})
	template := mcp.NewResourceTemplate(
		"tdrive:///{id}",
		"File Contents",
		mcp.WithTemplateDescription("Returns file content"),
		mcp.WithTemplateMIMEType("text/plain"),
	)

	mcpServer.AddResourceTemplate(template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		fileID := extractIDFromURI(request.Params.URI)
		if fileID == "" {
			return nil, fmt.Errorf("invalid file ID in URI: %s", request.Params.URI)
		}

		res, err := readFile(ctx, client, token, fileID)
		if err != nil {
			return nil, err
		}
		if res.Category == Text {
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      res.URI,
					MIMEType: res.MIMEType,
					Text:     res.Content,
				},
			}, nil
		} else {
			return []mcp.ResourceContents{
				mcp.BlobResourceContents{
					URI:      res.URI,
					MIMEType: res.MIMEType,
					Blob:     res.Content,
				},
			}, nil
		}
	})
	mux := http.NewServeMux()
	mux.Handle("/mcp/http", server.NewStreamableHTTPServer(mcpServer))
	mux.Handle("/mcp/sse", server.NewSSEServer(mcpServer))

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

type category int

const (
	Image category = iota
	Audio
	Text
)

type fileContent struct {
	URI      string
	MIMEType string
	Content  string
	Category category
}

func readFile(ctx context.Context, client *api.Client, token, fileID string) (*fileContent, error) {
	file, err := client.FilesStream(ctx, api.FilesStreamParams{
		ID:          fileID,
		AccessToken: api.NewOptString(token),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to stream file: %w", err)
	}
	f, ok := file.(*api.FilesStreamOKHeaders)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}
	bytes, err := io.ReadAll(f.Response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}
	content := ""
	var category category
	switch {
	case strings.HasPrefix(f.GetContentType(), "text/") || strings.HasPrefix(f.GetContentType(), "application/json"):
		content = string(bytes)
		category = Text

	case strings.HasPrefix(f.GetContentType(), "audio/"):
		content = base64.StdEncoding.EncodeToString(bytes)
		category = Audio

	case strings.HasPrefix(f.GetContentType(), "image/"):
		content = base64.StdEncoding.EncodeToString(bytes)
		category = Image

	default:
		return nil, fmt.Errorf("unsupported content type")
	}
	return &fileContent{URI: "tdrive:///" + fileID, MIMEType: f.GetContentType(), Content: content,
		Category: category}, nil

}
func extractIDFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	parts := strings.Split(u.Path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
