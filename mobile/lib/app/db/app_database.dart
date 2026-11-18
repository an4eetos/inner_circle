import 'package:drift/drift.dart';
import 'package:drift_flutter/drift_flutter.dart';

part 'app_database.g.dart';

/// Local database: offline cache and the check-in outbox (tables arrive in Phase 2).
/// Regenerate with `make gen-mobile` after changing tables.
@DriftDatabase(tables: [])
class AppDatabase extends _$AppDatabase {
  AppDatabase([QueryExecutor? executor])
    : super(executor ?? driftDatabase(name: 'inner_circle'));

  @override
  int get schemaVersion => 1;
}
