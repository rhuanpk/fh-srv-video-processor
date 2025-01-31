package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"extractor/internal/config"
	"extractor/internal/resouce/processor"
	"extractor/internal/resouce/zipper"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type service struct {
	zip zipper.Service
}

func NewService(zip zipper.Service) processor.Service {
	return &service{zip: zip}
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
		kwArgs["q:v"] = config.FrameHighQuality
	}

	var wg sync.WaitGroup
	for index, videoPath := range videosPaths {
		wg.Add(1)
		go func(index int, videoPath string) {
			defer wg.Done()
			err := ffmpeg.Input(videoPath).
				Output(filepath.Join(
					outputsDirs[index],
					strings.Replace(config.FrameName, "*", "%04d", 1)),
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

func (s service) Process(videosPaths []string, frameInterval int, highQuality bool) ([]string, error) {
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

	if err := s.zip.Create(zipsPaths, framesDirs); err != nil {
		return nil, err
	}

	return zipsPaths, nil
}
