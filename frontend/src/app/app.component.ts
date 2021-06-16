import { Component } from '@angular/core';
import * as Wails from '@wailsapp/runtime';

@Component({
  selector: '[id="app"]',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'my-app';

  ngOnInit() {
    Wails.Events.On("ytd:track", payload => console.log(payload))
  }
}
