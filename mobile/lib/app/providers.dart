import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:inner_circle_api/inner_circle_api.dart';

import 'config.dart';
import 'db/app_database.dart';

/// Build-time config. Overridden in tests.
final appConfigProvider = Provider<AppConfig>(
  (ref) => AppConfig.fromEnvironment(),
);

/// Shared HTTP client. Auth and refresh interceptors are added in Phase 1.
final dioProvider = Provider<Dio>((ref) {
  final config = ref.watch(appConfigProvider);
  final dio = Dio(
    BaseOptions(
      baseUrl: config.apiBaseUrl,
      connectTimeout: const Duration(seconds: 5),
      receiveTimeout: const Duration(seconds: 10),
    ),
  );
  ref.onDispose(dio.close);
  return dio;
});

/// Generated API client (packages/api_client), sharing [dioProvider].
final apiProvider = Provider<InnerCircleApi>(
  (ref) => InnerCircleApi(dio: ref.watch(dioProvider), interceptors: const []),
);

/// Local drift database.
final appDatabaseProvider = Provider<AppDatabase>((ref) {
  final db = AppDatabase();
  ref.onDispose(db.close);
  return db;
});
