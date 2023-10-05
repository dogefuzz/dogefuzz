package job

import (
	"time"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/bus"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"go.uber.org/zap"
)

type transactionsTimeoutCheckerJob struct {
	cfg                   *config.Config
	logger                *zap.Logger
	transactionService    interfaces.TransactionService
	taskInputRequestTopic interfaces.Topic[bus.TaskInputRequestEvent]
}

func NewTransactionsTimeoutCheckerJob(e Env) *transactionsTimeoutCheckerJob {
	return &transactionsTimeoutCheckerJob{
		cfg:                   e.Config(),
		logger:                e.Logger(),
		transactionService:    e.TransactionService(),
		taskInputRequestTopic: e.TaskInputRequestTopic(),
	}
}

func (j *transactionsTimeoutCheckerJob) Id() string         { return "transactions-timeout" }
func (j *transactionsTimeoutCheckerJob) CronConfig() string { return "*/5 * * * * *" }

func (j *transactionsTimeoutCheckerJob) Handler() {
	now := common.Now()
	threshold := now.Add(-1 * j.cfg.FuzzerConfig.TransactionTimeout)
	transactions, err := j.transactionService.FindRunningAndCreatedBeforeThreshold(threshold)
	if err != nil {
		j.logger.Sugar().Errorf("an error occured when retrieving transactions that are still running and expired: %v", err)
		return
	}

	if len(transactions) == 0 {
		return
	}

	for _, transaction := range transactions {
		transaction.Status = common.TRANSACTION_TIMEOUT
	}

	maxRetries := 5
	for retries := 0; retries < maxRetries; retries++ {
		err = j.transactionService.BulkUpdate(transactions)
		if err != nil {
			j.logger.Sugar().Warnf("an error occured when updating transactions: %v", err)
			j.logger.Sugar().Warnf("retrying...%d", retries)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
}
