# Security and Secrets

## Privileged operator access

The Deployer API is an administrative interface intended only for trusted operators. Deployments must restrict access to trusted operator identities and protect API traffic from untrusted parties. The current security model does not distinguish operator roles; an operator with API access is trusted to use all supported operations, including APIs that return secret-bearing runtime inputs.

Access to deployer-managed host state, runtime metadata, and configured secret providers is part of the same trusted operational boundary.

## Registry credentials

Registry credentials are referenced from global registry profiles and used only for registry interactions. Exact registry auth mechanics are specified in [Registry Auth and Tag Discovery](../../specs/0_registry_auth_and_tag_discovery.md).

## Deploy tokens

Deploy tokens are per-runtime operational proof tokens used by readiness and verification checks. During routed verification, they prove that the response came from the expected runtime.

Deploy tokens should be treated as sensitive and must not be logged raw.

## Backups

Backup artifacts are sensitive recovery data. The canonical backup path, retention, file mode, and capacity policy are in [Appendix - Queue and Backup Policy](../../appendix/06_queue_and_backup_policy.md).

## Logging redaction

Never log registry credentials, bearer tokens, authorization headers, raw connection strings containing credentials, or deploy tokens.

Redaction rules are canonical in [Appendix - Logging Contract](../../appendix/10_logging_contract.md).

## Runtime inputs

Configured and materialized runtime files and environment values may contain sensitive configuration or secrets. Their host storage is private deployer state: configured-input scopes and runtime materializations are accessible only to the deployer service account. Materialized files that are mounted into application containers remain protected on the host by their enclosing deployer-private directory while being readable by the application process. Exact host ownership and file modes are defined in [Appendix – Directory Structure](../../appendix/07_directory_structure.md).

Configured or materialized runtime-input values and file contents must not be written to deployer logs or diagnostic output. The canonical redaction rule is in [Appendix – Logging Contract](../../appendix/10_logging_contract.md).

The materialization fingerprint is operational runtime metadata, not a secret or authentication value. Access to deployer-managed runtime metadata is part of the trusted operational boundary. Consumers treat the fingerprint as opaque, but its derivation does not need to resist offline guessing by someone already trusted to inspect that metadata. A keyed fingerprint is therefore not required.

The canonical runtime-input model and persisted shapes are defined in [State and Persistence](state_and_persistence.md#configured-runtime-input-and-materialization) and [Appendix - Persistence Schemas](../../appendix/03_persistence_schemas.md#runtime-input-requirements-canonical). The deployment procedure is described in [Runtime Input Materialization](../../specs/1.2_runtime_input_materialization.md).
