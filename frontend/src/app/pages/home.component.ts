import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { FormControl } from '@angular/forms';
import { AudioPlayerService } from 'app/components/audio-player/audio-player.service';
import { Track, Entry } from '@models';
import { AppConfig, AppState } from '../models/app-state';
import * as Wails from '@wailsapp/runtime';
import { MatDialog } from '@angular/material/dialog';
import { SettingsComponent } from 'app/components';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class HomeComponent implements OnInit {
  public searchInput: FormControl;

  public entries: Entry[] = [];

  public inPlayback: Track = null;

  public get inPlaybackTrackId(): string {
    if(!this.inPlayback) {
      return null;
    }
    return this.inPlayback.id;
  }

  constructor(
    private _cdr: ChangeDetectorRef,
    private _dialog: MatDialog,
    private _snackbar: MatSnackBar,
    private _audioPlayerService: AudioPlayerService
  ) {
    this.searchInput = new FormControl('');
  }

  ngOnInit(): void {
    this.entries = (window.APP_STATE as AppState).entries;
    console.log(this.entries)


    Wails.Events.On("ytd:track", (payload: Entry) => {
      console.log("ytd:track", payload)
      const entry = this.entries.find(e => e.track.id === payload.track.id);
      if(entry) {
        entry.track = payload.track;
      } else {
        this.entries.unshift(payload);
      }
      this._cdr.detectChanges();
    });

    Wails.Events.On("ytd:track:progress", ({ id, progress }) => {
      const entry = this.entries.find(e => e.track.id === id);
      entry.track.downloadProgress = progress;
      this._cdr.detectChanges();
    });

    Wails.Events.On("ytd:playlist", payload => console.log(payload))

    this._audioPlayerService.onPlayCmdTrack.subscribe(track => {
      this.inPlayback = track;
      this._cdr.detectChanges();
    });

    this._audioPlayerService.onStopCmdTrack.subscribe(track => {
      this.inPlayback = null;
      this._cdr.detectChanges();
    });
  }

  clearSearch(): void {
    this.searchInput.setValue('');
  }

  openSettings(): void {
    const dialogRef = this._dialog.open(SettingsComponent, {
      autoFocus: false,
      panelClass: ['settings-dialog',  'with-header-dialog'],
      width: '600px',
      maxHeight: '700px',
      data: { config: window.APP_STATE.config }
    });

    dialogRef.afterClosed().subscribe(async (result: { config: AppConfig }) => {
      if(!result) {
        return;
      }

      const { config } = result;
      // put all for into try/catch block and open snackbar with error if something fails
      try {
        for (const [key, value] of Object.entries(config)) {
          switch(key) {
            case 'BaseSaveDir':
              await window.backend.saveSettingValue(key, value as string);
            break;

            case 'ClipboardWatch':
            case 'DownloadOnCopy':
            case 'ConcurrentDownloads':
            case 'ConcurrentPlaylistDownloads':
              await window.backend.saveSettingBoolValue(key, value as boolean);
            break;

            case 'MaxParrallelDownloads':
              await window.backend.saveSettingValue(key, `${value}`);
            break;
          }
        }

        // if no errors save new config without retrieve it from backend again (we load app state only once when app is launched)
        window.APP_STATE.config = config;
        this._snackbar.open("Settings has been saved");
      } catch(e) {
        this._snackbar.open("An error occured while saving settings");
      }
    });
  }

  trackById(idx: number, entry: Entry): string {
    if(entry.playlist.id) {
      return entry.playlist.id;
    }
    return entry.track.id;
  }

  getBgUrl(entry: Entry): string {
    return `url(${entry.track.thumbnails ? entry.track.thumbnails[4] ? entry.track.thumbnails[4] : entry.track.thumbnails[3] : entry.playlist.thumbnail})`;
  }

  onMouseEnter($event: Event, entry: Entry): void {
    console.log('onMouseEnter', $event, entry);
    ($event.target as HTMLDivElement).classList.toggle('onHover')
  }

  onMouseLeave($event: Event, entry: Entry): void {
    console.log('onMouseLeave', $event, entry);
    ($event.target as HTMLDivElement).classList.toggle('onHover')
  }

  playback(entry: Entry): void {
    this.inPlayback = entry.track;
    this._audioPlayerService.onPlaybackTrack.next(entry.track);
  }

  stop(entry: Entry): void {
    this.inPlayback = null;
    this._audioPlayerService.onStopTrack.next(entry.track);
  }
}
