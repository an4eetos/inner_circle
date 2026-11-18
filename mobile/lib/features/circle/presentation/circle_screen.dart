import 'package:flutter/material.dart';

import '../../../app/theme.dart';
import '../../../l10n/gen/app_localizations.dart';

/// Circle tab. The feed arrives in Phase 4.
class CircleScreen extends StatelessWidget {
  const CircleScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return SafeArea(
      child: ListView(
        padding: const EdgeInsets.fromLTRB(20, 24, 20, 24),
        children: [
          Text(
            l10n.circleTitle,
            style: text.headlineLarge?.copyWith(fontWeight: FontWeight.w700),
          ),
          const SizedBox(height: 8),
          Text(
            l10n.circleEmpty,
            style: text.bodyLarge?.copyWith(color: AppColors.textMuted),
          ),
        ],
      ),
    );
  }
}
