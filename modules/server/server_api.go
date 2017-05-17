package server

import (
	"github.com/amlun/linda/linda"
	"github.com/amlun/linda/linda/core"
	"github.com/twinj/uuid"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

type api struct {
	linda *linda.Linda
}

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// GET /ping
func (a *api) ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, "pong")
	}
}

// POST /task push task and push a job of task
func (a *api) pushTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var task core.Task
		err := c.Bind(&task)
		if err != nil {
			panic(err)
		}
		if task.TaskId == "" {
			task.TaskId = uuid.NewV4().String()
		}
		if task.Func == "" {
			panic("Func can not be empty")
		}
		if task.Period > 0 && task.Period < 60 {
			panic("Period too quickly, at least one minute")
		}
		err = a.linda.PushTask(task)
		if err != nil {
			panic(err)
		}
		var job = core.Job{
			JobId: uuid.NewV4().String(),
			Task:  task,
		}
		err = a.linda.PushJob(job)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, Result{
			Code: 200,
			Msg:  "ok",
			Data: task,
		})
	}
}

func (a *api) getJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		var result = Result{
			Code: 404,
			Msg:  "not found",
		}
		queue := c.Query("queue")
		if queue == "" {
			panic("queue can not be empty")
		}
		job := a.linda.GetJob(queue)
		if job.JobId != "" {
			result.Code = 200
			result.Msg = "ok"
			result.Data = job
		}
		c.JSON(http.StatusOK, result)
	}
}

func (a *api) queuesStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		queueStatus := a.linda.MonitorQueues()
		c.JSON(http.StatusOK, Result{
			Code: 200,
			Msg:  "ok",
			Data: queueStatus,
		})
	}
}

func (a *api) tasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var taskList core.TaskList
		taskList.State = c.Query("state")
		if err := a.linda.TaskList(&taskList); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, Result{
			Code: 200,
			Msg:  "ok",
			Data: taskList,
		})
	}
}
