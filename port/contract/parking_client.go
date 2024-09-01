package contract

import "booking/port/dto"

type ParkingClientInterface interface {
	Login() error
	Book(payload *dto.BookParkingPayload) error
	GetAccessToken() string
}
