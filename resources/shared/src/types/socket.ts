export interface BackendWsMessage<TPayload = unknown> {
  event: string;
  data: TPayload;
}
