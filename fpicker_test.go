package fpicker

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

const (
	t_root string = "./test_path"
	t_path string = "./test_path/folder_one"
)

func TestHandleFolderPicker(t *testing.T) {
	absPath := filepath.ToSlash(filepath.Join(currentDir, t_path))
	url := FolderPickerUrl + "?path=" + t_path + "&root=" + t_root + "&hide=false"
	Setup(nil)

	rr, err := getResponse(url, handleFolderPicker)
	if err != nil {
		t.Fatal(err)
	}

	checkStatus(rr, t)

	body := readResponseBody(rr, t)

	// Page needs 4 folder icons:
	// ../ , folder11 , folder12 , *bottom bar selection*
	checkCount(body, `<use xlink:href="#icon-folder">`, 4, t)

	// Selected folder is stored in compilation time in a constant with value = initial absolute path
	checkString(body, `const selected = "`+absPath+`"`, t)

	// Checks folder names rendering
	checkString(body, `<span class="panel-files-label">folder11</span>`, t)
	checkString(body, `<span class="panel-files-label">folder12</span>`, t)
}

func TestHandleFilePicker(t *testing.T) {
	url := FilePickerUrl + "?path=" + t_path + "&root=" + t_root + "&hide=false"
	Setup(nil)

	rr, err := getResponse(url, handleFilePicker)
	if err != nil {
		t.Fatal(err)
	}

	checkStatus(rr, t)

	bodyStr := readResponseBody(rr, t)

	// Page needs 3 file icons:
	// file11.txt , file12.txt , *bottom bar selection*
	checkCount(bodyStr, `<use xlink:href="#icon-file">`, 3, t)

	// Selected file is stored in a variable at runtime, initial value = ""
	checkString(bodyStr, `let selected = ""`, t)

	// Checks file names render
	checkString(bodyStr, `<span class="panel-files-label">file11.txt</span>`, t)
	checkString(bodyStr, `<span class="panel-files-label">file12.txt</span>`, t)
}

/*
* --- HELPER FUNCTIONS ---
 */

func checkCount(body string, subString string, want int, t *testing.T) {
	got := strings.Count(body, subString)
	if got != want {
		t.Errorf("Expected %s count in body %d; got %d\nBody:\n%s", subString, want, got, body)
	}
}

func checkStatus(rr *httptest.ResponseRecorder, t *testing.T) {
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}
}

func checkString(body string, subString string, t *testing.T) {
	if !strings.Contains(body, subString) {
		t.Errorf("Body should contain: %s\nBody:\n%s", subString, body)
	}
}

func getResponse(url string, fn func(http.ResponseWriter, *http.Request)) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(fn)
	handler.ServeHTTP(rr, req)

	return rr, nil
}

func readResponseBody(rr *httptest.ResponseRecorder, t *testing.T) string {
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("The response body could not be read: %s", err)
	}

	return string(b)
}
