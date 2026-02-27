export function RequireConnected() {
  return (_target: object, _propertyKey: string | symbol): any => {
    const storageKey = Symbol(String(_propertyKey));
    Object.defineProperty(_target, _propertyKey, {
      configurable: true,
      enumerable: false,
      get(this: Record<string | symbol, unknown>) {
        return this[storageKey];
      },
      set(
        this: { isConnected?: boolean } & Record<string | symbol, unknown>,
        value: unknown,
      ) {
        if (typeof value === "function") {
          const original = value as (...args: unknown[]) => unknown;
          this[storageKey] = (...args: unknown[]) => {
            if (!this.isConnected) {
              return undefined;
            }
            return original.apply(this, args);
          };
          return;
        }

        this[storageKey] = value;
      },
    });
  };
}
