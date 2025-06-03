package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

var (
	print_queue,
	status_queue chan string
)

func pdf_print(pdf_path string) error {

	app := "./PDFtoPrinter"
	cmd := exec.Command(app, pdf_path)
	err = cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()

	return err

}

func do_print_queue(print_queue chan string) {
	for {
		select {
		case new_file, ok := <-print_queue:
			if ok {
				log.Println("Info: Printing: ", new_file)
				err := pdf_print(new_file)
				if err != nil {
					log.Println(err)
				}
			} else {
				return
			}
		}
	}
}

func show_status(message string) {
	message = fmt.Sprintf("%s\t\t%s", time.Now().Format("15:04:05.000"), message)
	status_queue <- message
}

func status_bar_show(message string, timer *time.Timer) {
	status_bar.SetText(message)
	select {
	case <-timer.C:
		status_bar.SetText("")
	}
}

func do_status_queue(status_queue chan string) {
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
