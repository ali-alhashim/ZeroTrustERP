console.log("Initializing WebSocket connection...");

class ZTERPWebSocket {
    constructor() {
        this.socket = null;
        this.reconnectDelay = 3000;
        this.onPresenceUpdate = null; // ✅ allow pages to hook into presence updates
        this.connect();
    }

    connect() {
        const protocol = window.location.protocol === "https:" ? "wss" : "ws";
        const wsUrl = `${protocol}://${window.location.host}/ws`;

        this.socket = new WebSocket(wsUrl);

        this.socket.onopen = () => {
            console.log("✅ WebSocket connected");
        };

        this.socket.onmessage = (event) => {
            console.log("WS Message:", event.data);

            let msg;
            try {
                msg = JSON.parse(event.data);
            } catch (e) {
                console.warn("Non-JSON WS message:", event.data);
                return;
            }

            // ✅ Handle presence update
            if (msg.type === "ONLINE_USERS_UPDATE") {
                console.log("✅ Presence update received:", msg.users);

                if (this.onPresenceUpdate) {
                    this.onPresenceUpdate(msg.users);
                }
            }
        };

        this.socket.onclose = () => {
            console.warn("❌ WS closed, reconnecting...");
            setTimeout(() => this.connect(), this.reconnectDelay);
        };

        this.socket.onerror = (err) => {
            console.error("WS error:", err);
            this.socket.close();
        };
    }
}

// ✅ Create global instance ONCE
window.ZTERPWebSocket = new ZTERPWebSocket();