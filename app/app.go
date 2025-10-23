package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"speed_violation_tracker/broker"
	"speed_violation_tracker/cat"
	"speed_violation_tracker/config"
	"speed_violation_tracker/dog"
	"speed_violation_tracker/processor"
	"speed_violation_tracker/repo"
	"syscall"
)

func MustRun() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer slog.Info("Tracker closed")

	cfg, err := config.MustLoad()
	if err != nil {
		return fmt.Errorf("config load failed: %w", err)
	}
	slog.Info("Config loaded", "host", cfg.Server.Host, "port", cfg.Server.Port)

	dog := dog.New()
	if err := dog.Connect(cfg.Dog.ConnString); err != nil {
		return fmt.Errorf("dog connection failed: %w", err)
	}
	defer func() {
		dog.Close()
		slog.Info("Dog closed")
	}()
	slog.Info("DB Dog connected")
	dbService := repo.New(dog)

	cat := cat.New()
	if err := cat.Connect(cfg.Cat.ConnString); err != nil {
		return fmt.Errorf("cat connection failed: %w", err)
	}
	defer func() {
		cat.Close()
		slog.Info("Cat closed")
	}()
	slog.Info("Broker Cat connected")
	brokerService := broker.New(cat)

	processor := processor.New(ctx, dbService, brokerService, cfg.MaxOffenders, cfg.CloseTimeout)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		cancel()
	}()

	return processor.Run()
}
