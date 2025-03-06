package service

import (
	"bad-code-review/common"
	"bad-code-review/model"
	"bad-code-review/pkg"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type UserServiceImpl struct {
	httpTimeoutSeconds    int
	httpClient            *http.Client
	userDetailsWorkerPoll pkg.UserDetailsWorkerPool
}

func NewUserService(httpTimeoutSeconds, numWorkersInt int) (UserService, error) {
	pool := pkg.NewUserDetailsWorkerPool(numWorkersInt)

	return UserServiceImpl{
		httpTimeoutSeconds: httpTimeoutSeconds,
		httpClient: &http.Client{
			Timeout: time.Duration(httpTimeoutSeconds) * time.Second,
		},
		userDetailsWorkerPoll: pool,
	}, nil
}
func (us UserServiceImpl) FetchUserDetailsBatch(ctx context.Context, orderIDs []string) ([]model.User, error) {
	return us.userDetailsWorkerPoll.ProcessOrdersIds(ctx, orderIDs, us.fetchUserDetails)
}

func (us UserServiceImpl) fetchUserDetails(ctx context.Context, userID string) (model.User, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://user-service:8081/user?id="+userID, nil)
	if err != nil {
		log.Println("error creating request", err)
		return model.User{}, err
	}
	resp, err := us.httpClient.Do(request)
	if err != nil {
		log.Println("error fetching user details", err)
		return model.User{}, common.ErrUserFetchDetailsError
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("error fetching user details", resp.StatusCode)
		return model.User{}, common.ErrUserFetchDetailsError
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println("error closing response body", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body", err)
		return model.User{}, err
	}
	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("error unmarshalling response body", err)
		return model.User{}, err
	}
	return user, nil
}
