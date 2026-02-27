export type SocketMessage<TData = unknown> = {
  event: string;
  data?: TData;
};

export type BaseExports = {
  sendMessage: <T = any>(event: string, data?: T) => void;
  close: () => void;
  registerHandler: <TData = unknown>(
    event: string,
    handler: (message: SocketMessage<TData>) => void,
  ) => void;
};

function getResourceExports<T>(resourceName: string): T {
  const runtime = globalThis as unknown as {
    exports?: Record<string, unknown>;
  };

  const resourceExports = runtime.exports?.[resourceName] as T | undefined;
  if (!resourceExports) {
    throw new Error(`Resource exports not found for '${resourceName}'`);
  }

  return resourceExports;
}

export function registerBaseHandler<TData = unknown>(
  event: string,
  handler: (message: SocketMessage<TData>) => void,
): void {
  getResourceExports<BaseExports>("base").registerHandler<TData>(
    event,
    handler,
  );
}

export function sendMessage<T = any>(event: string, data?: T): void {
  getResourceExports<BaseExports>("base").sendMessage(event, data);
}
