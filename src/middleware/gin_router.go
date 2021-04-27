package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lingye-gin/src/rest"
	"strings"
)

type GinRouter struct {
}

func (v GinRouter) Init(r *gin.Engine) {
	// 路由组映射关系
	groupMap := make(map[string]*gin.RouterGroup)

	for _, currentRa := range rest.Urls {
		// 判断是否在某一组下
		if strings.Compare(currentRa.GroupName, "") == 0 {
			// get
			if strings.Compare(currentRa.Mode, "get") == 0 {
				r.GET(currentRa.RelativePath, currentRa.HandleFunction)
				continue
			}
			// post
			if strings.Compare(currentRa.Mode, "post") == 0 {
				r.POST(currentRa.RelativePath, currentRa.HandleFunction)
				continue
			}
			// delete
			if strings.Compare(currentRa.Mode, "delete") == 0 {
				r.DELETE(currentRa.RelativePath, currentRa.HandleFunction)
				continue
			}
			// put
			if strings.Compare(currentRa.Mode, "put") == 0 {
				r.PUT(currentRa.RelativePath, currentRa.HandleFunction)
			}
		} else {
			// 分组名称
			groupName := fmt.Sprintf("/%s", currentRa.GroupName)
			// 分组不存在
			if groupMap[currentRa.GroupName] == nil {
				if currentRa.GroupHandleFunction == nil {
					// 不存在分组过滤器
					groupMap[currentRa.GroupName] = r.Group(groupName)
				} else {
					groupMap[currentRa.GroupName] = r.Group(groupName, currentRa.GroupHandleFunction)
				}
			}
			// get
			if strings.Compare(currentRa.Mode, "get") == 0 {
				groupMap[currentRa.GroupName].GET(currentRa.RelativePath, currentRa.HandleFunction)
				continue
			}
			// post
			if strings.Compare(currentRa.Mode, "post") == 0 {
				groupMap[currentRa.GroupName].POST(currentRa.RelativePath, currentRa.HandleFunction)
				continue
			}
			// delete
			if strings.Compare(currentRa.Mode, "delete") == 0 {
				groupMap[currentRa.GroupName].DELETE(currentRa.RelativePath, currentRa.HandleFunction)
				continue
			}
			// put
			if strings.Compare(currentRa.Mode, "put") == 0 {
				groupMap[currentRa.GroupName].PUT(currentRa.RelativePath, currentRa.HandleFunction)
			}
		}
	}
}
