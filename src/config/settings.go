package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type FiltersSettings struct {
	BlackList   []string `json:"black_list" yaml:"black_list"`
	WhiteList   []string `json:"white_list" yaml:"white_list"`
	CurrentList string   `json:"current_list" yaml:"current_list"`
}

type Settings struct {
	Filters FiltersSettings `json:"filters" yaml:"filters"`
}

func LoadSettings() (*Settings, error) {
	data, err := os.ReadFile("settings.json")
	if os.IsNotExist(err) {
		return &Settings{
			Filters: FiltersSettings{
				BlackList:   []string{},
				WhiteList:   []string{},
				CurrentList: "Black list",
			},
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("[LoadSettings] failed read file: %s", err)
	}

	settings := &Settings{}
	if err := json.Unmarshal(data, settings); err != nil {
		return nil, fmt.Errorf("[LoadSettings] failed parse data: %s", err)
	}

	if settings.Filters.CurrentList == "" {
		settings.Filters.CurrentList = "Black list"
	}

	return settings, nil
}

func (settings *Settings) Dump() error {
	file, err := os.OpenFile("settings.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("[Settings] [Dump] failed open file: %s", err)
	}
	if err := json.NewEncoder(file).Encode(settings); err != nil {
		return fmt.Errorf("[Settings] [Dump] failed dump data: %s", err)
	}
	return nil
}
