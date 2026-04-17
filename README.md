# Broccoli Game Platform

A real-time, local-authoritative multiplayer game engine.

## Download & Play
**[Download this release for Mac and Windows here](https://github.com/bermantanner/broccoli-game/releases/latest)**

*(Note: You do not need to download this repository or have Go installed to play. Just download the .zip file for your operating system from the link above, extract it, and follow the included instructions.txt file)*

I am a huge fan of Jackbox-style party games and have a ton of ideas for my own games. I wanted to build an extremely lightweight and free platform where anyone can simply download, host, and play. Whether you are all sitting on the same couch or streaming the host screen across the world on Discord, players should be able to instantly join from their phones/tablets/laptops.

One interesting engineering challenge was when designing the networking architecture, it initially seemed like there were only two realistic options for multiplayer hosting:
1. **Fully Hosted Cloud Servers:** Fast and highly accessible, but expensive to maintain 24/7.
2. **Traditional Port-Forwarding:** Free, but tedious and impossible to expect casual hosts to set up.

I wanted to engineer a third option which was free, fast, and completely accessible. The solution is a locally-hosted Go engine that runs entirely on the Host's machine. To bypass NAT and router restrictions without port-forwarding, the platform automatically utilizes **Ngrok** to punch a secure TCP tunnel out to the public web. 
* The host gets zero local latency.
* The server is completely free.
* Remote players can join instantly via the web.

Currently the only game is 'Drawn Together' which is very bare bones, but I will improve it and make more games.

I plan to change the ngrok situation to something with shorter links and no install. Either creating a custom reverse-proxy or a using a different tunneling service.

---

## Prerequisites (Important)
To host a game so players on the public internet can join via their phones, the Host MUST have **Ngrok** installed on their computer to create the secure tunnel. *(Note: Players joining on their phones/tablets/laptops do not need to install anything).*

**Step 1: Install Ngrok**
* **Mac:** Open your terminal and run `brew install ngrok`, or download it from [ngrok.com/download](https://ngrok.com/download).
* **Windows:** Download the installer from [ngrok.com/download](https://ngrok.com/download).

**Step 2: Authenticate**
Create a free account at ngrok.com. Copy your personal auth token from your dashboard, open your terminal, and run the command they provide (e.g., `ngrok config add-authtoken YOUR_TOKEN`).

---

## How to Play

**For Mac Users:**
1. Double-click `start_game.command`. *(If your Mac says permission denied, open your terminal and run `chmod +x start_game.command` first).*
2. The game will automatically boot up the local server, establish the secure tunnel, and open the Host screen.
3. Tell your friends to connect using the URL displayed on your screen!

**For Windows Users:**
1. Double-click `start_game.bat`.
2. A terminal window will open to establish the secure Ngrok tunnel.
3. The Host screen will open automatically. Give your friends the Ngrok URL displayed on the screen to join!

## Quitting the Game
When you are done playing, simply close the terminal windows to shut down the local server and sever the tunnel.

---

⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡦⠑⢒⠒⠡⠠⠌⣀⠂⠍⡐⠀⠷⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡴⢺⠅⠨⠐⡈⠡⢂⠐⠠⢈⠐⣠⠝⠢⠞⢋⠙⣦⣀⠀⣀⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣤⣤⡴⠟⠡⠚⡁⢂⠁⡂⠌⠠⢁⠲⠋⢉⠠⢁⠂⠤⠈⠄⠎⠛⠟⠲⠟⣟⠪⣐⢆⠄⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⣴⣫⠍⢒⠒⠒⡈⢁⠂⡐⢂⠐⠠⢈⠐⠊⡁⠐⡀⠂⠄⡈⠄⠡⠈⠄⠡⠘⡀⠒⣈⢩⡉⠗⣕⠄⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣰⠟⢄⠊⠄⠌⠡⠐⢂⠡⠐⠠⢈⠐⠠⢈⠐⡀⠡⢀⠁⢂⠐⡈⣄⠡⢈⠂⢡⠠⢁⠄⢒⣈⡆⠽⣷⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢀⡽⠜⡈⢉⠰⢈⠡⠈⢄⠂⣁⠢⠈⠄⡁⠂⡐⢀⠡⠀⠌⡀⢂⠐⡈⢦⡂⠌⡀⢂⠄⠚⠒⠩⢘⠛⣙⢢⡀⠀⠀⠀⠀
⠀⠀⠀⠀⢀⣶⣿⠆⡈⠀⠇⢰⢀⡾⠉⠈⠰⠀⡀⠁⠆⡀⠁⠀⠆⢀⠁⠆⠰⠀⠆⠰⢀⢹⠀⢰⠀⡈⠰⢁⠆⡁⠎⡰⠇⣷⡀⠀⠀⠀
⠀⠀⠀⠀⢿⣜⡍⠒⠔⣉⡐⢂⣼⠁⠄⡁⢂⠡⢀⠁⠂⠄⠡⢈⡰⢄⡂⡌⢤⢣⣈⠐⣘⡡⢌⠀⢂⠤⢁⠢⡈⢆⣉⣒⣌⢾⡁⠀⠀⠀
⠀⠀⠀⠀⣾⢡⢌⠡⢒⡀⠔⡘⢁⠂⡐⠈⠄⠂⠄⡈⢜⡣⣍⣧⣜⣣⡜⣙⠦⣓⠬⣓⢣⡓⢮⣵⣢⠰⣌⣶⠟⠋⠉⠉⠉⠙⠛⢤⡀⠀
⠀⠀⠀⠀⢙⡎⡄⠃⢆⡘⡐⠌⣀⠢⢀⠡⢈⠐⠠⡰⢃⡳⡽⢋⣝⣀⡈⠓⠉⠄⠛⡀⠄⡘⠱⣎⠼⡱⢾⡁⠀⠀⠀⠀⠀⠀⠀⢀⠰⠀
⠀⠀⠀⠀⠘⣷⢈⠜⠤⢐⠠⠊⠄⡐⡀⢢⠓⣎⢣⢓⡭⣎⣵⣿⡿⣿⣿⣷⣬⢀⠐⡀⠐⡀⣁⣬⣶⣬⣿⠀⢰⠀⠠⡏⠀⡀⠀⠈⠀⡆
⠀⠀⠀⠀⣿⣒⣮⣼⣌⣆⡡⢍⡰⢐⡠⢇⡽⣈⠧⣉⢲⢼⣿⡿⠛⠉⠻⣟⣿⣷⡀⢁⢂⣼⣿⠿⠛⢿⣿⣷⣄⠳⣄⠣⣀⣱⡄⠀⠀⠁
⠀⠀⠀⠚⡟⠉⠀⠀⠀⠈⠙⢷⡜⣥⠚⣆⠞⢭⡲⣄⠎⣿⣿⠃⠀⠀⠀⠘⣿⣽⣿⣲⣾⣿⠇⣀⠀⠈⢿⣿⡾⣍⣀⠉⡨⠃⢀⠀⠀⠀
⠀⠀⠀⠀⠃⠀⠀⠀⠴⣁⣀⣨⣿⣐⠫⡴⠏⣤⠙⣌⠳⣿⣿⠀⡞⣿⡆⠀⣿⣟⣾⣟⣯⣿⢰⣻⣇⠀⣺⣿⠍⢻⡌⢣⠖⠒⠻⡄⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠱⠀⡿⠶⣧⡃⠜⡈⠳⣌⠀⢿⣿⡄⣿⣿⡇⢀⣿⣿⣳⣯⣟⣿⡼⣿⡏⢀⣿⡿⠂⢸⡇⠰⠃⠀⠀⡧⠀⡀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠆⡃⠀⠀⣻⠖⠓⠓⢾⣆⠸⣿⣿⣜⡛⢁⣼⣿⣷⠟⠋⠉⢻⣿⣮⣴⣾⣿⠃⠀⠸⡇⠀⠀⠀⠀⢁⡴⠁
⣴⣶⣷⣄⡀⠀⠀⠀⠀⠀⣨⠴⠓⢶⣶⠧⢤⣀⠀⠀⠈⠀⠙⠿⣿⣿⣿⠟⠋⠁⠀⠀⣠⠟⠛⠿⠿⠛⠁⠀⠀⢸⠓⠦⠴⠶⠒⠋⠀⠀
⡼⣿⢯⣿⣿⣷⠒⠂⠀⠀⠁⠀⢀⣼⠃⠀⠀⠈⠃⠀⠀⢀⡀⠀⠀⠀⠀⠀⣴⢤⣀⣀⣀⣀⡀⠀⠀⠀⠀⠀⠀⢼⠀⠀⠀⠀⠀⠀⠀⠀
⠛⣿⣿⣿⣿⣿⣷⣤⣤⢤⠤⡶⠟⠁⠈⠳⠄⠀⠀⠀⠀⣾⠁⠀⠀⠀⠀⠀⣿⢀⡀⠀⣨⠟⠁⠀⠀⠀⠀⠀⠀⡿⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠸⣿⣿⣿⣿⣿⣿⣿⣷⡈⢶⡀⠀⠀⠀⠀⠀⠀⢀⣾⠁⠀⠀⠀⠀⠀⠀⠸⣆⠀⣰⠇⠀⠀⠀⠀⠀⡀⠀⢀⡇⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠹⣽⣿⣿⣿⣿⣿⣿⣿⣶⣮⣑⡒⠀⣀⣠⣴⣿⣧⠀⠀⠀⠀⢻⡀⠀⠀⠮⠭⠁⠀⠀⠀⠀⠀⢰⠇⠀⢸⠁⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⡈⠮⠻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⢷⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡟⠀⢀⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⡀⠄⠠⠑⢉⡻⣿⣿⣿⡿⠿⠛⢻⣿⣿⣿⣿⣿⣿⣿⡧⠀⠀⠀⢸⡄⠀⠀⠀⠀⠀⠀⠀⠀⣼⠁⠀⣸⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠈⠙⠒⠒⠒⠉⠈⠉⠀⠀⠀⠀⠘⠿⠿⠿⠟⠛⠉⠁⡗⠀⠀⠀⠈⡇⠀⠀⠀⠀⠀⠀⠀⣸⠇⠀⠠⡋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣇⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⢀⡠⠟⠂⠁⡸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠓⠒⠲⠦⠧⠄⠀⠀⠀⠀⠀⠈⠀⠀⡠⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀