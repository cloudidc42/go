// ===== ตัวอย่าง: File I/O Operations =====
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("===== File I/O ใน Go =====\n")

	// ===== Part 1: Writing Files =====
	fmt.Println("--- 1. Writing Files ---")

	// วิธีที่ 1: ioutil.WriteFile (ง่ายที่สุด)
	content1 := []byte("Hello, Go!\nThis is a test file.\n")
	err := ioutil.WriteFile("test1.txt", content1, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	} else {
		fmt.Println("✓ Created test1.txt")
	}

	// วิธีที่ 2: os.Create + Write
	file, err := os.Create("test2.txt")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("Line 1\n")
	_, err = file.WriteString("Line 2\n")
	_, err = file.WriteString("Line 3\n")
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	} else {
		fmt.Println("✓ Created test2.txt")
	}

	// วิธีที่ 3: Using bufio.Writer (สำหรับไฟล์ใหญ่)
	file3, err := os.Create("test3.txt")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file3.Close()

	writer := bufio.NewWriter(file3)
	for i := 1; i <= 10; i++ {
		writer.WriteString(fmt.Sprintf("Line %d\n", i))
	}
	writer.Flush() // สำคัญ! ต้อง flush buffer
	fmt.Println("✓ Created test3.txt\n")

	// ===== Part 2: Reading Files =====
	fmt.Println("--- 2. Reading Files ---")

	// วิธีที่ 1: ioutil.ReadFile (อ่านทั้งไฟล์)
	data, err := ioutil.ReadFile("test1.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	} else {
		fmt.Println("Content of test1.txt:")
		fmt.Println(string(data))
	}

	// วิธีที่ 2: os.Open + Read
	file2, err := os.Open("test2.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file2.Close()

	buffer := make([]byte, 100)
	n, err := file2.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading file: %v\n", err)
	} else {
		fmt.Printf("Read %d bytes from test2.txt:\n", n)
		fmt.Println(string(buffer[:n]))
	}

	// วิธีที่ 3: Line by line with bufio.Scanner
	file4, err := os.Open("test3.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file4.Close()

	fmt.Println("Reading test3.txt line by line:")
	scanner := bufio.NewScanner(file4)
	lineNum := 1
	for scanner.Scan() {
		fmt.Printf("  %d: %s\n", lineNum, scanner.Text())
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v\n", err)
	}
	fmt.Println()

	// ===== Part 3: Append to File =====
	fmt.Println("--- 3. Appending to Files ---")

	fileAppend, err := os.OpenFile("test1.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file for append: %v\n", err)
		return
	}
	defer fileAppend.Close()

	_, err = fileAppend.WriteString("Appended line 1\n")
	_, err = fileAppend.WriteString("Appended line 2\n")
	if err != nil {
		fmt.Printf("Error appending to file: %v\n", err)
	} else {
		fmt.Println("✓ Appended to test1.txt")
	}

	// Show updated content
	data, _ = ioutil.ReadFile("test1.txt")
	fmt.Println("Updated content:")
	fmt.Println(string(data))

	// ===== Part 4: File Info =====
	fmt.Println("--- 4. File Information ---")

	fileInfo, err := os.Stat("test1.txt")
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
	} else {
		fmt.Printf("File name: %s\n", fileInfo.Name())
		fmt.Printf("Size: %d bytes\n", fileInfo.Size())
		fmt.Printf("Permissions: %v\n", fileInfo.Mode())
		fmt.Printf("Modified: %v\n", fileInfo.ModTime())
		fmt.Printf("Is Directory: %v\n\n", fileInfo.IsDir())
	}

	// ===== Part 5: Check File Exists =====
	fmt.Println("--- 5. Check File Exists ---")

	files := []string{"test1.txt", "nonexistent.txt", "test2.txt"}
	for _, filename := range files {
		if fileExists(filename) {
			fmt.Printf("✓ %s exists\n", filename)
		} else {
			fmt.Printf("✗ %s does not exist\n", filename)
		}
	}
	fmt.Println()

	// ===== Part 6: Copy File =====
	fmt.Println("--- 6. Copy File ---")

	err = copyFile("test1.txt", "test1_copy.txt")
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
	} else {
		fmt.Println("✓ Copied test1.txt to test1_copy.txt\n")
	}

	// ===== Part 7: Directory Operations =====
	fmt.Println("--- 7. Directory Operations ---")

	// Create directory
	err = os.Mkdir("testdir", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating directory: %v\n", err)
	} else {
		fmt.Println("✓ Created directory: testdir")
	}

	// Create nested directories
	err = os.MkdirAll("testdir/subdir1/subdir2", 0755)
	if err != nil {
		fmt.Printf("Error creating nested directories: %v\n", err)
	} else {
		fmt.Println("✓ Created nested directories")
	}

	// List directory contents
	entries, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
	} else {
		fmt.Println("\nCurrent directory contents:")
		for _, entry := range entries {
			if entry.IsDir() {
				fmt.Printf("  [DIR]  %s\n", entry.Name())
			} else {
				fmt.Printf("  [FILE] %s (%d bytes)\n", entry.Name(), entry.Size())
			}
		}
	}
	fmt.Println()

	// ===== Part 8: Walk Directory Tree =====
	fmt.Println("--- 8. Walk Directory Tree ---")

	fmt.Println("Walking directory tree from current directory:")
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip certain directories
		if info.IsDir() && (info.Name() == ".git" || info.Name() == "node_modules") {
			return filepath.SkipDir
		}

		// Only show .txt files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			fmt.Printf("  %s (%d bytes)\n", path, info.Size())
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
	fmt.Println()

	// ===== Part 9: Temporary Files =====
	fmt.Println("--- 9. Temporary Files ---")

	// Create temporary file
	tmpFile, err := ioutil.TempFile("", "example-*.txt")
	if err != nil {
		fmt.Printf("Error creating temp file: %v\n", err)
	} else {
		fmt.Printf("✓ Created temp file: %s\n", tmpFile.Name())

		// Write to temp file
		tmpFile.WriteString("Temporary data\n")
		tmpFile.Close()

		// Clean up
		defer os.Remove(tmpFile.Name())
	}

	// Create temporary directory
	tmpDir, err := ioutil.TempDir("", "example-dir-")
	if err != nil {
		fmt.Printf("Error creating temp directory: %v\n", err)
	} else {
		fmt.Printf("✓ Created temp directory: %s\n", tmpDir)

		// Clean up
		defer os.RemoveAll(tmpDir)
	}
	fmt.Println()

	// ===== Part 10: Practical Examples =====
	fmt.Println("--- 10. Practical Examples ---")

	// Example 1: Count lines, words, and characters
	stats := getFileStats("test3.txt")
	fmt.Printf("File statistics for test3.txt:\n")
	fmt.Printf("  Lines: %d\n", stats.Lines)
	fmt.Printf("  Words: %d\n", stats.Words)
	fmt.Printf("  Characters: %d\n\n", stats.Chars)

	// Example 2: Search in file
	searchTerm := "Line 5"
	matches := searchInFile("test3.txt", searchTerm)
	fmt.Printf("Search results for '%s' in test3.txt:\n", searchTerm)
	for lineNum, line := range matches {
		fmt.Printf("  Line %d: %s\n", lineNum, line)
	}
	fmt.Println()

	// Example 3: CSV-like data
	csvData := [][]string{
		{"Name", "Age", "City"},
		{"Alice", "25", "Bangkok"},
		{"Bob", "30", "Chiang Mai"},
		{"Carol", "28", "Phuket"},
	}

	writeCSV("data.csv", csvData)
	fmt.Println("✓ Created data.csv")

	readData := readCSV("data.csv")
	fmt.Println("CSV Data:")
	for _, row := range readData {
		fmt.Printf("  %v\n", row)
	}
	fmt.Println()

	// Example 4: Config file
	config := map[string]string{
		"host":     "localhost",
		"port":     "8080",
		"database": "mydb",
		"user":     "admin",
	}

	writeConfig("config.txt", config)
	fmt.Println("✓ Created config.txt")

	loadedConfig := readConfig("config.txt")
	fmt.Println("Loaded config:")
	for key, value := range loadedConfig {
		fmt.Printf("  %s = %s\n", key, value)
	}

	// ===== Cleanup =====
	fmt.Println("\n--- Cleanup ---")
	cleanupFiles := []string{
		"test1.txt", "test2.txt", "test3.txt", "test1_copy.txt",
		"data.csv", "config.txt",
	}

	for _, filename := range cleanupFiles {
		os.Remove(filename)
	}
	os.RemoveAll("testdir")

	fmt.Println("✓ Cleaned up test files and directories")
}

// ===== Helper Functions =====

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

type FileStats struct {
	Lines int
	Words int
	Chars int
}

func getFileStats(filename string) FileStats {
	stats := FileStats{}

	file, err := os.Open(filename)
	if err != nil {
		return stats
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		stats.Chars += len(line)
		stats.Words += len(strings.Fields(line))
	}

	return stats
}

func searchInFile(filename, searchTerm string) map[int]string {
	matches := make(map[int]string)

	file, err := os.Open(filename)
	if err != nil {
		return matches
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchTerm) {
			matches[lineNum] = line
		}
		lineNum++
	}

	return matches
}

func writeCSV(filename string, data [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, row := range data {
		line := strings.Join(row, ",")
		writer.WriteString(line + "\n")
	}

	return writer.Flush()
}

func readCSV(filename string) [][]string {
	var data [][]string

	file, err := os.Open(filename)
	if err != nil {
		return data
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		data = append(data, fields)
	}

	return data
}

func writeConfig(filename string, config map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for key, value := range config {
		writer.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}

	return writer.Flush()
}

func readConfig(filename string) map[string]string {
	config := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return config
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config[parts[0]] = parts[1]
		}
	}

	return config
}

// วิธีรัน:
// go run 05_file_io.go
//
// Key Points:
// - ใช้ defer file.Close() เสมอ
// - ตรวจสอบ errors ทุกครั้ง
// - ioutil.ReadFile/WriteFile สำหรับไฟล์เล็ก
// - bufio สำหรับไฟล์ใหญ่ (มี buffer)
// - os.OpenFile สำหรับ advanced options
//
// File Modes:
// - os.O_RDONLY: Read only
// - os.O_WRONLY: Write only
// - os.O_RDWR: Read and write
// - os.O_APPEND: Append mode
// - os.O_CREATE: Create if not exists
// - os.O_TRUNC: Truncate when opening
//
// Permissions (Unix):
// - 0644: rw-r--r-- (owner read/write, others read)
// - 0755: rwxr-xr-x (owner all, others read/execute)
// - 0666: rw-rw-rw- (all read/write)
//
// Best Practices:
// - Always close files with defer
// - Check errors after every I/O operation
// - Use bufio for better performance
// - Use filepath package for path operations
// - Use io.Copy for efficient file copying
