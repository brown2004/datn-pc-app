package store

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Identity struct {
	PCAgentID string    `json:"pc_agent_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Credentials struct {
	PCAgentID   string    `json:"pc_agent_id"`
	AgentSecret string    `json:"agent_secret"`
	SavedAt     time.Time `json:"saved_at"`
}

type Store struct {
	dir string
}

func New() (*Store, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(base, "datn-pc-app")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}

	return &Store{dir: dir}, nil
}

func (s *Store) LoadOrCreatePCAgentID() (string, error) {
	var identity Identity
	if err := readJSON(s.identityPath(), &identity); err == nil {
		pcAgentID := strings.TrimSpace(identity.PCAgentID)
		if pcAgentID != "" {
			return pcAgentID, nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	pcAgentID, err := generateUUID()
	if err != nil {
		return "", err
	}

	identity = Identity{
		PCAgentID: pcAgentID,
		CreatedAt: time.Now().UTC(),
	}
	if err := writeJSON(s.identityPath(), identity); err != nil {
		return "", err
	}

	return pcAgentID, nil
}

func (s *Store) LoadCredentials() (*Credentials, error) {
	var credentials Credentials
	if err := readJSON(s.credentialsPath(), &credentials); err != nil {
		return nil, err
	}
	if credentials.PCAgentID == "" || credentials.AgentSecret == "" {
		return nil, os.ErrNotExist
	}

	return &credentials, nil
}

func (s *Store) SaveCredentials(credentials Credentials) error {
	credentials.SavedAt = time.Now().UTC()
	return writeJSON(s.credentialsPath(), credentials)
}

func (s *Store) ClearCredentials() error {
	return removeIfExists(s.credentialsPath())
}

func (s *Store) credentialsPath() string {
	return filepath.Join(s.dir, "credentials.json")
}

func (s *Store) identityPath() string {
	return filepath.Join(s.dir, "identity.json")
}

func generateUUID() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}

	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16],
	), nil
}

func readJSON(path string, target any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(target)
}

func writeJSON(path string, value any) error {
	tempPath := path + ".tmp"
	file, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}

	encodeErr := json.NewEncoder(file).Encode(value)
	closeErr := file.Close()
	if encodeErr != nil {
		_ = os.Remove(tempPath)
		return encodeErr
	}
	if closeErr != nil {
		_ = os.Remove(tempPath)
		return closeErr
	}

	return os.Rename(tempPath, path)
}

func removeIfExists(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}
