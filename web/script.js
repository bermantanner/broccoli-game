cards = document.querySelectorAll('.game-card');

cards.forEach(card => {
    card.addEventListener('click', function() {
        if (card.classList.contains('selected')) {
            card.classList.remove('selected');
        }
        else {
            cards.forEach(c => {
                c.classList.remove('selected');
            });
            card.classList.add('selected');
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
            errorMsg.innerText = "Disconnected from host.";
        };
    });
}