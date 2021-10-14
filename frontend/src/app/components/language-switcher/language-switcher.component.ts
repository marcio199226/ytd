import { Component, OnInit, ChangeDetectionStrategy, ChangeDetectorRef, ViewEncapsulation, ViewChild, ApplicationRef } from '@angular/core';
import { Subscription } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { MatSelect, MatSelectChange } from '@angular/material/select';
import * as Wails from '@wailsapp/runtime';
import to from 'await-to-js';
import { SnackbarService } from 'app/services/snackbar.service';

@Component({
  selector: 'language-switcher',
  templateUrl: './language-switcher.component.html',
  styleUrls: ['./language-switcher.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None
})
export class LanguageSwitcherComponent {

  public langs: string[];

  public model: any = {
    lang: null
  }

  public langTransMap: any = {
    'en': 'LANGUAGE.EN',
    'pl': 'LANGUAGE.PL',
    'it': 'LANGUAGE.IT',
  }

  @ViewChild('matSelect')
  private _matSelectRef: MatSelect = null;

  constructor(
    private _appRef: ApplicationRef,
    private _cdr: ChangeDetectorRef,
    private _translateService: TranslateService,
    private _snackbar: SnackbarService
  ) {
    this.langs = this._translateService.langs;
    this.model.lang = this._translateService.currentLang;
  }

  triggerSelect(): void {
    this._matSelectRef.open();
  }

  public changeLanguage(change: MatSelectChange): void {
    const lang = change.value;
    this._translateService.getTranslation(lang).subscribe(async (translations) => {
      this._translateService.use(lang);

      const [error, saved] = await to(window.backend.main.AppState.SaveSettingValue('Language', lang));
      if(error) {
        this._snackbar.openError(this._translateService.instant("LANGUAGE.SAVE_ERROR"));
        return;
      }

      const [errorReload, _] = await to(window.backend.main.AppState.ReloadNewLanguage());
      this._appRef.tick();
    })
  }
}
