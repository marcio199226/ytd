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
import { AppConfig, NgrokState } from '@models';
import { MatTabGroup } from '@angular/material/tabs';

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

    if(this.data.tab) {
      const tab = this._matTabs._tabs.find(tab => tab.textLabel === this.data.tab);
      this._matTabs.selectedIndex = tab.position;
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

  private _getQrCodeData(): string {
    const pwa = window.APP_STATE.pwaUrl;
    if(this.data.isNgrokRunning) {
      const url = new URL(pwa);
      url.searchParams.append("url", this.data.ngrok.url);

      if(this.model.PublicServer.Ngrok.Auth.Enabled) {
        url.searchParams.append("username", this.model.PublicServer.Ngrok.Auth.Username);
        url.searchParams.append("password", this.model.PublicServer.Ngrok.Auth.Password);
      } else {
        url.searchParams.delete("username");
        url.searchParams.delete("password");
      }

      if(this.model.PublicServer.VerifyAppKey) {
        url.searchParams.append("api_key", this.model.PublicServer.AppKey);
      } else {
        url.searchParams.delete("api_key");
      }
      return url.toString();
    }
    return "{}";
  }
}
