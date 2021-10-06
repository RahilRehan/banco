package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/RahilRehan/banco/db/util"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testQueries *Queries
var testDB *sql.DB
var dataSource string

func TestMain(m *testing.M) {
	var err error

	cfg, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalln(err)
	}

	dataSource, err = CreateTestDBContainer(*cfg)
	if err != nil {
		log.Fatalln("Cannot create test postgres container")
	}

	testDB, err = sql.Open(cfg.DRIVER_NAME, dataSource)
	if err != nil {
		log.Fatalf("Cannot connect to db %v", err)
	}

	_, err = exec.Command("migrate", "-database", dataSource, "-path", "migrations", "up").Output()
	if err != nil {
		log.Fatalf("Cannot run migrations %v", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}

func CreateTestDBContainer(cfg util.Config) (string, error) {

	timeout, err := strconv.Atoi(cfg.TIMEOUT)
	if err != nil {
		log.Fatalln(err)
	}

	var env = map[string]string{
		"POSTGRES_PASSWORD": viper.GetString("DB_PASSWORD"),
		"POSTGRES_USER":     cfg.DB_USER,
		"POSTGRES_DB":       cfg.DB_NAME,
	}
	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DB_USER,
			viper.GetString("DB_PASSWORD"),
			cfg.DB_HOST,
			cfg.DB_PORT,
			cfg.DB_NAME,
			cfg.SSL_MODE,
		)
	}
	natPort := nat.Port(cfg.DB_PORT)
	ctx := context.Background()

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{cfg.DB_PORT},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          env,
			WaitingFor:   wait.ForSQL(natPort, "postgres", dbURL).Timeout(time.Second * time.Duration(timeout)),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to start container: %s", err)
	}

	mappedPort, err := container.MappedPort(ctx, natPort)
	if err != nil {
		return "", fmt.Errorf("failed to get container external port: %s", err)
	}

	log.Println("postgres container ready and running at port: ", mappedPort)

	return dbURL(nat.Port(mappedPort)), nil
}
