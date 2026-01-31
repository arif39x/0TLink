package auth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func JoinMesh(apiURL, token, nodeName string) error {
	if apiURL[:8] != "https://" {
		return errors.New("provisioning must use HTTPS")
	}

	privPEM, csrDER, err := GenerateClientIdentity(nodeName)
	if err != nil {
		return fmt.Errorf("failed to generate identity: %w", err)
	}

	req := ProvisionRequest{
		Version:        1,
		BootstrapToken: token,
		CSR:            base64.StdEncoding.EncodeToString(csrDER),
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiURL+"/provision",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("server unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("provisioning failed: %s", resp.Status)
	}

	var provResp ProvisionResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&provResp); err != nil {
		return fmt.Errorf("invalid server response: %w", err)
	}

	configDir := filepath.Join(os.Getenv("HOME"), ".local/share/sidecar-net")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	if err := writeAtomic(filepath.Join(configDir, "node.key"), privPEM, 0600); err != nil {
		return err
	}
	if err := writeAtomic(filepath.Join(configDir, "node.crt"), []byte(provResp.Certificate), 0644); err != nil {
		return err
	}
	if err := writeAtomic(filepath.Join(configDir, "ca.crt"), []byte(provResp.CACert), 0644); err != nil {
		return err
	}

	exp := time.Unix(provResp.ExpiresAt, 0)
	fmt.Printf("Successfully joined mesh. Certificate expires at %s\n", exp.Format(time.RFC3339))

	return nil
}

func writeAtomic(path string, data []byte, perm os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
