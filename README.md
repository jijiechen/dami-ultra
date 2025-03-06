# Service Catalog Microservice Design Document

## Overview

This document outlines the design and implementation of a service catalog microservice that enables users to view and manage services within their organization.

## User Story & Requirements

### Business Context
The service catalog is designed to provide users with a comprehensive overview of available services in their organization, enhancing service discovery and management.

### User Story
As a user, I need to view and manage services in my organization.

### Acceptance Criteria
1. Users can view service details including:
   - Service name
   - Brief description
   - Available versions
2. Users can navigate to specific services via service cards
3. Users can search for services using search functionality

### UI Design Reference
Design mockups are available in Figma:
[Service Card List Design](https://www.figma.com/file/zeaWiePnc3OCe34I4oZbzN/Service-Card-List?node-id=0%3A1)

## Technical Architecture

### High-Level Design
The service follows a three-layer architecture pattern:
```
┌───────────────┐
│ API Layer     │ → Request handling & routing
├───────────────┤
│ Business Layer│ → Core business logic
├───────────────┤
│ Repository    │ → Data persistence
└───────────────┘
```

### Component Details

#### 1. API Layer
- Located in `internal` package
- Implements a context-based API signature: `func(ctx context.Context, r R) (T, error)`
- Framework-agnostic design allowing easy adaptation to different API frameworks
- Current implementation: Gin framework

#### 2. Business Logic Layer
- Interface-based design enabling easy implementation swapping
- Dependency injection handled via `dig` in the `startup` package
- Clear separation of concerns from other layers

#### 3. Repository Layer
- Interface-based data access layer
- Current implementation: SQLite3 with GORM
- Designed for multi-tenancy support

### Multi-tenancy Support
- Organization separation through different databases
- Organization identification via `x-organization-id` header
- Middleware handles:
  - Database connection management per request
  - Connection pooling
  - Context injection

## Implementation Details

### API Endpoints

| Method | Endpoint        | Description           |
|--------|-----------------|-----------------------|
| GET    | /service/list   | List all services     |
| POST   | /service        | Create a new service  |
| GET    | /service/:id    | Get service by ID     |
| PUT    | /service        | Update service        |

### Technical Considerations
- Configuration management currently uses environment variables
- Authentication is not implemented in the current version
- Database middleware uses a simple connection management strategy
- Future improvements planned for connection pooling
- External configuration management integration planned

## Testing Strategy

### Behavior Tests
- Concurrent test execution
- Integration tests covering API endpoints
- Test suites:
  1. ServiceList
  2. ServiceCreateAndGet
  3. ServiceUpdate

### Running Tests
```bash
# Run all behavior tests
GO_DEBUG=1 go test -v ./behavior_tests/... -run TestBehavior

# Run specific test suite
GO_DEBUG=1 go test -v ./behavior_tests/... -run TestBehavior/ServiceList
```

## Future improvements

- Authentication and Authorization
- Implement connection pooling Add external configuration management 
- Additional Test Coverage
- Unit tests for business logic

## Technical Debt and Considerations

- Current lack of authentication system
- Simple database connection management
- Limited test coverage
- Hardcoded configurations
