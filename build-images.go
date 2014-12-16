package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const dockerfileTmpl = `FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y curl git build-essential

ENV PECO_GOVERSION  {{ .GoVersion }}
ENV PATH /opt/local/go/bin:$PATH

RUN mkdir /go-src
COPY go/* /go-src/

ENV PECO_GOFILENAME go-src/go{{ .GoVersion }}.linux-amd64.tar.gz
RUN tar xvzf $PECO_GOFILENAME

RUN rm -rf /go-src

RUN mkdir -p /opt/local
RUN mv go /opt/local/

RNN rm /go*

ENV PATH /work/bin:/opt/local/go/bin:$PATH
ENV GOROOT /opt/local/go
ENV GOPATH /work

RUN go get -u github.com/laher/goxc
RUN goxc -t

RUN apt-get install -y unzip
RUN curl -LO https://github.com/tcnksm/ghr/releases/download/v0.2.0/ghr_v0.2.0_linux_amd64.zip
RUN unzip ghr_v0.2.0_linux_amd64.zip
RUN mv ghr /usr/local/bin
RUN rm -rf ghr_v0.2.0_linux_amd64.zip

COPY script/test-docker.sh /test-docker.sh
COPY script/build-docker.sh /build-docker.sh
COPY script/release-docker.sh /release-docker.sh

VOLUME ["/work/src/github.com/peco/peco"]

CMD echo "peco-docker built for go version {{ .GoVersion }}"
`

func main() {
	noCache := flag.Bool("rebuild", false, "pass --no-cache to docker build")
	flag.Parse()

	supportedGoVersions := []string{
		"1.4",
	}

	helperScripts := []string{
		"script/test-docker.sh",
		"script/build-docker.sh",
		"script/release-docker.sh",
	}

	curDir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current directory: %s", err)
		return
	}

	// First, copy download all files, so that it's easier to cache
	// docker intermediate images
	for _, ver := range supportedGoVersions {
		fn := "go" + ver + ".linux-amd64.tar.gz"
		localFn := curDir + "/go/" + fn
		helperScripts = append(helperScripts, localFn)

		if _, err := os.Stat(localFn); err == nil {
			log.Printf("File '%s' already exists...", localFn)
			continue
		}
		url := "https://storage.googleapis.com/golang/" + fn
		log.Printf("Fetching %s...", url)
		res, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to fetch file: %s", err)
			return
		}
		if res.StatusCode != 200 {
			log.Printf("Failed to fetch file: got code %d", res.StatusCode)
			return
		}
		f, err := os.OpenFile(localFn, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Failed to open file %s for wriging: %s", localFn, err)
			return
		}
		io.Copy(f, res.Body)
		f.Close()
	}

	t := template.Must(template.New("docker").Parse(dockerfileTmpl))
	v := struct {
		GoVersion string
	}{}
	for i, ver := range supportedGoVersions {
		dir, err := ioutil.TempDir("", "peco-docker-go-")
		if err != nil {
			log.Printf("Failed to create temporary directory: %s", err)
			return
		}
		defer os.RemoveAll(dir)
		log.Printf("Work directory is '%s'...", dir)

		// Copy some helper scripts..
		for _, fn := range helperScripts {
			var relPath string
			if ! filepath.IsAbs(fn) {
				relPath = fn
			} else {
				relPath, err = filepath.Rel(curDir, fn)
				if err != nil {
					log.Printf("Failed to get relative path of %s (against %s): %s", fn, curDir, err)
					return
				}
			}
			newFn := filepath.Join(dir, relPath)
			log.Printf("Copying %s to %s", fn, newFn)

			parentDir := filepath.Dir(newFn)
			if _, err := os.Stat(parentDir); err != nil {
				log.Printf("Creating directory %s", parentDir)
				if err := os.MkdirAll(parentDir, 0755); err != nil {
					log.Printf("Failed to create directory %s: %s", parentDir, err)
					return
				}
			}

			if err := os.Link(fn, newFn); err != nil {
				log.Printf("Error copying %s: %s", fn, err)
				return
			}
		}

		fn := filepath.Join(dir, "Dockerfile")
		f, err := os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Failed to create Dockerfile: %s", err)
			return
		}

		v.GoVersion = ver
		if err := t.Execute(f, v); err != nil {
			log.Printf("Failed to create Dockerfile: %s", err)
			f.Close()
			return
		}

		f.Close()

		tag := "peco-docker:go" + ver

		var cmd *exec.Cmd
		if i == 0 && *noCache {
			cmd = exec.Command("docker", "build", "--no-cache", "-t", tag, ".")
		} else {
			cmd = exec.Command("docker", "build", "-t", tag, ".")
		}
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to build image for go version %s: %s", ver, err)
			return
		}

		log.Printf("Built image for go %s\n", ver)
	}
}