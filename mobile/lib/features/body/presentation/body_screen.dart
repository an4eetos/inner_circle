import 'package:flutter/material.dart';

import '../../../app/theme.dart';
import '../../../l10n/gen/app_localizations.dart';
import '../../system/presentation/connection_status_card.dart';

/// Body tab. Phase 2 adds the 2D body map; Phase 6 the 3D one.
class BodyScreen extends StatelessWidget {
  const BodyScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    final text = Theme.of(context).textTheme;
    return SafeArea(
      child: ListView(
        padding: const EdgeInsets.fromLTRB(20, 24, 20, 24),
        children: [
          Text(
            l10n.bodyTitle,
            style: text.headlineLarge?.copyWith(fontWeight: FontWeight.w700),
          ),
          const SizedBox(height: 8),
          Text(
            l10n.bodyEmpty,
            style: text.bodyLarge?.copyWith(color: AppColors.textMuted),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              for (final track in Track.values) ...[
                _TrackChip(track: track),
                const SizedBox(width: 8),
              ],
            ],
          ),
          const SizedBox(height: 24),
          const ConnectionStatusCard(),
        ],
      ),
    );
  }
}

class _TrackChip extends StatelessWidget {
  const _TrackChip({required this.track});

  final Track track;

  @override
  Widget build(BuildContext context) {
    final label = switch (track) {
      Track.health => 'Health',
      Track.biohacking => 'Biohacking',
      Track.longevity => 'Longevity',
    };
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(999),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 8,
            height: 8,
            decoration: BoxDecoration(
              color: track.color,
              shape: BoxShape.circle,
            ),
          ),
          const SizedBox(width: 6),
          Text(label, style: Theme.of(context).textTheme.labelLarge),
        ],
      ),
    );
  }
}
