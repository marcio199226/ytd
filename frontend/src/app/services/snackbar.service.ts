import { Injectable, TemplateRef, EmbeddedViewRef } from '@angular/core';
import { MatSnackBar, MatSnackBarRef, SimpleSnackBar, MatSnackBarConfig } from '@angular/material/snack-bar';

@Injectable({providedIn: 'root'})
export class SnackbarService extends MatSnackBar {
  public openSuccess(text: string, action: string = null, opts: MatSnackBarConfig = {}): MatSnackBarRef<SimpleSnackBar> {
    const { panelClass, ...options } = opts;
    const classes: string[] = [panelClass as string, 'success'];
    return this.open(text, action, { panelClass: classes, ...options });
  }

  public openError(text: string, action: string  = null, opts: MatSnackBarConfig = {}): MatSnackBarRef<SimpleSnackBar> {
    const { panelClass, ...options } = opts;
    const classes: string[] = [panelClass as string, 'error'];
    return this.open(text, action, { panelClass: classes, ...options });
  }

  public openInfo(text: string, action: string  = null, opts: MatSnackBarConfig = {}): MatSnackBarRef<SimpleSnackBar> {
    const { panelClass, ...options } = opts;
    const classes: string[] = [panelClass as string, 'info'];
    return this.open(text, action, { panelClass: classes, ...options });
  }

  public openWarning(text: string, action: string  = null, opts: MatSnackBarConfig = {}): MatSnackBarRef<SimpleSnackBar> {
    const { panelClass, ...options } = opts;
    const classes: string[] = [panelClass as string, 'warning'];
    return this.open(text, action, { panelClass: classes, ...options });
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
