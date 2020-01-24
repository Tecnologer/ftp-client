package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/cheggaaa/pb/v3"
	"github.com/jlaffaye/ftp"
)

type fileError struct {
	err   error
	file  string
	index int
}

func downloadContent(path string) error {
	content, err := c.List(path)

	if err != nil {
		return err
	}

	for _, element := range content {
		elementPath := fmt.Sprintf("%s/%s", path, element.Name)

		if strings.HasSuffix(elementPath, "..") || strings.HasSuffix(elementPath, ".") {
			continue
		}

		if element.Type == ftp.EntryTypeFolder {
			// logrus.Info("new folder found ", elementPath)
			if err = downloadContent(elementPath); err != nil {
				return err
			}
		}

		if element.Type != ftp.EntryTypeFile {
			continue
		}

		filesToDownload.Add(elementPath)

		// if err = writeFile(elementPath); err != nil {
		// 	return err
		// }
	}
	return nil
}

func downloadMarkedFiles() error {
	count := filesToDownload.Len()

	if count == 0 {
		return fmt.Errorf("no files to download")
	}

	fmt.Printf("\n >>>> found %d files <<<<\n\n", count)

	bar := getProgressBar(count)
	defer bar.Finish()

	filesChannel := []chan string{}
	results := make(chan *fileError)

	for w := 1; w <= 3; w++ {
		files := make(chan string)
		filesChannel = append(filesChannel, files)
		go worker(w, files, results, w-1, bar)
	}

	//initial workers
	for _, ch := range filesChannel {
		if !filesToDownload.HasFiles() {
			break
		}

		ch <- filesToDownload.GetNext()
	}

	for filesToDownload.HasFiles() {
		select {
		case result := <-results:
			if result.err != nil {
				logrus.Debugf("error donwloading file %s. Error: %v\n", result.file, result.err)
			}
			filesChannel[result.index] <- filesToDownload.GetNext()
		}
	}

	//close channels
	for _, ch := range filesChannel {
		close(ch)
	}

	return nil
}

func worker(id int, files <-chan string, results chan<- *fileError, index int, bar *pb.ProgressBar) {
	for file := range files {
		err := writeFile(file)
		bar.Increment()
		results <- &fileError{
			err:   err,
			file:  file,
			index: index,
		}
	}
}

func writeFile(filename string) error {
	// logrus.Info("downloading file ", filename)
	r, err := c.Retr(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	buf, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir("."+filename), 0700)
	if err != nil {
		return err
	}
	totalBytes += int64(len(buf))
	err = ioutil.WriteFile("."+filename, buf, 0644)

	if err != nil {
		return err
	}

	// logrus.Info("file ", filename, " download sucessfully")
	fileCount++
	return nil
}

func getProgressBar(count int) *pb.ProgressBar {
	// create bar
	bar := pb.New(count)

	// refresh info every second (default 200ms)
	// bar.SetRefreshRate(time.Second)

	// force set io.Writer, by default it's os.Stderr
	bar.SetWriter(os.Stdout)

	// bar will format numbers as bytes (B, KiB, MiB, etc)
	// bar.Set(pb.Byte, true)

	// bar use SI bytes prefix names (B, kB) instead of IEC (B, KiB)
	bar.Set(pb.SIBytesPrefix, true)

	// start bar
	bar.Start()
	return bar
}
