// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'readiness.dart';

// **************************************************************************
// CopyWithGenerator
// **************************************************************************

abstract class _$ReadinessCWProxy {
  Readiness status(ReadinessStatusEnum status);

  Readiness checks(ReadinessChecks checks);

  /// This function **does support** nullification of nullable fields. All `null` values passed to `non-nullable` fields will be ignored. You can also use `Readiness(...).copyWith.fieldName(...)` to override fields one at a time with nullification support.
  ///
  /// Usage
  /// ```dart
  /// Readiness(...).copyWith(id: 12, name: "My name")
  /// ````
  Readiness call({ReadinessStatusEnum status, ReadinessChecks checks});
}

/// Proxy class for `copyWith` functionality. This is a callable class and can be used as follows: `instanceOfReadiness.copyWith(...)`. Additionally contains functions for specific fields e.g. `instanceOfReadiness.copyWith.fieldName(...)`
class _$ReadinessCWProxyImpl implements _$ReadinessCWProxy {
  const _$ReadinessCWProxyImpl(this._value);

  final Readiness _value;

  @override
  Readiness status(ReadinessStatusEnum status) => this(status: status);

  @override
  Readiness checks(ReadinessChecks checks) => this(checks: checks);

  @override
  /// This function **does support** nullification of nullable fields. All `null` values passed to `non-nullable` fields will be ignored. You can also use `Readiness(...).copyWith.fieldName(...)` to override fields one at a time with nullification support.
  ///
  /// Usage
  /// ```dart
  /// Readiness(...).copyWith(id: 12, name: "My name")
  /// ````
  Readiness call({
    Object? status = const $CopyWithPlaceholder(),
    Object? checks = const $CopyWithPlaceholder(),
  }) {
    return Readiness(
      status: status == const $CopyWithPlaceholder()
          ? _value.status
          // ignore: cast_nullable_to_non_nullable
          : status as ReadinessStatusEnum,
      checks: checks == const $CopyWithPlaceholder()
          ? _value.checks
          // ignore: cast_nullable_to_non_nullable
          : checks as ReadinessChecks,
    );
  }
}

extension $ReadinessCopyWith on Readiness {
  /// Returns a callable class that can be used as follows: `instanceOfReadiness.copyWith(...)` or like so:`instanceOfReadiness.copyWith.fieldName(...)`.
  // ignore: library_private_types_in_public_api
  _$ReadinessCWProxy get copyWith => _$ReadinessCWProxyImpl(this);
}

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Readiness _$ReadinessFromJson(Map<String, dynamic> json) =>
    $checkedCreate('Readiness', json, ($checkedConvert) {
      $checkKeys(json, requiredKeys: const ['status', 'checks']);
      final val = Readiness(
        status: $checkedConvert(
          'status',
          (v) => $enumDecode(
            _$ReadinessStatusEnumEnumMap,
            v,
            unknownValue: ReadinessStatusEnum.unknownDefaultOpenApi,
          ),
        ),
        checks: $checkedConvert(
          'checks',
          (v) => ReadinessChecks.fromJson(v as Map<String, dynamic>),
        ),
      );
      return val;
    });

Map<String, dynamic> _$ReadinessToJson(Readiness instance) => <String, dynamic>{
  'status': _$ReadinessStatusEnumEnumMap[instance.status]!,
  'checks': instance.checks.toJson(),
};

const _$ReadinessStatusEnumEnumMap = {
  ReadinessStatusEnum.ok: 'ok',
  ReadinessStatusEnum.unknownDefaultOpenApi: 'unknown_default_open_api',
};
