package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/natefinch/lumberjack.v2"
)

// newWriter : 根据指定输入，生成对应的 io Writer
func newWriter(name string, cfg map[string]string) (io.Writer, error) {
	switch name {
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	case "file":
		return newRotateFileWriter(cfg)
	}

	return nil, fmt.Errorf("[%s] writer not supported", name)
}

// newRotateFileWriter : 生成支持日志文件滚动的 io Writer
func newRotateFileWriter(cfg map[string]string) (w io.Writer, err error) {
	// 文件路径 & 文件名配置
	filename, ok := cfg["filename"]
	if !ok || filename == "" {
		return nil, errors.New("the writer config must provide the non-empty filename setting")
	}
	// 检查目录是否存在
	dir := filepath.Dir(filename)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("the writer config - filename wrong, dir(%s) not found", dir)
	}

	// 单个文件容量大小，单位 MB，默认 100 MB
	maxSize := 100
	if cfg["maxsize"] != "" {
		if maxSize, err = strconv.Atoi(cfg["maxsize"]); err != nil {
			return nil, fmt.Errorf("the writer config - maxsize(%s) wrong, must be an integer", cfg["maxsize"])
		}
	}

	// 备份数量，默认 5 份
	maxBackups := 5
	if cfg["maxbackups"] != "" {
		if maxBackups, err = strconv.Atoi(cfg["maxbackups"]); err != nil {
			return nil, fmt.Errorf("the writer config - maxbackups(%s) wrong, must be an integer", cfg["maxbackups"])
		}
	}

	// 备份文件最长留存时间, 单位 天，默认为 0 ，即不删除
	maxAge := 0
	if cfg["maxage"] != "" {
		if maxAge, err = strconv.Atoi(cfg["maxage"]); err != nil {
			return nil, fmt.Errorf("the writer config - maxage(%s) wrong, must be an integer", cfg["maxage"])
		}
	}

	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		LocalTime:  true,
		Compress:   true,
	}, nil
}
