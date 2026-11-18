import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../app/theme.dart';
import '../../../l10n/gen/app_localizations.dart';
import '../data/health_repository.dart';

/// Shows whether the app can reach the API (Phase 0 smoke check).
class ConnectionStatusCard extends ConsumerWidget {
  const ConnectionStatusCard({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l10n = AppLocalizations.of(context);
    final status = ref.watch(serverStatusProvider);

    final (
      Color dot,
      String title,
      String? detail,
      bool canRetry,
    ) = switch (status) {
      AsyncData(value: ServerOnline(:final version)) => (
        AppColors.mint,
        l10n.connectionConnected,
        l10n.connectionConnectedDetail(version),
        false,
      ),
      AsyncData(value: ServerOffline(:final reason)) => (
        AppColors.danger,
        l10n.connectionOffline,
        reason,
        true,
      ),
      AsyncError() => (AppColors.danger, l10n.connectionOffline, null, true),
      _ => (AppColors.textMuted, l10n.connectionChecking, null, false),
    };

    return Card(
      key: const Key('connection-status'),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Container(
              width: 10,
              height: 10,
              decoration: BoxDecoration(color: dot, shape: BoxShape.circle),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: Theme.of(context).textTheme.titleSmall),
                  if (detail != null)
                    Text(
                      detail,
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: AppColors.textMuted,
                      ),
                    ),
                ],
              ),
            ),
            if (canRetry)
              TextButton(
                onPressed: () => ref.invalidate(serverStatusProvider),
                child: Text(l10n.connectionRetry),
              ),
          ],
        ),
      ),
    );
  }
}
