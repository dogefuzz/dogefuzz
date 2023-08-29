package config

import (
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/solidity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestLoadConfig_WhenPassingAValidPath_ReturnAValidConfig() {
	expectedConfig := Config{
		StorageFolder: path.Join(os.TempDir(), "dogefuzz"),
		DatabaseName:  "<databaseName>",
		ServerPort:    1001,
		GethConfig: GethConfig{
			NodeAddress:           "<nodeAddress>",
			ChainID:               1002,
			MinGasLimit:           700000000,
			DeployerPrivateKeyHex: "<deployerPrivateKeyHex>",
			AgentPrivateKeyHex:    "<agentPrivateKeyHex>",
		},
		FuzzerConfig: FuzzerConfig{
			CritialInstructions: []string{"<critical1>"},
			BatchSize:           1003,
			SeedsSize:           1004,
			TransactionTimeout:  10 * time.Second,
			Seeds: map[common.TypeIdentifier][]string{
				solidity.BOOL: {"true", "false"},
			},
		},
		VandalConfig: VandalConfig{
			Endpoint: "<endpoint>",
		},
		JobConfig: JobConfig{
			EnabledJobs: []string{"<job1>"},
		},
		EventConfig: EventConfig{
			EnabledListeners: []string{"<listener1>"},
		},
		ReporterConfig: ReporterConfig{
			Type:            common.CONSOLE_REPORTER,
			WebhookEndpoint: "<endpoint>",
			WebhookTimeout:  30 * time.Second,
		},
	}

	cfg, err := LoadConfig("../test/resources")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedConfig.StorageFolder, cfg.StorageFolder)
	assert.Equal(s.T(), expectedConfig.DatabaseName, cfg.DatabaseName)
	assert.Equal(s.T(), expectedConfig.ServerPort, cfg.ServerPort)
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.GethConfig, cfg.GethConfig))
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.FuzzerConfig, cfg.FuzzerConfig))
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.VandalConfig, cfg.VandalConfig))
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.JobConfig, cfg.JobConfig))
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.EventConfig, cfg.EventConfig))
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.ReporterConfig, cfg.ReporterConfig))
}

func (s *ConfigTestSuite) TestLoadConfig_WhenPassingAInvalidPath_ReturnErrConfigNotFound() {
	cfg, err := LoadConfig(gofakeit.Word())

	assert.Equal(s.T(), ErrConfigNotFound, err)
	assert.Nil(s.T(), cfg)
}

func (s *ConfigTestSuite) TestLoadConfig_WhenPassingHaveAnEnvVariable_ReturnAValidConfigWithOverwritedValueByEnvVariable() {
	expectedDatabaseName := "<newDatabaseName>"
	expectedGethNodeAddress := "<newNodeAddress>"
	expectedFuzzerSeedsBool := []string{"false", "true", "false"}
	expectedConfig := Config{
		StorageFolder: path.Join(os.TempDir(), "dogefuzz"),
		DatabaseName:  expectedDatabaseName,
		ServerPort:    1001,
		GethConfig: GethConfig{
			NodeAddress:           expectedGethNodeAddress,
			ChainID:               1002,
			MinGasLimit:           700000000,
			DeployerPrivateKeyHex: "<deployerPrivateKeyHex>",
			AgentPrivateKeyHex:    "<agentPrivateKeyHex>",
		},
		FuzzerConfig: FuzzerConfig{
			CritialInstructions: []string{"<critical1>"},
			BatchSize:           1003,
			SeedsSize:           1004,
			TransactionTimeout:  10 * time.Second,
			Seeds: map[common.TypeIdentifier][]string{
				solidity.BOOL: expectedFuzzerSeedsBool,
			},
		},
		VandalConfig: VandalConfig{
			Endpoint: "<endpoint>",
		},
		JobConfig: JobConfig{
			EnabledJobs: []string{"<job1>"},
		},
		EventConfig: EventConfig{
			EnabledListeners: []string{"<listener1>"},
		},
		ReporterConfig: ReporterConfig{
			Type:            common.CONSOLE_REPORTER,
			WebhookEndpoint: "<endpoint>",
			WebhookTimeout:  30 * time.Second,
		},
	}

	os.Setenv("DOGEFUZZ_DATABASENAME", expectedDatabaseName)
	os.Setenv("DOGEFUZZ_GETH_NODEADDRESS", expectedGethNodeAddress)
	os.Setenv("DOGEFUZZ_FUZZER_SEEDS_BOOL", strings.Join(expectedFuzzerSeedsBool, ","))

	cfg, err := LoadConfig("../test/resources")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedConfig.StorageFolder, cfg.StorageFolder)
	assert.Equal(s.T(), expectedConfig.DatabaseName, cfg.DatabaseName)
	assert.Equal(s.T(), expectedConfig.ServerPort, cfg.ServerPort)
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.GethConfig, cfg.GethConfig))
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.FuzzerConfig, cfg.FuzzerConfig))
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.VandalConfig, cfg.VandalConfig))
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.JobConfig, cfg.JobConfig))
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.EventConfig, cfg.EventConfig))
	assert.True(s.T(), reflect.DeepEqual(expectedConfig.ReporterConfig, cfg.ReporterConfig))
}
