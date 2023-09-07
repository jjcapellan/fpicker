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
	cssVars    string = ""
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
}

func ResetCssVars() {
	cssVars = ""
}

func SetCssVars(vars string) {
	cssVars = vars
}

// Setup configures route handlers for file and folder selection functionalities
// provided by the fpicker package. Call this function to register the corresponding
// HTTP route handlers either on an existing HTTP multiplexer (http.ServeMux) or, if
// 'mux' is nil, directly on Go's global HTTP server (http.DefaultServeMux).
//
// When users make GET requests to the registered URLs, such as FilePickerUrl and
// FolderPickerUrl, they can retrieve the respective HTML content for file and folder
// pickers.
//
// Example usage:
//   mux := http.NewServeMux()
//   fpicker.Setup(mux)
//   http.ListenAndServe(":8080", mux)
//
// Parameters:
//   mux: The HTTP multiplexer where route handlers will be registered. If nil, the
//        handlers will be registered on http.DefaultServeMux.
func Setup(mux *http.ServeMux) {
	if mux == nil {
		http.HandleFunc(FilePickerUrl, handleFilePicker)
		http.HandleFunc(FolderPickerUrl, handleFolderPicker)
		return
	}
	mux.HandleFunc(FilePickerUrl, handleFilePicker)
	mux.HandleFunc(FolderPickerUrl, handleFolderPicker)
}
