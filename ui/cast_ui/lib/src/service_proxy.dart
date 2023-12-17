import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter/services.dart' show rootBundle;
import 'package:flutter/widgets.dart';

import 'data/device.dart';
import 'data/media.dart';

ServiceProxy proxy = ServiceProxy();

class ServiceProxy {
  String? serviceUrl;

  List<Device> allDevices = [];
  String selectedDeviceId = "";

  List<Media> allMedia = [];
  Map<String, Media> allSubtitles = {};

  Map<String, PlaybackState> allPlaybacks = {};

  void addMedia({
    required String mediaId,
    required String uri,
  }) {
    var fileName = uri.split('/').last;
    var name = fileName.split('.').first;
    var ext = fileName.split('.').last;
    var media = Media(allMedia.length, mediaId, uri, name);
    if (ext == 'srt') {
      allSubtitles[uri] = media;
      return;
    }
    allMedia.add(media);
  }

  void addDevice({
    required String deviceId,
  }) {
    var device = Device(allMedia.length, deviceId);
    allDevices.add(device);
  }

  Media getMedia(String id) {
    return allMedia[int.parse(id)];
  }

  StreamController<bool> mediaChangedStreamController =
      StreamController<bool>.broadcast();

  StreamController<bool> devicesChangedStreamController =
      StreamController<bool>.broadcast();

  StreamController<String> errorChangedStreamController =
      StreamController<String>.broadcast();

  Future fetchServiceUrl() async {
    if (serviceUrl != null) {
      return serviceUrl;
    }
    WidgetsFlutterBinding.ensureInitialized();
    serviceUrl = await rootBundle.loadString('cast-service-url.txt');
    return serviceUrl;
  }

  void setServiceUrl(String url) {
    serviceUrl = url;
  }

  Future fetchMedia() async {
    final serviceUrl = await fetchServiceUrl();

    allMedia.clear();

    final response = await http.get(Uri.parse('$serviceUrl/media'));
    if (response.statusCode != 200) {
      errorChangedStreamController.add(response.body);
      return;
    }

    final j = jsonDecode(response.body) as Map<String, dynamic>;
    if (!j.containsKey('data')) {
      errorChangedStreamController.add(response.body);
      return;
    }

    for (var media in j['data']) {
      addMedia(
        mediaId: media['id'] as String,
        uri: media['uri'] as String,
      );
    }

    mediaChangedStreamController.add(true);
    errorChangedStreamController.add("Ok");
  }

  Future fetchDevices({bool rescan = false}) async {
    final serviceUrl = await fetchServiceUrl();

    allDevices.clear();

    final response =
        await http.get(Uri.parse('$serviceUrl/devices?rescan=$rescan'));
    if (response.statusCode != 200) {
      errorChangedStreamController.add(response.body);
      return;
    }

    final j = jsonDecode(response.body) as Map<String, dynamic>;
    if (!j.containsKey('data')) {
      errorChangedStreamController.add(response.body);
      return;
    }

    for (var dev in j['data']) {
      addDevice(
        deviceId: dev['id'] as String,
      );
    }

    if (allDevices.isNotEmpty) {
      selectedDeviceId = allDevices.first.deviceId;
    }

    devicesChangedStreamController.add(true);
    errorChangedStreamController.add("Ok");
  }

  Future loadMedia(String mediaId, String uri) async {
    final serviceUrl = await fetchServiceUrl();

    List<String> uriParts = uri.split(".");
    uriParts.removeLast();
    String uriSubtitle = '${uriParts.join(".")}.srt';
    String subtitleMediaId = "";
    if (allSubtitles.containsKey(uriSubtitle)) {
      subtitleMediaId = allSubtitles[uriSubtitle]!.mediaId;
    }
    final response = await http.post(Uri.parse(
        '$serviceUrl/load?mediaId=$mediaId&&subtitleId=$subtitleMediaId&&deviceId=$selectedDeviceId'));
    if (response.statusCode != 200) {
      errorChangedStreamController.add(response.body);
      return;
    }
    errorChangedStreamController.add("Ok");
  }

  Future controlPlayback(String command, {double amount = 0}) async {
    final serviceUrl = await fetchServiceUrl();

    final response = await http.post(
        Uri.parse('$serviceUrl/playback?command=$command&amount=$amount'));
    if (response.statusCode != 200) {
      errorChangedStreamController.add(response.body);
      return;
    }
    errorChangedStreamController.add("Ok");
  }

  Future getPlaybackProgress() async {
    final serviceUrl = await fetchServiceUrl();

    final response = await http.get(Uri.parse('$serviceUrl/playback'));
    if (response.statusCode != 200) {
      // this is not necessarily an error, it appears also when the playback is stopped
      return 0.0;
    }

    final j = jsonDecode(response.body) as Map<String, dynamic>;
    if (!j.containsKey('data')) {
      errorChangedStreamController.add(response.body);
      return 0.0;
    }

    final playbackInfo = j['data'];
    final position = playbackInfo['position'];
    final duration = playbackInfo['duration'];
    if (duration == 0) {
      errorChangedStreamController.add(response.body);
      return 0.0;
    }
    errorChangedStreamController.add("Ok");
    return position / duration;
  }

  List<Media> get openMedia => [
        ...allMedia.where((media) => allPlaybacks.containsKey(media.mediaId)),
      ];
}
