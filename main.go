package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/aki017/assetbundle"
)

var Root string

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	if len(req.URL.Path) <= 2 {
		http.Redirect(w, req, "/_/index.html", http.StatusFound)
	} else if req.URL.Path[1:2] == "_" {
		logrus.Infof("nothing")
	} else if req.URL.Path[1:2] == "*" {
		path := Root + req.URL.Path[2:]

		if info, err := os.Stat(path); os.IsNotExist(err) {
			io.WriteString(w, path)
			return
		} else if info.IsDir() {
			io.WriteString(w, "-1")
			return
		} else {
			ab := assetbundle.DecodeFile(path)
			if len(ab.Bodies) != 1 {
				io.WriteString(w, "-1")
				return
			}
			crc := ab.Bodies[0].CRC()
			b, err := json.Marshal(map[string]interface{}{
				"crc":    crc,
				"crchex": "0x" + strconv.FormatUint(uint64(crc), 16),
				"body":   ab.Bodies[0],
			})
			if err != nil {
				logrus.Warnf("AssetBundle Marshal error %s", err)
			}
			w.Write(b)
		}
	}
}

func TreeServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "[")
	i := 0
	filepath.Walk(Root, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, ".git") {
			return nil
		}
		if i != 0 {
			io.WriteString(w, ",")
		}
		i++
		j, _ := json.Marshal(map[string]interface{}{
			"path":     strings.Replace(path, Root, "", 1),
			"name":     f.Name(),
			"size":     f.Size(),
			"mode":     f.Mode(),
			"mod_time": f.ModTime(),
			"is_dir":   f.IsDir(),
			"err":      err,
		})
		w.Write(j)
		return nil
	})

	io.WriteString(w, "]")
}

func GitStatusServer(w http.ResponseWriter, req *http.Request) {
	result := make(map[string]string)
	cmd := exec.Command("git", "status", "--porcelain", "-u")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		logrus.Warnf("Git Status Err %s", err)
		return
	}

	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		l := scanner.Text()
		status := l[0:2]
		path := l[3:]
		result[path] = status
	}

	cmd.Wait()
	b, err := json.Marshal(result)
	if err != nil {
		logrus.Warnf("JSON Marshal error %s", err)
	}
	w.Write(b)
}

func main() {
	logrus.Info("Starting Server")
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	Root = string(out[0:len(out)-1]) + "/"
	if err != nil {
		logrus.Fatalf("Cannot Get Root Directory %s", err)
	}
	originalpath, err := filepath.Abs(".")
	if err != nil {
		logrus.Fatalf("Cannot Get CWD %s", err)
	}
	defer os.Chdir(originalpath)
	os.Chdir(Root)

	logrus.Infof("Server Started at `%s`", Root)

	http.Handle("/_/", http.StripPrefix("/_/", http.FileServer(assetFS())))

	//http.Handle("/_/", http.StripPrefix("/_/", http.FileServer(http.Dir("./web"))))
	http.Handle("/", http.HandlerFunc(HelloServer))
	http.Handle("/$", http.HandlerFunc(TreeServer))
	http.Handle("/@", http.HandlerFunc(GitStatusServer))
	err = http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
