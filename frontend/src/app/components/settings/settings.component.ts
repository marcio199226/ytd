import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
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
  };

  constructor(
    private _cdr: ChangeDetectorRef,
    private _dialogRef: MatDialogRef<SettingsComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { config: AppConfig }
  ) {
  }

  ngOnInit(): void {
    this.model = { ...this.data.config };
  }

  changeBaseSaveDir(): void {

  }

  save(): void {
    this._dialogRef.close({ config: this.model });
  }

  close(): void {
    this._dialogRef.close();
  }
}
