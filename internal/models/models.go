package models

type PRRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Changes     []struct {
		FilePath string `json:"file_path"`
		ChangeType string `json:"change_type"` // add, modify, delete
	} `json:"changes"`
}

type PRResponse struct {
	Description string `json:"description"`
}

