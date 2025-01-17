package version

import (
	"encoding/json"
	"fmt"
	"runtime"

	// CLI 中以表格形式展示数据
	"github.com/gosuri/uitable"
)

var (
	// 语义化版本号
	GitVersion = "v0.0.0-master+$Format:%h$"
	// ISO8601 格式时间
	BuildDate = "1970-01-01T00:00:00Z"
	// commit 的 SHA1 值 = git rev-parse HEAD
	GitCommit = "$Format:%H$"
	// git 仓库状态 clean or diry
	GitTreeState = ""
)

type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

func (info Info) String() string {
	if s, err := info.Text(); err == nil {
		return string(s)
	}

	return info.GitVersion
}

// json 格式信息
func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)

	return string(s)
}

// 表格化信息
func (info Info) Text() ([]byte, error) {
	table := uitable.New()

	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "

	table.AddRow("gitVersion:", info.GitVersion)
	table.AddRow("gitCommit:", info.GitCommit)
	table.AddRow("gitTreeState:", info.GitTreeState)
	table.AddRow("buildDate:", info.BuildDate)
	table.AddRow("goVersion:", info.GitVersion)
	table.AddRow("compiler:", info.Compiler)
	table.AddRow("platform:", info.Platform)

	return table.Bytes(), nil
}

// 结构化信息
func Get() Info {
	return Info{
		GitVersion: GitVersion,
		GitCommit: GitCommit,
		GitTreeState: GitTreeState,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Compiler: runtime.Compiler,
		Platform: fmt.Sprintf("%s%s", runtime.GOOS, runtime.GOARCH),
	}
}
