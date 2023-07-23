package toffu

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/term"
)

// The function `getToken` retrieves a Woffu token from a configuration file and returns it, or returns
// an error if no token is found.
func getToken() (token string, err error) {
	token = viper.GetString("woffu_token")
	if token == "" {
		return "", errors.New("no Woffu token found! Run `toffu token` to generate a token")
	}
	return
}

// The function `extractTokenConfig` extracts and decodes a token, then unmarshalls the JSON data into
// a map.
func extractTokenConfig() (map[string]interface{}, error) {
	// Get Token
	if viper.GetString("woffu_token") == "" {
		return nil, errors.New("no Woffu token found! Run `toffu token` to generate a token")
	}
	token := viper.GetString("woffu_token")

	// Decode token
	parts := strings.Split(token, ".")
	jsonStr, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(parts[1])
	if err != nil {
		panic(err)
	}

	// Unmarshall JSON
	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON from token: %v", err)
	}

	return data, nil
}

// The function `getUserId` extracts the user ID from a token and returns it as an integer.
func getUserId() (int, error) {
	// Get Token Configuration
	data, err := extractTokenConfig()
	if err != nil {
		return 0, err
	}

	// Extract user id
	userIdStr, ok := data["UserId"].(string)
	if !ok {
		return 0, errors.New("error extracting User Id from token")
	}

	// Convert to integer
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return 0, fmt.Errorf("error converting User Id to integer: %v", err)
	}

	return userId, nil
}

// The function `getTokenExpiration` extracts the expiration date from a token configuration and
// returns it as a Unix timestamp.
func getTokenExpiration() (int64, error) {
	// Get Token Configuration
	data, err := extractTokenConfig()
	if err != nil {
		return 0, err
	}

	// Extract expiration date
	unixEpoch, ok := data["exp"].(float64)
	if !ok {
		return 0, errors.New("error extracting expiration date from token")
	}

	return time.Unix(int64(unixEpoch), 0).Unix(), nil
}

// The function `getCredentials“ prompts the user for a username and password and returns them as
// strings.
func getCredentials() (username, password string) {
	// Username
	fmt.Print("Username: ")
	fmt.Scanln(&username)
	// Password
	fmt.Print("Password: ")
	if runtime.GOOS == "windows" { // TODO: https://github.com/golang/go/issues/16552
		fmt.Scanln(&password)
		return
	}
	bytePassword, _ := term.ReadPassword(0)
	password = string(bytePassword)
	return
}

// The function `secondsToHumanReadable“ takes a duration in seconds and returns a human-readable
// string representation of the duration in the format "Xh Xm Xs".
func secondsToHumanReadable(duration time.Duration) string {
	// Hours
	hours := duration / time.Hour
	// Minutes
	duration -= hours * time.Hour
	minutes := duration / time.Minute
	// Seconds
	duration -= minutes * time.Minute
	seconds := duration / time.Second

	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}

func runConcurrently(funcs ...func() error) (err error) {
	var wg sync.WaitGroup
	wg.Add(len(funcs))
	errChan := make(chan error, len(funcs))

	// Run goroutines
	for _, f := range funcs {
		go func(fn func() error) {
			defer wg.Done()
			if err := fn(); err != nil {
				errChan <- err
			}
		}(f)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Collect any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
