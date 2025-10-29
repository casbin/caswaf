#!/bin/bash

# Demonstration script for CasWAF auto-recovery mechanism
# This script creates a simple test program that crashes and shows how the supervisor restarts it

echo "=== CasWAF Auto-Recovery Demonstration ==="
echo ""

# Create a temporary directory for the demo
DEMO_DIR=$(mktemp -d)
echo "Working directory: $DEMO_DIR"

# Create a simple test program that crashes after a few seconds
cat > "$DEMO_DIR/demo_crash.go" << 'EOF'
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	pid := os.Getpid()
	startTime := time.Now()
	
	fmt.Printf("[%s] Process started (PID: %d)\n", startTime.Format("15:04:05"), pid)
	
	// Simulate some work
	for i := 1; i <= 3; i++ {
		fmt.Printf("[%s] Working... (%d/3)\n", time.Now().Format("15:04:05"), i)
		time.Sleep(1 * time.Second)
	}
	
	// Simulate a crash (only if we're supervised)
	if os.Getenv("CASWAF_SUPERVISED") == "1" {
		fmt.Printf("[%s] CRASH! (simulated for demo)\n", time.Now().Format("15:04:05"))
		panic("Simulated crash for demonstration")
	} else {
		fmt.Printf("[%s] Exiting cleanly (not supervised)\n", time.Now().Format("15:04:05"))
	}
}
EOF

# Create a simple supervisor based on CasWAF's implementation
cat > "$DEMO_DIR/demo_supervisor.go" << 'EOF'
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	EnvSupervisorKey = "CASWAF_SUPERVISED"
	MaxRestarts = 3  // Using 3 for demo (shorter runtime), actual supervisor uses 5
	RestartDelay = 2 * time.Second
)

func main() {
	// If we're the supervised process, run the actual program
	if os.Getenv(EnvSupervisorKey) == "1" {
		// This would be the main application code
		// For demo, we just exit here and let the external program run
		return
	}

	fmt.Println("=== Starting Supervisor ===")
	fmt.Printf("Max restarts: %d\n", MaxRestarts)
	fmt.Printf("Restart delay: %v\n", RestartDelay)
	fmt.Println("")

	restartCount := 0
	
	for {
		if restartCount >= MaxRestarts {
			fmt.Printf("\n=== Max restarts (%d) reached. Stopping supervisor. ===\n", MaxRestarts)
			return
		}
		
		// Get the program to supervise from command line args
		if len(os.Args) < 2 {
			fmt.Println("Usage: supervisor <program>")
			return
		}
		
		programPath := os.Args[1]
		
		// Start the child process
		cmd := exec.Command(programPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(), EnvSupervisorKey+"=1")
		
		if err := cmd.Start(); err != nil {
			fmt.Printf("Failed to start process: %v\n", err)
			return
		}
		
		startTime := time.Now()
		if restartCount == 0 {
			fmt.Printf("Supervisor: Started process (PID: %d)\n\n", cmd.Process.Pid)
		} else {
			fmt.Printf("Supervisor: Restarted process (PID: %d) - restart #%d\n\n", cmd.Process.Pid, restartCount)
		}
		
		// Setup signal handling
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		
		// Wait for process completion or signal
		doneChan := make(chan error, 1)
		go func() {
			doneChan <- cmd.Wait()
		}()
		
		select {
		case err := <-doneChan:
			uptime := time.Since(startTime)
			
			if err != nil {
				fmt.Printf("\n--- Supervisor: Process crashed after %v ---\n", uptime.Round(time.Millisecond))
				restartCount++
				
				if restartCount < MaxRestarts {
					fmt.Printf("Supervisor: Waiting %v before restart...\n\n", RestartDelay)
					time.Sleep(RestartDelay)
					continue
				}
			} else {
				fmt.Printf("\nSupervisor: Process exited cleanly after %v\n", uptime.Round(time.Millisecond))
				return
			}
			
		case sig := <-sigChan:
			fmt.Printf("\nSupervisor: Received signal %v, shutting down...\n", sig)
			if cmd.Process != nil {
				cmd.Process.Signal(sig)
			}
			<-doneChan
			return
		}
	}
}
EOF

# Build the programs
echo "Building demo programs..."
# Build and check exit status explicitly
if ! go build -o "$DEMO_DIR/demo_crash" "$DEMO_DIR/demo_crash.go" 2>&1 | grep -v "go: downloading"; then
    if [ ! -f "$DEMO_DIR/demo_crash" ]; then
        echo "Failed to build demo_crash"
        exit 1
    fi
fi

if ! go build -o "$DEMO_DIR/demo_supervisor" "$DEMO_DIR/demo_supervisor.go" 2>&1 | grep -v "go: downloading"; then
    if [ ! -f "$DEMO_DIR/demo_supervisor" ]; then
        echo "Failed to build demo_supervisor"
        exit 1
    fi
fi

echo "Build successful!"
echo ""
echo "=== Running demonstration (will auto-stop after 3 restarts) ==="
echo ""

# Run the supervisor with the crash program
"$DEMO_DIR/demo_supervisor" "$DEMO_DIR/demo_crash"

echo ""
echo "=== Demonstration Complete ==="
echo ""
echo "Summary:"
echo "- The supervisor started the process"
echo "- Each time the process crashed, the supervisor restarted it"
echo "- After 3 restarts, the supervisor stopped to prevent infinite loops"
echo ""
echo "This demonstrates CasWAF's auto-recovery mechanism!"

# Cleanup
rm -rf "$DEMO_DIR"
