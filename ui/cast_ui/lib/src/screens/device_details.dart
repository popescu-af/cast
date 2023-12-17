import 'package:flutter/material.dart';

import '../data.dart';

class DeviceDetailsScreen extends StatelessWidget {
  final Device device;

  const DeviceDetailsScreen({
    super.key,
    required this.device,
  });

  @override
  Widget build(BuildContext context) => Scaffold(
        appBar: AppBar(
          title: Text(device.deviceId),
        ),
        body: const Center(),
      );
}
