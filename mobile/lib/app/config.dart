/// Build-time configuration, injected with `--dart-define-from-file=config/<flavor>.json`.
class AppConfig {
  const AppConfig({required this.environment, required this.apiBaseUrl});

  /// Reads the values compiled into this build.
  factory AppConfig.fromEnvironment() {
    const env = String.fromEnvironment('APP_ENV', defaultValue: 'dev');
    const url = String.fromEnvironment(
      'API_BASE_URL',
      defaultValue: 'http://localhost:18080',
    );
    return const AppConfig(environment: env, apiBaseUrl: url);
  }

  /// `dev`, `staging` or `prod`.
  final String environment;

  /// Base URL of the Go API, without a trailing slash.
  final String apiBaseUrl;

  bool get isProd => environment == 'prod';
}
