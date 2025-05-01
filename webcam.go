package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/liyue201/goqr"
	"gocv.io/x/gocv"
)

func ReadFromWebcam() ([]*goqr.QRData, error) {
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		return nil, fmt.Errorf("failed to open camera: %v", err)
	}
	defer webcam.Close()

	webcam.Set(gocv.VideoCaptureFPS, 30)
	window := gocv.NewWindow("Original")
	img := gocv.NewMat()
	defer img.Close()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		img_proper, err := img.ToImage()
		if err == nil {
			result, err := goqr.Recognize(img_proper)
			if err == nil {
				return result, nil
			}
		}
		window.WaitKey(1)
	}

}

func SingleReadFromWebcam(webcam *gocv.VideoCapture, img *gocv.Mat) ([]*goqr.QRData, error) {
	var qrCodes []*goqr.QRData
	webcam.Read(img)
	img_proper, err := img.ToImage()
	if err == nil {
		qrCodes, err = goqr.Recognize(img_proper)

	}
	return qrCodes, err

}

func DoReadFromWebcam(wg *sync.WaitGroup, qr_text chan string, webcam_done chan bool) {
	log.Printf("DoReadFromWebcam: \n")
	defer wg.Done()

	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Printf("failed to open camera: %v", err)
		return
	}
	defer webcam.Close()

	webcam.Set(gocv.VideoCaptureFPS, 30)

	img := gocv.NewMat()
	defer img.Close()
	defer close(qr_text)

	for {
		select {
		case <-webcam_done:
			return
		default:
			qrCodes, err := SingleReadFromWebcam(webcam, &img)
			if err == nil {
				for _, qrCode := range qrCodes {
					qr_text <- string(qrCode.Payload)
				}
				return
			}
		}
	}
}
