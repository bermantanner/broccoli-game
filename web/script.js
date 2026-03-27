const cards = document.querySelectorAll('.game-card');
const optionsBox = document.getElementById('options-box');
const optionsContent = document.getElementById('options-content');

cards.forEach(card => {
    card.addEventListener('click', function() {
        // if clicking the already selected card, deselect it
        if (card.classList.contains('selected')) {
            card.classList.remove('selected');
            optionsBox.classList.add('hidden');
            optionsContent.innerHTML = '';
        }
        else {
            // deselect all others, select this one
            cards.forEach(c => c.classList.remove('selected'));
            card.classList.add('selected');

            // inject specific options based on the game name
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
                // for other games right now, leave options empty and hidden
                optionsBox.classList.add('hidden');
                optionsContent.innerHTML = '';
            }
        }
    });
});

if (window.location.pathname.includes("host.html")) {
    let roomCode = "";
    let socket = null;

    // ask Go server for a new room code
    fetch("/create-room", { 
        method: "POST",
        headers: {
            "ngrok-skip-browser-warning": "true" 
        }
    })
        .then(response => response.json())
        .then(data => {
            roomCode = data.room;
            
            // update the UI with the real code
            document.querySelector(".code").innerText = roomCode;
            
            // open the WebSocket connection as the Host
            connectWebSocket();
        })
        .catch(err => console.error("Error creating room:", err));

    function connectWebSocket() {
        // build the WS URL with the new code and role
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${wsProtocol}//${window.location.host}/ws?room=${roomCode}&name=HostScreen&role=host`;
        socket = new WebSocket(wsUrl);

        socket.onopen = () => {
            console.log("Host connected to room:", roomCode);
        };

        // listen for messages from the Go Hub
        socket.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            console.log("Received:", msg);

            // if Go server sends a lobby snapshot, update the UI
            if (msg.type === "lobby") {
                updatePlayerList(msg.players);
            }
        };

        socket.onclose = () => {
            console.log("Disconnected from server");
        };
    }

    // helper function to render the players in the HTML
    function updatePlayerList(players) {
        const playerList = document.getElementById("player-list");
        playerList.innerHTML = ""; // clear hardcoded list

        if (players.length === 0) {
            playerList.innerHTML = "<li>waiting for players...</li>";
            return;
        }

        players.forEach(playerName => {
            const li = document.createElement("li");
            li.innerText = playerName;
            playerList.appendChild(li);
        });
    }

    const startBtn = document.getElementById("startGameBtn");

    if (startBtn) {
        startBtn.addEventListener("click", () => {
            // make sure a game is actually selected
            const selectedGame = document.querySelector('.game-card.selected');
            if (!selectedGame) {
                console.log("no game selected!");
                return; // prevent starting without a game
            }

            if (socket && socket.readyState === WebSocket.OPEN) {
                const gameName = selectedGame.innerText.trim();
                const payload = {
                    type: "start_game",
                    game: gameName
                };

                // grab the specific settings for selected game
                if (gameName === 'drawn together') {
                    payload.rounds = parseInt(document.getElementById('rounds').value) || 3;
                    payload.duration = parseInt(document.getElementById('duration').value) || 3;
                }

                socket.send(JSON.stringify(payload));
                console.log("Sent start signal:", payload);
            } else {
                console.error("Cannot start game: WebSocket is not connected.");
            }
        });
    }
}

// player join logic
if (window.location.pathname.includes("join.html")) {
    const joinBtn = document.getElementById("joinBtn");
    const codeInput = document.getElementById("codeInput");
    const nameInput = document.getElementById("nameInput");
    const errorMsg = document.getElementById("error-msg");
    
    const joinUi = document.getElementById("join-ui");
    const waitingUi = document.getElementById("waiting-ui");

    let playerSocket = null;

    joinBtn.addEventListener("click", () => {
        const roomCode = codeInput.value.toUpperCase();
        const playerName = nameInput.value.trim();

        if (roomCode.length !== 4 || playerName === "") {
            errorMsg.innerText = "please enter a valid 4-letter code and name.";
            return;
        }

        errorMsg.innerText = "connecting...";

        // build the WebSocket URL specifically for a player
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${wsProtocol}//${window.location.host}/ws?room=${roomCode}&name=${encodeURIComponent(playerName)}&role=player`;
        
        playerSocket = new WebSocket(wsUrl);

        // trap 1 -> connection opens successfully
        playerSocket.onopen = () => {
            console.log("Connected to room:", roomCode);
            // Hide the form, show the waiting screen
            joinUi.classList.add("hidden");
            waitingUi.classList.remove("hidden");
        };

        // trap 2 -> Go Hub sends us a message
        playerSocket.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            console.log("Player received:", msg);
            // later, will listen for {"type": "start_game"} here to switch to the drawing canvas!
        };

        // trap 3 -> connection drops
        playerSocket.onclose = () => {
            console.log("Disconnected from server.");
            joinUi.classList.remove("hidden");
            waitingUi.classList.add("hidden");
            errorMsg.innerText = "disconnected from host.";
        };
    });
}