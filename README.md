# fpicker

This package provides functionality for selecting files and folders in a web application. It allows users to explore and select full path of files and folders in the local file system.  

**Security Notice:**
Due to the fact that this package accesses the local file system for read-only purposes, for security reasons, it should only be used for local web applications without external access.

## Usage
This package provides the following api endpoints to interact with the file or folder picker in your web application:
- (GET)  **/fpicker/file-picker** : Fetch file picker page. fpicker registers its own handler to the DefaultServeMux.
- (GET)  **/fpicker/folder-picker** : Fetch folder picker page. fpicker registers its own handler to the DefaultServeMux.
- (POST) **/fpicker/selected-file?path={full path of selected file}** : this request is sent by "select" button of the picker dialog.
- (POST) **/fpicker/selected-folder?path={full path of selected folder}** : this request is sent by "select" button of the picker dialog.

## Example
This basic example demonstrates how to open the file picker in a new window or tab when a link on the web page is clicked. When the user makes a selection in the file picker, the selection is displayed on the main web page using Server-Sent Events (SSE). This approach provides a simple way to interact with the local file system in a Go web application.
### Content of */public/index.html*

```html
<!-- Existing HTML code goes here -->

<!-- Add a link to open the file picker in a new window or tab -->
<a href="/fpicker/file-picker" target="_blank">Open file picker</a>
<!-- An element to display the selection -->
<span id="content"></span>

<script>
	const content = document.getElementById("content");
	// Set up an EventSource to receive real-time updates
	const eventSource = new EventSource("/sse");
	eventSource.addEventListener("file", (evt) => {
		// Display the selection on the main web page when an event arrives
		content.innerText = evt.data;
	});
</script>

<!-- More existing HTML code -->
```
### Content of */main.go*

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/jjcapellan/fpicker"
)

var ch chan string = make(chan string)

func main() {
	// Serve static files from the "/public" folder
	fs := http.FileServer(http.Dir("./public"))

	// Route that receives the selected file path
	http.HandleFunc(fpicker.SelectedFileUrl, handleSelectFile)

	// Route for Server-Sent Events (SSE)
	http.HandleFunc("/sse", handleSSE)

	// Default route to serve static files
	http.Handle("/", fs)

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleSelectFile(w http.ResponseWriter, r *http.Request) {
	// Get the path of the selected file
	filePath := r.URL.Query().Get("path")

	// Send the path to SSE using the "ch" channel
	eventStr := "event: file\ndata: " + filePath + "\n\n"
	ch <- eventStr
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	select {
	case eventStr := <-ch:
		// Send a new event to the web frontend
		fmt.Fprint(w, eventStr)
		w.(http.Flusher).Flush()
	case <-r.Context().Done():
		return
	}
}
```

### Screenshot of the file-picker dialog
[<img src="readme_imgs/fpicker_src1.png" width=500/>](readme_imgs/fpicker_src1.png)