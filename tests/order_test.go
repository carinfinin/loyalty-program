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
		Login:    "tes",
		Password: "333",
	}

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	client := &http.Client{
		Jar: jar,
	}

	userJSON, err := json.Marshal(user)
	assert.NoError(t, err)
	buffer := strings.NewReader(string(userJSON))

	//register
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/user/register", buffer)
	assert.NoError(t, err)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	assert.NoError(t, err)
	fmt.Println(response.StatusCode)
	response.Body.Close()

	//create order
	buffer = strings.NewReader("0109")
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
