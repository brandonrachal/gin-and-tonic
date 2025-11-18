package controllers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brandonrachal/gin-and-tonic/controllers"
	"github.com/brandonrachal/gin-and-tonic/db"
	"github.com/brandonrachal/gin-and-tonic/internal"
	"github.com/brandonrachal/gin-and-tonic/models"
	"github.com/brandonrachal/gin-and-tonic/models/api"
	"github.com/brandonrachal/go-toolbox/jsonutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var (
	router   *gin.Engine
	dbClient *db.Client
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	// Calling os.Exit() kills the program immediately and doesn't call any deferred methods.
	// So the easiest way to fix this is to call another method.
	exitCode := testMain(m)
	os.Exit(exitCode)
}

func testMain(m *testing.M) int {
	ctx := context.Background()
	migrationClient, migrationClientErr := internal.TestDBMigrationClient()
	if migrationClientErr != nil {
		log.Printf("Error: Couldn't get the migration client - %s\n", migrationClientErr)
		return 1
	}
	defer func() {
		if downAllErr := migrationClient.Reset(ctx); downAllErr != nil {
			log.Printf("Error: Couldn't down all migrations - %s\n", downAllErr)
		}
		if closeErr := migrationClient.Close(); closeErr != nil {
			log.Printf("Error: Couldn't close migration client - %s\n", closeErr)
		}
	}()
	if upAllErr := migrationClient.UpAll(ctx); upAllErr != nil {
		log.Printf("Error: Couldn't up all migrations - %s\n", upAllErr)
		return 1
	}
	var dbClientErr error
	dbClient, dbClientErr = internal.TestDBClient()
	if dbClientErr != nil {
		log.Printf("Error: Couldn't retrieve the db client - %s\n", dbClientErr)
		return 1
	}
	defer func() {
		if dbCloseErr := dbClient.Close(); dbCloseErr != nil {
			log.Printf("Error: Couldn't close the db client - %s\n", dbCloseErr)
		}
	}()
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	router = controllers.GetRouter(logger, dbClient)
	exitCode := m.Run()
	return exitCode
}

func TestPing(t *testing.T) {
	r := require.New(t)
	resp := callRequest(r, "GET", "/ping", nil)
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusOK, resp.StatusCode)
	bodyBytes, bodyBytesErr := io.ReadAll(resp.Body)
	r.NoError(bodyBytesErr)
	r.JSONEq(`{"status": "ok"}`, string(bodyBytes))
}

func TestCreateUserAction(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()
	_, deleteResultErr := dbClient.DeleteAllUsers(ctx)
	r.NoError(deleteResultErr)
	// Bad user data
	newUser, newUserErr := GetFirstNewUser()
	r.NoError(newUserErr)
	newUser.FirstName = ""
	newUser.LastName = ""
	newUser.Email = ""
	badResp := callRequest(r, "POST", "/v1.0/user", newUser)
	defer func() {
		closeErr := badResp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusBadRequest, badResp.StatusCode)
	// Create a user
	newUser, newUserErr = GetFirstNewUser()
	r.NoError(newUserErr)
	goodResp := callRequest(r, "POST", "/v1.0/user", newUser)
	defer func() {
		closeErr := goodResp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusOK, goodResp.StatusCode)
	user, userErr := dbClient.GetFirstUser(ctx)
	r.NoError(userErr)
	r.NotNil(user)
	// No duplicate user
	dupResp := callRequest(r, "POST", "/v1.0/user", newUser)
	defer func() {
		closeErr := dupResp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusInternalServerError, dupResp.StatusCode)
}

func TestGetUserAction(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()
	user, userErr := getFirstOrCreateUser(ctx)
	r.NoError(userErr)
	userId := models.GetIdUser(user.Id)
	// Get User
	userResp := callRequest(r, "GET", "/v1.0/user", userId)
	defer func() {
		closeErr := userResp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusOK, userResp.StatusCode)
	// Get User Bad Data
	userIdBadData := "{\"user\": nil}"
	userBadResp := callRequest(r, "GET", "/v1.0/user", userIdBadData)
	defer func() {
		closeErr := userBadResp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusBadRequest, userBadResp.StatusCode)
}

func TestUpdateUserAction(t *testing.T) {
	ctx := context.Background()
	r := require.New(t)
	user, userErr := getFirstOrCreateUser(ctx)
	r.NoError(userErr)
	// Update user good
	newFirstName := "HELLO!"
	user.FirstName = newFirstName
	goodResp := callRequest(r, "PUT", "/v1.0/user", user)
	defer func() {
		closeErr := goodResp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusOK, goodResp.StatusCode)
	user, userErr = dbClient.GetUser(ctx, user.Id)
	r.NoError(userErr)
	r.NotNil(user)
	r.Equal(newFirstName, user.FirstName)
}

func TestDeleteUserAction(t *testing.T) {
	ctx := context.Background()
	r := require.New(t)
	user, userErr := getFirstOrCreateUser(ctx)
	r.NoError(userErr)
	// Delete user
	userId := models.GetIdUser(user.Id)
	goodResp := callRequest(r, "DELETE", "/v1.0/user", userId)
	defer func() {
		closeErr := goodResp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusOK, goodResp.StatusCode)
	user, userErr = dbClient.GetUser(ctx, user.Id)
	r.Error(userErr)
	r.True(errors.Is(userErr, sql.ErrNoRows))
}

func TestGetUsersAction(t *testing.T) {
	ctx := context.Background()
	r := require.New(t)
	createUsers(ctx, r)
	// Get users
	goodResp := callRequest(r, "GET", "/v1.0/users", nil)
	defer func() {
		closeErr := goodResp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusOK, goodResp.StatusCode)
	bodyBytes, bodyBytesErr := io.ReadAll(goodResp.Body)
	r.NoError(bodyBytesErr)
	var apiUsers api.UsersMessage
	jsonErr := json.Unmarshal(bodyBytes, &apiUsers)
	r.NoError(jsonErr)
	r.Equal(3, len(apiUsers.Users))
}

func TestGetUsersWithAgeAction(t *testing.T) {
	ctx := context.Background()
	r := require.New(t)
	createUsers(ctx, r)
	// Get Users with age
	goodResp := callRequest(r, "GET", "/v1.0/users_with_age", nil)
	defer func() {
		closeErr := goodResp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusOK, goodResp.StatusCode)
	bodyBytes, bodyBytesErr := io.ReadAll(goodResp.Body)
	r.NoError(bodyBytesErr)
	var apiUsers api.UsersWithAgeMessage
	jsonErr := json.Unmarshal(bodyBytes, &apiUsers)
	r.NoError(jsonErr)
	r.Equal(3, len(apiUsers.Users))
}

func TestGetAgeStatsAction(t *testing.T) {
	ctx := context.Background()
	r := require.New(t)
	createUsers(ctx, r)
	// Users with age
	goodResp := callRequest(r, "GET", "/v1.0/age_stats", nil)
	defer func() {
		closeErr := goodResp.Body.Close()
		if closeErr != nil {
			log.Printf("Error: Couldn't close the response body - %s\n", closeErr)
		}
	}()
	r.Equal(http.StatusOK, goodResp.StatusCode)
	bodyBytes, bodyBytesErr := io.ReadAll(goodResp.Body)
	r.NoError(bodyBytesErr)
	var ageStats api.AgeStatsMessage
	jsonErr := json.Unmarshal(bodyBytes, &ageStats)
	r.NoError(jsonErr)
}

// Helper methods

func GetFirstNewUser() (*models.CreateUser, error) {
	return models.GetCreateUser("Testy", "McTesterson", "testy.mctesterson@gmail.com", "1996-06-06")
}

func GetSecondNewUser() (*models.CreateUser, error) {
	return models.GetCreateUser("John", "Doe", "john.doe@gmail.com", "2000-12-16")
}

func GetThirdNewUser() (*models.CreateUser, error) {
	return models.GetCreateUser("Jane", "Doe", "jane.doe@gmail.com", "2003-03-16")
}

func getFirstOrCreateUser(ctx context.Context) (*models.User, error) {
	user, userErr := dbClient.GetFirstUser(ctx)
	if userErr != nil && !errors.Is(userErr, sql.ErrNoRows) {
		return nil, userErr
	} else if user != nil {
		return user, nil
	}
	newUser, newUserErr := GetFirstNewUser()
	if newUserErr != nil {
		return nil, newUserErr
	}
	userId, userIdErr := createUser(ctx, newUser)
	if userIdErr != nil {
		return nil, userIdErr
	}
	user, userErr = dbClient.GetUser(ctx, *userId)
	if userErr != nil {
		return nil, userErr
	} else if user != nil {
		return user, nil
	}
	return user, errors.New("user could not found or created")
}

func createUser(ctx context.Context, user *models.CreateUser) (*int64, error) {
	createResult, createResultErr := dbClient.CreateUser(ctx, user.FirstName, user.LastName, user.Email, user.Birthday.ToTime())
	if createResultErr != nil {
		return nil, createResultErr
	}
	userId, userIdErr := createResult.LastInsertId()
	if userIdErr != nil {
		return nil, userIdErr
	}
	return &userId, nil
}

func createUsers(ctx context.Context, r *require.Assertions) {
	_, deleteResultErr := dbClient.DeleteAllUsers(ctx)
	r.NoError(deleteResultErr)

	newUser, newUserErr := GetFirstNewUser()
	r.NoError(newUserErr)
	_, createErr := createUser(ctx, newUser)
	r.NoError(createErr)

	newUser, newUserErr = GetSecondNewUser()
	r.NoError(newUserErr)
	_, createErr = createUser(ctx, newUser)
	r.NoError(createErr)

	newUser, newUserErr = GetThirdNewUser()
	r.NoError(newUserErr)
	_, createErr = createUser(ctx, newUser)
	r.NoError(createErr)
}

func callRequest(r *require.Assertions, method, url string, data any) *http.Response {
	jsonBodyReader, jsonBodyReaderErr := jsonutils.ToJsonReader(data)
	r.NoError(jsonBodyReaderErr)
	w := httptest.NewRecorder()
	req, reqErr := http.NewRequest(method, url, jsonBodyReader)
	r.NoError(reqErr)
	router.ServeHTTP(w, req)
	return w.Result()
}
