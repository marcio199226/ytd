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
    Telegram: {
      Share: null,
      Username: null
    }
  };

  public isFfmpegAvailable: boolean = false;

  constructor(
    private _cdr: ChangeDetectorRef,
    private _dialogRef: MatDialogRef<SettingsComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { config: AppConfig }
  ) {
  }

  async ngOnInit(): Promise<void> {
    this.model = { ...this.data.config };
    const [err, path] = await to(window.backend.main.AppState.IsFFmpegInstalled());
    if(path) {
      this.isFfmpegAvailable = true;
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
}
