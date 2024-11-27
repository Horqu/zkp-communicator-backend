# zkp-communicator-backend

## Repo structe:
zkp-communicator-backend/
├── api/                     # API definitions
│   ├── openapi.yaml         # API specification in OpenAPI format
│   └── docs/                # API documentation
├── cmd/                     # Main applications (entry points)
│   ├── api-gateway/         # Entry point for API Gateway
│   ├── auth-service/        # Entry point for Auth Service
│   ├── messaging-service/   # Entry point for Messaging Service
│   ├── zkp-service/         # Entry point for ZKP Service
│   ├── contacts-service/    # Entry point for Contacts Service
│   └── monitoring-service/  # Entry point for Monitoring Service
├── configs/                 # Configuration files
│   ├── dev/                 # Configurations for the development environment
│   ├── staging/             # Configurations for the staging environment
│   └── prod/                # Configurations for the production environment
├── docs/                    # Project documentation
│   ├── architecture.md      # Architecture documentation
│   ├── setup.md             # Application setup instructions
│   └── zkp-explained.md     # Explanation of ZKP functionality in the application
├── internal/                # Modules not exposed externally
│   ├── encryption/          # Encryption module
│   ├── auth/                # Login/registration module
│   ├── messaging/           # Messaging handling module
│   └── zkp/                 # ZKP module
├── pkg/                     # Shared modules between services
│   ├── utils/               # Utility functions and tools
│   ├── logger/              # Logger for the entire project
│   └── middleware/          # Middleware for gRPC/HTTP
├── deployments/             # Deployment configurations
│   ├── k8s/                 # Kubernetes configurations
│   ├── terraform/           # Terraform scripts for AWS
│   └── docker-compose.yaml  # Compose file for local setup
├── test/                    # Tests
│   ├── integration/         # Integration tests
│   ├── performance/         # Performance tests
│   └── unit/                # Unit tests
├── Makefile                 # Automation for building and running the application
├── go.mod                   # Go dependency management file
├── go.sum                   # Dependency checksum file for Go
├── README.md                # Main project documentation
└── .github/                 # GitHub configuration
    ├── workflows/           # CI/CD workflows
    └── dependabot.yml       # Automatic dependency updates