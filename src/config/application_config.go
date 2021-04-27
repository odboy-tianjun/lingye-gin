package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// 配置文件
type AppConfig struct {
	// 环境
	Env string `yaml:"env"`
	// 运行模式
	Mode string `yaml:"mode"`
	// 应用名称
	Name string `yaml:"name"`
	// 加密盐
	Secret string `yaml:"secret"`
}
type ServerConfig struct {
	// 服务端口
	Port int `yaml:"port"`
	// API相关配置
	Api ServerApiConfig `yaml:"api"`
}
type ServerApiConfig struct {
	// 过期时间
	Expiry string `yaml:"expiry"`
}
type LogConfig struct {
	// 日志文件配置
	File LogFileConfig `yaml:"file"`
	// 日志级别(debug, info, warn, error)
	Level string `yaml:"level"`
}
type LogFileConfig struct {
	// 日志路径
	Path string `yaml:"path"`
	// 日志名称
	Name string `yaml:"name"`
}

type RedisConfig struct {
	// 连接地址
	Addr string `yaml:"addr"`
	// 密码
	Passwd string `yaml:"passwd"`
	// 库索引
	Database int `yaml:"database"`
}

type ApplicationProperties struct {
	// 应用配置
	App AppConfig `yaml:"app"`
	// Gin服务配置
	Server ServerConfig `yaml:"server"`
	// 日志配置
	Log LogConfig `yaml:"log"`
	// redis配置
	Redis RedisConfig `yaml:"redis"`
}

// 根据路径读取yaml文件
func readYaml(path string) ApplicationProperties {
	var result ApplicationProperties
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("File reading error, application.yml not exist!")
	}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	return result
}

// 根据环境配置取配置明细
func (v *ApplicationProperties) Init() {
	result := readYaml(fmt.Sprintf("%s/application.yml", GetCurrentPath()))
	// 判断环境
	if strings.Compare(result.App.Env, "dev") == 0 {
		AppProps = readYaml(fmt.Sprintf("%s/application-%s.yml", GetCurrentPath(), "dev"))
	}
}

// 获取当前路径，比如：d:/abc
func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
