package table

import (
	"errors"
	"main/utils"
	"sort"
	"time"
)

// TableReservation struct for a table reservation
type TableReservation struct {
	CustomerID          string        // Customer ID
	TableID             int64         // Table ID
	NumberOfPeople      int           // Number of people
	MobileNumber        string        // Customer's mobile number
	ReservationDatetime time.Time     // Date and time of the reservation
	StartTime           time.Time     // Start time of the reservation
	EndTime 			time.Time 	  // End time of the reservation
}

type Table struct {
	table_id int
	capacity int
	reservedAt time.Time
	reservedFrom time.Time
	reservedTill time.Time
};

var tables []Table;
var reservations map[int] []Table

func InitializeTables(size int, capacity int){
	tables = make([]Table, size)
	for i := 0; i < size; i++ {
		tables[i] = Table{i, capacity, time.Now(), time.Now(), time.Now()}
	}
	reservations = make(map[int] []Table)
}

func ReserveTable(reservation * TableReservation, restaurantOpeningTime time.Time, restaurantClosingTime time.Time) ([] int, error){

	// This function will reserve a table for a customer
	// The reservation will be stored in the appropriate data structure

	// we will use hashmap for this purpose

	// hashmap key - table id
	// hashmap value - list of booked Tables entries for that table ID

	// error handling
	// 1. The start time and end time should be in the future
	// 2. The end time should be after the start time
	// 2. The start time and end time must not exceed the restaurant's opening hours
	// number of people must be > 0

	if( reservation.EndTime.Sub(reservation.StartTime) <= 0 ){ // end time should be after the start time
		return []int{}, errors.New("end time should be after the start time")
	}else if( reservation.StartTime.Before ( reservation.ReservationDatetime ) ){ // reservation start time should be after the reservation datetime
		return []int{}, errors.New("reservation start time should be after the current time of reservation")
	}else if( reservation.StartTime.Before(restaurantOpeningTime) || reservation.EndTime.After(restaurantClosingTime) ){ // reservation time should be within the restaurant's opening hours

		checkValidity := utils.IsWithinOpenHours(reservation.StartTime, restaurantOpeningTime, restaurantClosingTime)
		if !checkValidity {
			return []int{}, errors.New("reservation time should be within the restaurant's opening hours")
		}

	}else if(reservation.NumberOfPeople <= 0){ // number of people must be > 0
		return []int{}, errors.New("number of people must be atleast 1")
	}


	booked_tables := findAvailableTableWithSizeK(reservation.NumberOfPeople, reservation.ReservationDatetime, reservation.StartTime, reservation.EndTime)

	if len(booked_tables) == 0 {
		return booked_tables, errors.New("no table available for the given number of people")
	}

	return booked_tables, nil
}

func findAvailableTableWithSizeK(k int, currTime time.Time, reserveFrom time.Time, reserveTill time.Time ) ([]int){

	totalCapacityAchieved := 0

	// sort each table in increasing order of capacity
	// sort in increasing order of capacity, so that we can find the first table with capacity >= k
	sort.SliceStable(tables, func (i, j int) bool {
		return tables[i].capacity < tables[j].capacity
	})

	booked_tables := make([]int, 0)

	// For each table
		// if reservation of this table is possible, then possibly reserve it, store the table_id somewhere, increase the totalCapacityAchieved, and store the table_id in the result
			// if totalCapacityAchieved >= k, then return the result
		// if reservation of this table is not possible, then continue to the next table
	// at the end, if toatlCapacityAchieved >= k,only then start adding all the table_ids in the hashmap
	totalTables := len(tables)

	for i := 0; i < totalTables; i++ {

		table_id := tables[i].table_id

		if totalCapacityAchieved >= k {
			break;
		}

		reserved_tables_for_table_id, ok := reservations[table_id]

		if !ok {
			// table is not reserved for anything
			totalCapacityAchieved += tables[i].capacity
			booked_tables = append(booked_tables, table_id)
			continue
		}

		for j := 0; j < len(reserved_tables_for_table_id); j++ {

			reserved_from := reserved_tables_for_table_id[j].reservedFrom
			reserved_till := reserved_tables_for_table_id[j].reservedTill

			want_to_reserve_from := reserveFrom
			want_to_reserve_till := reserveTill

			// 3 main conditions to check
			if( want_to_reserve_till.After(reserved_from) && (want_to_reserve_till.Before(reserved_till) || want_to_reserve_till.Equal(reserved_till) ) ){ // want to reserve table from 10 to 12, but table is reserved from 11 to 1
				// table is not available
				continue
			}else if( ( want_to_reserve_from.After(reserved_from) || want_to_reserve_from.Equal(reserved_from) ) && ( want_to_reserve_till.Before(reserved_till) || want_to_reserve_till.Equal(reserved_till) ) ){ // want to reserve table from 11 to 12, but table is reserved from 10 to 1
				// table is not available
				continue
			}else if( (want_to_reserve_from.After(reserved_from) || want_to_reserve_from.Equal(reserved_from)) && want_to_reserve_from.Before(reserved_till) ){ // want to reserve table from 11 to 12, but table is reserved from 10 to 11
				// table is not available
				continue
			}

			// we can reserve this table

			// increase the totalCapacityAchieved
			totalCapacityAchieved += tables[i].capacity

			// store the table_id in the result
			booked_tables = append(booked_tables, table_id)

			break;
		}

	}

	// if totalCapacityAchieved < k, then return empty list and don't add any reservation to the hashmap
	if totalCapacityAchieved < k {
		return []int{}
	}

	for i := 0; i < len(booked_tables); i++ {
		table_id := booked_tables[i]
		newTable := Table{table_id, tables[table_id].capacity, currTime, reserveFrom, reserveTill}
		reservations[table_id] = append(reservations[table_id], newTable)
	}

	return booked_tables

}