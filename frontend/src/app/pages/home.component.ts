import {
  ApplicationRef,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  ElementRef,
  Inject,
  NgZone,
  OnDestroy,
  OnInit,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import { FormControl } from '@angular/forms';
import { AudioPlayerService } from 'app/components/audio-player/audio-player.service';
import { Track, Entry, UpdateRelease, ReleaseEventPayload } from '@models';
import { AppConfig, AppState } from '../models/app-state';
import * as Wails from '@wailsapp/runtime';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { SettingsComponent, UpdaterComponent } from 'app/components';
import { debounceTime, distinctUntilChanged, filter } from 'rxjs/operators';
import { MatMenu } from '@angular/material/menu';
import { DOCUMENT } from '@angular/common';
import { SnackbarService } from 'app/services/snackbar.service';
import { MediaMatcher } from '@angular/cdk/layout';
import { OfflinePlaylist } from 'app/models/offline-playlist';
import { AddToPlaylistComponent, CreatePlaylistComponent } from 'app/components/playlist';
import { ConfirmationDialogComponent } from 'app/components/confirmation-dialog/confirmation-dialog.component';
import to from 'await-to-js';

const minMax = (min: number, max: number) => {
  return Math.floor(Math.random() * (max - min)) + min;
}

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class HomeComponent implements OnInit, OnDestroy {
  public searchInput: FormControl;

  public urlInput: FormControl;

  public entries: Entry[] = [];
  public filteredEntries: Entry[] = [];

  public offlinePlaylists: OfflinePlaylist[] = [];

  public onHoverEntry: Entry = null;

  public inPlayback: Track = null;

  public newUpdateInfo: UpdateRelease = null;

  public get inPlaybackTrackId(): string {
    if(!this.inPlayback) {
      return null;
    }
    return this.inPlayback.id;
  }

  public get isClipboardWatchEnabled(): boolean {
    return window.APP_STATE.config.ClipboardWatch;
  }

  public get isSearching(): boolean {
    if(!this.searchNativeInput) {
      return false;
    }
    return this.searchNativeInput.nativeElement === document.activeElement || !!this.searchInput.value;
  }

  @ViewChild('searchNativeInput')
  public searchNativeInput: ElementRef<HTMLInputElement> = null;

  @ViewChild('pasteWrapper')
  public pasteWrapper: ElementRef<HTMLDivElement> = null;

  @ViewChild('pasteInput')
  public pasteInput: ElementRef<HTMLInputElement> = null;

  @ViewChild('menu')
  public matMenu: MatMenu = null;

  public menuIsOpened: boolean = false;

  public mobileQuery: MediaQueryList;

  constructor(
    private _cdr: ChangeDetectorRef,
    private _appRef: ApplicationRef,
    private _ngZone: NgZone,
    @Inject(DOCUMENT) private _document: Document,
    private _dialog: MatDialog,
    private _snackbar: SnackbarService,
    private _media: MediaMatcher,
    private _audioPlayerService: AudioPlayerService
  ) {
    this.searchInput = new FormControl('');
    this.urlInput = new FormControl('');

    this.mobileQuery = this._media.matchMedia('(max-width: 1024px)');
    this._mobileQueryListener = () => this._cdr.detectChanges();
    this.mobileQuery.addEventListener("change", this._mobileQueryListener);
  }

  ngOnInit(): void {
    this.entries = window.APP_STATE.entries;
    this.filteredEntries = this.entries;
    this.offlinePlaylists = window.APP_STATE.offlinePlaylists;

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

    Wails.Events.On("ytd:app:config", (config) => {
      this._ngZone.run(() => {
        window.APP_STATE.config = config;
        this._snackbar.openSuccess("Settings has been saved");
        this._cdr.detectChanges();
      });
    });

    Wails.Events.On("ytd:show:dialog:settings", () => {
      this._ngZone.run(() => {
        this.openSettings();
      })
    });

    Wails.Events.On("ytd:app:update:available", (release: ReleaseEventPayload) => {
      console.log('ON UPDATE', release)
      this._ngZone.run(() => {
        this.newUpdateInfo = UpdateRelease.fromJSON(release);
        this._cdr.detectChanges();
      })
    });

    Wails.Events.On("ytd:app:update:changelog", (release: ReleaseEventPayload) => {
      console.log('ON UPDATE', release)
      this._ngZone.run(() => {
        this.newUpdateInfo = UpdateRelease.fromJSON(release);
        this.openUpdate();
        this._cdr.detectChanges();
      })
    });

    Wails.Events.On("ytd:app:update:apply", (release: ReleaseEventPayload) => {
      console.log('ON UPDATE', release)
      this._ngZone.run(() => {
        this.newUpdateInfo = UpdateRelease.fromJSON(release);
      })
    });

    // players commands
    this._audioPlayerService.onPlayCmdTrack.subscribe(track => {
      this.inPlayback = track;
      this._cdr.detectChanges();
    });

    this._audioPlayerService.onStopCmdTrack.subscribe(track => {
      this.inPlayback = null;
      this._cdr.detectChanges();
    });

    this._audioPlayerService.onPrevTrackCmd.pipe(filter(track => track !== null)).subscribe(currentTrack => {
      const trackIdx = this.entries.findIndex(e => e.track && e.track.id === currentTrack.id);

      if(trackIdx - 1 < 0) {
        this._snackbar.openWarning("Cannot playback prev track")
        return;
      }

      let idx = 1;
      let prevTrack;
      do {
        prevTrack = this.entries[trackIdx - idx].track;
        idx++;
      } while(
        window.wails.System.Platform() !=='darwin' ||
        window.wails.System.Platform() ==='darwin' && !prevTrack.isConvertedToMp3
      )

      this.inPlayback = prevTrack;
      this._audioPlayerService.onPlaybackTrack.next(prevTrack);
    });

    this._audioPlayerService.onNextTrackCmd.pipe(filter(track => track !== null)).subscribe(currentTrack => {
      const trackIdx = this.entries.findIndex(e => e.track && e.track.id === currentTrack.id);
      if(trackIdx + 1 >= this.entries.length) {
        this._snackbar.openWarning("Cannot playback next track")
        return;
      }

      let idx = 1;
      let nextTrack;
      do {
        nextTrack = this.entries[trackIdx + idx].track;
        idx++;
      } while(
        window.wails.System.Platform() !=='darwin' ||
        window.wails.System.Platform() ==='darwin' && !nextTrack.isConvertedToMp3
      )

      this.inPlayback = nextTrack;
      this._audioPlayerService.onPlaybackTrack.next(nextTrack);
    });

    this._audioPlayerService.onShuffleTrackCmd.pipe(filter(track => track !== null)).subscribe(currentTrack => {
      const tracks = this.entries.filter(e => e.type === 'track');
      const idx = minMax(0, tracks.length)
      this.inPlayback = tracks[idx].track;
      this._audioPlayerService.onPlaybackTrack.next(tracks[idx].track);
    });


    this._onSearch();
  }

  clearSearch(event: any): void {
    this.searchNativeInput.nativeElement.focus();
    setTimeout(() => this.searchInput.setValue(''))
  }

  openSettings(): MatDialogRef<SettingsComponent, any> {
    console.log('window.APP_STATE.config', window.APP_STATE.config)
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
              await window.backend.main.AppState.SaveSettingValue(key, value as string);
            break;

            case 'ClipboardWatch':
            case 'DownloadOnCopy':
            case 'ConcurrentDownloads':
            case 'ConcurrentPlaylistDownloads':
            case 'ConvertToMp3':
            case 'CleanWebmFiles':
              await window.backend.main.AppState.SaveSettingBoolValue(key, value as boolean);
            break;

            case 'MaxParrallelDownloads':
              await window.backend.main.AppState.SaveSettingValue(key, `${value}`);
            break;

            case 'Telegram':
              await window.backend.main.AppState.SaveSettingValue(key, JSON.stringify(value));
            break;
          }
        }

        // if no errors save new config without retrieve it from backend again (we load app state only once when app is launched)
        window.APP_STATE.config = config;
        this._snackbar.openSuccess("Settings has been saved");
        Wails.Events.Emit("ytd:app:tray:update")
        this._cdr.detectChanges();
      } catch(e) {
        console.log(e)
        this._snackbar.openError("An error occured while saving settings");
      }
    });

    return dialogRef
  }

  openUpdate(): void {
    const dialogRef = this._dialog.open(UpdaterComponent, {
      autoFocus: false,
      panelClass: ['updater-dialog',  'with-header-dialog'],
      width: '600px',
      maxHeight: '700px',
      data: { release: this.newUpdateInfo, oldVersion: window.APP_STATE.appVersion }
    });

    dialogRef.afterClosed().subscribe(async (result) => {
      if(!result) {
        return;
      }

      const { action } = result;

      switch(action) {
        case 'UpdateAndReplace':
          await window.backend.main.AppState.Update(false);
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
    this.onHoverEntry = entry;
    if(this.menuIsOpened) {
      return;
    }
    ($event.target as HTMLDivElement).classList.add('onHover')
  }

  onMouseLeave($event: Event, entry: Entry): void {
    this.onHoverEntry = null;
    if(this.menuIsOpened) {
      return;
    }
    ($event.target as HTMLDivElement).classList.remove('onHover')
  }

  playback(entry: Entry): void {
    if(window.wails.System.Platform() ==='darwin' && !entry.track.isConvertedToMp3) {
      const ref = this._snackbar.openWarning("Cannot playback on MacOs, you should enable \"Convert to mp3\" option", "Open settings")
      ref.onAction().subscribe(() => this.openSettings());
      return;
    }

    this.inPlayback = entry.track;
    this._audioPlayerService.onPlaybackTrack.next(entry.track);
  }

  stop(entry: Entry): void {
    this.inPlayback = null;
    this._audioPlayerService.onStopTrack.next(entry.track);
  }

  menuOpened(): void {
    this.menuIsOpened = true;
  }

  menuClosed(): void {
    this.menuIsOpened = false;
    const onHoveredEntry = this._document.querySelector('.entry.onHover');
    if(onHoveredEntry && !this.onHoverEntry) {
      onHoveredEntry.classList.toggle('onHover')
    }
  }

  async addToDownload(): Promise<void> {
    const url = this.urlInput.value;
    if(!url) {
      return;
    }

    const isSupported = await window.backend.main.AppState.IsSupportedUrl(url);
    if(!isSupported) {
      this._snackbar.openWarning('Unsupported url');
      return
    }

    try {
      await window.backend.main.AppState.AddToDownload(url, false)
      this.urlInput.setValue('');
      this.pasteInput.nativeElement.blur();
      this.pasteWrapper.nativeElement.classList.remove('focused');
      this._snackbar.openSuccess('Track added to download');
    } catch(e) {
      console.log(e);
      this._snackbar.openError('Error while adding track to downloads');
    }
  }

  async startDownload(entry: Entry): Promise<void> {
    try {
      await window.backend.main.AppState.StartDownload(entry);
      this._snackbar.openSuccess("Started downloading");
    } catch(e) {
      this._snackbar.openError(e);
    }
  }

  addToPlaylist(entry: Entry): void {
    const foundInPlaylists = this.offlinePlaylists.filter(op => op.tracksIds.indexOf(entry.track.id) > -1);
    const dialogRef = this._dialog.open(AddToPlaylistComponent, {
      autoFocus: false,
      panelClass: ['add-to-playlist-dialog',  'with-header-dialog'],
      maxWidth: '500px',
      maxHeight: '500px',
      data: { entry, playlists: this.offlinePlaylists, foundInPlaylists }
    });

    dialogRef.afterClosed().subscribe(async (result) => {
      if(!result) {
        return;
      }

      const { action, selectedPlaylists } = result;
      if(!action) {
        console.log('selectedPlaylists', selectedPlaylists)
        const playlists: OfflinePlaylist[] = selectedPlaylists.map((uuid: string) => this.offlinePlaylists.find(op => op.uuid === uuid))
        playlists.forEach(p => p.tracksIds.push(entry.track.id));
        console.log(playlists)
        const [err, added] = await to(window.backend.main.OfflinePlaylistService.AddTrackToPlaylist(playlists))

        if(err) {
          this._snackbar.openError("Error while adding track to playlist");
          return;
        }

        return;
      }

      switch(action) {
        case 'createNew':
          this.createPlaylist();
      }

    });
  }

  createPlaylist(): void {
    const dialogRef = this._dialog.open(CreatePlaylistComponent, {
      autoFocus: false,
      panelClass: ['create-playlist-dialog',  'with-header-dialog'],
      width: '300px',
      maxHeight: '500px',
      data: { playlists: this.offlinePlaylists }
    });

    dialogRef.afterClosed().subscribe(async (result) => {
      if(!result) {
        return;
      }

      const { playlist } = result;

      const [err, createdPlaylist] = await to(window.backend.main.OfflinePlaylistService.CreateNewPlaylist(playlist.name))
      console.log('createdPlaylist', createdPlaylist)
      if(err) {
        this._snackbar.openError("Error while creating playlist");
        return
      }

      this._snackbar.openSuccess("Playlist created");
      this.offlinePlaylists.push(createdPlaylist);
      this._cdr.detectChanges();
    });
  }

  removePlaylist(playlist: OfflinePlaylist): void {
    const dialogRef = this._dialog.open(ConfirmationDialogComponent, {
      autoFocus: false,
      panelClass: ['with-header-dialog'],
      width: '300px',
      data: {
        title: 'Delete playlist',
        text: `Are you sure you would to remove <strong>${playlist.name}</strong> playlist?`
       }
    });

    dialogRef.afterClosed().subscribe(async (result) => {
      if(!result) {
        return;
      }

      // call remove from backend and update list
      console.log("remove playlist", playlist)
      const [err, isRemoved] = await to(window.backend.main.OfflinePlaylistService.RemovePlaylist(playlist.uuid))
      if(err || !isRemoved) {
        console.log(err, isRemoved)
        this._snackbar.openError("Error while removing playlist");
        return
      }

      this._snackbar.openSuccess("Playlist has been deleted");
      this.offlinePlaylists = this.offlinePlaylists.filter(p => p.uuid !== playlist.uuid)
      this._cdr.detectChanges();

    });
  }

  async exportPlaylist(playlist: OfflinePlaylist): Promise<void> {
    const [err, dir] = await to(window.backend.main.OfflinePlaylistService.ExportPlaylist(playlist.uuid))
    console.log('exportPlaylist', err, dir)
  }

  async remove(entry: Entry, i: number): Promise<void> {
    try {
      await window.backend.main.AppState.RemoveEntry(entry);
      this._snackbar.openSuccess(`${entry.type} has been removed`);
      const idx = this.entries.findIndex(e => {
        if(entry.type === 'playlist') {
          return e.playlist.id === entry.playlist.id;
        }
        return e.track.id === entry.track.id;
      })
      this.entries.splice(idx, 1);
      this.filteredEntries = this.entries;
      //this._cdr.detectChanges();
    } catch(e) {
      this._snackbar.openError(`Cannot delete`);
    }
  }

  private _onSearch(): void {
    this.searchInput.valueChanges
    .pipe(
      debounceTime(300),
      distinctUntilChanged(),
    )
    .subscribe((searchText: string) => {
      if(!!searchText === false) {
        this.filteredEntries = this.entries;
        this._cdr.detectChanges();
        return;
      }

      searchText = searchText.toLowerCase();
      this.filteredEntries = this.entries.filter(e => {
        if(e.type === 'playlist') {
          return e.playlist.name.toLowerCase().includes(searchText);
        }
        return e.track.author.toLowerCase().includes(searchText) || e.track.name.toLowerCase().includes(searchText)
      });
      this._cdr.detectChanges();
    });
  }

  private _mobileQueryListener: () => void;

  ngOnDestroy(): void {
    this.mobileQuery.removeEventListener("change", this._mobileQueryListener);
  }
}
