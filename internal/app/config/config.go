package config

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/goccy/go-yaml"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	Port     int    `yaml:"port"`
	SSLMode  bool   `yaml:"sslmode"`
}

type AppConfig struct {
	Database DB     `yaml:"db"`
	Port     string `yaml:"port"`
}

type Configs struct {
	DatabaseConnection *gorm.DB
	AppPort            string
}

func getProjectRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return "", fmt.Errorf("project root directory not found")
}

func LoadConfig(filename string) *Configs {
	rootDir, err := getProjectRootDir()
	if err != nil {
		log.Fatal(err)
	}

	app_config_file, err := os.ReadFile(rootDir + "/env/" + filename)

	if err != nil {
		log.Fatal(err)
	}

	// Read Config
	var app_config AppConfig
	err = yaml.Unmarshal(app_config_file, &app_config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Config DB
	db_conn, err := DBConfig(app_config.Database)
	if err != nil {
		panic(err)
	}

	// Config Server
	// ServerConfig(app_config.Port, mux)

	return &Configs{
		DatabaseConnection: db_conn,
		AppPort:            app_config.Port,
	}
}

func DBConfig(database DB) (*gorm.DB, error) {
	// if any of the database config is empty, return error
	if database.Host == "" || database.User == "" || database.Password == "" || database.DBName == "" || database.Port == 0 {
		return nil, errors.New("database config is not complete")
	}

	db_host := "host=" + database.Host
	db_user := "user=" + database.User
	db_password := "password=" + database.Password
	db_name := "dbname=" + database.DBName
	db_port := "port=" + strconv.Itoa(database.Port)

	var db_sslmode string

	if database.SSLMode {
		db_sslmode = "sslmode=enable"
	} else {
		db_sslmode = "sslmode=disable"
	}

	dsn := fmt.Sprintf("%s %s %s %s %s %s", db_host, db_user, db_password, db_name, db_port, db_sslmode)
	db_conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db_conn, err
}

func ServerConfig(port string, mux *http.ServeMux) {
	log.Println("Starting server at port " + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
