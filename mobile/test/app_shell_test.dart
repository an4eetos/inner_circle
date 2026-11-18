import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:inner_circle/app/app.dart';
import 'package:inner_circle/app/config.dart';
import 'package:inner_circle/app/providers.dart';
import 'package:inner_circle/features/system/data/health_repository.dart';

class _FakeHealthRepository implements HealthRepository {
  _FakeHealthRepository(this.status);

  ServerStatus status;
  int calls = 0;

  @override
  Future<ServerStatus> check() async {
    calls++;
    return status;
  }
}

Widget _app(HealthRepository repo) => ProviderScope(
  overrides: [
    appConfigProvider.overrideWithValue(
      const AppConfig(environment: 'test', apiBaseUrl: 'http://api.test'),
    ),
    healthRepositoryProvider.overrideWithValue(repo),
  ],
  child: const InnerCircleApp(),
);

void main() {
  testWidgets('shell shows Body · Circle · + · Me and the connected state', (
    tester,
  ) async {
    await tester.pumpWidget(
      _app(_FakeHealthRepository(const ServerOnline('v-test'))),
    );
    await tester.pumpAndSettle();

    for (final label in ['Body', 'Circle', 'Add', 'Me']) {
      expect(find.text(label), findsOneWidget, reason: '$label tab');
    }
    expect(find.text('Your body'), findsOneWidget);
    expect(find.text('Connected'), findsOneWidget);
    expect(find.text('Server v-test'), findsOneWidget);
  });

  testWidgets('tabs switch screens and "+" opens the quick add sheet', (
    tester,
  ) async {
    await tester.pumpWidget(
      _app(_FakeHealthRepository(const ServerOnline('v'))),
    );
    await tester.pumpAndSettle();

    await tester.tap(find.text('Circle'));
    await tester.pumpAndSettle();
    expect(
      find.textContaining('Your inner circle will live here'),
      findsOneWidget,
    );

    await tester.tap(find.text('Me'));
    await tester.pumpAndSettle();
    expect(find.text('http://api.test'), findsOneWidget);
    expect(find.text('test'), findsOneWidget);

    await tester.tap(find.text('Add'));
    await tester.pumpAndSettle();
    expect(find.text('Quick add'), findsOneWidget);
  });

  testWidgets('offline state offers a retry that re-checks the server', (
    tester,
  ) async {
    final repo = _FakeHealthRepository(const ServerOffline('DioException'));
    await tester.pumpWidget(_app(repo));
    await tester.pumpAndSettle();

    expect(find.text("Can't reach the server"), findsOneWidget);
    final callsBefore = repo.calls;

    repo.status = const ServerOnline('v2');
    await tester.tap(find.text('Retry'));
    await tester.pumpAndSettle();

    expect(repo.calls, greaterThan(callsBefore));
    expect(find.text('Connected'), findsOneWidget);
  });
}
