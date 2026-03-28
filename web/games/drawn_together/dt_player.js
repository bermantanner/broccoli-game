let myRole = "";
let isDrawing = false;
let lastX = 0;
let lastY = 0;
let ctx = null;

// Local coordinate space settings
const CANVAS_WIDTH = 300;
const CANVAS_HEIGHT = 250;
const BOUNDARY_ZONE = 40;
const OFFSET = CANVAS_HEIGHT - BOUNDARY_ZONE; // 360

window.handleGameMessage = function(msg) {
    // ---------------------------------------------------------
    // STATE 1: TEAM SELECT
    // ---------------------------------------------------------
    if (msg.type === "team_update") {
        document.getElementById("connection-ui").classList.add("hidden");
        const gameRoot = document.getElementById("game-root");

        if (!document.getElementById("dt-player-app")) {
            gameRoot.innerHTML = `
                <div id="dt-player-app">
                    <h2>pick a team!</h2>
                    <p>status: <span id="dt-status">unassigned</span></p>
                    <div class="team-buttons">
                        <button id="dt-btn-team1">join team 1</button>
                        <button id="dt-btn-team2">join team 2</button>
                    </div>
                </div>
            `;

            document.getElementById("dt-btn-team1").addEventListener("click", () => {
                if (window.appSocket) window.appSocket.send(JSON.stringify({ type: "join_team", team: 1 }));
            });
            document.getElementById("dt-btn-team2").addEventListener("click", () => {
                if (window.appSocket) window.appSocket.send(JSON.stringify({ type: "join_team", team: 2 }));
            });
        }

        const statusSpan = document.getElementById("dt-status");
        
        // use (msg.team1 || []) to default to an empty array if Go sends null
        if ((msg.team1 || []).includes(window.playerName)) {
            statusSpan.innerText = "team 1";
            statusSpan.style.color = "var(--color-blue)";
        } else if ((msg.team2 || []).includes(window.playerName)) {
            statusSpan.innerText = "team 2";
            statusSpan.style.color = "var(--color-blue)";
        } else {
            statusSpan.innerText = "unassigned";
            statusSpan.style.color = "var(--color-black)";
        }
    }

    // ---------------------------------------------------------
    // STATE 2: DRAWING PHASE
    // ---------------------------------------------------------
    if (msg.type === "drawing_started") {
        myRole = msg.assignments[window.playerName];
        const gameRoot = document.getElementById("game-root");
        
        // inject Canvas UI
        gameRoot.innerHTML = `
            <div id="dt-drawing-app">
                <h2>draw the ${myRole.replace("_", " ")}!</h2>
                <div class="canvas-container">
                    ${getBoundaryOverlays(myRole)}
                    <canvas id="dt-canvas" width="${CANVAS_WIDTH}" height="${CANVAS_HEIGHT}"></canvas>
                </div>
            </div>
        `;

        const canvas = document.getElementById("dt-canvas");
        ctx = canvas.getContext("2d");
        ctx.lineWidth = 4;
        ctx.lineCap = "round";
        ctx.strokeStyle = "black";

        // fill background white so base64 export isn't transparent
        ctx.fillStyle = "white";
        ctx.fillRect(0, 0, CANVAS_WIDTH, CANVAS_HEIGHT);

        setupCanvasListeners(canvas);
    }

    // ---------------------------------------------------------
    // receiving a neighbor's line
    // ---------------------------------------------------------
    if (msg.type === "draw_line") {
        if (ctx) {
            drawLine(ctx, msg.startX, msg.startY, msg.endX, msg.endY);
        }
    }

    // ---------------------------------------------------------
    // THE COLLECTOR: time is up
    // ---------------------------------------------------------
    if (msg.type === "time_up") {
        const canvas = document.getElementById("dt-canvas");
        if (canvas) {
            const imageData = canvas.toDataURL("image/png");
            window.appSocket.send(JSON.stringify({
                type: "submit_drawing",
                image: imageData
            }));
            
            document.getElementById("dt-drawing-app").innerHTML = `
                <h2>time's up!</h2>
                <p>look at the host screen...</p>
            `;
        }
    }
    // ---------------------------------------------------------
    // STATE 3: THE REVEAL 
    // ---------------------------------------------------------
    if (msg.type === "reveal_started") {
        const gameRoot = document.getElementById("game-root");
        gameRoot.innerHTML = `
            <div id="dt-player-app">
                <h2>round ${msg.currentRound} complete!</h2>
                <p>look at host screen.</p>
            </div>
        `;
    }
};

// --- DRAWING HELPER FUNCTIONS ---

function drawLine(context, x1, y1, x2, y2) {
    context.beginPath();
    context.moveTo(x1, y1);
    context.lineTo(x2, y2);
    context.stroke();
    context.closePath();
}

function setupCanvasListeners(canvas) {
    const startDrawing = (e) => {
        isDrawing = true;
        const pos = getMousePos(canvas, e);
        lastX = pos.x;
        lastY = pos.y;
    };

    const draw = (e) => {
        if (!isDrawing) return;
        e.preventDefault(); // prevent scrolling for mobile

        const pos = getMousePos(canvas, e);
        const currentX = pos.x;
        const currentY = pos.y;

        if (currentX < 0 || currentX > CANVAS_WIDTH || currentY < 0 || currentY > CANVAS_HEIGHT) {
            stopDrawing();
            return;
        }

        // draw locally
        drawLine(ctx, lastX, lastY, currentX, currentY);

        const neighbors = getNeighbors(myRole);

        // if drawing in the bottom zone, send DOWN
        if (neighbors.down && (lastY > OFFSET || currentY > OFFSET)) {
            window.appSocket.send(JSON.stringify({
                type: "draw_line", targetRole: neighbors.down,
                startX: lastX, startY: lastY - OFFSET,
                endX: currentX, endY: currentY - OFFSET
            }));
        }

        // if drawing in the top zone, send UP
        if (neighbors.up && (lastY < BOUNDARY_ZONE || currentY < BOUNDARY_ZONE)) {
            window.appSocket.send(JSON.stringify({
                type: "draw_line", targetRole: neighbors.up,
                startX: lastX, startY: lastY + OFFSET,
                endX: currentX, endY: currentY + OFFSET
            }));
        }

        lastX = currentX;
        lastY = currentY;
    };

    const stopDrawing = () => { isDrawing = false; };

    canvas.addEventListener("mousedown", startDrawing);
    canvas.addEventListener("mousemove", draw);
    canvas.addEventListener("mouseup", stopDrawing);
    canvas.addEventListener("mouseout", stopDrawing);

    canvas.addEventListener("touchstart", startDrawing, { passive: false });
    canvas.addEventListener("touchmove", draw, { passive: false });
    canvas.addEventListener("touchend", stopDrawing);
}

function getMousePos(canvas, evt) {
    const rect = canvas.getBoundingClientRect();
    const scaleX = canvas.width / rect.width;
    const scaleY = canvas.height / rect.height;
    
    // Handle touch vs mouse
    const clientX = evt.touches ? evt.touches[0].clientX : evt.clientX;
    const clientY = evt.touches ? evt.touches[0].clientY : evt.clientY;

    return {
        x: (clientX - rect.left) * scaleX,
        y: (clientY - rect.top) * scaleY
    };
}

function getNeighbors(role) {
    if (role === 'head') return { down: 'body' };
    if (role === 'body') return { up: 'head', down: 'legs' };
    if (role === 'legs') return { up: 'body' };
    if (role === 'top_half') return { down: 'bottom_half' };
    if (role === 'bottom_half') return { up: 'top_half' };
    return {};
}

// generates the CSS dashed lines purely as visual overlays
function getBoundaryOverlays(role) {
    let html = "";
    if (role === 'body' || role === 'legs' || role === 'bottom_half') {
        html += `<div class="boundary-top"></div>`;
    }
    if (role === 'head' || role === 'body' || role === 'top_half') {
        html += `<div class="boundary-bottom"></div>`;
    }
    return html;
}