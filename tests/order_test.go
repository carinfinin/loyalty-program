package tests

import (
	"encoding/json"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
)

func TestOrder(t *testing.T) {

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := &http.Client{
		Jar: jar,
	}

	//create order err
	buffer := strings.NewReader("0505")
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/orders", buffer)
	assert.NoError(t, err)
	response, err := client.Do(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)

	user := models.User{
		Login:    "tes02",
		Password: "33302",
	}

	userJSON, err := json.Marshal(&user)
	assert.NoError(t, err)
	buffer = strings.NewReader(string(userJSON))

	//register
	request, err = http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", buffer)
	assert.NoError(t, err)
	request.Header.Add("Content-Type", "application/json")
	response, err = client.Do(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	response.Body.Close()

	//register
	buffer = strings.NewReader(string(userJSON))
	request, err = http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", buffer)
	assert.NoError(t, err)
	request.Header.Add("Content-Type", "application/json")
	response, err = client.Do(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, response.StatusCode)
	response.Body.Close()

	//login
	buffer = strings.NewReader(string(userJSON))
	request, err = http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/", buffer)
	assert.NoError(t, err)
	request.Header.Add("Content-Type", "application/json")
	response, err = client.Do(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	response.Body.Close()

	//create order
	buffer = strings.NewReader("0604")
	request, err = http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/orders", buffer)
	assert.NoError(t, err)
	response, err = client.Do(request)
	assert.NoError(t, err)
	response.Body.Close()
	assert.Equal(t, http.StatusAccepted, response.StatusCode)

	// get orders
	request, err = http.NewRequest(http.MethodGet, "http://localhost:8080/api/user/orders", nil)
	assert.NoError(t, err)
	response, err = client.Do(request)
	assert.NoError(t, err)
	orders := make([]models.Order, 0)
	err = json.NewDecoder(response.Body).Decode(&orders)
	assert.NoError(t, err)
	assert.NotEmpty(t, orders)
	response.Body.Close()
}
