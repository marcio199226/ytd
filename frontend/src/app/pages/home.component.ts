import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnInit,
  ViewEncapsulation,
} from '@angular/core';
import { FormControl } from '@angular/forms';
import { AudioPlayerService } from 'app/components/audio-player/audio-player.service';
import { AppState } from '../models/app-state';
import { Entry } from '../models/entry';

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

  constructor(
    private _cdr: ChangeDetectorRef,
    private _audioPlayerService: AudioPlayerService
  ) {
    this.searchInput = new FormControl('');
  }

  ngOnInit(): void {
    this.entries = (window.APP_STATE as AppState).entries;
    console.log(this.entries)
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
    this._audioPlayerService.onPlaybackTrack.next(entry.track);
  }
}
