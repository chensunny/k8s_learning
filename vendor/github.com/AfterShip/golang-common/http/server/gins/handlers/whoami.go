package handlers

import (
	"net/http"

	"github.com/AfterShip/golang-common/http/server/gins"
	"github.com/AfterShip/golang-common/whoami"
	"github.com/gin-gonic/gin"
)

const (
	whoamiPath = "/whoami"
)

type whoAmI struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Build   Build  `json:"build"`
	Commit  Commit `json:"commit"`
}

type Build struct {
	Number   string `json:"number"`
	Datetime string `json:"datetime"`
}

type Commit struct {
	Hash   string `json:"hash"`
	Branch string `json:"branch"`
}

func RegisterWhoamiHandler(group *gin.RouterGroup) {
	group.GET(whoamiPath, func(c *gin.Context) {
		gins.ResponseSuccess(c, http.StatusOK, whoAmI{
			Service: whoami.Name(),
			Version: whoami.Version(),
			Build: Build{
				Number:   whoami.BuildNumber(),
				Datetime: whoami.BuildAt(),
			},
			Commit: Commit{
				Hash:   whoami.CommitHash(),
				Branch: whoami.CommitBranch(),
			},
		})
	})
}
