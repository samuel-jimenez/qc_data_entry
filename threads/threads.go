package threads

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/samuel-jimenez/windigo"
)

var (
	PRINT_QUEUE,
	STATUS_QUEUE chan string

	Status_bar *windigo.StatusBar
)

func pdf_print(pdf_path string) error {

	app := "./PDFtoPrinter"
	cmd := exec.Command(app, pdf_path)
	err := cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()

	return err

}

func Do_print_queue(print_queue chan string) {
	for {
		select {
		case new_file, ok := <-print_queue:
			if ok {
				log.Println("Info: Printing: ", new_file)
				err := pdf_print(new_file)
				if err != nil {
					log.Printf("Error: %q: %s\n", err, "Do_print_queue")
				}
			} else {
				return
			}
		}
	}
}

func Show_status(message string) {
	if STATUS_QUEUE != nil {
		message = fmt.Sprintf("%s\t\t%s", time.Now().Format("15:04:05.000"), message)
		STATUS_QUEUE <- message
	} else {
		log.Println("Warn: Status queue not configured. Call threads.Do_status_queue to set up.")
	}
}

func status_bar_show(message string, timer *time.Timer) {
	if Status_bar != nil {
		Status_bar.SetText(message)
		select {
		case <-timer.C:
			Status_bar.SetText("")
		}
	} else {
		log.Println("Warn: Status bar not configured. Call windigo.NewStatusBar to set up.")
		log.Println("Info: Status bar message:", message)
	}
}

func Do_status_queue(status_queue chan string) {
	var display_timeout_timer *time.Timer
	display_timeout := 2 * time.Second

	for {
		select {
		case message, ok := <-status_queue:
			if ok {
				display_timeout_timer = time.NewTimer(display_timeout)
				status_bar_show(message, display_timeout_timer)
			} else {
				return
			}
		}
	}
}
