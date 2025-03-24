package tests

import (
	"encoding/json"
	"fmt"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
)

func TestOrder(t *testing.T) {
	user := models.User{
		Login:    "test12",
		Password: "333333",
	}

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := &http.Client{
		Jar: jar,
	}

	userJSON, err := json.Marshal(user)
	assert.NoError(t, err)
	buffer := strings.NewReader(string(userJSON))

	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", buffer)
	assert.NoError(t, err)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	assert.NoError(t, err)
	fmt.Println(response.StatusCode)
	response.Body.Close()

	request, err = http.NewRequest(http.MethodGet, "http://localhost:8080/api/user/orders", nil)
	assert.NoError(t, err)
	response, err = client.Do(request)
	assert.NoError(t, err)
	response.Body.Close()

	fmt.Println(response.StatusCode)
}
