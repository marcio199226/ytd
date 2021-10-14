import { Component, OnInit, OnDestroy, TemplateRef, ChangeDetectionStrategy, ChangeDetectorRef, QueryList, ViewChildren, ViewChild, AfterViewInit, Inject } from '@angular/core';
import { Subscription } from 'rxjs';
import { LoaderService, LoaderState } from 'app/services/loader.service';
import { DOCUMENT } from '@angular/common';

@Component({
  selector: 'loader-overlay',
  templateUrl: './loader.component.html',
  styleUrls: ['./loader.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class LoaderOverlayComponent implements OnInit, AfterViewInit, OnDestroy {

  public loadingState: LoaderState;

  public defaultTpl: "default" | "card" = 'card';

  public templatesMap: {[key: string]: TemplateRef<any>} = {};

  private subscription: Subscription;

  @ViewChild('defaultLoaderTpl')
  private _defaultLoaderTpl: TemplateRef<any>;

  @ViewChild('cardLoaderTpl')
  private _cardLoaderTpl: TemplateRef<any>;

  constructor(
    private _cdr: ChangeDetectorRef,
    private _loaderService: LoaderService,
    @Inject(DOCUMENT) private _document: Document
  ) { }

  ngOnInit(): void {
    this.subscription = this._loaderService.loaderState$.subscribe(
      state => {
        this.loadingState = state;
        if(state.isLoading) {
          this._document.body.classList.add('loading');
        } else {
          this._document.body.classList.remove('loading');
        }
        this._cdr.detectChanges();
      }
    );
  }

  ngAfterViewInit(): void {
    this.templatesMap = {
      'default': this._defaultLoaderTpl,
      'card': this._cardLoaderTpl
    }
  }

  isHidden(): boolean {
    return !this.loadingState.isLoading || !this.loadingState.show;
  }

  ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }
}
