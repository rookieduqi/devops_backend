package models

// 定义请求的参数结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ServerNode struct {
	ID         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name" binding:"required"`
	Host       string `db:"host" json:"host" binding:"required"`
	Port       string `db:"port" json:"port"`
	Account    string `db:"account" json:"account" binding:"required"`
	Password   string `db:"password" json:"password" binding:"required"`
	Status     bool   `db:"status" json:"status"`
	Remark     string `db:"remark" json:"remark"`
	CreateTime string `db:"create_time" json:"create_time"`
	UpdateTime string `db:"update_time" json:"update_time"`
}

type NodeView struct {
	ID           string `db:"id" json:"id"`
	NodeID       string `db:"node_id" json:"node_id"`
	Weather      string `db:"weather" json:"weather"`
	Name         string `db:"name" json:"name"`
	LastSuccess  string `db:"last_success" json:"lastSuccess,omitempty"`
	LastFailure  string `db:"last_failure" json:"lastFailure,omitempty"`
	LastDuration string `db:"last_duration" json:"lastDuration,omitempty"`
	CreateTime   string `db:"create_time" json:"createTime,omitempty"`
	Type         string `db:"type" json:"type"`
}

type RequestData struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type NodeViewT struct {
	ID           string `json:"id"`
	NodeID       string `json:"node_id"`
	Name         string `json:"name"`
	Type         string `json:"type"` // 新增字段，区分文件夹 or Job
	Weather      string `json:"weather"`
	LastSuccess  string `json:"last_success"`
	LastFailure  string `json:"last_failure"`
	LastDuration string `json:"last_duration"`
	CreateTime   string `json:"create_time"`
}

type RequestJobData struct {
	NodeID   string `form:"nodeId" binding:"required"`
	ViewID   string `form:"viewId" binding:"required"`
	Host     string `form:"host" binding:"required"`
	Port     string `form:"port" binding:"required"`
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
	JobName  string `form:"jobname"`
}

// Jenkins Job 数据结构
type JenkinsJob struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Color string `json:"color"`
}

// Jenkins API 响应数据结构
type JenkinsResponse struct {
	Jobs []JenkinsJob `json:"jobs"`
}

type StartJobRequest struct {
	ViewID   string `json:"viewId" binding:"required"`
	ViewName string `json:"viewName"`
	NodeId   string `json:"nodeId"`
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}
