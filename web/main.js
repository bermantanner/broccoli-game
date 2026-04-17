// make sockets globally available so game modules can use them later
window.appSocket = null;
window.playerName = "";

// ==========================================
// HOST LOGIC
// ==========================================
if (window.location.pathname.includes("host.html")) {
    let roomCode = "";

    // UI Logic for the Lobby
    const cards = document.querySelectorAll('.game-card');
    const optionsBox = document.getElementById('options-box');
    const optionsContent = document.getElementById('options-content');

    cards.forEach(card => {
        card.addEventListener('click', function() {
            if (card.classList.contains('selected')) {
                card.classList.remove('selected');
                optionsBox.classList.add('hidden');
                optionsContent.innerHTML = '';
            } else {
                cards.forEach(c => c.classList.remove('selected'));
                card.classList.add('selected');

                const gameName = card.innerText.trim();
                if (gameName === 'drawn together') {
                    optionsBox.classList.remove('hidden');
                    optionsContent.innerHTML = `
                        <div class="option-row">
                            <label for="rounds">rounds (1-3):</label>
                            <input type="number" id="rounds" min="1" max="3" value="3">
                        </div>
                        <div class="option-row">
                            <label for="duration">duration (1-30m):</label>
                            <input type="number" id="duration" min="1" max="30" value="3">
                        </div>
                    `;
                } else {
                    optionsBox.classList.add('hidden');
                    optionsContent.innerHTML = '';
                }
            }
        });
    });

    // Network Logic
    fetch("/create-room", { 
        method: "POST",
        headers: { "ngrok-skip-browser-warning": "true" }
    })
    .then(response => response.json())
    .then(data => {
        roomCode = data.room;
        document.querySelector(".code").innerText = roomCode;
        
        // Grab the Ngrok URL and display it safely in the dedicated div
        const urlParams = new URLSearchParams(window.location.search);
        const tunnelUrl = urlParams.get('tunnel');
        if (tunnelUrl) {
            const linkDisplay = document.getElementById("join-link-display");
            if (linkDisplay) {
                linkDisplay.innerHTML = `go to: <br><strong class="public-join-link">${tunnelUrl}/web/join.html</strong>`;            
            }
        }

        connectHostWebSocket();
    })
    .catch(err => console.error("Error creating room:", err));

    function connectHostWebSocket() {
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${wsProtocol}//${window.location.host}/ws?room=${roomCode}&name=HostScreen&role=host`;
        window.appSocket = new WebSocket(wsUrl);

        window.appSocket.onopen = () => {
            console.log("Host connected to room:", roomCode);
            
            // sending a pulse every 30 seconds, hopefully this keeps server alive 
            setInterval(() => {
                if (window.appSocket.readyState === WebSocket.OPEN) {
                    window.appSocket.send(JSON.stringify({ type: "ping" }));
                }
            }, 30000); 
        };

        // if the socket dies, this message takes over the screen
        window.appSocket.onclose = () => {
            document.body.innerHTML = `
                <div style="display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; background: #ffebee; color: #c62828; font-family: monospace;">
                    <h1>⚠️ CONNECTION LOST</h1>
                    <p>The server dropped the connection.</p>
                    <button onclick="location.reload()" style="padding: 10px 20px; font-size: 1.2rem; cursor: pointer;">Refresh Page</button>
                </div>
            `;
        };

        window.appSocket.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            console.log("Host Received:", msg);

            if (msg.type === "lobby") {
                document.getElementById("connection-ui").classList.remove("hidden");
                document.getElementById("game-root").innerHTML = "";
                updatePlayerList(msg.players);
            } else if (msg.type === "error") {
                alert(msg.message);
            } else {
                // ROUTING... if it's not a lobby message, then pass it to the active game module
                if (typeof window.handleGameMessage === "function") {
                    window.handleGameMessage(msg);
                }
            }
        };
    }

    function updatePlayerList(players) {
        const playerList = document.getElementById("player-list");
        playerList.innerHTML = "";
        if (players.length === 0) {
            playerList.innerHTML = "<li>waiting for players...</li>";
            return;
        }
        players.forEach(name => {
            const li = document.createElement("li");
            li.innerText = name;
            playerList.appendChild(li);
        });
    }

    const startBtn = document.getElementById("startGameBtn");
    if (startBtn) {
        startBtn.addEventListener("click", () => {
            const selectedGame = document.querySelector('.game-card.selected');
            if (!selectedGame) {
                alert("please select a game first!");
                return;
            }
            if (window.appSocket && window.appSocket.readyState === WebSocket.OPEN) {
                const gameName = selectedGame.innerText.trim();
                const payload = { type: "start_game", game: gameName };

                if (gameName === 'drawn together') {
                    payload.rounds = parseInt(document.getElementById('rounds').value) || 3;
                    payload.duration = parseInt(document.getElementById('duration').value) || 3;
                }
                window.appSocket.send(JSON.stringify(payload));
            }
        });
    }
}

// ==========================================
// PLAYER LOGIC
// ==========================================
if (window.location.pathname.includes("join.html")) {
    const joinBtn = document.getElementById("joinBtn");
    const codeInput = document.getElementById("codeInput");
    const nameInput = document.getElementById("nameInput");
    const errorMsg = document.getElementById("error-msg");
    const joinUi = document.getElementById("join-ui");
    const waitingUi = document.getElementById("waiting-ui");

    joinBtn.addEventListener("click", () => {
        const roomCode = codeInput.value.toUpperCase();
        window.playerName = nameInput.value.trim();

        if (roomCode.length !== 4 || window.playerName === "") {
            errorMsg.innerText = "please enter a valid 4-letter code and name.";
            return;
        }

        errorMsg.innerText = "connecting...";

        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${wsProtocol}//${window.location.host}/ws?room=${roomCode}&name=${encodeURIComponent(window.playerName)}&role=player`;
        
        window.appSocket = new WebSocket(wsUrl);

        window.appSocket.onopen = () => {
            joinUi.classList.add("hidden");
            waitingUi.classList.remove("hidden");
        };

        window.appSocket.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            console.log("Player Received:", msg);
            
            if (msg.type === "lobby") {
                document.getElementById("connection-ui").classList.remove("hidden");
                document.getElementById("game-root").innerHTML = "";
                // phones go back to the waiting screen
                document.getElementById("join-ui").classList.add("hidden");
                document.getElementById("waiting-ui").classList.remove("hidden");
                return; // Stop routing!
            }

            // ROUTER: Pass messages to the active game module
            if (typeof window.handleGameMessage === "function") {
                window.handleGameMessage(msg);
            }
        };

        window.appSocket.onclose = () => {
            joinUi.classList.remove("hidden");
            waitingUi.classList.add("hidden");
            errorMsg.innerText = "disconnected from host.";
        };
    });
}