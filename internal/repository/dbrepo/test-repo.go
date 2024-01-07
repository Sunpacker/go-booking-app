package dbrepo

import (
	"errors"
	"github.com/Sunpacker/go-booking-app/internal/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(dto models.Reservation) (int, error) {
	_ = dto
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	_ = r
	return nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	_ = start
	_ = end
	_ = roomID
	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	_ = start
	_ = end
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	_ = id
	var user models.User
	return user, nil
}
func (m *testDBRepo) UpdateUser(user models.User) error {
	_ = user
	return nil
}
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	_ = email
	_ = testPassword
	return 0, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	_ = id
	var reservation models.Reservation
	return reservation, nil
}
func (m *testDBRepo) UpdateReservation(reservation models.Reservation) error {
	_ = reservation
	return nil
}
func (m *testDBRepo) DeleteReservation(id int) error {
	_ = id
	return nil
}
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	_ = id
	_ = processed
	return nil
}
