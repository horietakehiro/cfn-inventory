package config

import (
	"fmt"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
)

type Credential struct {
	ProfileName string
}

type Tag struct {
	Key   string
	Value string
}

type Filters struct {
	Regions         []string
	StackTags       []Tag
	StackNamePrefix string
}

type RootConfig struct {
	Credential Credential
	Filters    Filters
}

type AccountConfig struct {
	Name       string
	Id         string
	Credential Credential
	Filters    Filters
}

type InventoryConfig struct {
	RootConfig       RootConfig
	ConfigByAccounts []AccountConfig
}

func init() {
	config.WithOptions()
	config.AddDriver(yamlv3.Driver)
}

func setDefaultConfig(config *InventoryConfig) {
	for i := range config.ConfigByAccounts {
		if config.ConfigByAccounts[i].Credential.ProfileName == "" {
			if config.RootConfig.Credential.ProfileName != "" {
				config.ConfigByAccounts[i].Credential.ProfileName = config.RootConfig.Credential.ProfileName
			} else {
				config.ConfigByAccounts[i].Credential.ProfileName = ""
			}
		}
		if len(config.ConfigByAccounts[i].Filters.Regions) == 0 {
			if len(config.RootConfig.Filters.Regions) != 0 {
				config.ConfigByAccounts[i].Filters.Regions = config.RootConfig.Filters.Regions
			} else {
				config.ConfigByAccounts[i].Filters.Regions = []string{"*"}

			}
		}
		if config.ConfigByAccounts[i].Filters.StackNamePrefix == "" {
			if config.RootConfig.Filters.StackNamePrefix != "" {
				config.ConfigByAccounts[i].Filters.StackNamePrefix = config.RootConfig.Filters.StackNamePrefix
			} else {
				config.ConfigByAccounts[i].Filters.StackNamePrefix = ""
			}
		}
		if len(config.ConfigByAccounts[i].Filters.StackTags) == 0 {
			if len(config.RootConfig.Filters.StackTags) != 0 {
				config.ConfigByAccounts[i].Filters.StackTags = config.RootConfig.Filters.StackTags
			} else {
				config.ConfigByAccounts[i].Filters.StackTags = []Tag{}
			}
		}
	}
}

func validate(config *InventoryConfig) error {
	err := []string{}
	for i, accountConfig := range config.ConfigByAccounts {
		if accountConfig.Id == "" {
			err = append(err, fmt.Sprintf("ConfigByAccounts[%v].Id is required", i))
		}
		if len(accountConfig.Filters.StackTags) == 0 && accountConfig.Filters.StackNamePrefix == "" {
			err = append(err, fmt.Sprintf("either ConfigByAccounts[%v].Filters.StackTags or ConfigByAccounts[%v].Filters.StackNamePrefix are required", i, i))
		}
	}

	if len(err) == 0 {
		return nil
	} else {
		return fmt.Errorf("%s", strings.Join(err, "; "))
	}
}

func GetConfig(filePath string) (*InventoryConfig, error) {
	inventoryConfig := &InventoryConfig{}

	err := config.LoadFiles(filePath)
	if err != nil {
		return inventoryConfig, err
	}

	err = config.BindStruct("", &inventoryConfig)
	if err != nil {
		return inventoryConfig, err
	}

	setDefaultConfig(inventoryConfig)
	err = validate(inventoryConfig)
	if err != nil {
		return inventoryConfig, err
	}

	return inventoryConfig, nil

}
