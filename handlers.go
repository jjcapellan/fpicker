package fpicker

import (
	"net/http"
	"path/filepath"

	"github.com/jjcapellan/fsinfo"
)

func handleFolderPicker(w http.ResponseWriter, r *http.Request) {
	handlePicker(w, r, false)
}

func handleFilePicker(w http.ResponseWriter, r *http.Request) {
	handlePicker(w, r, true)
}

func handlePicker(w http.ResponseWriter, r *http.Request, isFilePicker bool) {
	path := r.URL.Query().Get("path")
	root := r.URL.Query().Get("root")
	if path == "" || root == "" {
		root := home
		show(w, root, home, isFilePicker)
		return
	}
	show(w, root, path, isFilePicker)
}

func show(w http.ResponseWriter, root string, path string, isFilePicker bool) {
	root = pathToAbs(root)
	path = pathToAbs(path)

	breadcrumb := makeBreadcrumb(root, path)
	currentFolder := filepath.Base(path)

	fi, _ := fsinfo.GetFolderInfo(path)

	data := &struct {
		Root              string
		Path              string
		Folder            string
		Home              string
		Breadcrumb        []string
		Drives            []fsinfo.DriveInfo
		Folders           []fsinfo.Folder
		Files             []fsinfo.File
		FilePickerUrl     string
		FolderPickerUrl   string
		SelectedFileUrl   string
		SelectedFolderUrl string
		IsFilePicker      bool
	}{
		root,
		path,
		currentFolder,
		home,
		breadcrumb,
		drives,
		fi.Folders,
		fi.Files,
		FilePickerUrl,
		FolderPickerUrl,
		SelectedFileUrl,
		SelectedFolderUrl,
		isFilePicker,
	}
	tmpl.ExecuteTemplate(w, "index.tmpl", data)
}
