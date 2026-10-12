# Readiness and Verification

This concept describes how the deployer checks whether a candidate app runtime is safe to keep moving through deployment.

These checks are deployer-facing HTTP endpoints exposed by the application. The deployer calls them after important deployment changes, such as preparing a candidate, migrating tenant state, and switching routing. The goal is to verify that the candidate still behaves correctly before the deployer moves on.

Exact endpoint request and response contracts are defined in [Candidate Readiness and Deployment Verification](../../specs/2_verification.md).

## Purpose

Readiness and verification prove different things at different phases.

| Phase | Purpose |
|---|---|
| Candidate readiness check | Prove the candidate app process is reachable and minimally ready without assuming migrated tenant state. |
| Post-migration internal verification | Prove app/env/tenant/domain sanity against migrated state before route switch. |
| Hostname verification | Prove real routed traffic reaches the intended runtime and tenant. |

The three checks deliberately overlap, but each adds a guarantee that could not safely be established at an earlier phase.

Before migration, the deployer calls the runtime directly over the internal network. At this point it only needs to establish that the process is reachable and that the application has the basic runtime configuration and inputs needed for deployment. The check cannot depend on tenant resolution or migrated tenant state.

After migration, the deployer can exercise the full application path for the tenant against the migrated tenant state. This is the first point where tenant-state-dependent checks are expected to work.

After the route switch, the deployer repeats that verification through a configured tenant hostname and normal request path. This adds proof that routing reaches the intended runtime and that the application resolves the expected tenant from the hostname. Reusing the same full verification keeps the application-side interface simple even though some checks are repeated.

## Candidate readiness check

This is a runtime-level check and does not require tenant context. It may check things that are safe before migration, such as expected mounts, environment values, naming expectations, and other basic runtime inputs. The check is meant as an early checkpoint for obvious candidate problems before the deployer performs backup or migration.

Because it runs before migration, it must not depend on tenant state already matching the target app version. Having a separate readiness check keeps these pre-migration assumptions distinct from the full post-migration verification path.

Docker healthchecks may exist as container-runtime signals, but they do not replace this application-level readiness check.

## Post-migration verification

This check establishes that the target app works with the tenant after migration. It should exercise the real application, tenant, and dependency paths needed to prove that the target app works with the migrated tenant state.

It is performed at two points in the deployment process: first directly against the runtime after migration and before the route switch, then again through a configured tenant hostname after the route switch.

For the internal pass, the deployer calls the runtime directly. Because there is no hostname routing context, tenant context may be supplied explicitly where the application needs it. The application should still resolve and validate that tenant through its normal application logic.

For the hostname pass, the deployer performs the same verification through a configured tenant hostname without supplying tenant context separately. The application resolves the tenant through its normal hostname/request handling, so this check also proves the tenant binding used by real routed traffic.

For a tenant with multiple hostnames, one successful hostname check is sufficient. This step verifies the routed end-to-end path rather than the availability of every configured hostname.

## Verification behavior

Deployment checks are observational and safe to repeat. They may inspect application state and dependencies, but they must not modify tenant state, start background work, or perform irreversible side effects.

Each invocation must perform the verification again. The application must not skip checks because an earlier invocation succeeded or because it recorded an internal "already verified" result.

The checks must remain usable while the tenant is frozen. In particular, post-migration verification must be able to exercise the required read and dependency paths without requiring tenant mutation.

## Deploy token

When the deployer creates a runtime, it generates a deploy token for that runtime and injects the token into the runtime configuration. The application returns this token in its readiness and verification responses.

During internal checks, matching the token confirms that the application received and can read that runtime input. During hostname verification, it additionally proves that routed traffic reached the expected runtime.

## Routing relationship

`verification` owns the HTTP checks and result interpretation. `routing` owns route rendering and route-target switching.

This keeps route changes separate from verification response validation.

## Verification after restart

Whenever the deployer restarts a runtime during deployment, it repeats the readiness and verification checks already completed for that runtime before continuing. For example, after a restart following post-migration internal verification, both the candidate readiness check and internal verification run again.

If the runtime contains a frozen tenant, the [restart verification invariant](tenant_freeze.md#restart-verification-invariant) for freeze state also applies.
