//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:copy_with_extension/copy_with_extension.dart';
import 'package:json_annotation/json_annotation.dart';

part 'readiness_checks.g.dart';

@CopyWith()
@JsonSerializable(
  checked: true,
  createToJson: true,
  disallowUnrecognizedKeys: false,
  explicitToJson: true,
)
class ReadinessChecks {
  /// Returns a new [ReadinessChecks] instance.
  ReadinessChecks({required this.database});

  @JsonKey(
    name: r'database',
    required: true,
    includeIfNull: false,
    unknownEnumValue: ReadinessChecksDatabaseEnum.unknownDefaultOpenApi,
  )
  final ReadinessChecksDatabaseEnum database;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is ReadinessChecks && other.database == database;

  @override
  int get hashCode => database.hashCode;

  factory ReadinessChecks.fromJson(Map<String, dynamic> json) =>
      _$ReadinessChecksFromJson(json);

  Map<String, dynamic> toJson() => _$ReadinessChecksToJson(this);

  @override
  String toString() {
    return toJson().toString();
  }
}

enum ReadinessChecksDatabaseEnum {
  @JsonValue(r'ok')
  ok(r'ok'),
  @JsonValue(r'unknown_default_open_api')
  unknownDefaultOpenApi(r'unknown_default_open_api');

  const ReadinessChecksDatabaseEnum(this.value);

  final String value;

  @override
  String toString() => value;
}
