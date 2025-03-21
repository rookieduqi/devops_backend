package controller

import (
	"bluebell/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

func GetNodeJobsT(c *gin.Context) {
	var reqData models.RequestJobData
	if err := c.ShouldBindQuery(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "参数绑定失败"})
		return
	}

	// 构造 Jenkins API 请求 URL
	jenkinsURL := fmt.Sprintf("http://%s:%s/me/my-views/view/all/job/%s/api/json",
		reqData.Host, reqData.Port, reqData.ViewID)

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
	//zap.L().Info("body", zap.ByteString("body", body))

	// 解析 JSON 数据
	var data models.JenkinsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "JSON 解析失败"})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data.Jobs})
}

// 获取 Jenkins View Jobs 数据
func GetNodeJobs(c *gin.Context) {
	var reqData models.RequestJobData
	if err := c.ShouldBindQuery(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "参数绑定失败"})
		return
	}

	// 构造 Jenkins API 请求 URL
	jenkinsURL := fmt.Sprintf("http://%s:%s/me/my-views/view/all/job/%s/api/json",
		reqData.Host, reqData.Port, reqData.ViewID)

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
	//zap.L().Info("body", zap.ByteString("body", body))

	// 解析 JSON 数据
	var data models.JenkinsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "JSON 解析失败"})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data.Jobs})
}

// 构建指定任务
func buildJob(ctx context.Context, jenkins *gojenkins.Jenkins, name string) (n int64) {
	var err error
	n, err = jenkins.BuildJob(ctx, name, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("构建完成：", n) // n是序号
	return
}

// 构建指定目录下的某个 Job
func buildJobInFolder(ctx context.Context, jenkins *gojenkins.Jenkins, folderName string, jobName string, params map[string]string) (int64, error) {
	// 获取指定目录 (Folder) 下的 Job
	job, err := jenkins.GetJob(ctx, jobName, folderName)
	if err != nil {
		return 0, fmt.Errorf("获取 Job [%s] 失败: %v", jobName, err)
	}
	fmt.Printf("成功获取 Job [%s] (URL: %s)\n", job.GetName(), job.GetDetails().URL)

	// 触发构建
	queueID, err := job.InvokeSimple(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("触发 Job [%s] 的构建失败: %v", jobName, err)
	}

	return queueID, nil
}

func StartNodeJobsT(c *gin.Context) {
	var reqData models.StartJobRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "参数绑定失败"})
		return
	}
	zap.L().Info("reqData", zap.Any("reqData", reqData))

	ctx := context.Background()
	// 创建 Jenkins 实例
	jenkinsURL := fmt.Sprintf("http://%s:%s", reqData.Host, reqData.Port)
	jenkins := gojenkins.CreateJenkins(nil, jenkinsURL, reqData.Account, reqData.Password)
	_, err := jenkins.Init(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "初始化 Jenkins 实例失败"})
		return
	}
	if reqData.ViewName != "" {
		_, _ = buildJobInFolder(ctx, jenkins, reqData.ViewID, reqData.ViewName, map[string]string{})
	} else {
		buildJob(ctx, jenkins, reqData.ViewID)
	}

	time.Sleep(2 * time.Second)
	// 请求已发起，立即返回成功响应
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "任务启动中，请稍后查看 Jenkins 构建状态"})
}

func stopBuildByJobLatest(ctx context.Context, jenkins *gojenkins.Jenkins, name string) {

	job, err := jenkins.GetJob(ctx, name)
	if err != nil {
		panic(err)
	}

	lastBuild, err := job.GetLastBuild(ctx)
	if err != nil {
		return
	}
	number := lastBuild.Raw.ID
	fmt.Println("准备停止：", number)
	stopped, err := lastBuild.Stop(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("是否停止：", stopped)
}

func StopNodeJobsT(c *gin.Context) {
	var reqData models.StartJobRequest
	if err := c.ShouldBindQuery(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "参数绑定失败"})
		return
	}

	ctx := context.Background()
	// 创建 Jenkins 实例
	jenkinsURL := fmt.Sprintf("http://%s:%s", reqData.Host, reqData.Port)

	// 创建 Jenkins 实例

	jenkins := gojenkins.CreateJenkins(nil, jenkinsURL, reqData.Account, reqData.Password)
	_, err := jenkins.Init(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "初始化 Jenkins 实例失败"})
		return
	}

	stopBuildByJobLatest(ctx, jenkins, reqData.ViewID)

	// 返回数据
	c.JSON(http.StatusOK, gin.H{"success": true, "data": ""})
}

// 启动 Jenkins Job (仅发起请求，不等待返回)
func StartNodeJobs(c *gin.Context) {
	var reqData models.StartJobRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "参数绑定失败"})
		return
	}

	// 构造 Jenkins API 请求 URL
	//jenkinsURL := fmt.Sprintf("http://%s:%s/job/%s/job/%s/build?delay=0sec",
	//	reqData.Host, reqData.Port, reqData.ViewID, reqData.JobName)
	jenkinsURL := fmt.Sprintf("http://%s:%s/job/%s/build",
		reqData.Host, reqData.Port, reqData.ViewID)

	fmt.Println("jenkinsURL===", jenkinsURL)

	// 异步触发 Jenkins 构建
	// 异步触发 Jenkins 构建
	go func() {
		client := &http.Client{Timeout: 5 * time.Second}
		req, _ := http.NewRequest("POST", jenkinsURL, nil) // ✅ 请求方法改为 POST

		// 设置 Basic Auth 认证
		req.SetBasicAuth(reqData.Account, reqData.Password)

		// 发送请求 (不处理返回结果)
		_, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ [StartNodeJobs] 请求失败: %v\n", err)
		} else {
			fmt.Println("✅ [StartNodeJobs] 任务成功触发")
		}
	}()
	time.Sleep(5 * time.Second)

	// 请求已发起，立即返回成功响应
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "任务启动中，请稍后查看 Jenkins 构建状态"})
}

func StopNodeJobs(c *gin.Context) {
	var reqData models.RequestJobData
	if err := c.ShouldBindQuery(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "参数绑定失败"})
		return
	}

	// 构造 Jenkins API 请求 URL
	jenkinsURL := fmt.Sprintf("http://%s:%s/me/my-views/view/all/job/%s/api/json",
		reqData.Host, reqData.Port, reqData.ViewID)

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

	// 解析 JSON 数据
	var data models.JenkinsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "JSON 解析失败"})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data.Jobs})
}
