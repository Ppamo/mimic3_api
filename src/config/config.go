package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"ppamo/api/common"
)

var config *ConfigStruct

type OutputFormat int32

const (
	OutputFormatUnknown OutputFormat = iota
	OutpuFormatDisabled
	OutputFormatLibOpus
	OutputFormatLibMp3Lame
)

type OutputOptions struct {
	Library   string   `json:"library"`
	Extension string   `json:"extension"`
	Options   []string `json:"options"`
}

type ConfigStruct struct {
	TempFolder string                        `json:"temp_folder"`
	Output     *OutputOptions                `json:"output"`
	Profiles   []common.ProfileOptionsStruct `json:"profiles"`
	Effects    []common.AudioEffectStruct    `json:"effects"`
}

func LoadConfigFromContent(content string) error {
	var (
		f   *os.File
		err error
	)
	f, err = os.CreateTemp("", "config.*.json")
	if err != nil {
		return fmt.Errorf("cf> Could not create config file\n%v", err)
	}
	defer os.Remove(f.Name())

	_, err = f.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("cf> Could not write to new file\n%v", err)
	}
	err = LoadConfig(f.Name())
	if err != nil {
		return err
	}
	return nil
}

func LoadConfig(path string) error {
	config = nil
	log.Printf("cf> Reading base config from: %s", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("cf> Error: failed to read base config file")
		return err
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Printf("cf> Error: failed to unmarshal content: %s", content)
		return err
	}
	if len(config.Profiles) == 0 {
		return fmt.Errorf("No profiles found")
	}
	if config.Output == nil {
		return fmt.Errorf("No output format defined")
	}
	log.Printf("cf> Loaded:\n%s", config)
	return nil
}

func GetConfig() *ConfigStruct {
	return config
}

func (c *ConfigStruct) GetProfileByName(name string) (*common.ProfileOptionsStruct, error) {
	if config == nil || config.Profiles == nil {
		return nil, fmt.Errorf("No profiles loaded")
	}
	for _, e := range config.Profiles {
		if e.Name == name {
			return &e, nil
		}
	}
	return nil, fmt.Errorf("No profile found")
}

func (c *ConfigStruct) GetEffectByName(name string) (*common.AudioEffectStruct, error) {
	if config == nil || config.Effects == nil {
		return nil, fmt.Errorf("No effects loaded")
	}
	for _, e := range config.Effects {
		if e.Name == name {
			return &e, nil
		}
	}
	return nil, fmt.Errorf("No effect found")
}

func (c *ConfigStruct) String() string {
	return fmt.Sprintf("temp_folder: %s; profiles: %d; effects: %d; output: %s - %s %v", c.TempFolder,
		len(c.Profiles), len(c.Effects), c.Output.Extension, c.Output.Library, c.Output.Options)
}
