package hashutil

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"
)

const userAgent = "derpack-edit/1.0 (modpack editor)"

// downloadClient has a long timeout for big jars.
var downloadClient = &http.Client{
	Timeout: 5 * time.Minute,
}

// SHA512OfURL fetches the URL and returns the SHA-512 hash of the response
// body as a lowercase hex string. Also returns the byte count.
func SHA512OfURL(url string) (hash string, size int64, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := downloadClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", 0, fmt.Errorf("download %s: %s", url, resp.Status)
	}

	h := sha512.New()
	n, err := io.Copy(h, resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("read body: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}
