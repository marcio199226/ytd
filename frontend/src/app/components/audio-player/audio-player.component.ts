import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  HostBinding,
  Input,
  OnInit,
  Output,
  ViewEncapsulation,
} from '@angular/core';
import { filter } from 'rxjs/operators'
import { Track } from '@models';
import { AudioPlayerService } from './audio-player.service';
import { MatSnackBar } from '@angular/material/snack-bar';

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
    private _snackbar: MatSnackBar,
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
    const src = `http://localhost:8080/youtube/${track.id}.webm`;
    this.audio = new Audio(src);

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
      this.track = null;
      this._snackbar.open("Cannot playback probably track's file does not exists");
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

  isReady(): Boolean {
    return this.audio && this.audio.readyState === 4;
  }

  ngOnDestroy(): void {
    this.audio.pause();
    this.audio = null;
  }
}
