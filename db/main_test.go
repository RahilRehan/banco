package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	host       = "localhost"
	name       = "test_db"
	user       = "test_user"
	password   = "passtemp"
	timeout    = 5
	driverName = "postgres"
)

var testQueries *Queries
var testDB *sql.DB
var dataSource string

func TestMain(m *testing.M) {
	var err error

	dataSource, err = CreateTestDBContainer()
	if err != nil {
		log.Fatalln("Cannot create test postgres container")
	}

	testDB, err = sql.Open(driverName, dataSource)
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

func CreateTestDBContainer() (string, error) {
	var env = map[string]string{
		"POSTGRES_PASSWORD": password,
		"POSTGRES_USER":     user,
		"POSTGRES_DB":       name,
	}
	var port = "5432/tcp"
	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port.Port(),
			name)
	}
	natPort := nat.Port(port)
	ctx := context.Background()

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{port},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          env,
			WaitingFor:   wait.ForSQL(natPort, "postgres", dbURL).Timeout(time.Second * timeout),
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
