import {
  ChangeDetectionStrategy,
  Component,
  HostListener,
  Inject,
  Input
} from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'confirmation-dialog',
  templateUrl: './confirmation-dialog.component.html',
  styleUrls: ['./confirmation-dialog.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ConfirmationDialogComponent {

  @Input()
  public cancelBtnLabel: string = null;

  @Input()
  public okBtnLabel: string = null;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: {
      text: string,
      title: string,
      cancelBtnLabel: string;
      okBtnLabel: string;
    },
    public dialogRef: MatDialogRef<ConfirmationDialogComponent>,
  ) { }

  @HostListener('keydown.esc')
  public onEsc(): void {
    this.cancel();
  }

  public delete(): void {
    this.dialogRef.close(true);
  }

  public cancel(): void {
    this.dialogRef.close(false);
  }

}
