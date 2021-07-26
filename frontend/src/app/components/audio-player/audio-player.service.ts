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

  public onTrackChanges: BehaviorSubject<string>;

  public get trackId(): string {
    if(!this.track) {
      return null;
    }
    return this.track.id;
  }

  constructor() {
    this.onPlaybackTrack  = new BehaviorSubject(null);
    this.onPlayCmdTrack   = new BehaviorSubject(null);
    this.onStopTrack      = new BehaviorSubject(null);
    this.onStopCmdTrack   = new BehaviorSubject(null);
    this.onTrackChanges   = new BehaviorSubject(null);

    this.onPlaybackTrack.subscribe(track => {
      this.action = 'play';
      this.track = track;
    });

    this.onStopTrack.subscribe(track => {
      this.action = 'stop';
    });
  }
}
