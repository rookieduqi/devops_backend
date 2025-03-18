package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"fmt"
)

func AddNode(n *models.ServerNode) (err error) {
	if err := mysql.AddNode(n); err != nil {
		return err
	}
	return nil
}

// GetServerNodes 获取节点 (支持 name 筛选)
func GetServerNodes(name string) (nodes []models.ServerNode, err error) {
	nodes, err = mysql.GetNodesByName(name)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func GetAllNodes() ([]models.ServerNode, error) {
	nodes, err := mysql.GetAllNodes()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if len(nodes) == 0 {
		fmt.Println("len(nodes) == 0")
		return []models.ServerNode{}, nil
	}

	return nodes, nil
}

func UpdateNode(id int, updatedNode models.ServerNode) error {

	err := mysql.UpdateNode(id, &updatedNode)
	if err != nil {
		return err
	}

	return nil
}

func DeleteNode(id int) error {

	err := mysql.DeleteNode(id)
	if err != nil {
		return err
	}
	return nil
}
