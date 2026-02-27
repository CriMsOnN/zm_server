import { EventListener, LocalEvent, RemoteEvent, Wrappers } from "@shared";
import { socket, Socket } from "src/socket/socket";

@EventListener()
export class Players extends Wrappers.Singleton<Players>() {
  private readonly identifiers: Map<number, Record<string, string>>;
  private readonly onlineUsers: Set<number>;
  private readonly socket: Socket;

  constructor() {
    super();
    this.identifiers = new Map();
    this.onlineUsers = new Set();
    this.socket = socket;
  }

  getIdentifiersForPlayer = (player: number | string) => {
    const source = player.toString();
    const identifierList = GetNumPlayerIdentifiers(source);
    const identifiers = {} as Record<string, string>;

    for (let i = 0; i < identifierList; i++) {
      const id = GetPlayerIdentifier(source, i);
      const key = id.replace(/\:\w+/, "");
      identifiers[key] = id;
    }

    this.identifiers.set(+source, identifiers);
    return identifiers;
  };

  @RemoteEvent("base:playerJoined")
  ClientPlayerJoined = (source: number) => {
    if (this.onlineUsers.has(source)) {
      console.log(`Player ${source} already joined the server`);
      return;
    }
    const playerName = GetPlayerName(source.toString());
    const identifiers = this.identifiers.get(source);
    if (identifiers) {
      const fivemID = identifiers.fivem;
      const license = identifiers.license;
      this.socket.send("user.upsert", {
        name: playerName,
        fivem: fivemID,
        license: license,
      });
      this.onlineUsers.add(source);
      this.socket.send("session.joined", {
        netID: source.toString(),
        name: playerName,
        identifier: fivemID,
      });
    } else {
      const identifiers = this.getIdentifiersForPlayer(source);
      const fivemID = identifiers.fivem;
      const license = identifiers.license;
      this.socket.send("user.upsert", {
        name: playerName,
        fivem: fivemID,
        license: license,
      });
      this.onlineUsers.add(source);
      this.socket.send("session.joined", {
        netID: source.toString(),
        name: playerName,
        identifier: identifiers?.fivem,
      });
    }
  };

  @LocalEvent("playerDropped")
  onPlayerDropped = async (reason: string) => {
    const src = +source;
    this.identifiers.delete(src);
    this.onlineUsers.delete(src);
    this.socket.send("session.dropped", {
      netID: source.toString(),
      reason: reason,
    });
  };

  @LocalEvent("playerJoining")
  onPlayerJoined = (oldSource: number) => {
    console.log(
      `Player ${source} with Name ${GetPlayerName(source.toString())} joined the server`,
    );
    const _nSource = +source;
    const identifiers = this.identifiers.get(oldSource);
    if (identifiers) {
      this.identifiers.delete(oldSource);
      this.identifiers.set(_nSource, identifiers);
      const fivemID = identifiers.fivem;
      const license = identifiers.license;
      this.socket.send("user.upsert", {
        name: GetPlayerName(source.toString()),
        fivem: fivemID,
        license: license,
      });
    }
  };
}

export const players = new Players();
