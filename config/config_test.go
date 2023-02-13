package config

import (
	"testing"

	"github.com/gookit/config/v2"
	"github.com/stretchr/testify/assert"
)

func TestConfig_valid(t *testing.T) {
	defer config.ClearAll()

	configPath := "./config_test_valid.yaml"

	inventoryConfig, err := GetConfig(configPath)
	assert.Nil(t, err)

	mainAccount := inventoryConfig.ConfigByAccounts[0]
	assert.Equal(t, "123456789012", mainAccount.Id)
	assert.Equal(t, "main-account", mainAccount.Name)
	assert.Equal(t, "test-profile", mainAccount.Credential.ProfileName)
	assert.Equal(t, "*", mainAccount.Filters.Regions[0])
	assert.Equal(t, "System", mainAccount.Filters.StackTags[0].Key)
	assert.Equal(t, "test-system", mainAccount.Filters.StackTags[0].Value)
	assert.Equal(t, "test-system-", mainAccount.Filters.StackNamePrefix)

	subAccount := inventoryConfig.ConfigByAccounts[1]
	assert.Equal(t, "210987654321", subAccount.Id)
	assert.Equal(t, "sub-account", subAccount.Name)
	assert.Equal(t, "sub-profile", subAccount.Credential.ProfileName)
	assert.Equal(t, "ap-northeast-1", subAccount.Filters.Regions[0])
	assert.Equal(t, "System", subAccount.Filters.StackTags[0].Key)
	assert.Equal(t, "sub-system", subAccount.Filters.StackTags[0].Value)
	assert.Equal(t, "sub-system-", subAccount.Filters.StackNamePrefix)

}

func TestConfig_invalid(t *testing.T) {
	defer config.ClearAll()
	configPath := "./config_test_invalid.yaml"

	_, err := GetConfig(configPath)
	assert.NotNil(t, err)

	assert.Contains(t, err.Error(), "ConfigByAccounts[0].Id is required")
	assert.Contains(t, err.Error(), "either ConfigByAccounts[0].Filters.StackTags or ConfigByAccounts[0].Filters.StackNamePrefix are required")

}
