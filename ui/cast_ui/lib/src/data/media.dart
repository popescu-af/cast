class Media {
  final int id;
  final String mediaId;
  final String uri;
  final String name;

  Media(this.id, this.mediaId, this.uri, this.name);
}

class PlaybackState {
  final int seekPosition;

  PlaybackState({
    required this.seekPosition,
  });
}
