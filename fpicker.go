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

// ResetCssVars resets CSS variables to default values.
func ResetCssVars() {
	cssVars = ""
}

// SetCssVars allows you to override default CSS variables within the template.
// In this way you can customize the color scheme.
//
// Parameters:
//   vars (string): A string containing CSS variable declarations with the format
//                 "--variable-name: value;", where "variable-name" is the name of the
//                 CSS variable to override, and "value" is the new value to set.
//
// These are the default css variables:
//   --color-dark: #2a363c;
//
//   /* Body */
//   --background-color: #a2a2a2;
//
//   /* Content */
//   --background-color-content: #f6f8f3;
//   --color-content: var(--color-dark);
//
//   /* Sidebar */
//   --background-color-sidebar: #e1e9ec;
//   --color-sidebar: var(--color-dark);
//   --color-sidebar-hover: #c51818;
//
//   /* Bars */
//   --background-color-bar: #546c78;
//   --color-bar: var(--background-color-sidebar);
//
//   /* Buttons */
//   --background-color-button: var(--color-dark);
//   --color-button: var(--background-color-sidebar);
//   --background-color-button-hover: #9fb2bc;
//   --color-button-hover: var(--color-dark);
//
// Example:
//   fpicker.SetCssVars("--background-color: #cc0000; --color-sidebar-hover: #c50000;")
//
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
