package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Name  string `yaml:"name"`
	Price string `yaml:"price"`
}

type BusinessConfig struct {
	Name              string    `yaml:"name"`
	Phone             string    `yaml:"phone"`
	Hours             string    `yaml:"hours"`
	WelcomeMessage    string    `yaml:"welcome_message"`
	OutOfHoursMessage string    `yaml:"out_of_hours_message"`
	Services          []Service `yaml:"services"`
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
