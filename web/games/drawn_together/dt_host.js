let countdownInterval = null;

window.handleGameMessage = function(msg) {
    // ---------------------------------------------------------
    // STATE 1: TEAM SELECT
    // ---------------------------------------------------------
    if (msg.type === "team_update") {
        document.getElementById("connection-ui").classList.add("hidden");
        const gameRoot = document.getElementById("game-root");

        if (!document.getElementById("dt-host-app")) {
            gameRoot.innerHTML = `
                <div id="dt-host-app">
                    <h2>sort into teams</h2>
                    <div class="team-columns">
                        <div class="col">
                            <h3>team 1</h3>
                            <ul id="dt-list-team1"></ul>
                        </div>
                        <div class="col">
                            <h3>unassigned</h3>
                            <ul id="dt-list-unassigned"></ul>
                        </div>
                        <div class="col">
                            <h3>team 2</h3>
                            <ul id="dt-list-team2"></ul>
                        </div>
                    </div>
                    <button id="dt-lock-btn">lock teams & start</button>
                </div>
            `;

            document.getElementById("dt-lock-btn").addEventListener("click", () => {
                if (window.appSocket) window.appSocket.send(JSON.stringify({ type: "confirm_teams" }));
            });
        }

        const renderList = (id, arr) => {
            const ul = document.getElementById(id);
            ul.innerHTML = "";
            (arr || []).forEach(name => {
                const li = document.createElement("li");
                li.innerText = name;
                ul.appendChild(li);
            });
        };

        renderList("dt-list-team1", msg.team1);
        renderList("dt-list-unassigned", msg.unassigned);
        renderList("dt-list-team2", msg.team2);
    }

    // ---------------------------------------------------------
    // STATE 2: CLOCK (DRAWING PHASE)
    // ---------------------------------------------------------
    if (msg.type === "drawing_started") {
        const gameRoot = document.getElementById("game-root");
        
        let timeLeft = (msg.duration || 3) * 60;

        gameRoot.innerHTML = `
            <div id="dt-host-timer-app">
                <h2>drawing in progress!</h2>
                <div id="dt-clock">${formatTime(timeLeft)}</div>
                <p>players, look at your phones.</p>
            </div>
        `;

        const clockElement = document.getElementById("dt-clock");

        // clear any existing intervals just to be safe
        if (countdownInterval) clearInterval(countdownInterval);

        countdownInterval = setInterval(() => {
            timeLeft--;
            clockElement.innerText = formatTime(timeLeft);

            if (timeLeft <= 10) {
                clockElement.style.color = "red";
            }

            if (timeLeft <= 0) {
                clearInterval(countdownInterval);
                clockElement.innerText = "0:00";
                
                // time is up, tell the Go server to collect the drawings.
                if (window.appSocket) {
                    window.appSocket.send(JSON.stringify({ type: "time_up" }));
                }

                // update UI while waiting for Go to process the images
                document.getElementById("dt-host-timer-app").innerHTML = `
                    <h2>time's up!</h2>
                    <p>collecting masterpiece data...</p>
                `;
            }
        }, 1000);
    }

    // ---------------------------------------------------------
    // STATE 3: THE REVEAL
    // ---------------------------------------------------------
    if (msg.type === "reveal_started") {
        const gameRoot = document.getElementById("game-root");
        
        // Helper function to stack images in the correct anatomical order
        const buildDrawingStack = (teamArray) => {
            const roleOrder = teamArray.length === 3 
                ? ['head', 'body', 'legs'] 
                : ['top_half', 'bottom_half'];

            let html = '<div class="drawing-stack">';
            
            // Loop through the correct anatomical order
            roleOrder.forEach(role => {
                // Find which player on this team had this role
                const playerName = teamArray.find(name => msg.assignments[name] === role);
                if (playerName && msg.drawings[playerName]) {
                    html += `<img src="${msg.drawings[playerName]}" class="canvas-slice" />`;
                }
            });
            
            html += '</div>';
            return html;
        };

        gameRoot.innerHTML = `
            <div id="dt-host-reveal-app">
                <h2>round ${msg.currentRound} of ${msg.totalRounds}</h2>
                <div class="reveal-container">
                    <div class="team-result">
                        <h3>team 1</h3>
                        ${buildDrawingStack(msg.team1)}
                    </div>
                    <div class="team-result">
                        <h3>team 2</h3>
                        ${buildDrawingStack(msg.team2)}
                    </div>
                </div>
                <p class="next-round-msg">
                   ${msg.currentRound < msg.totalRounds ? "next round starting soon..." : "game over! returning to lobby..."}
                </p>
            </div>
        `;
    }
};

// helper to turn seconds into M:SS format
function formatTime(seconds) {
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${m}:${s < 10 ? '0' : ''}${s}`;
}