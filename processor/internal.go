package processor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"speed_violation_tracker/models"
)

func (p *Processor) processMessage(workerID int, msg Message) error {
	msgAcked := true
	defer func() {
		slog.Info(fmt.Sprintf("Msg ack: %v", msgAcked))
		if msgAcked {
			msg.Ack()
		}
	}()

	var passage models.Passage
	if err := json.Unmarshal(msg.Bytes(), &passage); err != nil {
		return err
	}
	defer slog.Info("Processed passage", "worker", workerID, "GRN", passage.LicenseNum)

	offender, ok := isOffenseDetected(passage)
	if !ok {
		return nil
	}

	id, err := p.db.InsertMessage(offender.GRN, msg.Bytes())
	if err != nil {
		msgAcked = false
		return &models.FatalError{Cause: err.Error()}
	}

	p.offendersChan <- struct{}{}

	slog.Info("Offense detected", "GRN", offender.GRN, "db id", id)
	return nil
}

func isOffenseDetected(passage models.Passage) (models.Offender, bool) {
	if len(passage.Track) == 0 {
		slog.Error("Zero track len")
		return models.Offender{}, false
	}

	maxTimeStamp := passage.Track[0].T
	for _, fixationPoint := range passage.Track[1:] {
		if fixationPoint.T > maxTimeStamp {
			maxTimeStamp = fixationPoint.T
		}
	}

	if maxTimeStamp%60 < 45 {
		return models.Offender{}, false
	}

	return models.Offender{GRN: passage.LicenseNum}, true
}
