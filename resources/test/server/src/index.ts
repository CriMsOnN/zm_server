import { registerBaseHandler, sendMessage } from "@shared";

registerBaseHandler<{ message: string }>("pong", (message) => {
  console.log("Received ping message", message.data?.message);
});

sendMessage<{ message: string }>("ping", { message: "Hello, world!" });
