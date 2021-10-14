import {
  ApplicationRef,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  NgZone,
  OnDestroy,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AudioPlayerService } from 'app/components/audio-player/audio-player.service';
import { Track, Entry, UpdateRelease, ReleaseEventPayload } from '@models';
import wailsapp__runtime, * as Wails from '@wailsapp/runtime';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { debounceTime, distinctUntilChanged, filter, takeUntil } from 'rxjs/operators';
import { MatMenu } from '@angular/material/menu';
import { DOCUMENT } from '@angular/common';
import { SnackbarService } from 'app/services/snackbar.service';
import { OfflinePlaylist } from 'app/models/offline-playlist';
import { ConfirmationDialogComponent } from 'app/components/confirmation-dialog/confirmation-dialog.component';
import to from 'await-to-js';
import { Subject, Subscription } from 'rxjs';
import { minMax } from 'app/common/fn';

@Component({
  selector: 'app-offline-playlist',
  templateUrl: './offline-playlist.component.html',
  styleUrls: ['./offline-playlist.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class OfflinePlaylistComponent implements OnInit, OnDestroy {
  public playlist: OfflinePlaylist = null;

  public tracks: Track[] = [];

  public inPlayback: Track = null;

  public inPauseTrack: Track = null;

  public get inPlaybackTrackId(): string {
    if(!this.inPlayback) {
      return null;
    }
    return this.inPlayback.id;
  }

  private _unsubscribe: Subject<any> = new Subject();

  constructor(
    private _cdr: ChangeDetectorRef,
    private _appRef: ApplicationRef,
    private _ngZone: NgZone,
    @Inject(DOCUMENT) private _document: Document,
    private _router: Router,
    private _route: ActivatedRoute,
    private _dialog: MatDialog,
    private _snackbar: SnackbarService,
    private _audioPlayerService: AudioPlayerService
  ) {
    const { playlist, tracks } = this._router.getCurrentNavigation().extras.state;
    this.playlist = playlist;
    this.tracks = tracks;
    this.inPlayback = this.tracks[0];
  }

  ngOnInit(): void {
    const audioPlayer = this._document.querySelector('audio-player .player');
    audioPlayer.setAttribute("view", "playlist");

    // players commands
    this._audioPlayerService.onPlayCmdTrack.pipe(filter(track => track !== null), takeUntil(this._unsubscribe)).subscribe(track => {
      this.inPlayback = track;
      this.inPauseTrack = null;
      this._cdr.detectChanges();
    });

    this._audioPlayerService.onStopCmdTrack.pipe(filter(track => track !== null), takeUntil(this._unsubscribe)).subscribe(track => {
      console.log("ON STOPPPPP", track)
      this.inPlayback = null;
      this.inPauseTrack = track;
      this._cdr.detectChanges();
    });

    this._audioPlayerService.onPrevTrackCmd.pipe(filter(track => track !== null), takeUntil(this._unsubscribe)).subscribe(currentTrack => {
      const trackIdx = this.tracks.findIndex(track => track.id === currentTrack.id);

      if(trackIdx - 1 < 0) {
        this._snackbar.openWarning("PLAYER.CANNOT_PLAY_PREV")
        return;
      }

      let idx = 1;
      let prevTrack;
      do {
        prevTrack = this.tracks[trackIdx - idx];
        idx++;
      } while(
        window.wails.System.Platform() !=='darwin' ||
        window.wails.System.Platform() ==='darwin' && !prevTrack.isConvertedToMp3
      )

      this.inPlayback = prevTrack;
      this._audioPlayerService.onPlaybackTrack.next(prevTrack);
      this._cdr.detectChanges();
    });

    this._audioPlayerService.onNextTrackCmd.pipe(filter(track => track !== null), takeUntil(this._unsubscribe)).subscribe(currentTrack => {
      const trackIdx = this.tracks.findIndex(track => track.id === currentTrack.id);
      if(trackIdx + 1 >= this.tracks.length) {
        this._snackbar.openWarning("PLAYER.CANNOT_PLAY_NEXT")
        return;
      }

      let idx = 1;
      let nextTrack;
      do {
        nextTrack = this.tracks[trackIdx + idx];
        idx++;
      } while(
        window.wails.System.Platform() !=='darwin' ||
        window.wails.System.Platform() ==='darwin' && !nextTrack.isConvertedToMp3
      )

      this.inPlayback = nextTrack;
      this._audioPlayerService.onPlaybackTrack.next(nextTrack);
      this._cdr.detectChanges();
    });

    this._audioPlayerService.onShuffleTrackCmd.pipe(filter(track => track !== null), takeUntil(this._unsubscribe)).subscribe(currentTrack => {
      const idx = minMax(0, this.tracks.length)
      this.inPlayback = this.tracks[idx];
      this._audioPlayerService.onPlaybackTrack.next(this.tracks[idx]);
      this._cdr.detectChanges();
    });

    this._audioPlayerService.onClosePlayer.pipe(filter(v => v !== null), takeUntil(this._unsubscribe)).subscribe(() => {
      this.close();
    })
  }

  playback(track: Track): void {
    if(window.wails.System.Platform() ==='darwin' && !track.isConvertedToMp3) {
      const ref = this._snackbar.openWarning("SETTINGS.MACOS_MISSING_CONVERT_TO_MP3_OPTIONS")
      return;
    }

    this.inPlayback = track;
    this.inPauseTrack = null;
    this._audioPlayerService.onPlaybackTrack.next(track);
  }

  stop(track?: Track): void {
    this._audioPlayerService.onStopTrack.next(this.inPlayback);
    this.inPlayback = null;
    this.inPauseTrack = track;
  }

  async removeTrackFromPlaylist(trackId: string, playlist: OfflinePlaylist): Promise<void> {
    const [err, updatedPlaylist] = await to(window.backend.main.OfflinePlaylistService.RemoveTrackFromPlaylist(trackId, playlist))
    if(err) {
      console.log(err, updatedPlaylist)
      this._snackbar.openError("PLAYLISTS.REMOVE_TRACK.KO");
      return
    }

    this._snackbar.openSuccess("PLAYLISTS.REMOVE_TRACK.OK");
    // notify go backend that playlist has been deleted, this updates backend state.offlinePlaylists
    // and then emit ytd:offline:playlists back to fe with offlinePlaylists to sync state between each other
    Wails.Events.Emit("ytd:offline:playlists:removedTrack");
  }

  trackById(idx: number, track: Track): string {
    return track ? track.id : undefined;
  }

  close(): void {
    this._router.navigate(['home']);
  }

  ngOnDestroy(): void {
    const audioPlayer = this._document.querySelector('audio-player .player');
    audioPlayer.removeAttribute("view");

    this._unsubscribe.next();
    this._unsubscribe.complete();
  }
}
