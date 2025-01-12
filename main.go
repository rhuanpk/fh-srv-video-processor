package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	frameName        = "frame_*.png"
	frameHighQuality = 1

	videoPath     = "video/sample.mp4"
	frameInterval = 5
)

var videosPaths = []string{"video/sample_1.mp4", "video/sample_2.mp4"}

func formatDuration(seconds int) string {
	duration := time.Duration(seconds) * time.Second
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	secs := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

func extractFrames(videosPaths, outputsDirs []string, frameInterval int, highQuality bool) error {
	errChan := make(chan error, 1)

	for _, outputDir := range outputsDirs {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return err
		}
	}

	kwArgs := ffmpeg.KwArgs{"vf": fmt.Sprintf("fps=1/%d", frameInterval)}
	if highQuality {
		kwArgs["q:v"] = frameHighQuality
	}

	var wg sync.WaitGroup
	for index, videoPath := range videosPaths {
		wg.Add(1)
		go func(index int, videoPath string) {
			defer wg.Done()
			err := ffmpeg.Input(videoPath).
				Output(filepath.Join(
					outputsDirs[index],
					strings.Replace(frameName, "*", "%04d", 1)),
					kwArgs,
				).
				// OverWriteOutput().RunWithResource(0.75, 1.0)
				OverWriteOutput().Run()
			if err != nil {
				errChan <- err
			}
		}(index, videoPath)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	if err, ok := <-errChan; ok {
		return err
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}

func createZip(zipsPaths, sourcesDirs []string) error {
	for index, zipPath := range zipsPaths {
		zipFile, err := os.Create(zipPath)
		if err != nil {
			return err
		}
		defer zipFile.Close()

		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		frames, err := filepath.Glob(filepath.Join(sourcesDirs[index], frameName))
		if err != nil {
			return err
		}

		for _, frame := range frames {
			if err := addFileToZip(zipWriter, frame); err != nil {
				return err
			}
		}
	}

	return nil
}

func processVideo(videosPaths []string, frameInterval int, highQuality bool) ([]string, error) {
	var videosDirs, framesDirs, videosNames, videosNamesWithoutExt, zipsPaths []string

	for _, videoPath := range videosPaths {
		videosDirs = append(videosDirs, filepath.Dir(videoPath))
	}

	for index, videoPath := range videosPaths {
		videoName := filepath.Base(videoPath)
		videosNames = append(videosNames, videoName)

		videoNameWithoutExt := strings.TrimSuffix(videoName, filepath.Ext(videoName))
		videosNamesWithoutExt = append(videosNamesWithoutExt, videoNameWithoutExt)

		unixTime := strconv.Itoa(int(time.Now().UnixNano()))
		zipName := videoNameWithoutExt + "_frames_" + unixTime + ".zip"
		zipsPaths = append(zipsPaths, filepath.Join(videosDirs[index], zipName))
	}

	for index, videoDir := range videosDirs {
		framesDir, err := os.MkdirTemp(videoDir, videosNamesWithoutExt[index]+"_frames_")
		if err != nil {
			return nil, err
		}

		if err := os.MkdirAll(framesDir, os.ModePerm); err != nil {
			return nil, err
		}
		framesDirs = append(framesDirs, framesDir)
	}

	defer func() {
		for _, framesDir := range framesDirs {
			if err := os.RemoveAll(framesDir); err != nil {
				log.Println("error in remove frames dir:", err)
			}
		}
	}()

	if err := extractFrames(videosPaths, framesDirs, frameInterval, highQuality); err != nil {
		return nil, err
	}

	if err := createZip(zipsPaths, framesDirs); err != nil {
		return nil, err
	}

	return zipsPaths, nil
}

func main() {
	defer func(start time.Time) {
		log.Println("execution time:", time.Since(start).String())
	}(time.Now())

	// zipsPaths, err := processVideo([]string{videoPath}, frameInterval, true)
	zipsPaths, err := processVideo(videosPaths, frameInterval, true)
	if err != nil {
		log.Println("error in process videos:", err)
		return
	}

	log.Println("success in process videos:", zipsPaths)
}
