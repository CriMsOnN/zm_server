import { socket } from "./socket/socket";
import "./players";
import "./queue";

socket.connect();

exports("sendMessage", socket.send);
exports("close", socket.close);
exports("registerHandler", socket.registerHandler);
