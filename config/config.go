package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

const (
	filePath string = "./config.yaml"
)

var conf *config

type Config interface {
	HttpDBType() (string, error)
	HttpDBDataSource() (string, error)
	GrpcHostAndPort() (string, int, error)
	GrpcDBType() (string, error)
	GrpcDBDataSource() (string, error)
}

type config struct {
	Http *httpServer `yaml:"httpserver"`
	Grpc *grpcServer `yaml:"grpcserver"`
}

type httpServer struct {
	DB *db `yaml:"db"`
}

type grpcServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	DB   *db    `yaml:"db"`
}

type db struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (c *config) HttpDBType() (string, error) {
	if c.Http == nil || c.Http.DB == nil {
		return "", errors.New("invaid config")
	}

	return c.Http.DB.Type, nil
}

func (c *config) HttpDBDataSource() (string, error) {
	if c.Http == nil || c.Http.DB == nil {
		return "", errors.New("invaid config")
	}

	db := c.Http.DB

	return getDBDataSource(db.User, db.Password, db.Host, db.Port, db.Database), nil
}

func (c *config) GrpcHostAndPort() (string, int, error) {
	if c.Grpc == nil {
		return "", 0, errors.New("invaid config")
	}

	return c.Grpc.Host, c.Grpc.Port, nil
}

func (c *config) GrpcDBType() (string, error) {
	if c.Grpc == nil || c.Grpc.DB == nil {
		return "", errors.New("invaid config")
	}

	return c.Grpc.DB.Type, nil
}

func (c *config) GrpcDBDataSource() (string, error) {
	if c.Grpc == nil || c.Grpc.DB == nil {
		return "", errors.New("invaid config")
	}

	db := c.Grpc.DB

	return getDBDataSource(db.User, db.Password, db.Host, db.Port, db.Database), nil
}

func getDBDataSource(user, pass, host string, port int, database string) string {
	if host == "localhost" {
		host = ""
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, host, port, database)
}

func LoadConfig() (err error) {
	conf, err = getConfig()
	return
}

func GetConfig() (Config, error) {
	var err error
	if conf == nil {
		err = LoadConfig()
	}

	return conf, err
}

func getConfig() (*config, error) {
	// 外部からconfの中身を参照できるようにする
	var c config

	filePath, err := getPath(filePath)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func getPath(filePath string) (string, error) {

	if path.IsAbs(filePath) {
		return filePath, nil
	}

	_, cf, _, ok := runtime.Caller(0)
	if !ok {
		return filePath, errors.New("invalid")
	}
	currentDir := filepath.Dir(cf)

	return path.Clean(filepath.Join(currentDir, filePath)), nil
}
