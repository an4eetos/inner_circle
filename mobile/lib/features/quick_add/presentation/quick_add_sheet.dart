import 'package:flutter/material.dart';

import '../../../app/theme.dart';
import '../../../l10n/gen/app_localizations.dart';

/// Opens the "+" sheet. Check-ins, meals and new programs start here from Phase 2.
Future<void> showQuickAddSheet(BuildContext context) {
  return showModalBottomSheet<void>(
    context: context,
    backgroundColor: AppColors.surface,
    showDragHandle: true,
    builder: (context) {
      final l10n = AppLocalizations.of(context);
      return SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(24, 0, 24, 24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                l10n.quickAddTitle,
                style: Theme.of(context).textTheme.titleLarge,
              ),
              const SizedBox(height: 8),
              Text(
                l10n.quickAddBody,
                style: const TextStyle(color: AppColors.textMuted),
              ),
            ],
          ),
        ),
      );
    },
  );
}
