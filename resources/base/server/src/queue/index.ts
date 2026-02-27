import { EventListener, LocalEvent, Wrappers } from "@shared";
import { DeferralsCard } from "./deferrals_card";
import { players } from "src/players";

@EventListener()
export class Queue extends Wrappers.Singleton<Queue>() {
  private queue: string[];

  constructor() {
    super();
    StopResource("hardcap");
    this.queue = [];
  }

  @LocalEvent("playerConnecting")
  playerConnecting = async (
    name: string,
    setKickReason: (reason: string) => void,
    deferrals: any,
  ) => {
    const _src = source;
    await deferrals.defer();
    const card = DeferralsCard(deferrals);
    const playerIdentifiers = players.getIdentifiersForPlayer(_src);
    console.log(JSON.stringify(playerIdentifiers, null, 2));
    card(`Welcome ${name}! Validating your rockstar license...`);
    const rockstartLicense = playerIdentifiers.license;
    if (!rockstartLicense) {
      deferrals.done("You need a valid rockstar license");
      return;
    }

    card(`Welcome ${name}! Validating your Fivem ID`);
    const fivemId = playerIdentifiers.fivem;
    if (!fivemId) {
      deferrals.done("You need a valid Fivem ID");
      return;
    }
    card(`Welcome ${name}, Soon you will join the server. Stay tight`);
    setTimeout(() => {
      deferrals.done();
    }, 10000);
  };
}

export const queue = new Queue();
