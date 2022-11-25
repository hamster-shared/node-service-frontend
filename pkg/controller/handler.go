package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hamster-shared/node-api/pkg/controller/parameters"
	"github.com/hamster-shared/node-api/pkg/dispatcher"
	"github.com/hamster-shared/node-api/pkg/model"
	"github.com/hamster-shared/node-api/pkg/service"
	"gopkg.in/yaml.v3"
	"strconv"
)

type HandlerServer struct {
	jobService      service.IJobService
	dispatch        dispatcher.IDispatcher
	templateService service.ITemplateService
}

func NewHandlerServer(jobService service.IJobService, dispatch dispatcher.IDispatcher, templateService service.ITemplateService) *HandlerServer {
	return &HandlerServer{
		jobService:      jobService,
		dispatch:        dispatch,
		templateService: templateService,
	}
}

// createPipeline create pipeline jon
func (h *HandlerServer) createPipeline(gin *gin.Context) {
	createData := parameters.CreatePipeline{}
	err := gin.BindJSON(&createData)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	var jobData model.Job
	err = yaml.Unmarshal([]byte(createData.Yaml), &jobData)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	err = h.jobService.SaveJob(createData.Name, &jobData)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}

	Success("", gin)
}

func (h *HandlerServer) updatePipeline(gin *gin.Context) {
	oldName := gin.Param("oldName")
	updateData := parameters.UpdatePipeline{}
	err := gin.BindJSON(&updateData)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	var jobData model.Job
	err = yaml.Unmarshal([]byte(updateData.Yaml), &jobData)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	err = h.jobService.UpdateJob(oldName, updateData.NewName, &jobData)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	Success("", gin)
}

// getPipeline get pipeline job
func (h *HandlerServer) getPipeline(gin *gin.Context) {
	name := gin.Param("name")
	pipelineData := h.jobService.GetJob(name)
	Success(pipelineData, gin)
}

// deletePipeline delete pipeline job and pipeline job detail
func (h *HandlerServer) deletePipeline(gin *gin.Context) {
	name := gin.Param("name")
	err := h.jobService.DeleteJob(name)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	Success("", gin)
}

// pipelineList get pipeline job list
func (h *HandlerServer) pipelineList(gin *gin.Context) {
	query := gin.Query("query")
	pageStr := gin.DefaultQuery("page", "1")
	sizeStr := gin.DefaultQuery("size", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	jobData := h.jobService.JobList(query, page, size)
	Success(jobData, gin)
}

// getPipelineDetail get pipeline job detail info
func (h *HandlerServer) getPipelineDetail(gin *gin.Context) {
	name := gin.Param("name")
	idStr := gin.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	jobDetailData := h.jobService.GetJobDetail(name, id)
	Success(jobDetailData, gin)
}

// deleteJobDetail delete job detail
func (h *HandlerServer) deleteJobDetail(gin *gin.Context) {
	name := gin.Param("name")
	idStr := gin.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	err = h.jobService.DeleteJobDetail(name, id)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	Success("", gin)
}

// getPipelineDetailList get pipeline job detail list
func (h *HandlerServer) getPipelineDetailList(gin *gin.Context) {
	name := gin.Param("name")
	pageStr := gin.DefaultQuery("page", "1")
	sizeStr := gin.DefaultQuery("size", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	jobDetailPage := h.jobService.JobDetailList(name, page, size)
	Success(jobDetailPage, gin)
}

// execPipeline exec pipeline job
func (h *HandlerServer) execPipeline(gin *gin.Context) {
	name := gin.Param("name")
	job := h.jobService.GetJobObject(name)
	jobDetail, err := h.jobService.ExecuteJob(name)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	node := h.dispatch.DispatchNode(job)
	h.dispatch.SendJob(jobDetail, node)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	Success("", gin)
}

// reExecuteJob re exec pipeline job detail
func (h *HandlerServer) reExecuteJob(gin *gin.Context) {
	name := gin.Param("name")
	idStr := gin.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	err = h.jobService.ReExecuteJob(name, id)
	job := h.jobService.GetJobObject(name)
	jobDetail := h.jobService.GetJobDetail(name, id)
	node := h.dispatch.DispatchNode(job)
	h.dispatch.SendJob(jobDetail, node)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	Success("", gin)
}

// stopJobDetail stop pipeline job
func (h *HandlerServer) stopJobDetail(gin *gin.Context) {
	name := gin.Param("name")
	idStr := gin.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}

	err = h.jobService.StopJobDetail(name, id)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	job := h.jobService.GetJobObject(name)
	jobDetail := h.jobService.GetJobDetail(name, id)
	node := h.dispatch.DispatchNode(job)
	h.dispatch.CancelJob(jobDetail, node)
	Success("", gin)
}

// getJobLog get pipeline job detail logs
func (h *HandlerServer) getJobLog(gin *gin.Context) {
	name := gin.Param("name")
	idStr := gin.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	data := h.jobService.GetJobLog(name, id)
	Success(data, gin)
}

// getJobStageLog get job detail stage logs
func (h *HandlerServer) getJobStageLog(gin *gin.Context) {
	name := gin.Param("name")
	idStr := gin.Param("id")
	stageName := gin.Param("stagename")
	startStr := gin.DefaultQuery("start", "0")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	start, _ := strconv.Atoi(startStr)
	data := h.jobService.GetJobStageLog(name, id, stageName, start)
	Success(data, gin)
}

// getTemplates get template list
func (h *HandlerServer) getTemplates(gin *gin.Context) {
	data := h.templateService.GetTemplates()
	Success(data, gin)
}

// getTemplateDetail get template detail
func (h *HandlerServer) getTemplateDetail(gin *gin.Context) {
	idStr := gin.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Fail(err.Error(), gin)
		return
	}
	data, _ := h.templateService.GetTemplateDetail(id)
	Success(data, gin)
}
