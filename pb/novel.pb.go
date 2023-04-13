package pb

type Novel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type NovelResource struct {
	ID          string `json:"id"`
	NovelID     string `json:"novel_id"`
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}
