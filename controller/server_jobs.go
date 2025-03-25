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

// 将 Jenkins 的 color 字段映射为天气图标
func getWeatherStatus(color string) string {
	switch color {
	case "blue", "green":
		return "☀️ Sunny"
	case "yellow":
		return "⛅ Cloudy"
	case "red":
		return "⛈️ Storm"
	case "notbuilt", "grey":
		return "🌫️ Not Built"
	case "aborted":
		return "💨 Windy"
	case "disabled":
		return "❌ Disabled"
	default:
		return "❓ Unknown"
	}
}

// 将 Jenkins 的 Health Score 映射为天气图标
func getWeatherIconByHealthReport(healthScore int64) string {
	switch {
	case healthScore >= 80:
		return "☀️ Sunny"
	case healthScore >= 60:
		return "🌤️ Partly Sunny"
	case healthScore >= 40:
		return "🌥️ Cloudy"
	case healthScore >= 20:
		return "🌧️ Rain"
	case healthScore >= 0:
		return "⛈️ Storm"
	default:
		return "❓ Unknown"
	}
}

// 时间转换函数 (将毫秒转换为可读的时间格式)
func formatDurationT(ms int64) string {
	duration := time.Duration(ms) * time.Millisecond
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d hr %d min", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%d min %d sec", minutes, seconds)
	}
	return fmt.Sprintf("%.1f sec", duration.Seconds())
}

// 获取指定目录 (Folder) 下的所有 Job
func getJobsInFolder(ctx context.Context, jenkins *gojenkins.Jenkins, folderName string) ([]models.JenkinsJob, error) {
	folder, err := jenkins.GetJob(ctx, folderName)
	if err != nil {
		return nil, fmt.Errorf("获取目录 [%s] 失败: %v", folderName, err)
	}

	jobs, err := folder.GetInnerJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取目录 [%s] 下的 Job 失败: %v", folderName, err)
	}

	var jobInfos []models.JenkinsJob
	for _, job := range jobs {
		// 获取健康评分
		healthScore := int64(0)
		if len(job.GetDetails().HealthReport) > 0 {
			healthScore = job.GetDetails().HealthReport[0].Score
		}
		jobInfo := models.JenkinsJob{
			Name:  job.GetName(),
			URL:   job.GetDetails().URL,
			Color: job.GetDetails().Color,
			//Weather:    getWeatherStatus(job.GetDetails().Color),
			Weather:    getWeatherIconByHealthReport(healthScore),
			CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		}

		lastSucess, _ := job.GetLastSuccessfulBuild(ctx)
		if lastSucess != nil {
			jobInfo.LastSuccess = fmt.Sprintf("#%d", lastSucess.GetBuildNumber())
			LastSuccessDuration := int64(lastSucess.GetDuration())
			jobInfo.LastSuccessDuration = formatDurationT(LastSuccessDuration)
		} else {
			jobInfo.LastSuccess = "无"
			jobInfo.LastSuccessDuration = "无"
		}

		lastFailed, _ := job.GetLastFailedBuild(ctx)
		if lastFailed != nil {
			jobInfo.LastFailure = fmt.Sprintf("#%d", lastFailed.GetBuildNumber())
			LastFailureDuration := int64(lastFailed.GetDuration())
			jobInfo.LastFailureDuration = formatDurationT(LastFailureDuration)
		} else {
			jobInfo.LastFailure = "无"
			jobInfo.LastFailureDuration = "无"
		}

		// 获取上次构建时长
		lastBuild, _ := job.GetLastBuild(ctx)
		if lastBuild != nil {
			LastBuildDuration := int64(lastBuild.GetDuration())
			jobInfo.LastDuration = formatDurationT(LastBuildDuration)
		} else {
			jobInfo.LastDuration = "无"
		}

		jobInfos = append(jobInfos, jobInfo)
	}

	return jobInfos, nil
}

func GetNodeJobsT(c *gin.Context) {
	var reqData models.RequestJobData
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

	// 获取 Folder 下的 Jobs
	jobInfos, err := getJobsInFolder(ctx, jenkins, reqData.ViewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{"success": true, "data": jobInfos})
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

// 取消指定 Job 的最新构建
func cancelLatestBuild(ctx context.Context, jenkins *gojenkins.Jenkins, folderName string, jobName string) error {
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

	// 停止构建
	_, err = lastBuild.Stop(ctx)
	if err != nil {
		return fmt.Errorf("停止 Job [%s] 的构建失败: %v", jobName, err)
	}

	fmt.Printf("成功停止 Job [%s] 的最新构建 (构建编号: %d)\n", jobName, lastBuild.GetBuildNumber())
	return nil
}

func StopNodeJobsT(c *gin.Context) {
	var reqData models.StopJobRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		zap.L().Error("err", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "参数绑定失败"})
		return
	}
	zap.L().Info("reqData", zap.Any("reqData", reqData))

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

	if reqData.ViewName != "" {
		_ = cancelLatestBuild(ctx, jenkins, reqData.ViewID, reqData.ViewName)
	} else {
		stopBuildByJobLatest(ctx, jenkins, reqData.ViewID)
	}

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
