package mysql

import (
	"bluebell/models"
	"fmt"
	"time"
)

func AddNode(node *models.ServerNode) (err error) {
	// 设置添加时间
	node.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	query := `
    INSERT INTO server_nodes (name, host, port, account, password, status, remark, create_time)
    VALUES (:name, :host, :port, :account, :password, :status, :remark, :create_time)
    `

	_, err = db.NamedExec(query, node)
	if err != nil {
		fmt.Println("mysql.Node", err)
		return err
	}

	return nil
}

// GetNodeByID 获取单个节点
func GetNodeByID(id int) (*models.ServerNode, error) {
	var node models.ServerNode
	query := `SELECT * FROM server_nodes WHERE id = ?`
	err := db.Get(&node, query, id)
	if err != nil {
		fmt.Println("mysql.GetNodeByID", err)
		return nil, err
	}
	return &node, nil
}

// GetNodesByName 根据 name 筛选节点
func GetNodesByName(name string) ([]models.ServerNode, error) {
	var nodes []models.ServerNode

	// 使用 TRIM() 去除可能的空格
	query := `SELECT * FROM server_nodes WHERE TRIM(name) LIKE ?`

	// 模糊匹配 (支持 '%关键字%' 格式)
	err := db.Select(&nodes, query, "%"+name+"%")
	if err != nil {
		fmt.Println("mysql.GetNodesByName 错误:", err)
		return nil, err
	}
	return nodes, nil
}

// GetAllNodes 获取所有节点
func GetAllNodes() ([]models.ServerNode, error) {
	var nodes []models.ServerNode
	query := `SELECT * FROM server_nodes`
	err := db.Select(&nodes, query)
	if err != nil {
		fmt.Println("mysql.GetAllNodes", err)
		return nil, err
	}
	return nodes, nil
}

// UpdateNode 更新节点
func UpdateNode(id int, node *models.ServerNode) error {
	query := `
    UPDATE server_nodes
    SET name = :name, host = :host, port = :port,
        account = :account, password = :password,
        status = :status, remark = :remark
    WHERE id = :id
    `

	node.ID = id
	_, err := db.NamedExec(query, node)
	if err != nil {
		fmt.Println("mysql.UpdateNode", err)
		return err
	}
	return nil
}

// DeleteNode 删除节点
func DeleteNode(id int) error {
	query := `DELETE FROM server_nodes WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		fmt.Println("mysql.DeleteNode", err)
		return err
	}

	return nil
}
