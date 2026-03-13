## ADDED Requirements

### Requirement: Cross-Platform Build
The system SHALL support building executable binaries for multiple operating systems (Linux, Windows, Darwin).

#### Scenario: Successful Windows build
- **WHEN** the user runs the build command targeting Windows
- **THEN** a `server.exe` binary is generated in the `bin/` directory

#### Scenario: Successful Linux build
- **WHEN** the user runs the build command targeting Linux
- **THEN** a `server` binary is generated in the `bin/` directory

### Requirement: Versioned Artifacts
The system SHALL include version information in the built binary.

#### Scenario: Binary version check
- **WHEN** the user runs the built binary with a `--version` or similar identifier (if implemented) or inspects properties
- **THEN** the binary reports the correct version and build commit hash
