package main

import (
	"os"
	"testing"
)

func TestProcessIPFlag(t *testing.T) {
	ipChan := make(chan string, 10)
	processIPFlag("192.168.1.1,192.168.1.2", ipChan)
	close(ipChan)

	expected := []string{"192.168.1.1", "192.168.1.2"}
	for _, exp := range expected {
		if ip := <-ipChan; ip != exp {
			t.Errorf("Expected %s, got %s", exp, ip)
		}
	}
}

func TestProcessFileFlag(t *testing.T) {
	ipChan := make(chan string, 10)
	processFileFlag("test_ips.txt", ipChan)
	close(ipChan)

	expected := []string{"192.168.1.1", "192.168.1.2"} // Assuming these are the IPs in the file
	for _, exp := range expected {
		if ip := <-ipChan; ip != exp {
			t.Errorf("Expected %s, got %s", exp, ip)
		}
	}
}

func TestProcessGarbageFileFlag(t *testing.T) {
	ipChan := make(chan string, 10)
	processGarbageFileFlag("test_garbage.txt", ipChan)
	close(ipChan)

	expected := []string{"192.168.1.1", "192.168.1.2"} // Assuming these are the IPs in the file
	for _, exp := range expected {
		if ip := <-ipChan; ip != exp {
			t.Errorf("Expected %s, got %s", exp, ip)
		}
	}
}

func TestProcessStdin(t *testing.T) {
	ipChan := make(chan string, 10)

	// Create a temporary file and write test data to it
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte("192.168.1.1\n192.168.1.2\n")); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Redirect standard input to the temporary file
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin, err = os.Open(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	processStdin(false, ipChan)
	close(ipChan)

	expected := []string{"192.168.1.1", "192.168.1.2"}
	for _, exp := range expected {
		if ip := <-ipChan; ip != exp {
			t.Errorf("Expected %s, got %s", exp, ip)
		}
	}
}
