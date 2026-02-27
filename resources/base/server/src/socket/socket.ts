import * as WebSocket from "ws";
import { RequireConnected, type BackendWsMessage } from "@shared";
import { Wrappers } from "@shared/types";
const backendSecret = GetConvar("backend_secret", "");
if (backendSecret === "") {
  throw new Error("backend_secret is not configured");
}
const socketURL = `ws://localhost:8080/ws/backend?secret=${backendSecret}`;
export class Socket extends Wrappers.Singleton<Socket>() {
  private socket: WebSocket.WebSocket | null;
  private isConnected: boolean;
  private messagesHandlers: Record<
    string,
    ((message: BackendWsMessage) => void)[]
  >;

  private retryCount: number;
  private retryInterval: NodeJS.Timeout | null;
  private url: string;

  constructor() {
    super();
    this.url = socketURL;
    this.isConnected = false;
    this.messagesHandlers = {};
    this.retryCount = 0;
    this.retryInterval = null;
    this.socket = null;
  }

  public connect() {
    if (this.isConnected) {
      console.log("Already connected to backend websocket");
      return;
    }
    this.socket = new WebSocket.WebSocket(this.url);
    this.socket.onopen = () => {
      this.isConnected = true;
      console.log("Connected to backend websocket");
      if (this.retryInterval) {
        clearTimeout(this.retryInterval);
        this.retryInterval = null;
        this.retryCount = 0;
      }
    };
    this.socket.onmessage = (event: WebSocket.MessageEvent) => {
      const data = JSON.parse(event.data.toString()) as BackendWsMessage;
      const handlers = this.messagesHandlers[data.event];
      if (handlers) {
        handlers.forEach((handler) => {
          handler(data);
        });
      }
    };
    this.socket.onerror = (error: WebSocket.ErrorEvent) => {
      console.error("WebSocket error:", error.error);
    };
    this.socket.onclose = () => {
      this.isConnected = false;
      this.retryConnect();
    };
  }

  private retryConnect() {
    if (this.retryInterval) {
      clearTimeout(this.retryInterval);
      this.retryInterval = null;
    }
    this.retryInterval = setTimeout(() => {
      console.log(
        `[WebSocket] Retrying to connect... (attempt: ${this.retryCount})`,
      );
      this.retryCount++;
      if (this.retryCount > 3) {
        console.log(
          `[WebSocket] Max retries reached, stopping retries (attempt: ${this.retryCount})`,
        );
        return;
      }

      this.connect();
    }, 5000);
  }

  public registerHandler = <T = unknown>(
    ev: string,
    handler: (message: BackendWsMessage<T>) => void,
  ) => {
    this.messagesHandlers[ev] = [
      ...(this.messagesHandlers[ev] || []),
      (message: BackendWsMessage<unknown>) =>
        handler(message as BackendWsMessage<T>),
    ];
  };

  @RequireConnected()
  public send = <T = unknown>(event: string, data?: T) => {
    const debugObject = { event, data };
    console.log("[Socket] Message sent: ", debugObject);
    if (this.socket) {
      this.socket.send(JSON.stringify({ event, data }));
    } else {
      console.error("Not connected to backend websocket");
    }
  };

  public close() {
    if (this.socket) {
      this.socket.close();
    } else {
      console.error("Not connected to backend websocket");
    }
  }
}

export const socket = new Socket();
