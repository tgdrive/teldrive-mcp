# Teldrive MCP Server Instructions

## Overview
The Teldrive MCP server enables interaction with your Teldrive storage through powerful tools. Key capabilities:
- Advanced file searching with multi-filter support
- Folder navigation and creation
- File content retrieval with automatic format handling

## Available Tools & Usage

### 1. Search Files Tool - Advanced Filtering
Find files/folders with combinable filters:
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>search_files</tool_name>
<arguments>
{
  "query": "project alpha",
  "searchType": "text",  // "text" (default) or "regex",
  "limit": 8,
  "page": 1,
  "limit": 20,
  "category": ["document", "image"],
  "type": "file"
}
</arguments>
</use_mcp_tool>
```

**Key Parameters**:
- `query`: Filename/text search
- `searchType`: `text` for literal matches (default), `regex` for regular expressions
- `category`: Multi-select: `["document","image","video","audio","archive","other"]`
- `type`: `file` or `folder`
- `limit/page`: Pagination control

### Examples:

#### Text Search (Default)
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>search_files</tool_name>
<arguments>
{
  "query": "Taylor Swift",
  "searchType": "text"
}
</arguments>
</use_mcp_tool>
```

#### Regex Search
Find all files with a 4-digit year prefix:
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>search_files</tool_name>
<arguments>
{
  "query": "\\d{4}-.*\\.txt",
  "searchType": "regex",
  "category": ["document"]
}
</arguments>
</use_mcp_tool>
```

#### Find specific file or folder
Find file or folder with exact name
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>search_files</tool_name>
<arguments>
{
  "name": "hello.pdf",
  "type": "file",
  "limit":1
}
</arguments>
<arguments>
{
  "name": "test",
  "type": "folder",
  "limit": 1
}
</arguments>
</use_mcp_tool>
```

### 2. List Files Tool - Folder Navigation
View folder contents (requires folder ID):
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>list_files</tool_name>
<arguments>
{
  "folder_id": "F:12345",
  "limit": 100
}
</arguments>
</use_mcp_tool>
```

### 3. Create Folder Tool
Create folders anywhere in your structure:
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>create_folder</tool_name>
<arguments>
{
  "path": "/Clients/Acme Corp",
  "name": "Contracts"
}
</arguments>
</use_mcp_tool>
```

### 4. Read File Tool
Retrieve content by file ID with automatic formatting:
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>read_file</tool_name>
<arguments>
{"file_id": "F:98765"}
</arguments>
</use_mcp_tool>
```
**Format Handling**:
- Text/JSON: Raw text
- Images: Base64 for direct embedding `![img](data:image/png;base64,...)`
- Audio: Base64 for playback via `<audio src="data:audio/mp3;base64,...">`
- Binary: Not directly readable (use with preview tools)

---

## Advanced Usage Examples

### Example 1: Multi-category Search
Find marketing PDFs and images:
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>search_files</tool_name>
<arguments>
{
  "query": "Q3 campaign",
  "category": ["document", "image"],
  "limit": 5
}
</arguments>
</use_mcp_tool>
```

### Example 2: Recent Video Search
Find newest videos (pagination example):
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>search_files</tool_name>
<arguments>
{
  "category": ["video"],
  "limit": 10,
  "page": 1
}
</arguments>
</use_mcp_tool>
```

### Example 3: Folder Creation + Verification
Create nested folders and confirm:
```xml
<!-- Create parent folder -->
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>create_folder</tool_name>
<arguments>
{"path": "/", "name": "Research"}
</arguments>
</use_mcp_tool>

<!-- Create subfolder -->
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>create_folder</tool_name>
<arguments>
{"path": "/Research", "name": "Competitor Analysis"}
</arguments>
</use_mcp_tool>

<!-- Verify structure -->
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>search_files</tool_name>
<arguments>
{"query": "Competitor Analysis", "type": "folder"}
</arguments>
</use_mcp_tool>
```

### Example 4: Folder Exploration
Navigate folder hierarchy:
```xml
<!-- Find project root -->
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>search_files</tool_name>
<arguments>
{
  "query": "Project Phoenix",
  "type": "folder",
  "limit": 1
}
</arguments>
</use_mcp_tool>

<!-- List contents (using ID from search) -->
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>list_files</tool_name>
<arguments>
{
  "folder_id": "F:54321",
  "limit": 50
}
</arguments>
</use_mcp_tool>
```

### Example 5: Archive Retrieval
Find and read archived reports:
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>search_files</tool_name>
<arguments>
{
  "query": "financial_audit",
  "category": ["archive"]
}
</arguments>
</use_mcp_tool>

<!-- Extract and read (using file ID from results) -->
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>read_file</tool_name>
<arguments>
{"file_id": "F:ZIP001"}
</arguments>
</use_mcp_tool>
```

### Example 6: Audio Batch Processing
Find and process podcast episodes:
```xml
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>search_files</tool_name>
<arguments>
{
  "query": "podcast_episode_*",
  "category": ["audio"],
  "limit": 3
}
</arguments>
</use_mcp_tool>

<!-- Process each found ID independently -->
<use_mcp_tool>
<server_name>teldrive</server_name>
<tool_name>read_file</tool_name>
<arguments>
{"file_id": "F:AUD123"}
</arguments>
</use_mcp_tool>
```

---

## Best Practices
1. **Precision Searching**: Combine `category+query` for focused results
   `Category Filtering Guide`:
   ┌─────────────┬──────────────────────────────────┐
   │ Category    │ Best For                         │
   ├─────────────┼──────────────────────────────────┤
   │ document    │ PDF, Word, Excel, Text files     │
   │ image       │ JPG, PNG, GIF visual content     │
   │ video       │ MP4, MOV, AVI videos             │
   │ audio       │ MP3, WAV, OGG sound files        │
   │ archive     │ ZIP, RAR compressed bundles      │
   │ other       │ Unclassified file types          │
   └─────────────┴──────────────────────────────────┘

2. **ID Handling**: Always reference files by ID (not names) after initial search/list

3. **Pagination**: For large results, use `limit` + `page` parameters
   `page=1 limit=20` → First 20 results
   `page=2 limit=20` → Next 20 results

4. **Resource Notation**: Reference files via URI in actions:
   `teldrive:///F:12345` (for downstream processing)

5. **File Size Awareness**:
   - Text files <1MB: Direct reading
   - Media files: Use preview/metadata extraction
   - >5MB files: Use chunked reading strategies

---
