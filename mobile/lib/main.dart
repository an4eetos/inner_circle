import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'app/app.dart';

/// Run with a flavor and its config, e.g. `make mobile-run FLAVOR=dev`.
void main() {
  runApp(const ProviderScope(child: InnerCircleApp()));
}
