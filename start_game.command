#!/bin/bash

# Move into the directory where the script is located
cd "$(dirname "$0")"

echo "========================================="
echo "BOOTING BROCCOLI GAME PLATFORM..."
echo "========================================="

# 1. Start the local Go server in the background
echo "-> Starting local Go engine on port 8090..."
./broccoli-engine &
GO_PID=$!

sleep 2

# 2. Start Ngrok in the background
echo "-> Establishing secure Cloud Tunnel via Ngrok..."
ngrok http 8090 > /dev/null &
NGROK_PID=$!

sleep 3

# 3. Ask Ngrok for the public URL
PUBLIC_URL=$(curl -s localhost:4040/api/tunnels | grep -o '"public_url":"https://[^"]*' | cut -d'"' -f4)

echo "========================================="
echo "TUNNEL ESTABLISHED SUCCESSFULLY"
echo " "
echo "HOST SCREEN (LOCAL): http://localhost:8090/web/host.html"
echo "PLAYERS JOIN (PUBLIC): $PUBLIC_URL/web/join.html"
echo " "
echo "========================================="
echo "Opening Host interface..."

# 4. Open the LOCAL path, passing the public URL as a variable to display
open "http://localhost:8090/web/host.html?tunnel=$PUBLIC_URL"

trap 'echo "Shutting down..."; kill $GO_PID; kill $NGROK_PID; exit' INT
wait