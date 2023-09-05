# fpicker

This golang package provides a file selection dialog for web applications, granting access to complete local file system paths. Unlike standard HTML <code><input type="file" ...></code> dialogs, which restrict access to full file paths, fpicker does not have this limitation.

**Important Note:**
Please be aware that this package does not perform any write, delete, or modification operations on local files. However, since it exposes the directory structure of the file system, it should only be used in secure, local environments. Recommended scenarios include integration into emulator frontends, system utilities, and similar use cases.  

<p align="center"><img src="readme_imgs/fpicker_capture.gif"></p>  
<br>

## Features
- [x] Web dialogs for file and folder selection
- [x] Quick access to the Home directory and disk drives
- [x] Hidden file filtering
- [x] Supports Linux and Windows (Other OS not tested)
- [ ] Visual styles customization  

<br>

## Installation
```bash
go get github.com/jjcapellan/fpicker
```
<br>

## Usage
1. Call the function <code>Setup(mux *http.ServeMux)</code> (if mux == nil then fpicker will use http.DefaultServeMux). *Setup* configures route handlers for file and folder selection functionalities.
2. Make a get request to **/fpicker/file-picker** or **/fpicker/folder-picker** to retrieve the file picker or folder picker respectively. The *Setup* function, in the preceding step, added these URLs to the multiplexer's routing configuration.
3. After the user makes the selection of the file or folder pressing the select button, a post request with the selection is sent to one of this urls: **/fpicker/selected-file?path={full path of selected file}** or **/fpicker/selected-folder?path={full path of selected folder}** respectively. You must implement handlers for this routes.
<br><br>

## Example
This example corresponds to the content of the animated gif in this readme. Clicking the button opens the file picker in a new window, and after selecting a file, this window closes, and the file path is displayed on the screen.  
<br>

### Content of file "/public/index.html"
For the sake of simplicity, all the CSS, JavaScript, and HTML code has been included in the same file.  

```html
<!-- some html code here[...] -->
<body>
    <div class="container">
        <button class="button" id="bt-open-file">Open filePicker</button>
        <span id="content">No file selected</span>
    </div>

    <script>
        const button = document.getElementById("bt-open-file");
        const content = document.getElementById("content");

        // The file picker sends the selected file path to the backend, 
		// which forwards it to this page through a server-side events (SSE) stream.
        const eventSource = new EventSource("/sse");
        eventSource.addEventListener("file", (evt) => {
            content.innerText = evt.data;
        });

        button.addEventListener("click", handler);

        function handler() {
            const width = 1100,
                  height = 800,
				  offset = 48, // Aprox.
                  left = window.innerWidth / 2 - width / 2,
                  top = window.innerHeight / 2 - height / 2 + offset;
			// The URL "/fpicker/file-picker" retrieves the file picker page
            window.open("/fpicker/file-picker", "Select file", `width=${width},height=${height},left=${left},top=${top}`);
        }
    </script>
</body>
</html>
```  
<br>

### Content of file "main.go"
The file picker doesn't directly send the selection back to the invoking page but instead utilizes the backend as an intermediary. To relay the file picker's selection back to the client, the backend inserts it into the server-sent events (SSE) stream being listened to by the client.  


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

	fs := http.FileServer(http.Dir("./public"))
	mux := http.NewServeMux()
	
	// Setup configures route handlers for file and folder selection functionalities
    // provided by the fpicker package. Call this function to register the corresponding
    // HTTP route handlers either on an existing HTTP multiplexer (http.ServeMux) or, if
    // 'mux' is nil, directly on Go's global HTTP server (http.DefaultServeMux).
    //
    // When users make GET requests to the registered URLs, such as FilePickerUrl and
    // FolderPickerUrl, they can retrieve the respective HTML content for file and folder
    // pickers.
	fpicker.Setup(mux)
    
	// You must handle the route where the selectionn is sent by the file picker
	mux.HandleFunc(fpicker.SelectedFileUrl, handleSelectFile)

	// In this example the route "/sse" is used by an events stream as way to send data to the client
	mux.HandleFunc("/sse", handleSSE)

	mux.Handle("/", fs)
    
	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// This function forwards the file picker selection (path) to the "ch" channel within a server-side event
func handleSelectFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	eventStr := "event: file\ndata: " + filePath + "\n\n"
	ch <- eventStr
}

// This function sends the event containing the file picker selection to the client
func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	select {
	case eventStr := <-ch:
		fmt.Fprint(w, eventStr)
		w.(http.Flusher).Flush()
	case <-r.Context().Done():
		return
	}
}
```
<br>

## Credits
Appreciation to the [Heroicons](https://heroicons.com/) team for their valuable icon collection used in this project.  
<br>

## License
**fpicker** is licensed under the terms of the [MIT](https://opensource.org/licenses/MIT) license.