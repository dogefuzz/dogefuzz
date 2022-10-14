package controller

import (
	"errors"
	"net/http"

	"github.com/dogefuzz/dogefuzz/dto"
	"github.com/dogefuzz/dogefuzz/pkg/oracle"
	"github.com/dogefuzz/dogefuzz/repo"
	"github.com/gin-gonic/gin"
)

type WeaknessesController interface {
	Create(c *gin.Context)
}

type weaknessesController struct {
	transactionRepo repo.TransactionRepo
	taskRepo        repo.TaskRepo
	taskOracleRepo  repo.TaskOracleRepo
	oracleRepo      repo.OracleRepo
}

func NewWeaknessesController(e Env) *weaknessesController {
	return &weaknessesController{
		transactionRepo: e.TransactionRepo(),
		taskRepo:        e.TaskRepo(),
		taskOracleRepo:  e.TaskOracleRepo(),
		oracleRepo:      e.OracleRepo(),
	}
}

func (ctrl *weaknessesController) Create(c *gin.Context) {
	var request dto.NewWeaknessDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := ctrl.transactionRepo.FindByBlockchainHash(request.TxHash)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	task, err := ctrl.taskRepo.Find(transaction.TaskId)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tasksOracles, err := ctrl.taskOracleRepo.FindByTaskId(task.Id)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var oracleIds []string
	for _, to := range tasksOracles {
		oracleIds = append(oracleIds, to.OracleId)
	}

	oracleEntities, err := ctrl.oracleRepo.FindByIds(oracleIds)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var oracleNames []string
	for _, entity := range oracleEntities {
		oracleNames = append(oracleNames, entity.Name)
	}
	oracles := oracle.GetOracles(oracleNames)
	snapshot := oracle.NewEventsSnapshot(request.OracleEvents)

	var weaknesses []string
	for _, o := range oracles {
		if o.Detect(snapshot) {
			weaknesses = append(weaknesses, o.Name())
		}
	}

	transaction.SetDetectedWeaknesses(weaknesses)
	ctrl.transactionRepo.Update(transaction)
	if err != nil {
		if errors.Is(err, repo.ErrNotExists) {
			c.AbortWithStatus(404)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(200)
}
