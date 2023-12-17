import 'package:flutter/material.dart';

import '../service_proxy.dart';

class MediaScreen extends StatefulWidget {
  final Widget child;
  final ValueChanged<int> onTap;
  final int selectedIndex;

  const MediaScreen({
    required this.child,
    required this.onTap,
    required this.selectedIndex,
    super.key,
  });

  @override
  State<MediaScreen> createState() => _MediaScreenState();
}

class _MediaScreenState extends State<MediaScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this)
      ..addListener(_handleTabIndexChanged);
  }

  @override
  void dispose() {
    _tabController.removeListener(_handleTabIndexChanged);
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    _tabController.index = widget.selectedIndex;
    return Scaffold(
      appBar: AppBar(
        title: const Text('Media'),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(
              text: 'All',
              icon: Icon(Icons.list_outlined),
            ),
            Tab(
              text: 'Continue Watching',
              icon: Icon(Icons.list_outlined),
            ),
          ],
        ),
      ),
      body: widget.child,
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          setState(() {
            proxy.fetchMedia();
          });
        },
        foregroundColor: Colors.white,
        backgroundColor: Colors.blue,
        shape: const CircleBorder(),
        child: const Icon(Icons.refresh),
      ),
    );
  }

  void _handleTabIndexChanged() {
    widget.onTap(_tabController.index);
  }
}
