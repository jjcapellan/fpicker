package fpicker

import (
	"embed"
	"path/filepath"
	"text/template"
	"time"

	"github.com/jjcapellan/fsinfo"
)

//go:embed templates
var templateFiles embed.FS

func initTemplates() (*template.Template, error) {
	return template.New("tmpl").Funcs(template.FuncMap{
		"formatBytes": fsinfo.FormatBytes,
		"formatTime":  formatTime,
		"getParent":   getParent,
	}).ParseFS(templateFiles, "templates/*")
}

func formatTime(t time.Time) string {
	return t.Format("02/01/2006 15:04")
}

func getParent(path string) string {
	path = filepath.Dir(path)
	path = filepath.ToSlash(path)
	return path
}
