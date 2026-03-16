package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Name  string `yaml:"name"`
	Price string `yaml:"price"`
}

type BusinessHour struct {
	Days  []int  `yaml:"days"`
	Open  string `yaml:"open"`
	Close string `yaml:"close"`
}

type BusinessConfig struct {
	Name              string         `yaml:"name"`
	Phone             string         `yaml:"phone"`
	BusinessHours     []BusinessHour `yaml:"business_hours"`
	WelcomeMessage    string         `yaml:"welcome_message"`
	OutOfHoursMessage string         `yaml:"out_of_hours_message"`
	Services          []Service      `yaml:"services"`
}

var businessConfig BusinessConfig

func LoadBusinessConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading business.yaml: %w", err)
	}

	if err := yaml.Unmarshal(data, &businessConfig); err != nil {
		return fmt.Errorf("error parsing business.yaml: %w", err)
	}

	return nil
}

func IsBusinessHour() bool {
	now := time.Now()
	weekday := int(now.Weekday()) // 0=Sunday, ..., 6=Saturday

	for _, bh := range businessConfig.BusinessHours {
		for _, day := range bh.Days {
			if day == weekday {
				openTime, err1 := time.Parse("15:04", bh.Open)
				closeTime, err2 := time.Parse("15:04", bh.Close)
				if err1 != nil || err2 != nil {
					return false
				}
				// Build today's times using hour and minute
				openToday := time.Date(now.Year(), now.Month(), now.Day(), openTime.Hour(), openTime.Minute(), 0, 0, now.Location())
				closeToday := time.Date(now.Year(), now.Month(), now.Day(), closeTime.Hour(), closeTime.Minute(), 0, 0, now.Location())
				if !now.Before(openToday) && !now.After(closeToday) {
					return true
				}
			}
		}
	}
	return false
}
