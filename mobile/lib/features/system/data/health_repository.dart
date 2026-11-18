import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:inner_circle_api/inner_circle_api.dart';

import '../../../app/providers.dart';

/// Result of a liveness check against the API.
sealed class ServerStatus {
  const ServerStatus();
}

/// The API answered `/v1/healthz`.
final class ServerOnline extends ServerStatus {
  const ServerOnline(this.version);

  final String version;
}

/// The API could not be reached or answered with an error.
final class ServerOffline extends ServerStatus {
  const ServerOffline(this.reason);

  final String reason;
}

/// Reads server liveness.
abstract interface class HealthRepository {
  Future<ServerStatus> check();
}

/// [HealthRepository] backed by the generated API client.
class ApiHealthRepository implements HealthRepository {
  ApiHealthRepository(this._api);

  final SystemApi _api;

  @override
  Future<ServerStatus> check() async {
    try {
      final res = await _api.getHealthz();
      final body = res.data;
      if (body == null || body.status != HealthStatusEnum.ok) {
        return const ServerOffline('unexpected response');
      }
      return ServerOnline(body.version);
    } on Exception catch (e) {
      return ServerOffline(e.runtimeType.toString());
    }
  }
}

final healthRepositoryProvider = Provider<HealthRepository>(
  (ref) => ApiHealthRepository(ref.watch(apiProvider).getSystemApi()),
);

/// Current server status; `ref.invalidate(serverStatusProvider)` re-checks.
final serverStatusProvider = FutureProvider.autoDispose<ServerStatus>(
  (ref) => ref.watch(healthRepositoryProvider).check(),
);
