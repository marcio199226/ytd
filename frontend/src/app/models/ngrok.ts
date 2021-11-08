export type NgrokStateEventPayload = NgrokState;

export interface NgrokState {
  status: 'running' | 'error' | 'timeout' | 'killed';
  errCode: string;
  url: string
}

