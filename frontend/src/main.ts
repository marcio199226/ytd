import 'core-js/stable';
import { enableProdMode } from '@angular/core';
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';

import { AppModule } from './app/app.module';
import { environment } from './environments/environment';

import 'zone.js'

import * as Wails from '@wails/runtime';
import { AppState, BackendCallbacks } from '@models';

declare global {
  interface Window {
    [key: string]: any,
    APP_STATE: AppState,
    backend: BackendCallbacks
  }
}

if (environment.production) {
  enableProdMode();
}

Wails.Init(() => {
  // fetch app data from backend before ng app will be bootstraped
  Wails.Events.On('ytd:onload', ({ entries, config, appVersion }) => {
    window.APP_STATE = { entries, config, appVersion };
    platformBrowserDynamic().bootstrapModule(AppModule)
      .then(ngModule => Wails.Events.Emit('frontend:ready')) // notify wails about angular ready state
      .catch(err => console.error(err));
  })
});
