import 'package:flutter/material.dart';

import '../data.dart';

class MediaList extends StatelessWidget {
  final List<Media> media;
  final ValueChanged<Media>? onTap;

  const MediaList({
    required this.media,
    this.onTap,
    super.key,
  });

  @override
  Widget build(BuildContext context) => ListView.builder(
        itemCount: media.length,
        itemBuilder: (context, index) => ListTile(
          title: Text(
            media[index].name,
          ),
          subtitle: Text(
            media[index].mediaId.substring(0, 7),
          ),
          onTap: onTap != null ? () => onTap!(media[index]) : null,
        ),
      );
}
