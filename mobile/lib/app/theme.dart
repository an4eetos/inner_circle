import 'package:flutter/material.dart';

/// Colour tokens taken from the concept screens (docs/screens). Dark-first.
abstract final class AppColors {
  static const background = Color(0xFF0D1015);
  static const surface = Color(0xFF161B22);
  static const surfaceRaised = Color(0xFF1D232C);
  static const border = Color(0xFF262D38);
  static const text = Color(0xFFF2F4F7);
  static const textMuted = Color(0xFF8B94A3);

  /// Primary actions, success, the Longevity track.
  static const mint = Color(0xFF4FE3AE);

  /// "Needs attention", nudges, the Health track.
  static const amber = Color(0xFFFFB547);

  /// The Biohacking track, recommendations.
  static const periwinkle = Color(0xFF8C9EFF);

  static const danger = Color(0xFFFF6B6B);
}

/// The three program tracks and their marker colours.
enum Track {
  health(AppColors.amber),
  biohacking(AppColors.periwinkle),
  longevity(AppColors.mint);

  const Track(this.color);

  final Color color;
}

/// Builds the app theme.
ThemeData buildTheme() {
  const scheme = ColorScheme.dark(
    surface: AppColors.background,
    primary: AppColors.mint,
    onPrimary: Color(0xFF07140F),
    secondary: AppColors.amber,
    onSecondary: Color(0xFF1A1205),
    tertiary: AppColors.periwinkle,
    error: AppColors.danger,
    onSurface: AppColors.text,
    onSurfaceVariant: AppColors.textMuted,
    outline: AppColors.border,
    surfaceContainer: AppColors.surface,
    surfaceContainerHigh: AppColors.surfaceRaised,
  );

  final base = ThemeData(
    useMaterial3: true,
    colorScheme: scheme,
    scaffoldBackgroundColor: AppColors.background,
  );

  return base.copyWith(
    textTheme: base.textTheme.apply(
      bodyColor: AppColors.text,
      displayColor: AppColors.text,
    ),
    cardTheme: CardThemeData(
      color: AppColors.surface,
      elevation: 0,
      margin: EdgeInsets.zero,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(20),
        side: const BorderSide(color: AppColors.border),
      ),
    ),
    filledButtonTheme: FilledButtonThemeData(
      style: FilledButton.styleFrom(
        minimumSize: const Size.fromHeight(52),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        textStyle: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
      ),
    ),
    navigationBarTheme: NavigationBarThemeData(
      backgroundColor: AppColors.background,
      indicatorColor: Colors.transparent,
      iconTheme: WidgetStateProperty.resolveWith(
        (states) => IconThemeData(
          color: states.contains(WidgetState.selected)
              ? AppColors.mint
              : AppColors.textMuted,
        ),
      ),
      labelTextStyle: WidgetStateProperty.resolveWith(
        (states) => TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w500,
          color: states.contains(WidgetState.selected)
              ? AppColors.mint
              : AppColors.textMuted,
        ),
      ),
    ),
    dividerTheme: const DividerThemeData(color: AppColors.border, space: 1),
  );
}
