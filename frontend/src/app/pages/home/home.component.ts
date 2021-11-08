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
  TemplateRef,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import { FormControl } from '@angular/forms';
import { AudioPlayerService } from 'app/components/audio-player/audio-player.service';
import { Track, Entry, UpdateRelease, ReleaseEventPayload, NgrokState } from '@models';
import { AppConfig, AppState } from '../../models/app-state';
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
import { ActivatedRoute, Router } from '@angular/router';
import { MatDrawer } from '@angular/material/sidenav';
import { minMax } from 'app/common/fn';
import { OfflinePlaylistComponent } from '../offline-playlist/offline-playlist.component';
import { LoaderService } from 'app/services/loader.service';
import { TranslateService } from '@ngx-translate/core';
import { fakeAsync } from '@angular/core/testing';

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

  public ngrok: NgrokState = null;

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

  public get isOpenedOfflinePlaylist(): boolean {
    return this._route.children.length !== 0;
  }

  public get isPlaylistInPlayback(): boolean {
    if(!this._offlinePlaylistComponent) {
      return false;
    }
    return this._offlinePlaylistComponent.inPlaybackTrackId !== null;
  }

  public get appVersion(): string {
    return window.APP_STATE.appVersion;
  }

  public get isNgrokRunning(): boolean {
    if(!this.ngrok) {
      return false;
    }
    return this.ngrok.status === 'running';
  }

  @ViewChild('searchNativeInput')
  public searchNativeInput: ElementRef<HTMLInputElement> = null;

  @ViewChild('pasteWrapper')
  public pasteWrapper: ElementRef<HTMLDivElement> = null;

  @ViewChild('pasteInput')
  public pasteInput: ElementRef<HTMLInputElement> = null;

  @ViewChild('menu')
  public matMenu: MatMenu = null;

  @ViewChild('drawer')
  public drawer: MatDrawer = null;

  public menuIsOpened: boolean = false;

  public mobileQuery: MediaQueryList;

  private _mobileQueryListener: (event: MediaQueryListEvent) => void;

  private _offlinePlaylistComponent: OfflinePlaylistComponent = null

  constructor(
    private _cdr: ChangeDetectorRef,
    private _appRef: ApplicationRef,
    private _ngZone: NgZone,
    @Inject(DOCUMENT) private _document: Document,
    private _router: Router,
    private _route: ActivatedRoute,
    private _dialog: MatDialog,
    private _snackbar: SnackbarService,
    private _loader: LoaderService,
    private _media: MediaMatcher,
    private _trans: TranslateService,
    private _audioPlayerService: AudioPlayerService
  ) {
    this.searchInput = new FormControl('');
    this.urlInput = new FormControl('');

    this.mobileQuery = this._media.matchMedia('(max-width: 1024px)');
    if(!this.mobileQuery.matches) {
      this._document.body.classList.add('fullscreen');
    }
    this._mobileQueryListener = (e) => {
      if(!e.matches) {
        this._document.body.classList.add('fullscreen');
      } else {
        this._document.body.classList.remove('fullscreen');
      }
      this._cdr.detectChanges();
    }
    this.mobileQuery.addEventListener("change", this._mobileQueryListener);
  }

  ngOnInit(): void {
    this.entries = window.APP_STATE.entries;
    this.filteredEntries = this.entries;
    this.offlinePlaylists = window.APP_STATE.offlinePlaylists;

    // subscribe to events coming from backend
    Wails.Events.On("ytd:track", (payload: Entry) => {
      this._ngZone.run(() => {
        console.log("ytd:track", payload)
        const entry = this.entries.find(e => e.track.id === payload.track.id);
        if(entry) {
          entry.track = payload.track;
        } else {
          this.entries.unshift(payload);
        }
        this._cdr.detectChanges();
      });
    });

    Wails.Events.On("ytd:track:progress", ({ id, progress }) => {
      const entry = this.entries.find(e => e.track.id === id);
      entry.track.downloadProgress = progress;
      this._cdr.detectChanges();
    });

    Wails.Events.On("ytd:playlist", payload => console.log(payload));

    Wails.Events.On("ytd:offline:playlists", offlinePlaylists => {
      this._ngZone.run(() => {
        console.log("UPDATE offlinePlaylists", offlinePlaylists)
        this.offlinePlaylists = offlinePlaylists;
        this._cdr.detectChanges();
      });
    });

    Wails.Events.On("ytd:app:config", (config) => {
      this._ngZone.run(() => {
        window.APP_STATE.config = config;
        this._snackbar.openSuccess("SETTINGS.SAVED");
        this._cdr.detectChanges();
      });
    });

    Wails.Events.On("ytd:show:dialog:settings", () => {
      this._ngZone.run(() => {
        this.openSettings();
      })
    });

    Wails.Events.On("ytd:app:update:available", (release: ReleaseEventPayload) => {
      this._ngZone.run(() => {
        this.newUpdateInfo = UpdateRelease.fromJSON(release);
        this._cdr.detectChanges();
      })
    });

    Wails.Events.On("ytd:app:update:changelog", (release: ReleaseEventPayload) => {
      this._ngZone.run(() => {
        this.newUpdateInfo = UpdateRelease.fromJSON(release);
        this.openUpdate();
        this._cdr.detectChanges();
      })
    });

    Wails.Events.On("ytd:app:update:apply", (release: ReleaseEventPayload) => {
      this._ngZone.run(() => {
        this.newUpdateInfo = UpdateRelease.fromJSON(release);
      })
    });

    Wails.Events.On("ytd:ngrok", (ngrok: NgrokState) => {
      this._ngZone.run(() => {
        this.ngrok = ngrok;
        console.log("ytd:ngrok", ngrok)
        this._cdr.detectChanges();
      })
    });

    // players commands
    this._audioPlayerService.onPlayCmdTrack.subscribe(track => {
      this.inPlayback = track;
      this._cdr.detectChanges();
    });

    this._audioPlayerService.onStopCmdTrack.subscribe(track => {
      console.log("ON STOPPP FROM HOMEEE", track)
      this.inPlayback = null;
      this._cdr.detectChanges();
    });

    this._audioPlayerService.onPrevTrackCmd.pipe(filter(track => track !== null)).subscribe(currentTrack => {
      const trackIdx = this.entries.findIndex(e => e.track && e.track.id === currentTrack.id);

      if(trackIdx - 1 < 0) {
        this._snackbar.openWarning("PLAYER.CANNOT_PLAY_PREV")
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
        this._snackbar.openWarning("PLAYER.CANNOT_PLAY_NEXT")
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

  onRouterOutletActivate(cmp: any): void {
    this._offlinePlaylistComponent = cmp;
    this._cdr.detectChanges();
  }

  onRouterOutletDeactivate(e: any): void {
    this._offlinePlaylistComponent = null;
    this._cdr.detectChanges();
  }

  clearSearch(event: any): void {
    this.searchNativeInput.nativeElement.focus();
    setTimeout(() => this.searchInput.setValue(''))
  }

  openSettings(tab?: string): MatDialogRef<SettingsComponent, any> {
    const dialogRef = this._dialog.open(SettingsComponent, {
      panelClass: ['settings-dialog',  'with-header-dialog'],
      width: '700px',
      maxHeight: '700px',
      data: { tab, config: window.APP_STATE.config, isNgrokRunning: this.isNgrokRunning, ngrok: this.ngrok }
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

            // struct configs
            case 'Telegram':
            case 'PublicServer':
              await window.backend.main.AppState.SaveSettingValue(key, JSON.stringify(value));
            break;
          }
        }

        // if no errors save new config without retrieve it from backend again (we load app state only once when app is launched)
        window.APP_STATE.config = config;
        this._snackbar.openSuccess("SETTINGS.SAVED");
        Wails.Events.Emit("ytd:app:tray:update")
        this._cdr.detectChanges();
      } catch(e) {
        console.log(e)
        this._snackbar.openError("SETTINGS.KO");
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
      const ref = this._snackbar.openWarning("SETTINGS.MACOS_MISSING_CONVERT_TO_MP3_OPTIONS", "SETTINGS.OPEN_DIALOG_BTN")
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
      this._snackbar.openWarning("SETTINGS.MANUAL_MODE.UNSUPPORTED_URL");
      return
    }

    try {
      await window.backend.main.AppState.AddToDownload(url, false)
      this.urlInput.setValue('');
      this.pasteInput.nativeElement.blur();
      this.pasteWrapper.nativeElement.classList.remove("focused");
      this._snackbar.openSuccess("SETTINGS.MANUAL_MODE.TRACK_ADDED");
    } catch(e) {
      console.log(e);
      this._snackbar.openError("SETTINGS.MANUAL_MODE.TRACK_ADD_KO");
    }
  }

  async startDownload(entry: Entry): Promise<void> {
    try {
      await window.backend.main.AppState.StartDownload(entry);
      this._snackbar.openSuccess("HOME.STARTED_DOWNLOAD_TRACK");
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
      console.log(selectedPlaylists)
      if(!action) {
        const playlists: OfflinePlaylist[] = selectedPlaylists.map((uuid: string) => this.offlinePlaylists.find(op => op.uuid === uuid))
        playlists.forEach(p => p.tracksIds.push(entry.track.id));
        const [err, added] = await to(window.backend.main.OfflinePlaylistService.AddTrackToPlaylist(playlists))

        if(err) {
          this._snackbar.openError("PLAYLISTS.ADD_TRACK_KO");
          return;
        }

        this._snackbar.openSuccess("PLAYLISTS.ADD_TRACK_OK");
      }

      switch(action) {
        case 'createNew':
          this.createPlaylist();
      }

    });
  }

  createPlaylist(): void {
    const dialogRef = this._dialog.open(CreatePlaylistComponent, {
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
      if(err) {
        this._snackbar.openError("PLAYLISTS.CREATED_KO");
        return
      }

      this._snackbar.openSuccess("PLAYLISTS.CREATED_OK");
      // notify go backend that new playlist has been created, this updates backend state.offlinePlaylists
      // and then emit ytd:offline:playlists back to fe with offlinePlaylists to sync state between each other
      Wails.Events.Emit("ytd:offline:playlists:created");
    });
  }

  removePlaylist(playlist: OfflinePlaylist): void {
    const dialogRef = this._dialog.open(ConfirmationDialogComponent, {
      autoFocus: false,
      panelClass: ['with-header-dialog'],
      width: '300px',
      data: {
        title: "PLAYLISTS.REMOVE_DIALOG.TITLE",
        text: this._trans.instant("PLAYLISTS.REMOVE_DIALOG.DESC", { name: playlist.name })
       }
    });

    dialogRef.afterClosed().subscribe(async (result) => {
      if(!result) {
        return;
      }

      // call remove from backend and update list
      const [err, isRemoved] = await to(window.backend.main.OfflinePlaylistService.RemovePlaylist(playlist.uuid))
      if(err || !isRemoved) {
        console.log(err, isRemoved)
        this._snackbar.openError("PLAYLISTS.REMOVE_DIALOG.KO");
        return
      }

      this._snackbar.openSuccess("PLAYLISTS.REMOVE_DIALOG.OK");
      // notify go backend that playlist has been deleted, this updates backend state.offlinePlaylists
      // and then emit ytd:offline:playlists back to fe with offlinePlaylists to sync state between each other
      Wails.Events.Emit("ytd:offline:playlists:removed");
    });
  }

  async exportPlaylist(playlist: OfflinePlaylist): Promise<void> {
    if(!playlist.tracksIds.length) {
      this._snackbar.openWarning("HOME.PLAYLISTS.CANNOT_EXPORT_EMPTY");
      return;
    }

    const [err, dirResult] = await to<string, Error>(window.wails.Dialog.Open({
      DefaultDirectory:     playlist.name,
      AllowFiles:           false,
      CanCreateDirectories: true,
      AllowDirectories:     true,
      Title:                this._trans.instant("HOME.PLAYLISTS.CHOOSE_EXPORT_DIRECTORY"),
    }))
    const dir: string[] = JSON.parse(dirResult);

    if(err) {
      console.log('err, dirResult', err, dirResult)
      this._snackbar.openError("HOME.PLAYLISTS.CHOOSE_EXPORT_DIRECTORY_KO", null, null, { error: err });
      return;
    }

    if(!dir.length) {
      return;
    }

    this._loader.show("HOME.PLAYLISTS.LOADER.EXPORTING");
    const [exportErr, exportResult] = await to(window.backend.main.OfflinePlaylistService.ExportPlaylist(playlist.uuid, dir[0]))
    this._loader.hide();

    if(exportErr) {
      console.log('exportPlaylist go call results in error', exportErr, exportResult);
      this._snackbar.openError("HOME.PLAYLISTS.EXPORTED_KO", null, null, { error: exportErr });
      return;
    }

    this._snackbar.openSuccess("HOME.PLAYLISTS.EXPORTED_OK");
  }

  async playbackPlaylist(playlist: OfflinePlaylist): Promise<void> {
    if(!playlist.tracksIds.length) {
      this._snackbar.openWarning("PLAYLISTS.CANNOT_PLAYBACK_EMPTY");
      return;
    }
    if(this.mobileQuery.matches) {
      await this.drawer.close();
    }
    const tracks = this.entries.filter(e => e.type === 'track' && playlist.tracksIds.indexOf(e.track.id) > -1).map(e => e.track);
    this._router.navigate(['playlist', playlist.uuid], { relativeTo: this._route,  state: { playlist, tracks } });
    this._audioPlayerService.onPlaybackTrack.next(tracks[0]);
    this._cdr.detectChanges();
  }

  stopPlaylist(playlisy: OfflinePlaylist): void {
    this._offlinePlaylistComponent.stop();
    this._cdr.detectChanges();
  }

  isPlaylistDetailOpen(playlist: OfflinePlaylist): boolean {
    if(!this._route.children.length) {
      return false;
    }
    return this._route.children[0].snapshot.paramMap.get('playlist') === playlist.uuid;
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

  async remove(entry: Entry, i: number): Promise<void> {
    try {
      await window.backend.main.AppState.RemoveEntry(entry);
      this._snackbar.openSuccess(`HOME.REMOVED.OK_${entry.type.toUpperCase()}`);
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
      this._snackbar.openError(`HOME.REMOVED.KO_${entry.type.toUpperCase()}`);
    }
  }

  async openInYt(url: string): Promise<void> {
    await window.backend.main.AppState.OpenUrl(url);
  }

  async openGithubProfile(): Promise<void> {
    await window.backend.main.AppState.OpenUrl('https://github.com/marcio199226/ytd');
  }

  async openReleasePage(): Promise<void> {
    await window.backend.main.AppState.OpenUrl(`https://github.com/marcio199226/ytd/releases/tag/${this.appVersion}`);
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

  ngOnDestroy(): void {
    this.mobileQuery.removeEventListener("change", this._mobileQueryListener);
  }
}
