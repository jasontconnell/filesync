package data

type SyncFile struct {
	RelativePath string `json:"relativePath"`
	Contents     string `json:"contents"`
	Delete       bool   `json:"delete"`
	Type         string `json:"type"`
}

type Client struct {
	Url string
}
