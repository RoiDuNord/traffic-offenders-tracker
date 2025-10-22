package processor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"speed_violation_tracker/models"
)

func (p *Processor) processMessage(workerID int, msg Message) error {
	// В случае проблем на нашей стороне сообщаем брокеру, что сообщение не потреблено
	msgAcked := true
	defer func() {
		slog.Info(fmt.Sprintf("msg ack: %v", msgAcked))
		if msgAcked {
			msg.Ack()
		}
	}()

	var passage models.Passage
	if err := json.Unmarshal(msg.Bytes(), &passage); err != nil {
		return err
	}
	defer slog.Info("Processed passage", "worker", workerID, "GRN", passage.LicenseNum)

	offender, ok := isViolationsDetected(passage)
	if !ok {
		return nil
	}

	driverInfo := fmt.Sprintf("driver_%s", offender.GRN)
	id, err := p.db.InsertMessage(driverInfo, []byte(offender.GRN))
	if err != nil {
		msgAcked = false
		return &models.FatalError{Reason: err.Error()}
	}

	p.offendersChan <- struct{}{}

	slog.Info("Violation detected", "license", offender.GRN, "db id", id)
	return nil
}

func isViolationsDetected(passage models.Passage) (models.Offender, bool) {
	if len(passage.Track) == 0 {
		slog.Error("zero track len")
		return models.Offender{}, false
	}

	maxTimeStamp := passage.Track[0].T
	for _, point := range passage.Track[1:] {
		if point.T > maxTimeStamp {
			maxTimeStamp = point.T
		}
	}

	if maxTimeStamp%60 < 45 {
		return models.Offender{}, false
	}

	return models.Offender{GRN: passage.LicenseNum}, true
}
