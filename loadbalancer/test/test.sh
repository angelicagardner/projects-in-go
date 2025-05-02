#!/bin/bash
set -e

HEALTH_CHECK_INTERVAL=10   # Must match your LB default
LB_PORT=9090
BACKEND_PORTS=(8080 8081 8082)
BIN_DIR="bin"
PIDS=()

mkdir -p "$BIN_DIR"

echo "Building binaries..."
go build -o "$BIN_DIR/backend" mock/backend.go
go build -o "$BIN_DIR/lb" main.go

echo
echo "Starting backend servers..."
for port in "${BACKEND_PORTS[@]}"; do
  PORT=$port "$BIN_DIR/backend" > "backend_$port.log" 2>&1 &
  pid=$!
  echo "→ Backend running on port $port (PID $pid)"
  PIDS+=($pid)
  sleep 0.5
done

echo
echo "Starting load balancer..."
BACKENDS_FLAG=$(printf "http://localhost:%s," "${BACKEND_PORTS[@]}" | sed 's/,$//')
"$BIN_DIR/lb" --healthcheck-interval="${HEALTH_CHECK_INTERVAL}s" \
  -backends="$BACKENDS_FLAG" \
  > lb.log 2>&1 &
lb_pid=$!
echo "→ Load balancer running on port $LB_PORT (PID $lb_pid)"
PIDS+=($lb_pid)

# Ensure cleanup on exit
cleanup() {
  echo
  echo "Cleaning up..."
  for pid in "${PIDS[@]}"; do
    kill -9 "$pid" 2>/dev/null || true
  done
  rm -rf "$BIN_DIR" urls.txt backend_*.log lb.log
  echo "Done!"
}
trap cleanup EXIT

sleep 2

echo
for i in {1..6}; do
  curl -s "http://localhost:$LB_PORT"
  echo
done

echo
echo "Killing backend on port 8081..."
kill -9 "${PIDS[1]}" || true

echo
echo "Waiting for health check to detect failure..."
sleep $((HEALTH_CHECK_INTERVAL + 2))

echo
echo "Testing routing after 8081 is down..."
for i in {1..6}; do
  curl -s "http://localhost:$LB_PORT"
  echo
done

echo
echo "Restarting backend on port 8081..."
PORT=8081 "$BIN_DIR/backend" > "backend_8081.log" 2>&1 &
new_pid=$!
echo "→ Backend restarted on port 8081 (PID $new_pid)"
PIDS+=($new_pid)

echo
echo "Waiting for health check to detect recovery..."
sleep $((HEALTH_CHECK_INTERVAL + 2))

echo
echo "Testing routing after 8081 is back..."
for i in {1..6}; do
  curl -s "http://localhost:$LB_PORT"
  echo
done

echo
echo "Testing concurrent requests..."
cat > urls.txt <<EOF
url = "http://localhost:$LB_PORT"
url = "http://localhost:$LB_PORT"
url = "http://localhost:$LB_PORT"
url = "http://localhost:$LB_PORT"
url = "http://localhost:$LB_PORT"
url = "http://localhost:$LB_PORT"
url = "http://localhost:$LB_PORT"
url = "http://localhost:$LB_PORT"
EOF

curl --parallel --parallel-immediate --parallel-max 3 --config urls.txt
echo
echo "Test script completed successfully."
