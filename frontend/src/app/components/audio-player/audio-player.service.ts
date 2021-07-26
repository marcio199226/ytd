import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { Track } from '@models';


@Injectable({providedIn: 'root'})
export class AudioPlayerService {

  public onPlaybackTrack: BehaviorSubject<Track>;

  public onTrackChanges: BehaviorSubject<string>;

  constructor() {
    this.onPlaybackTrack  = new BehaviorSubject(null);
    this.onTrackChanges   = new BehaviorSubject(null);
  }
}