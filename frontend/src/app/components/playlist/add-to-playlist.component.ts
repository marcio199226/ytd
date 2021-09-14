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
import { Entry, UpdateRelease } from '@models';
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
  public a: any = {};
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
      this.a[playlist.uuid] = {action: 'selected'};
      this.selectedPlaylists.push(playlist.uuid);
      console.log(this.a)
      return;
    }
    const idx = this.selectedPlaylists.indexOf(playlist.uuid)
    this.selectedPlaylists.splice(idx, 1);
    this.a[playlist.uuid] = {action: 'deselected'};
    console.log(this.a)
  }

  createNew(): void {
    this._dialogRef.close({ action: 'createNew' });
  }

  apply(): void {
    this._dialogRef.close({ selectedPlaylists: this.selectedPlaylists });
  }

  close(): void {
    this._dialogRef.close();
  }
}
