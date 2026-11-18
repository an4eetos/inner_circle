// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appTitle => 'Inner Circle';

  @override
  String get navBody => 'Body';

  @override
  String get navCircle => 'Circle';

  @override
  String get navMe => 'Me';

  @override
  String get navAdd => 'Add';

  @override
  String get bodyTitle => 'Your body';

  @override
  String get bodyEmpty =>
      'No programs yet. Your organ programs will appear on the body map.';

  @override
  String get circleTitle => 'Circle';

  @override
  String get circleEmpty =>
      'Your inner circle will live here: friends\' progress, cheers, nudges and asks.';

  @override
  String get meTitle => 'Me';

  @override
  String get meEnvironment => 'Environment';

  @override
  String get meApiBaseUrl => 'API';

  @override
  String get meServerVersion => 'Server version';

  @override
  String get quickAddTitle => 'Quick add';

  @override
  String get quickAddBody =>
      'Check-ins, meals and new programs will start here.';

  @override
  String get connectionChecking => 'Checking server…';

  @override
  String get connectionConnected => 'Connected';

  @override
  String connectionConnectedDetail(String version) {
    return 'Server $version';
  }

  @override
  String get connectionOffline => 'Can\'t reach the server';

  @override
  String get connectionRetry => 'Retry';
}
