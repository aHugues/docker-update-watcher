package remotedocker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RemoteImage struct {
	Architecture string `json:"architecture"`
	Digest       string `json:"digest"`
	OS           string `json:"os"`
	LastPushed   string `json:"last_pushed"`
}

type RemoteTag struct {
	LastUpdated string        `json:"last_updated"`
	Name        string        `json:"name"`
	Images      []RemoteImage `json:"images"`
}

type RemoteRes struct {
	Results []RemoteTag `json:"results"`
}

func GetRemote(ctx context.Context, clt *http.Client, namespace, image string) ([]RemoteTag, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://hub.docker.com/v2/repositories/"+namespace+"/"+image+"/tags/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := clt.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected returncode %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	fmt.Print(string(content))

	var res RemoteRes
	if err := json.Unmarshal(content, &res); err != nil {
		return nil, fmt.Errorf("failed to parse body: %w", err)
	}

	return res.Results, nil
}
