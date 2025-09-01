# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [1.1.1](https://github.com/pisondev/rest-todo-api/compare/v1.1.0...v1.1.1) (2025-09-01)


### Bug Fixes

* ensure dueDate.Valid attribute is false if it is nil ([a503485](https://github.com/pisondev/rest-todo-api/commit/a503485fb36c2247b88437668fdb72a4e5bb20c2))

## [1.1.0](https://github.com/pisondev/rest-todo-api/compare/v1.0.0...v1.1.0) (2025-09-01)


### Features

* add DueDate field for Update Feature ([cde8c91](https://github.com/pisondev/rest-todo-api/commit/cde8c91dab6a37c60f2e7e1fcea4a8def1962d4f))
* add logger from sirupsen/logrus ([a681f78](https://github.com/pisondev/rest-todo-api/commit/a681f787939e1bda8308a8ad3af56edb29d72150))
* add logger in each layer


### Bug Fixes

* Ensure JSON response is an empty array
* handle sql.ErrNoRows properly ([497251e](https://github.com/pisondev/rest-todo-api/commit/497251e5d9860860e677221f137eb3dfe0ba575e))

## [1.0.0] - 2025-08-17

### Added
- Initial release of the REST To-Do API.
- User authentication system with JWT for Register and Login.
- Full CRUD functionality for tasks with soft deletes.
- Task filtering by status and due date.
- Layered architecture with Controller, Service, and Repository layers.
- OpenAPI specification for API documentation.