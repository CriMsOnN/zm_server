import * as WebSocket from "ws";
import { RequireConnected, type BackendWsMessage } from "@shared";

export class Socket {
  private socket: WebSocket.WebSocket;
  private isConnected: boolean;
  private messagesHandlers: Record<
    string,
    ((message: BackendWsMessage) => void)[]
  >;

  constructor(url: string) {
    this.socket = new WebSocket.WebSocket(url);
    this.isConnected = false;
    this.messagesHandlers = {};
  }

  public connect() {
    this.socket.onopen = () => {
      this.isConnected = true;
      console.log("Connected to backend websocket");
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
      console.log("WebSocket connection closed");
    };
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
    console.log("Event:", event);
    console.log("Data:", data);
    this.socket.send(JSON.stringify({ event, data }));
  };

  public close() {
    this.socket.close();
  }
}
