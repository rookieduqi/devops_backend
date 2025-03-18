package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// AddServerNode 新增节点
func AddServerNode(c *gin.Context) {
	node := new(models.ServerNode)
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	err := logic.AddNode(node)
	if err != nil {
		ResponseError(c, CodeInvalidNode)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "节点添加成功", "success": true, "data": ""})
}

// GetNameServerNodes 获取所有节点
func GetNameServerNodes(c *gin.Context) {
	name := c.Query("name") // 获取查询参数
	nodes, err := logic.GetServerNodes(name)
	if err != nil {
		ResponseError(c, CodeInvalidGetNode)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "查询成功", "success": true, "data": nodes})
}

// GetServerNodes 获取所有节点
func GetServerNodes(c *gin.Context) {
	nodes, err := logic.GetAllNodes()
	if err != nil {
		ResponseError(c, CodeInvalidGetNode)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "查询成功", "success": true, "data": nodes})
}

// UpdateServerNode 更新节点
func UpdateServerNode(c *gin.Context) {
	var updatedNode models.ServerNode
	if err := c.ShouldBindJSON(&updatedNode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	id := updatedNode.ID
	err := logic.UpdateNode(id, updatedNode)
	if err != nil {
		ResponseError(c, CodeInvalidUpdateNode)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "节点更新成功", "success": true})
}

// DeleteServerNode 删除节点
func DeleteServerNode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID无效"})
		return
	}
	err = logic.DeleteNode(id)
	if err != nil {
		ResponseError(c, CodeInvalidDeleteNode)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "节点删除成功", "success": true})
}
