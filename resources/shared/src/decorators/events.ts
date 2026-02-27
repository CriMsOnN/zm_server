import "reflect-metadata";

export const RemoteEvent = (eventName: string) => {
  return function (target: any, key: any) {
    if (!Reflect.hasMetadata("events", target)) {
      Reflect.defineMetadata("events", [], target);
    }

    const netEvents = Reflect.getMetadata("events", target) as any[];

    netEvents.push({
      name: eventName,
      net: true,
      key,
    });

    Reflect.defineMetadata("events", netEvents, target);
  };
};

export const LocalEvent = (eventName: string) => {
  return function (target: any, key: any) {
    if (!Reflect.hasMetadata("events", target)) {
      Reflect.defineMetadata("events", [], target);
    }

    const netEvents = Reflect.getMetadata("events", target) as any[];

    netEvents.push({
      name: eventName,
      net: false,
      key,
    });
    Reflect.defineMetadata("events", netEvents, target);
  };
};

export const EventListener = () => {
  return function <T extends { new (...args: any[]): any }>(constructor: T) {
    return class extends constructor {
      constructor(...args: any[]) {
        super(...args);

        if (!Reflect.hasMetadata("events", this)) {
          Reflect.defineMetadata("events", [], this);
        }

        const events = Reflect.getMetadata("events", this) as any[];
        events.forEach((event) => {
          if (event.net) {
            onNet(event.name, (...args: any[]) => {
              this[event.key](source, ...args);
            });
          } else {
            on(event.name, (...args: any[]) => {
              this[event.key](...args);
            });
          }
        });
      }
    };
  };
};
