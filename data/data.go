package data

type Mode int

type SyncFile struct {
	RelativePath string `json:"relativePath"`
	Contents     string `json:"contents"`
}

type Client struct {
	Url string
}
