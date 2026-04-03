document.addEventListener("DOMContentLoaded", () => {
    console.log("Initializing WebSocket connection...");



class ZTERPWebSocket {
    constructor() {
        this.socket = null;
        this.heartbeatInterval = null;
        this.reconnectDelay = 3000; // 3 seconds
        this.connect();
    }

    connect() {
        const protocol = window.location.protocol === "https:" ? "wss" : "ws";
        const wsUrl = `${protocol}://${window.location.host}/ws`;

        this.socket = new WebSocket(wsUrl);

        this.socket.onopen = () => {
            console.log("WebSocket connected");
            this.startHeartbeat();

               
        };

      this.socket.onmessage = (event) => {
            console.log("WS Message:", event.data);

            let msg;
            try {
                msg = JSON.parse(event.data);
            } catch (e) {
                console.warn("Received non‑JSON WS message:", event.data);
                return;
            }

            if (msg.type === "ONLINE_USERS_UPDATE") {
                if (this.onPresenceUpdate) {
                    this.onPresenceUpdate(msg.users);
                }
            }
        };

        this.socket.onclose = () => {
            console.warn("WebSocket closed, reconnecting...");
            this.stopHeartbeat();
            setTimeout(() => this.connect(), this.reconnectDelay);
        };

        this.socket.onerror = (err) => {
            console.error("WebSocket error:", err);
            this.stopHeartbeat();
            this.socket.close();
        };
    }

    startHeartbeat() {
        this.heartbeatInterval = setInterval(() => {
            if (this.socket.readyState === WebSocket.OPEN) {
                this.socket.send(JSON.stringify({ type: "heartbeat" }));
            }
        }, 5000); // every 5 seconds
    }

    stopHeartbeat() {
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
            this.heartbeatInterval = null;
        }
    }

  



} //close class

window.ZTERPWebSocket = new ZTERPWebSocket();

});