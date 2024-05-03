package service

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert" // for better assertions
)

func TestSetFile(t *testing.T) {
	f := NewFile()
	index := 10
	collectionHash := "test_hash"
	serverPath := "/tmp/test_folder" // a temporary location for testing

	// Clean up any leftover files from previous tests
	defer os.RemoveAll(serverPath)

	err := f.SetFile(index, collectionHash, serverPath)

	// Assert that there is no error
	assert.Nil(t, err, "Unexpected error creating file")

	// Assert that the directory was created
	_, err = os.Stat(filepath.Join(serverPath, collectionHash))
	assert.NoError(t, err, "Directory not created")

	// Assert that the file was created
	_, err = os.Stat(f.FilePath)
	assert.NoError(t, err, "File not created")

	// Test data for writing
	testData := []byte("This is some test data to write")

	err = f.Write(testData)
	assert.NoError(t, err, "Error writing data to file")

	// Check if the data is written to the file (optional)
	writtenData, err := os.ReadFile(f.FilePath)
	assert.NoError(t, err, "Error reading written data")
	assert.Equal(t, testData, writtenData, "Written data does not match")

	// Close the file (important for proper cleanup)
	defer f.Close()
}

// Data for testing
var testData []byte

func init() {
	// Generate some sample data for testing
	testData = make([]byte, 1024*1024) // Adjust size as needed (e.g., 1MB)
	for i := range testData {
		testData[i] = byte(i % 256) // Fill with some pattern
	}
}

func BenchmarkWrite(b *testing.B) {
	f := NewFile()
	err := f.SetFile(1, "test_hash", "/tmp/test_folder") // Adjust paths as needed
	if err != nil {
		b.Fatal(err)
	}
	defer f.OutputFile.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		f.Write(testData)
	}
}

func BenchmarkWriteFromReader(b *testing.B) {
	f := NewFile()
	err := f.SetFile(1, "test_hash", "/tmp/test_folder") // Adjust paths as needed
	if err != nil {
		b.Fatal(err)
	}
	defer f.OutputFile.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		reader := bytes.NewReader(testData)

		f.WriteFromReader(reader)
		// reader.Seek(0, io.SeekStart)
	}
}
