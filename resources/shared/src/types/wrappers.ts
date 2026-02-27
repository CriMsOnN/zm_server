export const Wrappers = {
  Singleton: <T>() => {
    return class {
      static instance: T;

      public static getInstance(): T {
        if (!this.instance) {
          this.instance = new this() as T;
        }
        return this.instance;
      }
    };
  },
};
