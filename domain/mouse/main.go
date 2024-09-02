package mouse

import (
	"booking/port/contract"
	"context"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	"math/rand"
	"time"
)

func Move(startTime, endTime string) error {
	start, err := time.Parse("02-01-2006 15:04:05", startTime)
	if err != nil {
		return fmt.Errorf("failed to parse start time: %v", err)
	}

	end, err := time.Parse("02-01-2006 15:04:05", endTime)
	if err != nil {
		return fmt.Errorf("failed to parse end time: %v", err)
	}

	location := time.Now().Location()
	start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), location)
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), location)

	now := time.Now()
	fmt.Printf("Current time: %v\n", now)
	fmt.Printf("Start time: %v (Zone: %v)\n", start, start.Location())
	fmt.Printf("End time: %v (Zone: %v)\n", end, end.Location())

	if now.Before(start) {
		fmt.Println(fmt.Sprintf("Waiting for start time. Starting in %v seconds", time.Until(start).Seconds()))
		time.Sleep(time.Until(start))
	} else if now.After(end) {
		return fmt.Errorf("end time has already passed")
	}

	screenWidth, screenHeight := robotgo.GetScreenSize()
	fmt.Printf("Screen resolution: %dx%d\n", screenWidth, screenHeight)

	for {
		now = time.Now()
		if now.After(end) {
			fmt.Println("End time reached. Stopping mouse movement.")
			break
		}

		x := rand.Intn(screenWidth)
		y := rand.Intn(screenHeight)

		robotgo.MoveSmooth(x, y, 1.0, 100.0)
		fmt.Printf("Mouse moved to: %d, %d\n", x, y)
		time.Sleep(10 * time.Second)
	}
	return nil
}

func MoveWithContext(ctx context.Context, startTime, endTime string, outputLabel *widget.Label, updateOutput contract.UpdateOutputFunc) {
	start, err := time.Parse("02-01-2006 15:04:05", startTime)
	if err != nil {
		updateOutput(outputLabel, fmt.Sprintf("[MouseMover][Error] failed to parse start time: %v", err))
		return
	}
	end, err := time.Parse("02-01-2006 15:04:05", endTime)
	if err != nil {
		updateOutput(outputLabel, fmt.Sprintf("[MouseMover][Error] failed to parse end time: %v", err))
		return
	}

	location := time.Now().Location()
	start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), location)
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), location)
	now := time.Now()

	updateOutput(outputLabel, fmt.Sprintf("[MouseMover] Current time: %v", now))
	updateOutput(outputLabel, fmt.Sprintf("[MouseMover] Start time: %v (Zone: %v)", start, start.Location()))
	updateOutput(outputLabel, fmt.Sprintf("[MouseMover] End time: %v (Zone: %v)", end, end.Location()))

	if now.Before(start) {
		updateOutput(outputLabel, fmt.Sprintf("[MouseMover] Waiting for start time. Starting in %v seconds", time.Until(start).Seconds()))
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Until(start)):
			// Continue execution
		}
	} else if now.After(end) {
		updateOutput(outputLabel, fmt.Sprintf("[MouseMover][Error] end time has already passed"))
		return
	}

	screenWidth, screenHeight := robotgo.GetScreenSize()
	updateOutput(outputLabel, fmt.Sprintf("[MouseMover] Screen resolution: %dx%d", screenWidth, screenHeight))

	for {
		select {
		case <-ctx.Done():
			updateOutput(outputLabel, "[MouseMover] Process stopped.")
			return
		default:
			now = time.Now()
			if now.After(end) {
				updateOutput(outputLabel, "[MouseMover] End time reached. Stopping mouse movement.")
				return
			}
			//
			//x := rand.Intn(screenWidth)
			//y := rand.Intn(screenHeight)

			robotgo.MoveSmooth(rand.Intn(200), rand.Intn(200), 5.0, 15.0)
		}
	}
}
