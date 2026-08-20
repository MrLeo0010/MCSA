package src

import (
	"io"
)

type StatusResponse struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int            `json:"max"`
		Online int            `json:"online"`
		Sample []PlayerSample `json:"sample"`
	} `json:"players"`
	Description any `json:"description"` // Может быть строкой или объектом Chat

	PirateStatus       string `json:"-"`
	PirateStatusReason string `json:"-"`
}
type ProgressReader struct {
	Reader     io.Reader
	Total      int64
	Downloaded int64
	Prefix     string
}

var BaseGitHubURL = "https://github.com/{OWNER}/{REPO}"
var BaseGitHubAPIURL = "https://api.github.com/repos/{OWNER}/{REPO}"
