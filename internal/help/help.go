package help

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func Method(method string, handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, method)
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}

// RecaptchaResponse represents the response from Google's reCAPTCHA API.
type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`       // Optional: for v3
	Action      string   `json:"action"`      // Optional: for v3
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

// VerifyRecaptcha verifies the reCAPTCHA response token with Google's servers.
// Returns true if the validation is successful.
func VerifyRecaptcha(secret, response string) (bool, error) {
	if secret == "" || response == "" {
		return false, fmt.Errorf("secret and response cannot be empty")
	}

	reqURL := "https://www.google.com/recaptcha/api/siteverify"

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.PostForm(reqURL, url.Values{
		"secret":   {secret},
		"response": {response},
	})
	if err != nil {
		return false, fmt.Errorf("failed to verify recaptcha: %w", err)
	}
	defer resp.Body.Close()

	var result RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode recaptcha response: %w", err)
	}

	// For reCAPTCHA v3, you might also want to check result.Score (e.g., > 0.5)
	// For now, we return Success which covers both v2 and v3 basic checks.
	return result.Success, nil
}
