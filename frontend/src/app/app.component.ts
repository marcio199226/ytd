import { ApplicationRef, ChangeDetectorRef, Component } from '@angular/core';
import * as Wails from '@wails/runtime';

@Component({
  selector: 'app-root,[id="app"]',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'frontend';

  constructor(
    private _cdr: ChangeDetectorRef,
    private _appRef: ApplicationRef,
  ) {}

  ngOnInit() {
    window.addEventListener('focus', () => Wails.Events.Emit("ytd:app:focused"));

    window.addEventListener('blur', () => Wails.Events.Emit("ytd:app:blurred"));

    document.addEventListener("visibilitychange", (event) => {
      Wails.Events.Emit("ytd:app:foreground", { isInForeground:  document.visibilityState === 'visible' })
      if(document.visibilityState === 'visible') {
        // run tick such us I encoutered some issues ex open settings from ytd tray -> settings when windows is not visible
        this._appRef.tick();
      }
    });
  }
}
