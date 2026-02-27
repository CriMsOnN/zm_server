import { Socket } from "./socket/socket";
const backendSecret = GetConvar("backend_secret", "");
if (backendSecret === "") {
  throw new Error("backend_secret is not configured");
}
const socketURL = `ws://localhost:8080/ws/backend?secret=${backendSecret}`;
const socket = new Socket(socketURL);
socket.connect();

exports("sendMessage", socket.send);
exports("close", socket.close);
exports("registerHandler", socket.registerHandler);
