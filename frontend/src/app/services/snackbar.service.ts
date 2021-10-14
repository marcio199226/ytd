import { LiveAnnouncer } from '@angular/cdk/a11y';
import { BreakpointObserver } from '@angular/cdk/layout';
import { Overlay } from '@angular/cdk/overlay';
import { Injectable, TemplateRef, EmbeddedViewRef, Injector, Optional, SkipSelf, Inject } from '@angular/core';
import { MatSnackBar, MatSnackBarRef, SimpleSnackBar, MatSnackBarConfig, MAT_SNACK_BAR_DEFAULT_OPTIONS } from '@angular/material/snack-bar';
import { TranslateService } from '@ngx-translate/core';

@Injectable({providedIn: 'root'})
export class SnackbarService extends MatSnackBar {
  constructor(
    private _trans: TranslateService,
    _overlay: Overlay,
    _live: LiveAnnouncer,
    _injector: Injector,
    _breakpointObserver: BreakpointObserver,
    @Optional() @SkipSelf() _parentSnackBar: MatSnackBar,
    @Inject(MAT_SNACK_BAR_DEFAULT_OPTIONS) _defaultConfig: MatSnackBarConfig,
  ) {
    super(_overlay, _live, _injector, _breakpointObserver, _parentSnackBar, _defaultConfig)
  }

  public openSuccess(text: string, action: string = null, opts: MatSnackBarConfig = {}, transPayload: any = {}): MatSnackBarRef<SimpleSnackBar> {
    const { panelClass, ...options } = opts;
    const classes: string[] = [panelClass as string, 'success'];
    return this.open(this._trans.instant(text, transPayload), this._trans.instant(action), { panelClass: classes, ...options });
  }

  public openError(text: string, action: string  = null, opts: MatSnackBarConfig = {}, transPayload: any = {}): MatSnackBarRef<SimpleSnackBar> {
    const { panelClass, ...options } = opts;
    const classes: string[] = [panelClass as string, 'error'];
    return this.open(this._trans.instant(text, transPayload), this._trans.instant(action), { panelClass: classes, ...options });
  }

  public openInfo(text: string, action: string  = null, opts: MatSnackBarConfig = {}, transPayload: any = {}): MatSnackBarRef<SimpleSnackBar> {
    const { panelClass, ...options } = opts;
    const classes: string[] = [panelClass as string, 'info'];
    return this.open(this._trans.instant(text, transPayload), this._trans.instant(action), { panelClass: classes, ...options });
  }

  public openWarning(text: string, action: string  = null, opts: MatSnackBarConfig = {}, transPayload: any = {}): MatSnackBarRef<SimpleSnackBar> {
    const { panelClass, ...options } = opts;
    const classes: string[] = [panelClass as string, 'warning'];
    return this.open(this._trans.instant(text, transPayload), this._trans.instant(action), { panelClass: classes, ...options });
  }

  public openTplSuccess(tpl: TemplateRef<any>, opts: MatSnackBarConfig = {}): MatSnackBarRef<EmbeddedViewRef<any>> {
    const { panelClass, ...options } = opts;
    const classes: string[] = [panelClass as string, 'success'];
    return this.openFromTemplate(tpl, { panelClass: classes, ...options });
  }

  public openTplError(tpl: TemplateRef<any>, opts: MatSnackBarConfig = {}): MatSnackBarRef<EmbeddedViewRef<any>> {
    const { panelClass, ...options } = opts;
    const classes: string[] = [panelClass as string, 'error'];
    return this.openFromTemplate(tpl, { panelClass: classes, ...options });
  }
}
