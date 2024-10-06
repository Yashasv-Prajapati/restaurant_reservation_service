package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	// "main/event"
	"main/table" // Importing the package "table" from the subfolder
	"main/utils"
	"strconv"
	"strings"
	"time"
)

var dateTimeFormat string
var timeFormat string

func init(){
	table.InitializeTables(100,5)
	dateTimeFormat = "2006-01-02::15:04"
	timeFormat = "15:04"
	// table.InitializeHalls(100)
}

// ClearConsole clears the console screen
func ClearConsole() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") // For Windows
	default:
		cmd = exec.Command("clear") // For Unix-like systems (Linux, macOS)
	}
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		utils.PrintConsole("Error clearing console:", err)
	}
}


func help(key []string){
	if len(key) == 0 {
		utils.PrintConsole("Please use the following format:")
		utils.PrintConsole("help table|event")
		return
	}

	switch key[0] {
	case "table":
		utils.PrintConsole("table {customer_name} {number_of_people} {mobile_number} {reservation_start_time} {reservation_end_time}")
	case "event":
		utils.PrintConsole("event {customer_name} {number_of_people} {mobile_number} {reservation_start_time} {reservation_end_time}")
	default:
		utils.PrintConsole("Please use the following format:")
		utils.PrintConsole("help table|event")
	}
}


func main() {

	// Set opening and closing times
	openingTime, err := time.Parse(timeFormat, "09:00")
	if err != nil {
		utils.PrintConsole("Error parsing opening time:", err)
		return
	}

	closingTime, err := time.Parse(timeFormat, "22:00") // 10:00 PM in 24-hour format
	if err != nil {
		utils.PrintConsole("Error parsing closing time:", err)
		return
	}

	utils.PrintConsole("Enter the reservation type and arguments in the following format:")
	utils.PrintConsole("table {customer_name} {number_of_people} {mobile_number} {reservation_start_time} {reservation_end_time}")


	for {

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("-> ")

		// Read an entire line of input
		input_args, err := reader.ReadString('\n')
		if err != nil {
			utils.PrintConsole("Error reading input:", err)
			return
		}

		// Trim the newline character from the input
		input_args = strings.TrimSpace(input_args)

		// Split the input by spaces
		args := strings.Split(input_args, " ")

		reservation_type := args[0]

		switch reservation_type {
			// table {number_of_people} {mobile_number} {reservation_end_time} {reservation_end_time}}
			case "table":
				// Parse the arguments and create a reservation
				r, err := parseTableBookingArguments(args[1:])
				if err != nil {
					utils.PrintConsole("Something went wrong while parsing:", err)
					return
				}
				r.ReservationDatetime = time.Now()
				handleTableReservation(&r, openingTime, closingTime)
			// case "event":
			// 	e, err := parseEventBookingArguments(args)
			// 	if err != nil {
			// 		utils.PrintConsole("Error in parsing event booking arguments")
			// 		return
			// 	}
			// 	break
			case "help":
				help(args[1:])
			case "clear":
				ClearConsole()
			case "exit":
				ClearConsole()
				os.Exit(0)
			default:
				utils.PrintConsole("Invalid reservation type")
		}

	}
}

func handleTableReservation(t* table.TableReservation, openingTime time.Time, closingTime time.Time) {

	table_ids, err := table.ReserveTable(t, openingTime, closingTime)
	if err != nil {
		utils.PrintConsole("Error reserving table:", err)
		return
	}
	utils.PrintConsole("Table reserved successfully. Table IDs:", table_ids)
}

// func parseEventBookingArguments(args[] string) (event.EventReservation, error){

// 	// Create a new reservation
// 	var r event.EventReservation

// 	event_id, err := strconv.ParseInt(args[1], 10, 64)
// 	if err != nil {
// 		utils.PrintConsole("Error in parsing event id")
// 		return r, err
// 	}

// 	number_of_people, err := strconv.ParseInt(args[2], 10, 64)
// 	if err != nil {
// 		utils.PrintConsole("Error in parsing number of people")
// 		return r, err
// 	}

// 	mobile_number := args[3]

// 	// Time format to parse the reservation datetime and start time
// 	dateTimeFormat := "2006-01-02 15:04"

// 	reservation_datetime, err := time.Parse(dateTimeFormat, args[4])
// 	if err != nil {
// 		utils.PrintConsole("Error in parsing reservation datetime")
// 		return r, err
// 	}

// 	start_time, err := time.Parse(dateTimeFormat, args[5])
// 	if err != nil {
// 		utils.PrintConsole("Error in parsing start time")
// 		return r, err
// 	}

// 	// Parse reservation duration
// 	reservation_duration_hours := args[6] + "h" // Assuming duration in hours
// 	reservation_duration, err := time.ParseDuration(reservation_duration_hours)
// 	if err != nil {
// 		utils.PrintConsole("Error in parsing reservation duration")
// 		return r, err
// 	}

// 	// hall_id, err := strconv.ParseInt(args[7], 10, 64)
// 	// if err != nil {
// 	// 	utils.PrintConsole("Error in parsing hall id")
// 	// 	return r, err
// 	// }

// 	r.CustomerID = "1"
// 	r.EventID = event_id
// 	r.NumberOfPeople = int(number_of_people)
// 	r.MobileNumber = mobile_number
// 	r.ReservationDatetime = reservation_datetime
// 	r.StartTime = start_time
// 	r.ReservationDuration = reservation_duration
// 	// r.HallID = hall_id

// 	return r, nil

// }

func parseTableBookingArguments(args[] string) (table.TableReservation, error){

	// Create a new reservation
	var r table.TableReservation

	if len(args) < 6 {
		return r, errors.New("invalid number of arguments")
	}

	customer_id := args[0]

	number_of_people, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		utils.PrintConsole(err)
		return r, errors.New("error in parsing number of people")
	}

	mobile_number := args[2]


	reservation_start_time, err := time.Parse(dateTimeFormat, args[3])
	if err != nil {
		return r, errors.New("error in parsing reservation start time")
	}

	// Parse reservation end time
	reservation_end_time, err := time.Parse(dateTimeFormat, args[4])
	if err != nil {
		return r, errors.New("error in parsing reservation end time")
	}


	r.CustomerID = customer_id
	// r.TableID = table_id
	r.NumberOfPeople = int(number_of_people)
	r.MobileNumber = mobile_number
	r.ReservationDatetime = time.Now()
	r.StartTime = reservation_start_time
	r.EndTime = reservation_end_time

	utils.PrintConsole("TableReservation parsed:")

	return r, nil
}