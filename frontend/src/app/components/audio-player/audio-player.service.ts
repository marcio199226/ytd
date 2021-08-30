import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { Track } from '@models';


@Injectable({providedIn: 'root'})
export class AudioPlayerService {
  public action:  'play' | 'stop';

  public track: Track = null;

  public onPlaybackTrack: BehaviorSubject<Track>;

  public onPlayCmdTrack: BehaviorSubject<Track>;

  public onStopTrack: BehaviorSubject<Track>;

  public onStopCmdTrack: BehaviorSubject<Track>;

  public onPrevTrack: BehaviorSubject<Track>;

  public onPrevTrackCmd: BehaviorSubject<Track>;

  public onNextTrack: BehaviorSubject<Track>;

  public onNextTrackCmd: BehaviorSubject<Track>;

  public onShuffleTrack: BehaviorSubject<Track>;

  public onShuffleTrackCmd: BehaviorSubject<Track>;

  public get trackId(): string {
    if(!this.track) {
      return null;
    }
    return this.track.id;
  }

  constructor() {
    this.onPlaybackTrack    = new BehaviorSubject(null);
    this.onPlayCmdTrack     = new BehaviorSubject(null);
    this.onStopTrack        = new BehaviorSubject(null);
    this.onStopCmdTrack     = new BehaviorSubject(null);
    this.onPrevTrack        = new BehaviorSubject(null);
    this.onPrevTrackCmd     = new BehaviorSubject(null);
    this.onNextTrack        = new BehaviorSubject(null);
    this.onNextTrackCmd     = new BehaviorSubject(null);
    this.onShuffleTrack     = new BehaviorSubject(null);
    this.onShuffleTrackCmd  = new BehaviorSubject(null);

    this.onPlaybackTrack.subscribe(track => {
      this.action = 'play';
      this.track = track;
    });

    this.onStopTrack.subscribe(track => {
      this.action = 'stop';
    });
  }
}
