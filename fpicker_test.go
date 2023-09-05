package fpicker

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleFolderPicker(t *testing.T) {

	root := "./test_path"
	path := "./test_path/folder_one"
	absPath := filepath.ToSlash(filepath.Join(currentDir, path))
	url := FolderPickerUrl + "?path=" + path + "&root=" + root + "&hide=false"
	Setup(nil)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleFolderPicker)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}

	b, _ := io.ReadAll(rr.Body)

	bodyStr := string(b)

	// Page needs 4 folder icons:
	// ../ , folder11 , folder12 , *bottom bar selection*
	htmlSubstr := `<use xlink:href="#icon-folder">`
	got := strings.Count(bodyStr, htmlSubstr)
	want := 4
	if got != want {
		t.Errorf("Expected #icon-folder count in body %d; got %d\nBody:\n%s", want, got, bodyStr)
	}

	// Selected folder is stored in compilation time in a constant with value = initial absolute path
	jsSubstr := `const selected = "` + absPath + `"`
	if !strings.Contains(bodyStr, jsSubstr) {
		t.Errorf("Body should contain: %s\nBody:\n%s", jsSubstr, bodyStr)
	}

	// Checks folder names rendering
	htmlSubstr = `<span class="panel-files-label">folder11</span>`
	if !strings.Contains(bodyStr, htmlSubstr) {
		t.Errorf("Body should contain: %s\nBody:\n%s", htmlSubstr, bodyStr)
	}
	htmlSubstr = `<span class="panel-files-label">folder12</span>`
	if !strings.Contains(bodyStr, htmlSubstr) {
		t.Errorf("Body should contain: %s\nBody:\n%s", htmlSubstr, bodyStr)
	}
}

func TestHandleFilePicker(t *testing.T) {

	root := "./test_path"
	path := "./test_path/folder_one"
	url := FilePickerUrl + "?path=" + path + "&root=" + root + "&hide=false"
	Setup(nil)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleFilePicker)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}

	b, _ := io.ReadAll(rr.Body)

	bodyStr := string(b)

	// Page needs 3 file icons:
	// file11.txt , file12.txt , *bottom bar selection*
	htmlSubstr := `<use xlink:href="#icon-file">`
	got := strings.Count(bodyStr, htmlSubstr)
	want := 3
	if got != want {
		t.Errorf("Expected #icon-file count in body %d; got %d\nBody:\n%s", want, got, bodyStr)
	}

	// Selected file is stored in a variable at runtime, initial value = ""
	jsSubstr := `let selected = ""`
	if !strings.Contains(bodyStr, jsSubstr) {
		t.Errorf("Body should contain: %s\nBody:\n%s", jsSubstr, bodyStr)
	}

	// Checks file names render
	htmlSubstr = `<span class="panel-files-label">file11.txt</span>`
	if !strings.Contains(bodyStr, htmlSubstr) {
		t.Errorf("Body should contain: %s\nBody:\n%s", htmlSubstr, bodyStr)
	}
	htmlSubstr = `<span class="panel-files-label">file12.txt</span>`
	if !strings.Contains(bodyStr, htmlSubstr) {
		t.Errorf("Body should contain: %s\nBody:\n%s", htmlSubstr, bodyStr)
	}
}
