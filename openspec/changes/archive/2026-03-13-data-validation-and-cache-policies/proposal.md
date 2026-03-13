## Why

To transition from a prototype to a production-grade system, we must ensure data integrity and resource efficiency. Currently, the backend lacks a strict validation layer for game stats (e.g., preventing impossible health values) and lacks explicit caching policies, which could lead to resource leakage in Redis.

## What Changes

- **Defensive Validation Layer**: Implement a centralized validation service using `go-playground/validator` to enforce min/max bounds on survival stats (Health, Hunger, Sanity, Warmth).
- **Inbound Request Hardening**: Update DTOs with validation tags to reject malformed or "cheated" data at the delivery layer.
- **Cache TTL Policies**: Define and implement explicit expiration times for cached player data in Redis to prevent orphaned records.
- **Domain-Level Integrity**: Move business rule validation into the Usecase layer to ensure consistency regardless of the delivery method.

## Capabilities

### New Capabilities
- `data-validation`: Centralized service for enforcing business rules and safety bounds on inbound data.

### Modified Capabilities
- `save-system`: Update save requirements to include statutory bounds check for survival metrics.
- `auth`: Add validation for character and origin selection during initialization.

## Impact

- **Security**: Basic anti-cheat protection against modified client memory values.
- **Stability**: Prevents "Garbage In, Garbage Out" scenarios that could crash or corrupt player saves.
- **Cost/Efficiency**: Redis memory usage is optimized through strict TTL policies.
