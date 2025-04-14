package app

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest/dto"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testServerAddr = "http://localhost:8070"
)

func TestIntegrationCreatePVZAndReception(t *testing.T) {
	t.Chdir("../..")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresContainer, dbURL, err := setupPostgresContainer(ctx, t)
	require.NoError(t, err, "failed to setup postgres container")
	defer func(postgresContainer testcontainers.Container, ctx context.Context, opts ...testcontainers.TerminateOption) {
		err := postgresContainer.Terminate(ctx, opts...)
		if err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	}(postgresContainer, ctx)

	setEnvironmentVariable(t, "DB_ADDRESS", dbURL)
	setEnvironmentVariable(t, "LOG_LEVEL", "DEBUG")
	setEnvironmentVariable(t, "TOKEN_SECRET", "test-secret")
	setEnvironmentVariable(t, "TOKEN_TTL", "24h")
	setEnvironmentVariable(t, "HTTP_ADDRESS", ":8070")
	setEnvironmentVariable(t, "HTTP_TIMEOUT", "5s")

	go func() {
		if err := Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("server failed: %v", err)
		}
	}()

	time.Sleep(5 * time.Second)

	moderatorToken, err := getDummyToken(t, "moderator")
	require.NoError(t, err, "failed to get moderator token")

	pvzID, err := createPVZ(t, moderatorToken)
	require.NoError(t, err, "failed to create PVZ")

	employeeToken, err := getDummyToken(t, "employee")
	require.NoError(t, err, "failed to get employee token")

	receptionID, err := createReception(t, employeeToken, pvzID)
	require.NoError(t, err, "failed to create reception")

	for i := 0; i < 50; i++ {
		err := createProduct(t, employeeToken, pvzID, receptionID)
		require.NoError(t, err, fmt.Sprintf("failed to create product %d", i))
	}

	err = closeReception(t, employeeToken, pvzID)
	require.NoError(t, err, "failed to close reception")

}

func setupPostgresContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string, error) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image: "postgres:17.4",
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "test_db",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		).WithDeadline(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get container port: %v", err)
	}

	dbURL := fmt.Sprintf("postgres://postgres:postgres@%s:%s/test_db?sslmode=disable", host, port.Port())

	return container, dbURL, nil
}

func getDummyToken(t *testing.T, role string) (string, error) {
	url := fmt.Sprintf("%s/dummyLogin", testServerAddr)
	payload := map[string]string{"role": role}
	jsonData, err := json.Marshal(payload)
	require.NoError(t, err, "failed to marshal login payload")

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err, "failed to send dummy login request")
	defer closeOrLog(t, resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK for dummy login")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")

	var token string
	err = json.Unmarshal(body, &token)
	require.NoError(t, err, "failed to unmarshal token")

	return token, nil
}

func createPVZ(t *testing.T, token string) (uuid.UUID, error) {
	url := fmt.Sprintf("%s/pvz", testServerAddr)
	id := uuid.New()
	ti := time.Now()
	pvz := dto.PVZ{
		Id:               &id,
		City:             dto.Москва,
		RegistrationDate: &ti,
	}
	jsonData, err := json.Marshal(pvz)
	require.NoError(t, err, "failed to marshal PVZ")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	require.NoError(t, err, "failed to create PVZ request")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "failed to send PVZ request")
	defer closeOrLog(t, resp.Body)

	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created for PVZ creation")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read PVZ response")

	var createdPVZ dto.PVZ
	err = json.Unmarshal(body, &createdPVZ)
	require.NoError(t, err, "failed to unmarshal created PVZ")
	require.NotNil(t, createdPVZ.Id, "failed to unmarshal created PVZ id")

	return *createdPVZ.Id, nil
}

func createReception(t *testing.T, token string, pvzID uuid.UUID) (uuid.UUID, error) {
	url := fmt.Sprintf("%s/receptions", testServerAddr)
	payload := map[string]string{"pvzId": pvzID.String()}
	jsonData, err := json.Marshal(payload)
	require.NoError(t, err, "failed to marshal reception payload")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	require.NoError(t, err, "failed to create reception request")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "failed to send reception request")
	defer closeOrLog(t, resp.Body)

	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created for reception creation")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read reception response")

	var createdReception models.Reception
	err = json.Unmarshal(body, &createdReception)
	require.NoError(t, err, "failed to unmarshal created reception")

	return createdReception.ID, nil
}

func createProduct(t *testing.T, token string, pvzID, receptionID uuid.UUID) error {
	url := fmt.Sprintf("%s/products", testServerAddr)
	payload := map[string]interface{}{
		"type":        "электроника",
		"pvzId":       pvzID.String(),
		"receptionId": receptionID.String(),
	}
	jsonData, err := json.Marshal(payload)
	require.NoError(t, err, "failed to marshal product payload")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	require.NoError(t, err, "failed to create product request")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "failed to send product request")
	defer closeOrLog(t, resp.Body)

	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created for product creation")

	return nil
}

func closeReception(t *testing.T, token string, pvzID uuid.UUID) error {
	url := fmt.Sprintf("%s/pvz/%s/close_last_reception", testServerAddr, pvzID.String())
	req, err := http.NewRequest("POST", url, nil)
	require.NoError(t, err, "failed to create close reception request")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "failed to send close reception request")
	defer closeOrLog(t, resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK for closing reception")

	return nil
}

func setEnvironmentVariable(t *testing.T, key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		t.Errorf("failed to set environment variable: %v", err)
		return
	}
}

func closeOrLog(t *testing.T, c io.Closer) {
	err := c.Close()
	if err != nil {
		t.Logf("failed to close connection: %v", err)
	}
}
