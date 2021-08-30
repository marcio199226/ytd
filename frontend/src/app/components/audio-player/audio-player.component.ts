import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  HostBinding,
  Inject,
  Input,
  OnInit,
  Output,
  Renderer2,
  ViewEncapsulation,
} from '@angular/core';
import { filter } from 'rxjs/operators'
import { Track } from '@models';
import { AudioPlayerService } from './audio-player.service';
import { SnackbarService } from 'app/services/snackbar.service';
import { MatSliderChange } from '@angular/material/slider';
import { DOCUMENT } from '@angular/common';

@Component({
  selector: 'audio-player',
  templateUrl: './audio-player.component.html',
  styleUrls: ['./audio-player.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class AudioPlayerComponent implements OnInit {

  public track: Track = null;

  public audio: HTMLAudioElement = null;

  public volume: number = 0.5;

  public isPlaying: Boolean = false;

  public canPlay: Boolean = false;

  public duration: number = null;

  public elapsedTime: string = null;

  public elapsedTimeProgress: number = null;

  public get trackCover(): string {
    return this.track.thumbnails[3];
  }

  @Output() public onPlayAudio: any = new EventEmitter<any>();

  constructor(
    private _cdr: ChangeDetectorRef,
    @Inject(DOCUMENT) private document: Document,
    private _snackbar: SnackbarService,
    private _audioPlayerService: AudioPlayerService
  ) {}

  ngOnInit(): void {
    this._audioPlayerService.onPlaybackTrack.pipe(filter(track => track !== null)).subscribe(track => {
      if(this.track && this.track.id === track.id) {
        this.play();
        return;
      }

      if(this.audio) {
        this.audio.pause();
        this.audio = null;
      }

      this.track = track;
      this._play(track);
    })

    this._audioPlayerService.onStopTrack.pipe(filter(track => track !== null)).subscribe(track => {
      this.stop();
    })
  }

  private _play(track: Track): void {
    const src = track.isConvertedToMp3 ? `http://localhost:8080/tracks/youtube/${track.id}.mp3` : `http://localhost:8080/tracks/youtube/${track.id}.webm`;
    this.audio = new Audio(src);
    this.audio.volume = this.volume;

    this.audio.ontimeupdate = (e) => {
      const s = parseInt((this.audio.currentTime % 60).toString(), 10);
      const m = parseInt(((this.audio.currentTime / 60) % 60).toString(), 10);
      this.duration = this.audio.duration;
      this.elapsedTimeProgress = +(
        (+this.audio.currentTime.toFixed(1) / +this.audio.duration.toFixed(1)) *
        100
      ).toFixed(0);
      this.elapsedTime = s < 10 ? m + ':0' + s : m + ':' + s;
      this._cdr.detectChanges();
    };

    this.audio.onended = (e) => {
      this.isPlaying = false;
      this.elapsedTimeProgress = 0;
      this._cdr.detectChanges();
    };

    this.audio.onerror = (e) => {
      console.log("Track playback error", e)
      this.track = null;
      this._snackbar.openError("Cannot playback probably track's file does not exists");
    };

    this.play();
  }

  reloadAudio(): void {
    this.audio.load();
  }

  ngDoCheck(): void {
    if (this.audio && this.audio.readyState === 4 && !this.canPlay) {
      this.canPlay = true;
      const s = parseInt((this.audio.duration % 60).toString(), 10);
      const m = parseInt(((this.audio.duration / 60) % 60).toString(), 10);
      this.elapsedTime = s < 10 ? m + ':0' + s : m + ':' + s;
      this._cdr.detectChanges();
    }
  }

  play(): void {
    this.document.body.classList.add('player-visible');
    this.onPlayAudio.emit();
    this.isPlaying = true;
    this.audio.play();
    this._audioPlayerService.onPlayCmdTrack.next(this.track);
  }

  stop(): void {
    this.isPlaying = false;
    this.audio.pause();
    this._audioPlayerService.onStopCmdTrack.next(this.track);

    this._cdr.detectChanges();
  }

  changeVolume(event: MatSliderChange): void {
    this.volume = event.value;
    this.audio.volume = event.value;
  }

  prev(): void {
    this._audioPlayerService.onPrevTrackCmd.next(this.track);
  }

  next(): void {
    this._audioPlayerService.onNextTrackCmd.next(this.track);
  }

  replay(): void {
    this.audio.currentTime = 0;
  }

  shuffle(): void {
    this._audioPlayerService.onShuffleTrackCmd.next(this.track);
  }

  isReady(): Boolean {
    return this.audio && this.audio.readyState === 4;
  }

  closePlayer(): void {
    this.audio.pause();
    this._audioPlayerService.onStopCmdTrack.next(this.track);
    this.audio = null;
    this.track = null;
    this.document.body.classList.remove('player-visible');
  }

  ngOnDestroy(): void {
    this.audio.pause();
    this.audio = null;
  }
}
