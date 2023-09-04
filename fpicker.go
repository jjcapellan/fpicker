package fpicker

import (
	"net/http"
	"os"
	"text/template"

	"github.com/jjcapellan/fsinfo"
)

var (
	drives     []fsinfo.DriveInfo
	home       string
	tmpl       *template.Template
	currentDir string
)

const (
	FilePickerUrl     string = "/fpicker/file-picker"
	FolderPickerUrl   string = "/fpicker/folder-picker"
	SelectedFileUrl   string = "/fpicker/selected-file"
	SelectedFolderUrl string = "/fpicker/selected-folder"
)

func init() {
	drives, _ = fsinfo.GetDrives()
	home, _ = fsinfo.GetHomePath()
	tmpl, _ = initTemplates()
	currentDir, _ = os.Getwd()

	http.HandleFunc(FilePickerUrl, handleFilePicker)
	http.HandleFunc(FolderPickerUrl, handleFolderPicker)
}
