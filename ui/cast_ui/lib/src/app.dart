import 'package:cast_ui/src/widgets/device_list.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'screens/device_details.dart';
import 'screens/devices.dart';
import 'screens/media.dart';
import 'screens/media_details.dart';
import 'screens/settings.dart';
import 'screens/scaffold.dart';
import 'widgets/fade_transition_page.dart';
import 'widgets/media_list.dart';

import 'service_proxy.dart';

final appShellNavigatorKey = GlobalKey<NavigatorState>(debugLabel: 'app shell');
final mediaNavigatorKey = GlobalKey<NavigatorState>(debugLabel: 'media shell');

class Cast extends StatefulWidget {
  const Cast({super.key});

  @override
  State<Cast> createState() => _CastState();
}

class _CastState extends State<Cast> {
  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      routerConfig: GoRouter(
        debugLogDiagnostics: true,
        initialLocation: '/media/all',
        redirect: (context, state) {
          return null;
        },
        routes: [
          ShellRoute(
            navigatorKey: appShellNavigatorKey,
            builder: (context, state, child) {
              return CastScaffold(
                selectedIndex: switch (state.uri.path) {
                  var p when p.startsWith('/media') => 0,
                  var p when p.startsWith('/devices') => 1,
                  var p when p.startsWith('/settings') => 2,
                  _ => 0,
                },
                child: child,
              );
            },
            routes: [
              ShellRoute(
                pageBuilder: (context, state, child) {
                  return FadeTransitionPage<dynamic>(
                    key: state.pageKey,
                    child: Builder(builder: (context) {
                      return MediaScreen(
                        onTap: (idx) {
                          GoRouter.of(context).go(switch (idx) {
                            0 => '/media/all',
                            1 => '/media/continue',
                            2 => '/settings',
                            _ => '/media/all',
                          });
                        },
                        selectedIndex: switch (state.uri.path) {
                          var p when p.startsWith('/media/all') => 0,
                          var p when p.startsWith('/media/continue') => 1,
                          _ => 0,
                        },
                        child: child,
                      );
                    }),
                  );
                },
                routes: [
                  GoRoute(
                    path: '/media/all',
                    pageBuilder: (context, state) {
                      return FadeTransitionPage<dynamic>(
                        key: state.pageKey,
                        child: StreamBuilder(
                          stream: proxy.mediaChangedStreamController.stream,
                          builder: (context, snapshot) {
                            return MediaList(
                              media: proxy.allMedia,
                              onTap: (media) {
                                GoRouter.of(context)
                                    .go('/media/all/media/${media.id}');
                              },
                            );
                          },
                        ),
                      );
                    },
                    routes: [
                      GoRoute(
                        path: 'media/:mediaId',
                        parentNavigatorKey: appShellNavigatorKey,
                        builder: (context, state) {
                          return MediaDetailsScreen(
                            media: proxy.getMedia(
                                state.pathParameters['mediaId'] ?? ''),
                          );
                        },
                      ),
                    ],
                  ),
                  GoRoute(
                    path: '/media/continue',
                    pageBuilder: (context, state) {
                      return FadeTransitionPage<dynamic>(
                        key: state.pageKey,
                        child: Builder(
                          builder: (context) {
                            return MediaList(
                              media: proxy.openMedia,
                              onTap: (media) {
                                GoRouter.of(context)
                                    .go('/media/continue/media/${media.id}');
                              },
                            );
                          },
                        ),
                      );
                    },
                    routes: [
                      GoRoute(
                        path: 'media/:mediaId',
                        parentNavigatorKey: appShellNavigatorKey,
                        builder: (context, state) {
                          return MediaDetailsScreen(
                            media: proxy.getMedia(
                                state.pathParameters['mediaId'] ?? ''),
                          );
                        },
                      ),
                    ],
                  ),
                ],
              ),
              ShellRoute(
                pageBuilder: (context, state, child) {
                  return FadeTransitionPage<dynamic>(
                    key: state.pageKey,
                    child: Builder(builder: (context) {
                      return DevicesScreen(
                        child: child,
                      );
                    }),
                  );
                },
                routes: [
                  GoRoute(
                    path: '/devices',
                    pageBuilder: (context, state) {
                      return FadeTransitionPage<dynamic>(
                        key: state.pageKey,
                        child: StreamBuilder(
                          stream: proxy.devicesChangedStreamController.stream,
                          builder: (context, snapshot) {
                            return DeviceList(
                              devices: proxy.allDevices,
                              onTap: (device) {
                                GoRouter.of(context)
                                    .go('/devices/${device.id}');
                              },
                            );
                          },
                        ),
                      );
                    },
                    routes: [
                      GoRoute(
                        path: ':deviceId',
                        builder: (context, state) {
                          final device = proxy.allDevices.firstWhere((device) =>
                              device.id ==
                              int.parse(state.pathParameters['deviceId']!));
                          return Builder(builder: (context) {
                            return DeviceDetailsScreen(
                              device: device,
                            );
                          });
                        },
                      )
                    ],
                  ),
                ],
              ),
              GoRoute(
                path: '/settings',
                pageBuilder: (context, state) {
                  return FadeTransitionPage<dynamic>(
                    key: state.pageKey,
                    child: const SettingsScreen(),
                  );
                },
              ),
            ],
          ),
        ],
      ),
    );
  }
}
