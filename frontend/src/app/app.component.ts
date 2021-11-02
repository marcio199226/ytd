import { ApplicationRef, ChangeDetectorRef, Component, NgZone } from '@angular/core';
import { MatIconRegistry } from '@angular/material/icon';
import { DomSanitizer } from '@angular/platform-browser';
import { TranslateService } from '@ngx-translate/core';
import * as Wails from '@wails/runtime';
import { RegisterCustomIcons } from './common/custom-icons';
import { LoaderService, LoaderEventBackend } from './services/loader.service';

@Component({
  selector: 'app-root,[id="app"]',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  constructor(
    private _ngZone: NgZone,
    private _cdr: ChangeDetectorRef,
    private _appRef: ApplicationRef,
    private _loader: LoaderService,
    private _matIconRegistry: MatIconRegistry,
    private _domSanitizer: DomSanitizer,
    private _translateService: TranslateService
  ) {
    RegisterCustomIcons(this._matIconRegistry, this._domSanitizer);

    this._translateService.addLangs(['en', 'it', 'pl']);
    this._translateService.use(window.APP_STATE.config.Language || 'en');
  }

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

    Wails.Events.On("ytd:loader:show", (payload: LoaderEventBackend) => {
      this._loader.show(payload.label, null , payload.templateName);
    });

    Wails.Events.On("ytd:loader:hide", () => {
      this._loader.hide();
    });

    Wails.Events.On("ytd:notification", payload => {
      console.log('New notification from backend', payload)
    })
  }
}
