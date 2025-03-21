package controller

import (
	"bluebell/models"
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

// 获取和删除指定目录下某个 Job 的最新构建
func getAndDeleteLatestBuildInFolder(ctx context.Context, jenkins *gojenkins.Jenkins, folderName string, jobName string) error {
	// 获取指定目录 (Folder) 下的 Job
	job, err := jenkins.GetJob(ctx, jobName, folderName)
	if err != nil {
		return fmt.Errorf("获取 Job [%s] 失败: %v", jobName, err)
	}

	// 获取最新构建
	lastBuild, err := job.GetLastBuild(ctx)
	if err != nil {
		return fmt.Errorf("获取 Job [%s] 的最新构建失败: %v", jobName, err)
	}

	// 打印最新构建信息
	buildNumber := lastBuild.GetBuildNumber()
	buildResult := lastBuild.GetResult()
	fmt.Printf("获取 Job [%s] 的最新构建成功:\n", jobName)
	fmt.Printf("   - 构建编号: %d\n", buildNumber)
	fmt.Printf("   - 构建状态: %s\n", buildResult)

	// 删除最新构建
	deleteURL := fmt.Sprintf("%s/%d/doDelete", job.Base, buildNumber)
	resp, err := jenkins.Requester.Post(ctx, deleteURL, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("删除 Job [%s] 的最新构建失败: %v", jobName, err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("删除 Job [%s] 的最新构建失败: HTTP 状态码 %d", jobName, resp.StatusCode)
	}

	fmt.Printf("✅ 成功删除 Job [%s] 的最新构建 (构建编号: %d)\n", jobName, buildNumber)
	return nil
}

func ConsoleBuildDelete(c *gin.Context) {
	var reqData models.RequestJobData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid JSON data"})
		return
	}

	ctx := context.Background()
	jenkinsURL := fmt.Sprintf("http://%s:%s", reqData.Host, reqData.Port)

	// 创建 Jenkins 实例
	jenkins := gojenkins.CreateJenkins(nil, jenkinsURL, reqData.Account, reqData.Password)
	_, err := jenkins.Init(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "初始化 Jenkins 实例失败"})
		return
	}

	_ = getAndDeleteLatestBuildInFolder(ctx, jenkins, reqData.ViewID, reqData.JobName)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": "删除成功"})
}

func ConsoleBuildPrevious(c *gin.Context) {
	var reqData models.RequestJobData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid JSON data"})
		return
	}

	zap.L().Info("reqData", zap.Any("reqData", reqData))
	// /job/GMB/job/GmbClient/lastSuccessfulBuild/pipeline-console/allSteps
	jenkinsURL := fmt.Sprintf("http://%s:%s/job/%s/lastSuccessfulBuild/pipeline-console/allSteps", reqData.Host, reqData.Port, reqData.ViewID)

	if reqData.JobName != "" {
		jenkinsURL = fmt.Sprintf("http://%s:%s/job/%s/job/%s/lastSuccessfulBuild/pipeline-console/allSteps", reqData.Host, reqData.Port, reqData.ViewID, reqData.JobName)
	}

	// 构造 Jenkins API URL
	//http://172.24.65.29:10001/job/GMB/job/GmbClient/lastBuild/consoleText

	// 构造 HTTP 请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", jenkinsURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "构造请求失败"})
		return
	}

	// 设置 Basic Auth 认证
	req.SetBasicAuth(reqData.Account, reqData.Password)

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "请求失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fmt.Sprintf("请求失败，状态码：%d", resp.StatusCode)})
		return
	}

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "读取响应失败"})
		return
	}
	zap.L().Info("body", zap.ByteString("body", body))

	c.JSON(http.StatusOK, gin.H{"success": true, "data": string(body)})
}

func ConsoleBuildNext(c *gin.Context) {
	var reqData models.RequestJobData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid JSON data"})
		return
	}

	zap.L().Info("reqData", zap.Any("reqData", reqData))
	// http://172.24.65.29:10001/job/GMB/job/GmbClient/lastBuild/
	jenkinsURL := "http://172.24.65.29:10001/job/GMB/job/GmbClient/20/"
	//jenkinsURL := fmt.Sprintf("http://%s:%s/job/%s/lastSuccessfulBuild/pipeline-console/allSteps", reqData.Host, reqData.Port, reqData.ViewID)
	//
	//if reqData.JobName != "" {
	//	jenkinsURL = fmt.Sprintf("http://%s:%s/job/%s/job/%s/lastSuccessfulBuild/pipeline-console/allSteps", reqData.Host, reqData.Port, reqData.ViewID, reqData.JobName)
	//}

	// 构造 Jenkins API URL
	//http://172.24.65.29:10001/job/GMB/job/GmbClient/lastBuild/consoleText

	// 构造 HTTP 请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", jenkinsURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "构造请求失败"})
		return
	}

	// 设置 Basic Auth 认证
	req.SetBasicAuth(reqData.Account, reqData.Password)

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "请求失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fmt.Sprintf("请求失败，状态码：%d", resp.StatusCode)})
		return
	}

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "读取响应失败"})
		return
	}
	zap.L().Info("body", zap.ByteString("body", body))

	c.JSON(http.StatusOK, gin.H{"success": true, "data": string(body)})
}

func GetConsolePipeConsole(c *gin.Context) {
	var reqData models.RequestJobData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid JSON data"})
		return
	}

	zap.L().Info("reqData", zap.Any("reqData", reqData))
	// /job/GMB/job/GmbClient/lastSuccessfulBuild/pipeline-console/allSteps
	jenkinsURL := fmt.Sprintf("http://%s:%s/job/%s/lastSuccessfulBuild/pipeline-console/allSteps", reqData.Host, reqData.Port, reqData.ViewID)

	if reqData.JobName != "" {
		jenkinsURL = fmt.Sprintf("http://%s:%s/job/%s/job/%s/lastSuccessfulBuild/pipeline-console/allSteps", reqData.Host, reqData.Port, reqData.ViewID, reqData.JobName)
	}

	// 构造 Jenkins API URL
	//http://172.24.65.29:10001/job/GMB/job/GmbClient/lastBuild/consoleText

	// 构造 HTTP 请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", jenkinsURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "构造请求失败"})
		return
	}

	// 设置 Basic Auth 认证
	req.SetBasicAuth(reqData.Account, reqData.Password)

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "请求失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fmt.Sprintf("请求失败，状态码：%d", resp.StatusCode)})
		return
	}

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "读取响应失败"})
		return
	}
	zap.L().Info("body", zap.ByteString("body", body))

	c.JSON(http.StatusOK, gin.H{"success": true, "data": string(body)})
}

func GetConsolePipeOverview(c *gin.Context) {
	var reqData models.RequestJobData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid JSON data"})
		return
	}

	zap.L().Info("reqData", zap.Any("reqData", reqData))
	jenkinsURL := fmt.Sprintf("http://%s:%s/job/%s/lastBuild/pipeline-graph/tree", reqData.Host, reqData.Port, reqData.ViewID)

	if reqData.JobName != "" {
		jenkinsURL = fmt.Sprintf("http://%s:%s/job/%s/job/%s/lastBuild/pipeline-graph/tree", reqData.Host, reqData.Port, reqData.ViewID, reqData.JobName)
	}

	// 构造 Jenkins API URL
	//http://172.24.65.29:10001/job/GMB/job/GmbClient/lastBuild/consoleText

	// 构造 HTTP 请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", jenkinsURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "构造请求失败"})
		return
	}

	// 设置 Basic Auth 认证
	req.SetBasicAuth(reqData.Account, reqData.Password)

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "请求失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fmt.Sprintf("请求失败，状态码：%d", resp.StatusCode)})
		return
	}

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "读取响应失败"})
		return
	}
	zap.L().Info("body", zap.ByteString("body", body))

	c.JSON(http.StatusOK, gin.H{"success": true, "data": string(body)})
}

func GetNodeConsole(c *gin.Context) {
	var reqData models.RequestJobData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid JSON data"})
		return
	}

	zap.L().Info("reqData", zap.Any("reqData", reqData))
	jenkinsURL := fmt.Sprintf("http://%s:%s/job/%s/lastBuild/consoleText", reqData.Host, reqData.Port, reqData.ViewID)

	if reqData.JobName != "" {
		jenkinsURL = fmt.Sprintf("http://%s:%s/job/%s/job/%s/lastBuild/consoleText", reqData.Host, reqData.Port, reqData.ViewID, reqData.JobName)
	}

	// 构造 Jenkins API URL
	//http://172.24.65.29:10001/job/GMB/job/GmbClient/lastBuild/consoleText

	// 构造 HTTP 请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", jenkinsURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "构造请求失败"})
		return
	}

	// 设置 Basic Auth 认证
	req.SetBasicAuth(reqData.Account, reqData.Password)

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "请求失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fmt.Sprintf("请求失败，状态码：%d", resp.StatusCode)})
		return
	}

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "读取响应失败"})
		return
	}
	zap.L().Info("body", zap.ByteString("body", body))

	c.JSON(http.StatusOK, gin.H{"success": true, "data": string(body)})
}
