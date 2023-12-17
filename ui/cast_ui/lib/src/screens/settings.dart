import 'package:flutter/material.dart';

import '../service_proxy.dart';

class SettingsScreen extends StatefulWidget {
  const SettingsScreen({super.key});

  @override
  State<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends State<SettingsScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings'),
      ),
      body: const SingleChildScrollView(
        child: Align(
          alignment: Alignment.center,
          child: SettingsContent(),
        ),
      ),
    );
  }
}

class SettingsContent extends StatelessWidget {
  const SettingsContent({
    super.key,
  });

  @override
  Widget build(BuildContext context) => Column(
        children: [
          Builder(
            builder: (context) {
              return TextField(
                decoration: const InputDecoration(
                  labelText: 'Service URL',
                ),
                controller: TextEditingController(
                  text: proxy.serviceUrl,
                ),
                onChanged: (value) {
                  proxy.setServiceUrl(value);
                },
              );
            },
          ),
        ],
      );
}
