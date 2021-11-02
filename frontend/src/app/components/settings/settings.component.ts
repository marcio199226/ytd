import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import to from 'await-to-js';
import { AppConfig } from '@models';

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

  constructor(
    private _cdr: ChangeDetectorRef,
    private _dialogRef: MatDialogRef<SettingsComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { config: AppConfig, isNgrokRunning: boolean }
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
    return JSON.stringify({
      // url: this.model.PublicServer.Ngrok.Url,
      auth: {
        username: this.model.PublicServer.Ngrok.Auth.Username,
        password: this.model.PublicServer.Ngrok.Auth.Password
      },
      ...(this.model.PublicServer.VerifyAppKey && {apiKey: this.model.PublicServer.AppKey})
    })
  }
}
