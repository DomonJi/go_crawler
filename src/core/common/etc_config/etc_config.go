package etc_config

import (
    "../config"
    "../util"
    "os"
)

func configpath() string {
    wd := os.Getenv("GOPATH")
    if wd == "" {
        panic("GOPATH is not setted in env.")
    }
    logpath := wd + "/etc/"
    filename := "main.conf"
    err := os.MkdirAll(logpath, 0755)
    if err != nil {
        panic("logpath error : " + logpath + "\n")
    }
    return logpath + filename
}

var conf *config.Config
var path string

func StartConf(configFilePath string) *config.Config {
    if configFilePath != "" && !util.IsFileExists(configFilePath) {
        panic("config path is not valiad:" + configFilePath)
    }

    path = configFilePath
    return Conf()
}

func Conf() *config.Config {
    if conf == nil {
        if path == "" {
            path = configpath()
        }
        conf = config.NewConfig().Load(path)
    }
    return conf
}
