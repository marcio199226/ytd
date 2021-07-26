import 'core-js/stable';
import { enableProdMode } from '@angular/core';
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';

import { AppModule } from './app/app.module';
import { environment } from './environments/environment';

import 'zone.js'

import * as Wails from '@wailsapp/runtime';

declare global {
  interface Window {
    [key: string]: any
  }
}

if (environment.production) {
  enableProdMode();
}

Wails.Init(() => {
  // fetch app data from backend before ng app will be bootstraped
  Wails.Events.Once('ytd:onload', ({ entries, config }) => {
    window.APP_STATE = { entries, config };
    //setTimeout(() => {
      platformBrowserDynamic().bootstrapModule(AppModule)
      .then(ngModule => Wails.Events.Emit('frontend:ready')) // notify wails about angular ready state
      .catch(err => console.error(err));
    //}, 500000)
  })
});