package bdd

import (
	"fmt"
	"os"
	"time"

	"github.com/cucumber/godog"

	"github.com/aholbreich/tl/cmd"
	"github.com/aholbreich/tl/internal/events"
)

// --- actor.feature support ------------------------------------------------

func initializeActorSteps(ctx *godog.ScenarioContext, w *world) {
	ctx.Step(`^environment variable "([^"]*)" is "([^"]*)"$`, w.setEnv)
	ctx.Step(`^the detected agent is "([^"]*)"$`, w.setDetectedAgent)
	ctx.Step(`^the claim expiry for "([^"]*)" is extended$`, w.claimExpiryIsExtended)
	ctx.Step(`^an event "([^"]*)" is recorded for "([^"]*)" by "([^"]*)"$`, w.eventRecordedForBy)
}

func (w *world) eventRecordedForBy(eventName, taskID, actor string) error {
	if w.cmdErr != nil {
		return w.cmdErr
	}
	journal, err := events.ReadAll(".tl")
	if err != nil {
		return err
	}
	found := false
	for _, event := range journal {
		if event.Event != eventName || event.TaskID != taskID {
			continue
		}
		found = true
		if event.Actor != actor {
			return fmt.Errorf("%s event for %s has actor %q, want %q", eventName, taskID, event.Actor, actor)
		}
	}
	if !found {
		return fmt.Errorf("no %s event for %s", eventName, taskID)
	}
	return nil
}

func (w *world) setEnv(key, value string) error {
	w.envOverrides = append(w.envOverrides, key)
	return os.Setenv(key, value)
}

func (w *world) setDetectedAgent(agent string) error {
	cmd.DetectedActor = func() string { return agent }
	return nil
}

func (w *world) claimExpiryIsExtended(id string) error {
	t, err := loadFixtureTask(id)
	if err != nil {
		return err
	}
	if t.Claim.ExpiresAt == nil {
		return fmt.Errorf("task %s has no claim expiry", id)
	}
	// Fixture tasks set UpdatedAt to fixtureTime (2026-05-16T12:00Z).
	// A successful renewal bumps UpdatedAt to time.Now().
	if t.UpdatedAt.Equal(fixtureTime) {
		return fmt.Errorf("task %s was not renewed (updated_at still at fixture time)", id)
	}
	if !t.Claim.ExpiresAt.After(time.Now().UTC()) {
		return fmt.Errorf("task %s claim expiry %v is not in the future", id, t.Claim.ExpiresAt)
	}
	return nil
}
