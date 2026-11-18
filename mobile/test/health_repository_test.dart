import 'dart:convert';
import 'dart:typed_data';

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:inner_circle/features/system/data/health_repository.dart';
import 'package:inner_circle_api/inner_circle_api.dart';

/// Answers every request with a fixed status and body.
class _StubAdapter implements HttpClientAdapter {
  _StubAdapter(this.status, this.body, {this.contentType = 'application/json'});

  final int status;
  final Object? body;
  final String contentType;
  RequestOptions? lastRequest;

  @override
  Future<ResponseBody> fetch(
    RequestOptions options,
    Stream<Uint8List>? requestStream,
    Future<void>? cancelFuture,
  ) async {
    lastRequest = options;
    return ResponseBody.fromString(
      jsonEncode(body),
      status,
      headers: {
        Headers.contentTypeHeader: [contentType],
      },
    );
  }

  @override
  void close({bool force = false}) {}
}

ApiHealthRepository _repo(_StubAdapter adapter) {
  final dio = Dio(BaseOptions(baseUrl: 'http://api.test'))
    ..httpClientAdapter = adapter;
  return ApiHealthRepository(
    InnerCircleApi(dio: dio, interceptors: const []).getSystemApi(),
  );
}

void main() {
  test('online when /v1/healthz answers ok', () async {
    final adapter = _StubAdapter(200, {
      'status': 'ok',
      'version': 'abc123',
      'time': '2026-11-19T10:00:00Z',
    });
    final status = await _repo(adapter).check();

    expect(adapter.lastRequest?.path, '/v1/healthz');
    expect(status, isA<ServerOnline>());
    expect((status as ServerOnline).version, 'abc123');
  });

  test('tolerates fields the server adds later', () async {
    final status = await _repo(
      _StubAdapter(200, {
        'status': 'ok',
        'version': 'v',
        'time': '2026-11-19T10:00:00Z',
        'region': 'added-in-a-future-release',
      }),
    ).check();
    expect(status, isA<ServerOnline>());
  });

  test('offline on a problem+json error', () async {
    final status = await _repo(
      _StubAdapter(503, {
        'type': 'about:blank',
        'title': 'Service Unavailable',
        'status': 503,
        'code': 'unavailable',
      }, contentType: 'application/problem+json'),
    ).check();
    expect(status, isA<ServerOffline>());
  });
}
