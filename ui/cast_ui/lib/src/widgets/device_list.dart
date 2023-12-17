import 'package:flutter/material.dart';

import '../data.dart';

class DeviceList extends StatelessWidget {
  final List<Device> devices;
  final ValueChanged<Device>? onTap;

  const DeviceList({
    required this.devices,
    this.onTap,
    super.key,
  });

  @override
  Widget build(BuildContext context) => ListView.builder(
        itemCount: devices.length,
        itemBuilder: (context, index) => ListTile(
          title: Text(
            devices[index].deviceId,
          ),
          onTap: onTap != null ? () => onTap!(devices[index]) : null,
        ),
      );
}
