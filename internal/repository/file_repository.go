package repository

import (
	"mock-api/internal/config"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type FileRepository interface {
	Find(filepath string, result interface{}) error
	FindFilePath(root, path, method string, queries map[string]string) string
}

type fileRepository struct {
	cfg *config.Config
}

func NewFileRepository(cfg *config.Config) FileRepository {
	return &fileRepository{cfg: cfg}
}

func (f *fileRepository) Find(filepath string, result interface{}) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}

func (f *fileRepository) FindFilePath(root, path, method string, queries map[string]string) string {
	method = strings.ToLower(method)
	queryString := getSortedQueryString(queries)

	path = strings.Trim(path, "/")
	if path == "" {
		path = "index"
	}

	var files []string
	if len(queryString) > 0 {
		files = append(files,
			filepath.Join(root, path+"_"+method+"_"+queryString+".json"),
			filepath.Join(root, path, "index_"+method+"_"+queryString+".json"),
		)
	} else {
		files = append(files,
			filepath.Join(root, path+"_"+method+".json"),
			filepath.Join(root, path, "index_"+method+".json"),
		)
	}

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}
	return ""
}

func getSortedQueryString(queries map[string]string) string {
	if len(queries) == 0 {
		return ""
	}
	query := url.Values{}
	for key, value := range queries {
		query.Add(key, value)
	}
	return strings.ReplaceAll(query.Encode(), "&", "_")
}
