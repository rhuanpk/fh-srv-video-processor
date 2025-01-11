package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	frameName        = "frame_*.png"
	frameHighQuality = 1

	videoPath     = "video/sample.mp4"
	frameInterval = 5
)

func formatDuration(seconds int) string {
	duration := time.Duration(seconds) * time.Second
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	secs := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

func extractFrames(videoPath, outputDir string, frameInterval int, highQuality bool) error {
	// errChan := make(chan error, 1)

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	kwArgs := ffmpeg.KwArgs{"vf": fmt.Sprintf("fps=1/%d", frameInterval)}
	if highQuality {
		kwArgs["q:v"] = frameHighQuality
	}

	// videoProbe, err := ffmpeg.Probe(videoPath)
	// if err != nil {
	// 	return err
	// }

	// // // videoSize, err := strconv.ParseInt(gjson.Get(videoProbe, "format.size").Str, 0, 64)
	// // videoSize, err := strconv.ParseFloat(gjson.Get(videoProbe, "format.size").Str, 64)
	// // if err != nil {
	// // 	return err
	// // }
	// // videoSizeMB := videoSize / (1 << 20)

	// videoDuration, err := strconv.ParseFloat(gjson.Get(videoProbe, "format.duration").Str, 64)
	// if err != nil {
	// 	return err
	// }

	// frameIntervalStep := frameInterval * 3
	// if frameIntervalStep >= 10 {
	// 	frameIntervalStep = frameIntervalStep / 2
	// }
	// goRoutinesCount := int(videoDuration) / frameIntervalStep

	// var wg sync.WaitGroup
	// for index := 0; index < goRoutinesCount; index++ {
	// 	wg.Add(1)
	// 	go func(index int, kwArgs ffmpeg.KwArgs) {
	// 		defer wg.Done()
	// 		startInterval := frameIntervalStep * index
	// 		kwArgs["ss"] = formatDuration(startInterval)
	// 		kwArgs["to"] = formatDuration(frameIntervalStep * (index + 1))
	// 		kwArgs["start_number"] = startInterval/frameInterval + 1
	// 		err := ffmpeg.Input(videoPath).
	// 			Output(filepath.Join(outputDir, strings.Replace(frameName, "*", "%04d", 1)), kwArgs).
	// 			// OverWriteOutput().RunWithResource(0.75, 1.0)
	// 			OverWriteOutput().Run()
	// 		if err != nil {
	// 			errChan <- err
	// 		}
	// 	}(index, kwArgs)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(errChan)
	// }()
	// if err, ok := <-errChan; ok {
	// 	return err
	// }

	err := ffmpeg.Input(videoPath).
		Output(filepath.Join(outputDir, strings.Replace(frameName, "*", "%04d", 1)), kwArgs).
		// OverWriteOutput().RunWithResource(0.75, 1.0)
		OverWriteOutput().Run()
	if err != nil {
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

func createZip(zipPath, sourceDir string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	frames, err := filepath.Glob(filepath.Join(sourceDir, frameName))
	if err != nil {
		return err
	}

	for _, frame := range frames {
		if err := addFileToZip(zipWriter, frame); err != nil {
			return err
		}
	}

	return nil
}

func processVideo(videoPath string, frameInterval int, highQuality bool) (string, error) {
	videoDir := filepath.Dir(videoPath)

	framesDir := filepath.Join(videoDir, "frames")
	if err := os.MkdirAll(framesDir, os.ModePerm); err != nil {
		return "", err
	}

	defer func() {
		if err := os.RemoveAll(framesDir); err != nil {
			log.Println("error in remove frames dir:", err)
		}
	}()

	videoName := filepath.Base(videoPath)
	zipName := strings.TrimSuffix(videoName, filepath.Ext(videoName)) + "_frames.zip"
	zipPath := filepath.Join(videoDir, zipName)

	if err := extractFrames(videoPath, framesDir, frameInterval, highQuality); err != nil {
		return "", err
	}

	if err := createZip(zipPath, framesDir); err != nil {
		return "", err
	}

	return zipPath, nil
}

func main() {
	defer func(start time.Time) {
		log.Println("execution time:", time.Since(start).String())
	}(time.Now())

	zipPath, err := processVideo(videoPath, frameInterval, true)
	if err != nil {
		log.Println("error in process video:", err)
		return
	}

	log.Println("success in process video:", zipPath)
}
