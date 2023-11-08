package cmd

import (
	"fmt"
)

// getVersionFmt 返回一个格式化的版本字符串，包括项目名称、版本号和构建时间。
func getVersionFmt() string {
	return fmt.Sprintf("%s %s %s", ProjectName, Version, BuildTime)
}

func Execute() {
	command := newCommandsBuilder().addAll().builder()
	rootCmd := command.GetCommand()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
