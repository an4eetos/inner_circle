import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../features/body/presentation/body_screen.dart';
import '../features/circle/presentation/circle_screen.dart';
import '../features/me/presentation/me_screen.dart';
import '../features/quick_add/presentation/quick_add_sheet.dart';
import '../l10n/gen/app_localizations.dart';
import 'theme.dart';

/// Route paths. Deep links (invites, push taps) resolve against these.
abstract final class Routes {
  static const body = '/body';
  static const circle = '/circle';
  static const me = '/me';
}

final routerProvider = Provider<GoRouter>((ref) {
  final router = GoRouter(
    initialLocation: Routes.body,
    routes: [
      StatefulShellRoute.indexedStack(
        builder: (context, state, shell) => AppShell(shell: shell),
        branches: [
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: Routes.body,
                builder: (context, state) => const BodyScreen(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: Routes.circle,
                builder: (context, state) => const CircleScreen(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: Routes.me,
                builder: (context, state) => const MeScreen(),
              ),
            ],
          ),
        ],
      ),
    ],
  );
  ref.onDispose(router.dispose);
  return router;
});

/// Bottom navigation from the screens: Body · Circle · + · Me.
/// "+" is an action (quick add sheet), not a tab.
class AppShell extends StatelessWidget {
  const AppShell({super.key, required this.shell});

  final StatefulNavigationShell shell;

  static const _addSlot = 2;

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context);
    // Visual slots: 0 Body, 1 Circle, 2 "+", 3 Me. Branches: 0 Body, 1 Circle, 2 Me.
    final selected = shell.currentIndex >= _addSlot
        ? shell.currentIndex + 1
        : shell.currentIndex;
    return Scaffold(
      body: shell,
      bottomNavigationBar: NavigationBar(
        selectedIndex: selected,
        onDestinationSelected: (slot) {
          if (slot == _addSlot) {
            showQuickAddSheet(context);
            return;
          }
          final branch = slot > _addSlot ? slot - 1 : slot;
          shell.goBranch(branch, initialLocation: branch == shell.currentIndex);
        },
        destinations: [
          NavigationDestination(
            icon: const Icon(Icons.accessibility_new_outlined),
            selectedIcon: const Icon(Icons.accessibility_new),
            label: l10n.navBody,
          ),
          NavigationDestination(
            icon: const Icon(Icons.group_outlined),
            selectedIcon: const Icon(Icons.group),
            label: l10n.navCircle,
          ),
          NavigationDestination(icon: const _AddButton(), label: l10n.navAdd),
          NavigationDestination(
            icon: const Icon(Icons.person_outline),
            selectedIcon: const Icon(Icons.person),
            label: l10n.navMe,
          ),
        ],
      ),
    );
  }
}

class _AddButton extends StatelessWidget {
  const _AddButton();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 52,
      height: 40,
      decoration: BoxDecoration(
        color: AppColors.text,
        borderRadius: BorderRadius.circular(14),
      ),
      child: const Icon(Icons.add, color: AppColors.background),
    );
  }
}
