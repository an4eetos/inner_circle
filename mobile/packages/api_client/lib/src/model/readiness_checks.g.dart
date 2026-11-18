// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'readiness_checks.dart';

// **************************************************************************
// CopyWithGenerator
// **************************************************************************

abstract class _$ReadinessChecksCWProxy {
  ReadinessChecks database(ReadinessChecksDatabaseEnum database);

  /// This function **does support** nullification of nullable fields. All `null` values passed to `non-nullable` fields will be ignored. You can also use `ReadinessChecks(...).copyWith.fieldName(...)` to override fields one at a time with nullification support.
  ///
  /// Usage
  /// ```dart
  /// ReadinessChecks(...).copyWith(id: 12, name: "My name")
  /// ````
  ReadinessChecks call({ReadinessChecksDatabaseEnum database});
}

/// Proxy class for `copyWith` functionality. This is a callable class and can be used as follows: `instanceOfReadinessChecks.copyWith(...)`. Additionally contains functions for specific fields e.g. `instanceOfReadinessChecks.copyWith.fieldName(...)`
class _$ReadinessChecksCWProxyImpl implements _$ReadinessChecksCWProxy {
  const _$ReadinessChecksCWProxyImpl(this._value);

  final ReadinessChecks _value;

  @override
  ReadinessChecks database(ReadinessChecksDatabaseEnum database) =>
      this(database: database);

  @override
  /// This function **does support** nullification of nullable fields. All `null` values passed to `non-nullable` fields will be ignored. You can also use `ReadinessChecks(...).copyWith.fieldName(...)` to override fields one at a time with nullification support.
  ///
  /// Usage
  /// ```dart
  /// ReadinessChecks(...).copyWith(id: 12, name: "My name")
  /// ````
  ReadinessChecks call({Object? database = const $CopyWithPlaceholder()}) {
    return ReadinessChecks(
      database: database == const $CopyWithPlaceholder()
          ? _value.database
          // ignore: cast_nullable_to_non_nullable
          : database as ReadinessChecksDatabaseEnum,
    );
  }
}

extension $ReadinessChecksCopyWith on ReadinessChecks {
  /// Returns a callable class that can be used as follows: `instanceOfReadinessChecks.copyWith(...)` or like so:`instanceOfReadinessChecks.copyWith.fieldName(...)`.
  // ignore: library_private_types_in_public_api
  _$ReadinessChecksCWProxy get copyWith => _$ReadinessChecksCWProxyImpl(this);
}

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ReadinessChecks _$ReadinessChecksFromJson(Map<String, dynamic> json) =>
    $checkedCreate('ReadinessChecks', json, ($checkedConvert) {
      $checkKeys(json, requiredKeys: const ['database']);
      final val = ReadinessChecks(
        database: $checkedConvert(
          'database',
          (v) => $enumDecode(
            _$ReadinessChecksDatabaseEnumEnumMap,
            v,
            unknownValue: ReadinessChecksDatabaseEnum.unknownDefaultOpenApi,
          ),
        ),
      );
      return val;
    });

Map<String, dynamic> _$ReadinessChecksToJson(ReadinessChecks instance) =>
    <String, dynamic>{
      'database': _$ReadinessChecksDatabaseEnumEnumMap[instance.database]!,
    };

const _$ReadinessChecksDatabaseEnumEnumMap = {
  ReadinessChecksDatabaseEnum.ok: 'ok',
  ReadinessChecksDatabaseEnum.unknownDefaultOpenApi: 'unknown_default_open_api',
};
