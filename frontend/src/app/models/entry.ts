import { Playlist } from "./playlist";
import { Track } from "./track";

export interface Entry {
    playlist: Playlist;
    source: 'youtube' | 'spotify';
    track: Track;
    type: 'track' | 'playlist';
}