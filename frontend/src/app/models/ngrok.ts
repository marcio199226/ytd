export type NgrokStateEventPayload = NgrokState;

export interface NgrokState {
  status: 'running' | 'error';
  errCode: string;
  url: string
}

