type EventKind = "net" | "local";

type MethodMeta = {
  propertyKey: string | symbol;
  eventName: string;
  kind: EventKind;
};

const eventStore = new WeakMap<object, MethodMeta[]>();
const instanceGuard = new WeakSet<object>();

function addMethodMeta(
  target: object,
  propertyKey: string | symbol,
  kind: EventKind,
  eventName: string,
): void {
  const existing = eventStore.get(target) ?? [];
  existing.push({ propertyKey, kind, eventName });
  eventStore.set(target, existing);
}

function setEvent(kind: EventKind, eventName: string): MethodDecorator {
  return (target, propertyKey) => {
    addMethodMeta(target, propertyKey, kind, eventName);
  };
}

function registerInstance(instance: object): void {
  if (instanceGuard.has(instance)) return;
  instanceGuard.add(instance);

  const prototype = Object.getPrototypeOf(instance);
  const methods = eventStore.get(prototype) ?? [];

  for (const method of methods) {
    const fn = (instance as Record<string | symbol, unknown>)[
      method.propertyKey
    ];
    if (typeof fn !== "function") continue;

    if (method.kind === "net") {
      onNet(method.eventName, (...args: unknown[]) => {
        const src = source;
        (fn as (...invokeArgs: unknown[]) => void).call(instance, src, ...args);
      });
      continue;
    }

    on(method.eventName, (...args: unknown[]) => {
      (fn as (...invokeArgs: unknown[]) => void).call(instance, ...args);
    });
  }
}

export function EventController(): ClassDecorator {
  return <TFunction extends Function>(target: TFunction): TFunction => {
    const Wrapped = class extends (target as unknown as new (
      ...args: unknown[]
    ) => object) {
      constructor(...args: unknown[]) {
        super(...args);
        registerInstance(this);
      }
    };

    return Wrapped as unknown as TFunction;
  };
}

export function OnNet(eventName: string): MethodDecorator {
  return setEvent("net", eventName);
}

export function OnLocal(eventName: string): MethodDecorator {
  return setEvent("local", eventName);
}

export const RemoteEvent = OnNet;
export const LocalEvent = OnLocal;
