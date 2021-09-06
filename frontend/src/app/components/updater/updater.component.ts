import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { UpdateRelease } from '@models';

@Component({
  selector: 'updater',
  templateUrl: './updater.component.html',
  styleUrls: ['./updater.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class UpdaterComponent implements OnInit {

  constructor(
    private _cdr: ChangeDetectorRef,
    private _dialogRef: MatDialogRef<UpdaterComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { release: UpdateRelease, oldVersion: string }
  ) {
  }

  ngOnInit(): void {

  }

  updateAndRestart(): void {
    this._dialogRef.close({ action: 'UpdateAndRestart' });
  }

  updateAndReplace(): void {
    this._dialogRef.close({ action: 'UpdateAndReplace' });
  }

  async updateManually(url: string): Promise<void> {
    await window.backend.main.AppState.OpenUrl(url);
  }

  close(): void {
    this._dialogRef.close();
  }
}
