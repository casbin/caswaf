// Copyright 2025 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

const (
	// EnvSupervisorKey is the environment variable key to detect if running under supervisor
	EnvSupervisorKey = "CASWAF_SUPERVISED"
	// EnvSupervisorMode indicates if this process is the supervisor itself
	EnvSupervisorMode = "CASWAF_SUPERVISOR_MODE"
	// MaxRestarts is the maximum number of restarts within the restart window
	MaxRestarts = 5
	// RestartWindow is the time window for counting restarts
	RestartWindow = 5 * time.Minute
	// RestartDelay is the delay before restarting after a crash
	RestartDelay = 2 * time.Second
)

// InitSelfGuard initializes the self-recovery mechanism
// On Windows: Launches a new CMD window with supervisor that monitors the main process
// On other platforms: Uses the current terminal to run supervisor
// If already supervised, it does nothing and returns
func InitSelfGuard() {
	// Check if we're already supervised (child process)
	if os.Getenv(EnvSupervisorKey) == "1" {
		// Already supervised, just return and continue normal execution
		fmt.Println("Running under supervisor, proceeding with normal startup...")
		return
	}
	
	// Check if we are the supervisor process itself
	if os.Getenv(EnvSupervisorMode) == "1" {
		// We are the supervisor, run supervisor logic
		err := runSupervisor()
		if err != nil {
			fmt.Printf("Supervisor error: %v\n", err)
			os.Exit(1)
		}
		// Supervisor exited cleanly
		os.Exit(0)
	}
	
	// We are the initial process, need to start supervisor
	if runtime.GOOS == "windows" {
		// On Windows, start supervisor in a new CMD window
		err := startWindowsSupervisor()
		if err != nil {
			fmt.Printf("Failed to start supervisor window: %v\n", err)
			os.Exit(1)
		}
		// Initial process exits after launching supervisor
		fmt.Println("Supervisor window started, this process will exit...")
		os.Exit(0)
	} else {
		// On non-Windows, run supervisor in current terminal
		err := runSupervisor()
		if err != nil {
			fmt.Printf("Supervisor error: %v\n", err)
			os.Exit(1)
		}
		// Supervisor exited cleanly
		os.Exit(0)
	}
}

// escapeWindowsArg escapes an argument for use in Windows command line
// Based on Windows command line parsing rules
func escapeWindowsArg(arg string) string {
	// If the argument doesn't contain special characters, return as-is
	if !containsSpecialChar(arg) {
		return arg
	}
	
	// Escape the argument with quotes and handle internal quotes/backslashes
	result := "\""
	for i := 0; i < len(arg); i++ {
		numBackslashes := 0
		
		// Count consecutive backslashes
		for i < len(arg) && arg[i] == '\\' {
			numBackslashes++
			i++
		}
		
		if i == len(arg) {
			// Backslashes at end of string - double them before closing quote
			for j := 0; j < numBackslashes*2; j++ {
				result += "\\"
			}
			break
		} else if arg[i] == '"' {
			// Backslashes followed by quote - double them and escape the quote
			for j := 0; j < numBackslashes*2; j++ {
				result += "\\"
			}
			result += "\\\""
		} else {
			// Backslashes followed by normal character - keep them as-is
			for j := 0; j < numBackslashes; j++ {
				result += "\\"
			}
			result += string(arg[i])
		}
	}
	result += "\""
	return result
}

// containsSpecialChar checks if a string contains characters that need escaping
func containsSpecialChar(s string) bool {
	for _, c := range s {
		if c == ' ' || c == '\t' || c == '\n' || c == '\v' || c == '"' {
			return true
		}
	}
	// Also need to escape if string ends with backslash (will be followed by closing quote)
	if len(s) > 0 && s[len(s)-1] == '\\' {
		return true
	}
	return false
}

// startWindowsSupervisor starts a new CMD window with the supervisor on Windows
func startWindowsSupervisor() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	
	// Get absolute path
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	
	// Build command arguments for the new CMD window
	// Using cmd /c start to create new window, then cmd /k to keep it open
	args := []string{"/c", "start", "CasWAF Supervisor", "cmd", "/k"}
	
	// Build the command line with properly escaped arguments
	cmdLine := escapeWindowsArg(exePath)
	for _, arg := range os.Args[1:] {
		cmdLine += " " + escapeWindowsArg(arg)
	}
	args = append(args, cmdLine)
	
	// Create the command
	cmd := exec.Command("cmd", args...)
	
	// Set environment variable to mark this as supervisor mode
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=1", EnvSupervisorMode))
	
	// Start the CMD window
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start supervisor window: %w", err)
	}
	
	fmt.Println("Starting CasWAF with supervisor in new CMD window...")
	return nil
}

// runSupervisor starts the supervisor that monitors and restarts the main process
func runSupervisor() error {
	fmt.Println("=====================================")
	fmt.Println("     CasWAF Supervisor Started")
	fmt.Println("=====================================")
	fmt.Println("This window monitors the CasWAF process")
	fmt.Println("and will automatically restart it if it crashes.")
	fmt.Println()
	
	restartTimes := []time.Time{}
	restartCount := 0
	
	for {
		// Clean up old restart times outside the window
		now := time.Now()
		validRestarts := []time.Time{}
		for _, t := range restartTimes {
			if now.Sub(t) < RestartWindow {
				validRestarts = append(validRestarts, t)
			}
		}
		restartTimes = validRestarts
		
		// Check if we've exceeded max restarts
		if len(restartTimes) >= MaxRestarts {
			fmt.Printf("\n[ERROR] Exceeded maximum restart limit (%d restarts in %v)\n", MaxRestarts, RestartWindow)
			fmt.Println("Stopping supervisor to prevent infinite restart loop.")
			fmt.Println("Please check the application logs for errors.")
			return fmt.Errorf("exceeded maximum restart limit (%d restarts in %v), stopping supervisor", MaxRestarts, RestartWindow)
		}
		
		// Start the child process
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		
		// Set environment variable to indicate supervised process
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=1", EnvSupervisorKey))
		
		// Start the process
		if err := cmd.Start(); err != nil {
			fmt.Printf("[ERROR] Failed to start process: %v\n", err)
			return err
		}
		
		processStartTime := time.Now()
		restartCount++
		fmt.Printf("\n[%s] Started CasWAF process (PID: %d, Start #%d)\n", 
			processStartTime.Format("2006-01-02 15:04:05"), cmd.Process.Pid, restartCount)
		
		// Setup signal handling to forward signals to child
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		
		// Wait for process completion or signal
		doneChan := make(chan error, 1)
		go func() {
			doneChan <- cmd.Wait()
		}()
		
		select {
		case err := <-doneChan:
			// Process exited
			exitTime := time.Now()
			uptime := exitTime.Sub(processStartTime)
			
			if err != nil {
				fmt.Printf("\n[%s] CasWAF process crashed after running for %v\n", 
					exitTime.Format("2006-01-02 15:04:05"), uptime)
				fmt.Printf("[ERROR] Exit error: %v\n", err)
				
				// Record this restart
				restartTimes = append(restartTimes, time.Now())
				
				fmt.Printf("[INFO] Waiting %v before restarting... (Restart %d/%d in last %v)\n", 
					RestartDelay, len(restartTimes), MaxRestarts, RestartWindow)
				time.Sleep(RestartDelay)
				
				// Continue to restart
				continue
			} else {
				// Clean exit
				fmt.Printf("\n[%s] CasWAF process exited cleanly after running for %v\n", 
					exitTime.Format("2006-01-02 15:04:05"), uptime)
				fmt.Println("Supervisor shutting down...")
				return nil
			}
			
		case sig := <-sigChan:
			// Received shutdown signal, forward to child
			fmt.Printf("\n[%s] Received signal %v, forwarding to CasWAF process...\n", 
				time.Now().Format("2006-01-02 15:04:05"), sig)
			if cmd.Process != nil {
				cmd.Process.Signal(sig)
			}
			// Wait for child to exit
			<-doneChan
			fmt.Println("Supervisor shutting down...")
			return nil
		}
	}
}
