import { Injectable, TemplateRef } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

export interface LoaderState {
  isLoading: boolean;
  show: boolean;
  label: string | null;
  templateRef: TemplateRef<any> | null;
  templateName?: "default" | "card";
  data?: {[key: string]: any};
}

export interface LoaderEvent {
  label: string | null;
  templateRef: TemplateRef<any> | null;
  templateName?: "default" | "card";
  data?: {[key: string]: any};
}

export interface LoaderEventBackend {
  label: string | null;
  templateName?: "default" | "card";
}

@Injectable({providedIn: 'root'})
export class LoaderService {

  private _showSubject: BehaviorSubject<LoaderState> = new BehaviorSubject<LoaderState>({
    isLoading: false,
    show: false,
    label: null,
    templateRef: null,
    templateName: 'default',
    data: null
  });

  constructor() {

  }

  public get loaderState$(): Observable<LoaderState> {
    return this._showSubject.asObservable();
  }

  public show(templateRef: TemplateRef<any> | string = null, data: any = null, templateName: 'default' | 'card' = 'card'): void {
    this._showSubject.next({
      isLoading: true,
      show: true,
      label: typeof templateRef === 'string' ? templateRef : null,
      templateRef: typeof templateRef === 'object' ? templateRef : null,
      templateName,
      ...data && { data }
    });
  }

  public loading(): void {
    this._showSubject.next({
      isLoading: true,
      show: false,
      label: null,
      templateRef: null,
    });
  }

  public hide(): void {
    this._showSubject.next({ isLoading: false, show: false, label: null, templateRef: null });
  }

}
