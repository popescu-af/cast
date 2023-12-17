import 'dart:async';

import 'package:flutter/material.dart';

import '../data.dart';
import '../service_proxy.dart';

class MediaDetailsScreen extends StatefulWidget {
  final Media? media;

  const MediaDetailsScreen({
    super.key,
    this.media,
  });

  @override
  State<MediaDetailsScreen> createState() => _MediaDetailsScreenState();
}

class _MediaDetailsScreenState extends State<MediaDetailsScreen>
    with SingleTickerProviderStateMixin {
  late Stream<double> _playbackProgressUpdater;
  late StreamSubscription<double> _playbackProgressSubscription;
  double _desiredProgress = 0;
  bool _seekInProgress = false;

  Stream<double> createPlaybackProgressStream() async* {
    yield* Stream.periodic(const Duration(seconds: 5), (_) {
      return proxy.getPlaybackProgress();
    }).asyncMap((event) async => await event);
  }

  @override
  void initState() {
    super.initState();
    _playbackProgressUpdater = createPlaybackProgressStream();
    _playbackProgressSubscription = _playbackProgressUpdater.listen((event) {
      if (!_seekInProgress) {
        setState(() {
          _desiredProgress = event * 100;
        });
      }
    });
  }

  @override
  void dispose() {
    _playbackProgressSubscription.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (widget.media == null) {
      return const Scaffold(
        body: Center(
          child: Text('No media found.'),
        ),
      );
    }
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.media!.name),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Align(
          alignment: Alignment.topLeft,
          child: Column(
            children: [
              Table(
                columnWidths: Map.fromIterables(
                  [0, 1],
                  [
                    const FractionColumnWidth(0.2),
                    const FractionColumnWidth(0.8),
                  ],
                ),
                children: [
                  TableRow(
                    children: [
                      Text(
                        'Name:',
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                      Text(
                        widget.media!.name,
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                    ],
                  ),
                  TableRow(
                    children: [
                      Text(
                        'ID:',
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                      Text(
                        widget.media!.mediaId,
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                    ],
                  ),
                  TableRow(
                    children: [
                      Text(
                        'URI:',
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                      Text(
                        widget.media!.uri,
                        style: Theme.of(context).textTheme.titleMedium,
                      ),
                    ],
                  ),
                ],
              ),
              const Divider(
                color: Colors.grey,
                height: 20,
                thickness: 1,
                indent: 0,
                endIndent: 0,
              ),
              Container(
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(10),
                  color: const Color.fromARGB(255, 64, 188, 255),
                ),
                child: Column(
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      children: [
                        IconButton(
                          icon: const Icon(Icons.cast),
                          onPressed: () {
                            proxy.loadMedia(
                                widget.media!.mediaId, widget.media!.uri);
                          },
                          tooltip: "load",
                        ),
                        IconButton(
                          icon: const Icon(Icons.play_arrow),
                          onPressed: () {
                            proxy.controlPlayback("play");
                          },
                          tooltip: "play",
                        ),
                        IconButton(
                          icon: const Icon(Icons.pause),
                          onPressed: () {
                            proxy.controlPlayback("pause");
                          },
                          tooltip: "pause",
                        ),
                        IconButton(
                          icon: const Icon(Icons.stop),
                          onPressed: () {
                            proxy.controlPlayback("stop");
                          },
                          tooltip: "stop",
                        ),
                      ],
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      children: [
                        IconButton(
                          icon: const Icon(Icons.keyboard_double_arrow_left),
                          onPressed: () {
                            proxy.controlPlayback("rev", amount: 30);
                          },
                          tooltip: "-30s",
                        ),
                        IconButton(
                          icon: const Icon(Icons.keyboard_arrow_left),
                          onPressed: () {
                            proxy.controlPlayback("rev", amount: 10);
                          },
                          tooltip: "-10s",
                        ),
                        Text(
                          "${_desiredProgress.toStringAsFixed(2)}%",
                          style: Theme.of(context).textTheme.titleMedium,
                        ),
                        IconButton(
                          icon: const Icon(Icons.keyboard_arrow_right),
                          onPressed: () {
                            proxy.controlPlayback("fwd", amount: 10);
                          },
                          tooltip: "+10s",
                        ),
                        IconButton(
                          icon: const Icon(Icons.keyboard_double_arrow_right),
                          onPressed: () {
                            proxy.controlPlayback("fwd", amount: 30);
                          },
                          tooltip: "+30s",
                        ),
                      ],
                    ),
                    Builder(
                      builder: (context) {
                        return Slider(
                          value: _desiredProgress,
                          max: 100,
                          divisions: null, // continuous
                          label: "seek",
                          onChanged: (double value) {
                            setState(() {
                              _seekInProgress = true;
                              _desiredProgress = value;
                            });
                          },
                          onChangeEnd: (double value) {
                            setState(() {
                              _seekInProgress = false;
                              proxy.controlPlayback("seek", amount: value);
                            });
                          },
                        );
                      },
                    ),
                  ],
                ),
              ),
              const Divider(
                color: Colors.grey,
                height: 20,
                thickness: 1,
                indent: 0,
                endIndent: 0,
              ),
              StreamBuilder<String>(
                stream: proxy.errorChangedStreamController.stream,
                builder: (context, snapshot) {
                  return Text(
                    snapshot.data ?? "",
                    style: Theme.of(context).textTheme.titleMedium,
                  );
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
