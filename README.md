# TelDrive MCP Server

A Model Context Protocol (MCP) server for TelDrive, enabling seamless integration with AI assistants for file management operations.

## Features

- **File Search**: Search and filter files by name, content, category, and type
- **File Listing**: List files within specific folders with pagination
- **Folder Creation**: Create new folders at specified paths
- **File Reading**: Read file contents with support for text, images, and audio
- **Resource Templates**: Automatic content loading via URI references
- **Multiple Content Types**: Support for text, image, and audio file formats

## Installation

### macOS/Linux
```sh
curl -sSL instl.vercel.app/tgdrive/teldrive-mcp | bash
```
### Windows
```powershell
powershell -c "irm https://instl.vercel.app/tgdrive/teldrive-mcp?platform=windows|iex"
```

## Configuration

### Starting the Server

##### Stdio Transport
To start the server with standard input/output transport, run:
```bash
./teldrive-mcp -base-url $BASE_URL -token $AUTH_TOKEN
```

#### HTTP Transport
To start the server with HTTP transport, run:
```bash
./teldrive-mcp -transport http -base-url $BASE_URL -token $AUTH_TOKEN 
```

The server will start on port 8080 with the following endpoints:
- HTTP Streaming: `http://localhost:8080/mcp/http`
- SSE: `http://localhost:8080/mcp/sse`
### Available Tools

#### 1. Search Files
Search for files and folders with various filters.

**Parameters:**
- `query` (string): File name or keyword to search
- `name` (string): Exact file name to find
- `limit` (number): Maximum results to return (default: 20)
- `page` (number): Page number to return (default: 1)
- `category` (array): Filter by category (document, image, video, audio, archive, other)
- `searchType` (string): "text" for literal matches, "regex" for regular expressions
- `type` (string): Filter by type (file or folder)

**Example:**
```json
{
  "name": "search_files",
  "arguments": {
    "query": "report",
    "category": ["document"],
    "limit": 10
  }
}
```

#### 2. List Files
List files within a specific folder.

**Parameters:**
- `folder_id` (string, required): ID of the folder to list
- `limit` (number): Maximum results to return (default: 50)
- `page` (number): Page number to return (default: 1)

**Example:**
```json
{
  "name": "list_files",
  "arguments": {
    "folder_id": "folder123",
    "limit": 20
  }
}
```

#### 3. Create Folder
Create a new folder at the specified path.

**Parameters:**
- `path` (string, required): Path to create the folder (e.g., /folder1/folder2)
- `name` (string, required): Name of the folder to create

**Example:**
```json
{
  "name": "create_folder",
  "arguments": {
    "path": "/documents",
    "name": "reports"
  }
}
```

#### 4. Read File
Read the content of a file by its ID.

**Parameters:**
- `file_id` (string, required): ID of the file to read

**Example:**
```json
{
  "name": "read_file",
  "arguments": {
    "file_id": "file456"
  }
}
```

### Resource Templates

The server provides a resource template for accessing file content:

- **URI Pattern**: `tdrive:///{id}`
- **Description**: Returns file content
- **Supported Types**: Text, images, and audio files

When an AI assistant references a URI like `tdrive:///file123`, the server automatically fetches and returns the file content.

## Response Schema

### File List Response
```json
{
  "files": [
    {
      "id": "file123",
      "name": "document.pdf",
      "mimeType": "application/pdf",
      "size": 1024,
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "count": 1,
    "totalPages": 1
  }
}
```

## Examples

### Searching for Documents
```json
{
  "name": "search_files",
  "arguments": {
    "query": "quarterly report",
    "category": ["document"],
    "limit": 5
  }
}
```

### Listing Folder Contents
```json
{
  "name": "list_files",
  "arguments": {
    "folder_id": "root",
    "limit": 10
  }
}
```

### Creating a New Folder
```json
{
  "name": "create_folder",
  "arguments": {
    "path": "/projects",
    "name": "2023"
  }
}
```

### Reading a File
```json
{
  "name": "read_file",
  "arguments": {
    "file_id": "file789"
  }
}
```
## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
