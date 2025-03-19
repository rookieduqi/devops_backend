package controller

import (
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 获取节点视图列表 (Mock 数据)
func GetNodeViews(c *gin.Context) {
	nodeID := c.Param("node_id")
	name := c.Query("name")

	mockData := []models.NodeView{
		{
			ID:           "1",
			NodeID:       nodeID,
			Weather:      "晴天",
			Name:         "视图1",
			LastSuccess:  "2025-03-15 10:00:00",
			LastFailure:  "2025-03-14 09:00:00",
			LastDuration: "00:10:00",
			CreateTime:   "2025-03-10 12:00:00",
		},
		{
			ID:           "2",
			NodeID:       nodeID,
			Weather:      "阴天",
			Name:         "视图2",
			LastSuccess:  "2025-03-16 11:00:00",
			LastFailure:  "2025-03-15 10:30:00",
			LastDuration: "00:15:00",
			CreateTime:   "2025-03-11 12:00:00",
		},
	}

	// 模拟带搜索条件的过滤
	var filteredData []models.NodeView
	for _, v := range mockData {
		if name == "" || v.Name == name {
			filteredData = append(filteredData, v)
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": filteredData})
}

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
