import { Component } from '@angular/core';
import * as Wails from '@wailsapp/runtime';

@Component({
  selector: 'app-root,[id="app"]',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'frontend';

  ngOnInit() {
/*     Wails.Events.On("ytd:onload", payload => {
      console.log(payload);
      window.APP_STATE = payload;
    }) */

    Wails.Events.On("ytd:track", payload => console.log(payload))

    Wails.Events.On("ytd:track:progress", payload => console.log("Progress of track download", payload))

    Wails.Events.On("ytd:playlist", payload => console.log(payload))
  }
}
