import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  OnInit,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import to from 'await-to-js';
import { AppConfig, getQrcodeData, NgrokState } from '@models';
import { MatTabGroup } from '@angular/material/tabs';
import { Clipboard } from '@angular/cdk/clipboard';
import { SnackbarService } from 'app/services/snackbar.service';

@Component({
  selector: 'settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class SettingsComponent implements OnInit {
  public model: AppConfig = {
    BaseSaveDir: null,
    ClipboardWatch: null,
    ConcurrentDownloads: null,
    ConcurrentPlaylistDownloads: null,
    DownloadOnCopy: null,
    MaxParrallelDownloads: null,
    ConvertToMp3: null,
    CleanWebmFiles: null,
	  RunInBackgroundOnClose: null,
	  CheckForUpdates: null,
	  StartAtLogin: null,
    Language: null,
    Telegram: {
      Share: null,
      Username: null
    },
    PublicServer: {
      Enabled: false,
      Ngrok: {
        Authtoken: null,
        Auth: {
          Enabled: false,
          Username: null,
          Password: null,
        },
      },
      VerifyAppKey: null,
      AppKey: null,
    }
  };

  public isFfmpegAvailable: boolean = false;

  public isNgrokAvailable: boolean = false;

  public publicServerQrcode: string = "{}";

  @ViewChild(MatTabGroup)
  private _matTabs: MatTabGroup;

  constructor(
    private _cdr: ChangeDetectorRef,
    private _dialogRef: MatDialogRef<SettingsComponent>,
    private _clipboard: Clipboard,
    private _snackbar: SnackbarService,
    @Inject(MAT_DIALOG_DATA) public data: { tab?: string, config: AppConfig, isNgrokRunning: boolean, ngrok: NgrokState },
  ) {
  }

  async ngOnInit(): Promise<void> {
    this.model = { ...this.data.config };
    const [errFfmpeg, ffmpegPath] = await to(window.backend.main.AppState.IsFFmpegInstalled());
    if(ffmpegPath) {
      this.isFfmpegAvailable = true;
      this._cdr.detectChanges();
    }

    const [errNgrok, ngrokPath] = await to(window.backend.main.NgrokService.IsNgrokInstalled());
    if(ngrokPath) {
      this.isNgrokAvailable = true;
      this._cdr.detectChanges();
    }

    if(this.model.PublicServer.Enabled) {
      this.reRenderQrcode();
    }
  }

  ngAfterViewInit(): void {
    if(this.data.tab) {
      const tab = this._matTabs._tabs.find(tab => tab.textLabel === this.data.tab);
      this._matTabs.selectedIndex = tab.position;
      this._cdr.detectChanges();
    }
  }

  async changeBaseSaveDir(): Promise<void> {
    const [err, path] = await to(window.backend.main.AppState.SelectDirectory());
    console.log(err, path)
  }

  async openUrl(url: string): Promise<void> {
    await window.backend.main.AppState.OpenUrl(url);
  }

  save(): void {
    this._dialogRef.close({ config: this.model });
  }

  close(): void {
    this._dialogRef.close();
  }

  reRenderQrcode(): void {
    this.publicServerQrcode = this._getQrCodeData();
    this._cdr.detectChanges();
  }

  copyUrlToClipboard(): void {
    this._clipboard.copy(this.publicServerQrcode);
    this._snackbar.openSuccess("SETTINGS.TABS.PUBLIC_SERVER.NGROK.URL_COPIED");
  }

  private _getQrCodeData(): string {
    const pwa = window.APP_STATE.pwaUrl;
    if(this.data.isNgrokRunning) {
      return getQrcodeData(this.data.ngrok.url);
    }
    return "{}";
  }
}
