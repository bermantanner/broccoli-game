@echo off
echo =========================================
echo BOOTING BROCCOLI GAME PLATFORM...
echo =========================================

echo -^> Starting local Go engine on port 8090...
start /B broccoli-engine.exe

timeout /t 2 /nobreak > NUL

echo -^> Establishing secure Cloud Tunnel via Ngrok...
start /B ngrok http 8090

timeout /t 3 /nobreak > NUL

echo =========================================
echo TUNNEL ESTABLISHED!
echo Check the Ngrok window for your public URL.
echo =========================================
echo Opening Host interface...

start http://localhost:8090/web/host.html