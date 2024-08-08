package clock

import (
	"fmt"
	"go-clock-app/database"
	"log"
	"sync"
	"time"
)

type Signal struct {
	TickMessage  string
	TockMessage  string
	BongMessage  string
	TickInterval time.Duration
	TockInterval time.Duration
	BongInterval time.Duration
}

type Manager struct {
	signals Signal
	Updates chan Signal
	db      *database.Database
	mu      sync.Mutex
}

func NewManager(db *database.Database) *Manager {
	initialSignals := Signal{
		TickMessage:  "tick",
		TockMessage:  "tock",
		BongMessage:  "bong",
		TickInterval: 1 * time.Second,
		TockInterval: 1 * time.Minute,
		BongInterval: 1 * time.Hour,
	}

	return &Manager{
		signals: initialSignals,
		Updates: make(chan Signal, 1),
		db:      db,
	}
}

func (cm *Manager) Run() {
	tickTicker := time.NewTicker(cm.signals.TickInterval)
	tockTicker := time.NewTicker(cm.signals.TockInterval)
	bongTicker := time.NewTicker(cm.signals.BongInterval)
	quit := time.After(3 * time.Hour)

	for {
		select {
		case <-tickTicker.C:
			cm.LogSignal(cm.signals.TickMessage)
		case <-tockTicker.C:
			cm.LogSignal(cm.signals.TockMessage)
		case <-bongTicker.C:
			cm.LogSignal(cm.signals.BongMessage)
		case newsignals := <-cm.Updates:
			cm.mu.Lock()
			cm.signals = newsignals
			if cm.signals.TickInterval > 0 {
				tickTicker.Reset(cm.signals.TickInterval)
			}
			if cm.signals.TockInterval > 0 {
				tockTicker.Reset(cm.signals.TockInterval)
			}
			if cm.signals.BongInterval > 0 {
				bongTicker.Reset(cm.signals.BongInterval)
			}
			cm.mu.Unlock()
		case <-quit:
			log.Println("Clock stopped after three hours.")
			return
		}
	}
}

func (cm *Manager) LogSignal(message string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	fmt.Println(message)
	cm.db.LogSignal(message)
}

func (cm *Manager) UpdateSignals(newsignals Signal) {
	cm.signals.TickMessage = newsignals.TickMessage
	cm.signals.TockMessage = newsignals.TockMessage
	cm.signals.BongMessage = newsignals.BongMessage

	select {
	case cm.Updates <- cm.signals:
		log.Printf("Signal update: tick %s, tock %s, bong %s\n",
			cm.signals.TickMessage, cm.signals.TockMessage, cm.signals.BongMessage)
	default:
		log.Println("Failed to queue signal update: channel is blocked")
	}
}

func (cm *Manager) UpdateIntervals(tickInterval, tockInterval, bongInterval string) error {
	parseDuration := func(interval string) (time.Duration, error) {
		return time.ParseDuration(interval)
	}

	tickDur, err := parseDuration(tickInterval)
	if err != nil {
		return err
	}

	tockDur, err := parseDuration(tockInterval)
	if err != nil {
		return err
	}

	bongDur, err := parseDuration(bongInterval)
	if err != nil {
		return err
	}

	cm.mu.Lock()
	cm.signals.TickInterval = tickDur
	cm.signals.TockInterval = tockDur
	cm.signals.BongInterval = bongDur
	cm.mu.Unlock()

	select {
	case cm.Updates <- cm.signals:
		log.Printf("Updated intervals to: tick %s, tock %s, bong %s\n",
			cm.signals.TickInterval, cm.signals.TockInterval, cm.signals.BongInterval)
	default:
		log.Println("Failed to queue interval update: channel is blocked")
	}

	return nil
}

func (cm *Manager) GetTickInterval() time.Duration {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.signals.TickInterval
}

func (cm *Manager) GetTockInterval() time.Duration {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.signals.TockInterval
}

func (cm *Manager) GetBongInterval() time.Duration {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.signals.BongInterval
}
