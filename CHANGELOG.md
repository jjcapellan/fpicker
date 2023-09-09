# v0.2.0
## Added
- Customization of default css variables using the functions `SetCssVars(vars string)` and `ResetCssVars()`.  

<br>

# v0.1.0
## Added
- Template for file-piker and folder-picker (css, js, and html)
- 4 api routes:
    * `FilePickerUrl`     "/fpicker/file-picker"
    * `FolderPickerUrl`   "/fpicker/folder-picker"
    * `SelectedFileUrl`   "/fpicker/selected-file"
    * `SelectedFolderUrl` "/fpicker/selected-folder"
- Implemented handlers for `FilePickerUrl` and `FolderPickerUrl`
- `Setup(mux *http.ServeMux)` function sets the servemux used to configure the api routes.
- Integration of `fsinfo` package.