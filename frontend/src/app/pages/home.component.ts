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
import { AppState } from '../models/app-state';

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

  public get isAudioPlayerPlaying(): boolean {
    console.log(this._audioPlayerService)
    return this._audioPlayerService.action === 'play';
  }

  constructor(
    private _cdr: ChangeDetectorRef,
    private _audioPlayerService: AudioPlayerService
  ) {
    this.searchInput = new FormControl('');
  }

  ngOnInit(): void {
    this.entries = (window.APP_STATE as AppState).entries;
    console.log(this.entries)

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
