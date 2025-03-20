package controller

import (
	"bluebell/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

// 获取节点视图列表 (Mock 数据)
func GetNodeViews(c *gin.Context) {
	var reqData models.RequestData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid JSON data"})
		return
	}

	// 构造 Jenkins API URL
	jenkinsURL := fmt.Sprintf("http://%s:%s/api/json?tree=jobs[name,lastSuccessfulBuild[timestamp],lastFailedBuild[timestamp],lastBuild[duration]]",
		reqData.Host, reqData.Port)

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
	var data struct {
		Jobs []struct {
			Name                string `json:"name"`
			LastSuccessfulBuild struct {
				Timestamp int64 `json:"timestamp"`
			} `json:"lastSuccessfulBuild"`
			LastFailedBuild struct {
				Timestamp int64 `json:"timestamp"`
			} `json:"lastFailedBuild"`
			LastBuild struct {
				Duration int64 `json:"duration"`
			} `json:"lastBuild"`
		} `json:"jobs"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "JSON 解析失败"})
		return
	}

	// 转换数据格式，匹配前端需求
	var nodeViews []models.NodeView
	for _, job := range data.Jobs {
		nodeViews = append(nodeViews, models.NodeView{
			ID:           job.Name,
			NodeID:       "1",  // 示例 Node ID
			Weather:      "未知", // Jenkins API 没有此字段，可根据 color 字段自行扩展
			Name:         job.Name,
			LastSuccess:  formatTimestamp(job.LastSuccessfulBuild.Timestamp),
			LastFailure:  formatTimestamp(job.LastFailedBuild.Timestamp),
			LastDuration: formatDuration(job.LastBuild.Duration),
			CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": nodeViews})
}

// 时间戳转换为字符串
func formatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return "N/A"
	}
	return time.Unix(timestamp/1000, 0).Format("2006-01-02 15:04:05")
}

// 构建时长转换为 "00:10:00" 格式
func formatDuration(duration int64) string {
	if duration == 0 {
		return "N/A"
	}
	minutes := duration / 60000
	seconds := (duration % 60000) / 1000
	return fmt.Sprintf("%02d:%02d:00", minutes, seconds)
}

//func GetNodeViews(c *gin.Context) {
//	nodeID := c.Param("node_id")
//	name := c.Query("name")
//
//	mockData := []models.NodeView{
//		{
//			ID:           "1",
//			NodeID:       nodeID,
//			Weather:      "晴天",
//			Name:         "视图1",
//			LastSuccess:  "2025-03-15 10:00:00",
//			LastFailure:  "2025-03-14 09:00:00",
//			LastDuration: "00:10:00",
//			CreateTime:   "2025-03-10 12:00:00",
//		},
//		{
//			ID:           "2",
//			NodeID:       nodeID,
//			Weather:      "阴天",
//			Name:         "视图2",
//			LastSuccess:  "2025-03-16 11:00:00",
//			LastFailure:  "2025-03-15 10:30:00",
//			LastDuration: "00:15:00",
//			CreateTime:   "2025-03-11 12:00:00",
//		},
//	}
//
//	// 模拟带搜索条件的过滤
//	var filteredData []models.NodeView
//	for _, v := range mockData {
//		if name == "" || v.Name == name {
//			filteredData = append(filteredData, v)
//		}
//	}
//
//	c.JSON(http.StatusOK, gin.H{"success": true, "data": filteredData})
//}

// 添加节点视图 (Mock 数据)
func AddNodeView(c *gin.Context) {
	var view models.NodeView
	if err := c.ShouldBindJSON(&view); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 模拟成功返回
	view.ID = "3"
	view.CreateTime = "2025-03-18 15:00:00"

	c.JSON(http.StatusOK, gin.H{
		"message": "添加成功",
		"success": true,
		"data":    view,
	})
}

// 更新节点视图 (Mock 数据)
func UpdateNodeView(c *gin.Context) {
	var view models.NodeView
	if err := c.ShouldBindJSON(&view); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 模拟更新数据
	view.Name = "已更新的视图"

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"success": true,
		"data":    view,
	})
}

// 删除节点视图 (Mock 数据)
func DeleteNodeView(c *gin.Context) {
	nodeID := c.Param("node_id")
	viewID := c.Param("view_id")

	// 模拟删除成功
	c.JSON(http.StatusOK, gin.H{
		"message":    "删除成功",
		"success":    true,
		"deleted_id": viewID,
		"node_id":    nodeID,
	})
}
