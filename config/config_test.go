package config

import (
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/reporter"
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
			DeployerPrivateKeyHex: "<deployerPrivateKeyHex>",
			AgentPrivateKeyHex:    "<agentPrivateKeyHex>",
		},
		FuzzerConfig: FuzzerConfig{
			CritialInstructions: []string{"<critical1>"},
			BatchSize:           1003,
			SeedsSize:           1004,
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
			Type:            reporter.CONSOLE_REPORTER,
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
