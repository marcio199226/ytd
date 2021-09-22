import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Entry } from '@models';
import { OfflinePlaylist } from 'app/models/offline-playlist';

@Component({
  selector: 'add-to-playlist',
  templateUrl: './add-to-playlist.component.html',
  styleUrls: ['./add-to-playlist.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class AddToPlaylistComponent implements OnInit {

  public selectedPlaylists: string[] = [];
  public alreadyInPlaylist: any = {};

  constructor(
    private _cdr: ChangeDetectorRef,
    private _dialogRef: MatDialogRef<AddToPlaylistComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { entry: Entry, playlists: OfflinePlaylist[],  foundInPlaylists: OfflinePlaylist[] }
  ) {
  }

  ngOnInit(): void {
    this.selectedPlaylists = this.data.foundInPlaylists.map(p => p.uuid);
    this.data.foundInPlaylists.forEach(p => this.alreadyInPlaylist[p.uuid] = true);
    this._cdr.detectChanges();
  }

  isSelectedPlaylist(uuid: string): boolean {
    return this.selectedPlaylists.findIndex(pUUID => pUUID === uuid) > -1;
  }

  add(change: MatCheckboxChange, playlist: OfflinePlaylist): void {
    if(change.checked) {
      this.selectedPlaylists.push(playlist.uuid);
      return;
    }
    const idx = this.selectedPlaylists.indexOf(playlist.uuid)
    this.selectedPlaylists.splice(idx, 1);
  }

  createNew(): void {
    this._dialogRef.close({ action: 'createNew' });
  }

  apply(): void {
    const uuids = Object.keys(this.alreadyInPlaylist);
    const diff = this.selectedPlaylists.filter(uuid => uuids.indexOf(uuid) === -1);
    this._dialogRef.close({ selectedPlaylists: diff });
  }

  close(): void {
    this._dialogRef.close();
  }
}
