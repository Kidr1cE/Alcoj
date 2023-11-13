package api

import "github.com/gin-gonic/gin"

type RunResponse struct {
	StatusMsg     string `json:"status_msg,omitempty"`
	Status        string `json:"status,omitempty"`
	StatusRuntime string `json:"status_runtime,omitempty"`
	Memory        string `json:"memory,omitempty"`
}

type RunRequest struct {
	TypedCode string `json:"typed_code,omitempty"`
}

func Check(c *gin.Context) {

}
