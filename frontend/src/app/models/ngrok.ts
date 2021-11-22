export type NgrokStateEventPayload = NgrokState;

export interface NgrokState {
  status: 'running' | 'error' | 'timeout' | 'killed';
  errCode: string;
  url: string
}

export const getQrcodeData = (ngrokUrl: string): string => {
  const pwa = window.APP_STATE.pwaUrl;
  const config = window.APP_STATE.config;
  const url = new URL(pwa);
  url.searchParams.append("url", ngrokUrl);

  if(config.PublicServer.Ngrok.Auth.Enabled) {
    url.searchParams.append("username", config.PublicServer.Ngrok.Auth.Username);
    url.searchParams.append("password", config.PublicServer.Ngrok.Auth.Password);
  } else {
    url.searchParams.delete("username");
    url.searchParams.delete("password");
  }

  if(config.PublicServer.VerifyAppKey) {
    url.searchParams.append("api_key", config.PublicServer.AppKey);
  } else {
    url.searchParams.delete("api_key");
  }
  return url.toString();
}

