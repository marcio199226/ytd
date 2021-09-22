import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Entry, UpdateRelease } from '@models';
import { OfflinePlaylist } from 'app/models/offline-playlist';

@Component({
  selector: 'create-playlist',
  templateUrl: './create.component.html',
  styleUrls: ['./create.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class CreatePlaylistComponent implements OnInit {

  public model: any = {
    name: null
  }

  constructor(
    private _cdr: ChangeDetectorRef,
    private _dialogRef: MatDialogRef<CreatePlaylistComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { entry: Entry, playlists: OfflinePlaylist[] }
  ) {
  }

  ngOnInit(): void {

  }

  create(): void {
    if(!this.model.name) {
      return;
    }
    this._dialogRef.close({ playlist: this.model });
  }

  close(): void {
    this._dialogRef.close();
  }
}
