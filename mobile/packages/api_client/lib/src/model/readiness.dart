//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:inner_circle_api/src/model/readiness_checks.dart';
import 'package:copy_with_extension/copy_with_extension.dart';
import 'package:json_annotation/json_annotation.dart';

part 'readiness.g.dart';

@CopyWith()
@JsonSerializable(
  checked: true,
  createToJson: true,
  disallowUnrecognizedKeys: false,
  explicitToJson: true,
)
class Readiness {
  /// Returns a new [Readiness] instance.
  Readiness({required this.status, required this.checks});

  @JsonKey(
    name: r'status',
    required: true,
    includeIfNull: false,
    unknownEnumValue: ReadinessStatusEnum.unknownDefaultOpenApi,
  )
  final ReadinessStatusEnum status;

  @JsonKey(name: r'checks', required: true, includeIfNull: false)
  final ReadinessChecks checks;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Readiness && other.status == status && other.checks == checks;

  @override
  int get hashCode => status.hashCode + checks.hashCode;

  factory Readiness.fromJson(Map<String, dynamic> json) =>
      _$ReadinessFromJson(json);

  Map<String, dynamic> toJson() => _$ReadinessToJson(this);

  @override
  String toString() {
    return toJson().toString();
  }
}

enum ReadinessStatusEnum {
  @JsonValue(r'ok')
  ok(r'ok'),
  @JsonValue(r'unknown_default_open_api')
  unknownDefaultOpenApi(r'unknown_default_open_api');

  const ReadinessStatusEnum(this.value);

  final String value;

  @override
  String toString() => value;
}
