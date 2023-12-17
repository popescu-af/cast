import 'package:flutter/material.dart';

import '../service_proxy.dart';

class DevicesScreen extends StatefulWidget {
  final Widget child;

  const DevicesScreen({
    required this.child,
    super.key,
  });

  @override
  State<DevicesScreen> createState() => _DevicesScreenState();
}

class _DevicesScreenState extends State<DevicesScreen>
    with SingleTickerProviderStateMixin {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Devices'),
      ),
      body: widget.child,
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          setState(() {
            proxy.fetchDevices(rescan: true);
          });
        },
        foregroundColor: Colors.white,
        backgroundColor: Colors.blue,
        shape: const CircleBorder(),
        child: const Icon(Icons.refresh),
      ),
    );
  }
}
