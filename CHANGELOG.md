# Change Log

## Redpanda Chart

### [Unreleased](https://github.com/redpanda-data/helm-charts/releases/tag/redpanda-FILLMEIN) - YYYY-MM-DD

#### Added

#### Changed
* `image.repository` longer needs to be the default value of
  `"docker.redpanda.com/redpandadata/redpanda"` to respect version checks of
  `image.tag`
  ([#1334](https://github.com/redpanda-data/helm-charts/issues/1334)).
* BREAKING: A variety of fields, primarily prefixed with `extra` are now arrays
  of structured values rather than strings and no longer supported nested
  template evaluation. ([#1408](https://github.com/redpanda-data/helm-charts/pull/1408), [#1384](https://github.com/redpanda-data/helm-charts/pull/1384))

    Affected fields:
    - `post_upgrade_job.extraEnv`
    - `post_upgrade_job.extraEnvFrom`
    - `statefulSet.initContainers.*.extraVolumeMounts`
    - `statefulSet.sideCars.*.extraVolumeMounts`
    - `statefulSet.initContainers.extraInitContainers`
    - `statefulSet.extraVolumes`

    Fields that utilize nested templating for string generation (e.g. `external.domain`) are unaffected.

    Updating these fields is typically a case of needing to remove `|-`'s from one's values file.

    Before:
    ```yaml
    post_upgrade_job:
      extraEnv: |-
      - name: SPECIAL_LEVEL_KEY
          valueFrom:
            configMapKeyRef:
              name: special-config
              key: special.how
    ```

    After:
    ```yaml
    post_upgrade_job:
      extraEnv:
      - name: SPECIAL_LEVEL_KEY
        valueFrom:
          configMapKeyRef:
            name: special-config
            key: special.how
    ```

    In practice, this ability to nest template evaluation (see below) was not
    generally useful due to needing to know the internals of the chart and made
    it quite easy to accidentally forget required keys or generate invalid
    YAML.

    If you were using a templated value and would like to see it added back,
    please [file us an
    issue](https://github.com/redpanda-data/helm-charts/issues/new/choose) and
    tell us about your use case!

    ```yaml
    statefulSet.extraVolumes: |-
      {{ include "redpanda.fullname" }}
    ```

#### Fixed
* Numeric node/broker configurations are now properly transcoded as numerics.

#### Removed

## Redpanda Operator Chart
### [Unreleased](https://github.com/redpanda-data/helm-charts/releases/tag/operator-FILLMEIN) - YYYY-MM-DD
#### Added
#### Changed
#### Fixed
#### Removed

## Connectors Chart
### [Unreleased](https://github.com/redpanda-data/helm-charts/releases/tag/connectors-FILLMEIN) - YYYY-MM-DD
#### Added
#### Changed
#### Fixed
#### Removed

## Console Chart
### [Unreleased](https://github.com/redpanda-data/helm-charts/releases/tag/console-FILLMEIN) - YYYY-MM-DD
#### Added
#### Changed
#### Fixed
#### Removed

## Kminion Chart
### [Unreleased](https://github.com/redpanda-data/helm-charts/releases/tag/console-FILLMEIN) - YYYY-MM-DD
#### Added
#### Changed
#### Fixed
#### Removed
