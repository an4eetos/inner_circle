//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:copy_with_extension/copy_with_extension.dart';
import 'package:json_annotation/json_annotation.dart';

part 'problem.g.dart';

@CopyWith()
@JsonSerializable(
  checked: true,
  createToJson: true,
  disallowUnrecognizedKeys: false,
  explicitToJson: true,
)
class Problem {
  /// Returns a new [Problem] instance.
  Problem({
    required this.type,

    required this.title,

    required this.status,

    this.detail,

    this.instance,

    required this.code,

    this.requestId,
  });

  /// URI reference identifying the problem type
  @JsonKey(name: r'type', required: true, includeIfNull: false)
  final String type;

  /// Short human-readable summary
  @JsonKey(name: r'title', required: true, includeIfNull: false)
  final String title;

  /// HTTP status code
  @JsonKey(name: r'status', required: true, includeIfNull: false)
  final int status;

  /// Human-readable explanation specific to this occurrence
  @JsonKey(name: r'detail', required: false, includeIfNull: false)
  final String? detail;

  /// URI reference identifying this occurrence
  @JsonKey(name: r'instance', required: false, includeIfNull: false)
  final String? instance;

  /// Stable machine-readable error code
  @JsonKey(name: r'code', required: true, includeIfNull: false)
  final String code;

  /// Request ID, also returned in the X-Request-ID header
  @JsonKey(name: r'request_id', required: false, includeIfNull: false)
  final String? requestId;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is Problem &&
          other.type == type &&
          other.title == title &&
          other.status == status &&
          other.detail == detail &&
          other.instance == instance &&
          other.code == code &&
          other.requestId == requestId;

  @override
  int get hashCode =>
      type.hashCode +
      title.hashCode +
      status.hashCode +
      detail.hashCode +
      instance.hashCode +
      code.hashCode +
      requestId.hashCode;

  factory Problem.fromJson(Map<String, dynamic> json) =>
      _$ProblemFromJson(json);

  Map<String, dynamic> toJson() => _$ProblemToJson(this);

  @override
  String toString() {
    return toJson().toString();
  }
}
