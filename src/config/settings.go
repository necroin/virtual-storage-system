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

type ReplicationSettings struct {
	SrcHostname string `json:"src_hostname" yaml:"src_hostname"`
	DstHostname string `json:"dst_hostname" yaml:"dst_hostname"`
	SrcPath     string `json:"src_path" yaml:"src_path"`
	DstPath     string `json:"dst_path" yaml:"dst_path"`
	Cron        string `json:"cron" yaml:"cron"`
}

type Settings struct {
	Filters     FiltersSettings       `json:"filters" yaml:"filters"`
	Replication []ReplicationSettings `json:"replication" yaml:"replication"`
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
			Replication: []ReplicationSettings{},
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

	if settings.Replication == nil {
		settings.Replication = []ReplicationSettings{}
	}

	settings.Dump()

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

func (replication *ReplicationSettings) String() string {
	return fmt.Sprintf("[Host] (%s -> %s): [Path] (%s -> %s)", replication.SrcHostname, replication.DstHostname, replication.SrcPath, replication.DstPath)
}
