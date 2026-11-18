import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../app/providers.dart';
import '../../../app/theme.dart';
import '../../../l10n/gen/app_localizations.dart';
import '../../system/data/health_repository.dart';

/// Me tab. Profile, privacy and data controls arrive in Phase 1.
class MeScreen extends ConsumerWidget {
  const MeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = AppLocalizations.of(context);
    final config = ref.watch(appConfigProvider);
    final serverVersion = switch (ref.watch(serverStatusProvider)) {
      AsyncData(value: ServerOnline(:final version)) => version,
      _ => '—',
    };
    return SafeArea(
      child: ListView(
        padding: const EdgeInsets.fromLTRB(20, 24, 20, 24),
        children: [
          Text(
            l10n.meTitle,
            style: Theme.of(
              context,
            ).textTheme.headlineLarge?.copyWith(fontWeight: FontWeight.w700),
          ),
          const SizedBox(height: 24),
          Card(
            child: Column(
              children: [
                _Row(label: l10n.meEnvironment, value: config.environment),
                const Divider(),
                _Row(label: l10n.meApiBaseUrl, value: config.apiBaseUrl),
                const Divider(),
                _Row(label: l10n.meServerVersion, value: serverVersion),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _Row extends StatelessWidget {
  const _Row({required this.label, required this.value});

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return ListTile(
      title: Text(label),
      trailing: Text(value, style: const TextStyle(color: AppColors.textMuted)),
    );
  }
}
