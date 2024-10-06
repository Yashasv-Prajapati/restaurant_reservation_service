package event

import "time"

type EventReservation struct {
	EventID             int64         // Event ID
	CustomerID          string        // Customer ID
	NumberOfPeople      int           // Number of people
	MobileNumber        string        // Customer's mobile number
	ReservationDatetime time.Time     // Date and time of the reservation
	StartTime           time.Time     // Start time of the reservation
	ReservationDuration time.Duration // Duration of the reservation
	HallID              int64         // Hall ID
}

type Hall struct {
	hall_id int
	capacity int
	available bool
	reservedAt time.Time
};

var halls []Hall;

func InitializeHalls(size int){
	halls = make([]Hall, size)
}