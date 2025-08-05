package parquet

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/segmentio/parquet-go"
)

// SimpleManager provides basic Parquet operations
type SimpleManager struct {
	baseDir string
}

// NewSimpleManager creates a new simple Parquet manager
func NewSimpleManager(baseDir string) *SimpleManager {
	if baseDir == "" {
		baseDir = "data/parquet"
	}
	return &SimpleManager{
		baseDir: baseDir,
	}
}

// ensureDir creates directory if it doesn't exist
func (m *SimpleManager) ensureDir() error {
	return os.MkdirAll(m.baseDir, 0755)
}

// WriteUsers writes user data to Parquet file with default settings
func (m *SimpleManager) WriteUsers(filename string, users []User) error {
	if err := m.ensureDir(); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(m.baseDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := parquet.NewGenericWriter[User](file)
	defer writer.Close()

	_, err = writer.Write(users)
	if err != nil {
		return fmt.Errorf("failed to write users: %w", err)
	}

	return nil
}

// ReadUsers reads user data from Parquet file
func (m *SimpleManager) ReadUsers(filename string) ([]User, error) {
	filePath := filepath.Join(m.baseDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := parquet.NewGenericReader[User](file)
	defer reader.Close()

	users := make([]User, reader.NumRows())
	n, err := reader.Read(users)
	if err != nil {
		return nil, fmt.Errorf("failed to read users: %w", err)
	}

	return users[:n], nil
}

// WriteProducts writes product data to Parquet file
func (m *SimpleManager) WriteProducts(filename string, products []Product) error {
	if err := m.ensureDir(); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(m.baseDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := parquet.NewGenericWriter[Product](file)
	defer writer.Close()

	_, err = writer.Write(products)
	if err != nil {
		return fmt.Errorf("failed to write products: %w", err)
	}

	return nil
}

// ReadProducts reads product data from Parquet file
func (m *SimpleManager) ReadProducts(filename string) ([]Product, error) {
	filePath := filepath.Join(m.baseDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := parquet.NewGenericReader[Product](file)
	defer reader.Close()

	products := make([]Product, reader.NumRows())
	n, err := reader.Read(products)
	if err != nil {
		return nil, fmt.Errorf("failed to read products: %w", err)
	}

	return products[:n], nil
}

// GetBasicFileInfo returns basic information about a Parquet file
func (m *SimpleManager) GetBasicFileInfo(filename string) (*BasicFileInfo, error) {
	filePath := filepath.Join(m.baseDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	pf, err := parquet.OpenFile(file, stat.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to open parquet file: %w", err)
	}

	return &BasicFileInfo{
		Filename: filename,
		FilePath: filePath,
		FileSize: stat.Size(),
		NumRows:  pf.NumRows(),
		Schema:   pf.Schema(),
	}, nil
}

// BasicFileInfo contains basic information about a Parquet file
type BasicFileInfo struct {
	Filename string
	FilePath string
	FileSize int64
	NumRows  int64
	Schema   *parquet.Schema
}

// ListFiles lists all Parquet files in the base directory
func (m *SimpleManager) ListFiles() ([]string, error) {
	if err := m.ensureDir(); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".parquet" {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// DeleteFile deletes a Parquet file
func (m *SimpleManager) DeleteFile(filename string) error {
	filePath := filepath.Join(m.baseDir, filename)
	return os.Remove(filePath)
}