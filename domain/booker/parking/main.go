package parking

import (
	"booking/domain/mouse"
	"booking/port/config"
	"booking/port/contract"
	"booking/port/dto"
	"context"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"sync"
	"time"
)

type Parking struct {
	parkingClient contract.ParkingClientInterface
	ctx           context.Context
	config        *config.Parking
	outputLabel   *widget.Label
	updateOutput  contract.UpdateOutputFunc
}

func NewParkingBooker(cfg *config.Config, parkingClient contract.ParkingClientInterface, ctx context.Context, widget *widget.Label, updateOutput contract.UpdateOutputFunc) *Parking {
	return &Parking{
		config:        &cfg.Parking,
		ctx:           ctx,
		outputLabel:   widget,
		updateOutput:  updateOutput,
		parkingClient: parkingClient,
	}
}

func (p *Parking) Start() {
	var wg sync.WaitGroup
	if p.config.StandBy.CanStandby {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mouse.MoveWithContext(p.ctx, p.config.StandBy.StartTime, p.config.StandBy.EndTime, p.outputLabel, p.updateOutput)
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.book()
	}()
	wg.Wait()
}

func (p *Parking) book() {
	start, err := time.Parse("02-01-2006 15:04:05", p.config.Schedule.Date)
	if err != nil {
		p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker][Error] failed to parse start time: %v", err))
		return
	}
	location := time.Now().Location()
	start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), location)
	now := time.Now()
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	daysSinceEpoch := int(start.Sub(unixEpoch).Hours() / 24)
	daysSinceEpoch += p.config.Schedule.Offset
	p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker] Current time: %v", now))
	p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker] Start time: %v (Zone: %v)", start, start.Location()))
	p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker] Days since epoch: %v", daysSinceEpoch))
	p.updateOutput(p.outputLabel, "[Booker] Logging user to booking App.")
	if now.Before(start) {
		p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker] Waiting for start time. Starting in %v seconds", time.Until(start).Seconds()))
		select {
		case <-p.ctx.Done():
			p.updateOutput(p.outputLabel, "[Booker] Process stopped.")
			return
		case <-time.After(time.Until(start)):
			// Continue execution
		}
	} else {
		p.updateOutput(p.outputLabel, "[Booker] Start time is in the past. Exiting")
		return
	}
	err = p.parkingClient.Login()
	if err != nil {
		p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker][Error] failed to login: %v", err))
		return
	}
	for time.Now().Before(time.Now().Add(5 * time.Minute)) {
		select {
		case <-p.ctx.Done():
			p.updateOutput(p.outputLabel, "[Booker] Process stopped.")
			return
		default:
			p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker][Error] Attempting to book parking space."))
			for _, parkingSpot := range p.config.ParkingSpot {
				select {
				case <-p.ctx.Done():
					p.updateOutput(p.outputLabel, "[Booker] Process stopped.")
					return
				default:
					if p.bookFirstAvailable(daysSinceEpoch, parkingSpot) {
						return
					}
				}
			}
			p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker]Retrying in : %v seconds", 3))
			time.Sleep(3 * time.Second)
		}
	}
}

func (p *Parking) bookFirstAvailable(daysSinceEpoch int, parkingSpot config.ParkingSpot) bool {
	p.updateOutput(p.outputLabel, fmt.Sprintf(
		"[Booker][Error] Attempting to book parking space: Days since epoch: %v, SpotID: %v, ParkingID %v",
		daysSinceEpoch,
		parkingSpot.SpotID,
		parkingSpot.ParkingID),
	)
	payload := dto.NewBookingPayload(daysSinceEpoch, parkingSpot.SpotID, parkingSpot.ParkingID)
	err := p.parkingClient.Book(payload)
	if err == nil {
		p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker] Successfully booked parking space for SpotID: %v, ParkingID: %v",
			parkingSpot.SpotID,
			parkingSpot.ParkingID,
		))
		return true
	}
	p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker][Error] failed to book. Reason: %v", err))
	p.updateOutput(p.outputLabel, fmt.Sprintf("[Booker]Retrying in : %v seconds", 1))
	time.Sleep(1 * time.Second)
	return false
}
