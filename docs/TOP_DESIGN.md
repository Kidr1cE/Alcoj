# Alcoj
## Overview
### Sequence diagrams
``` mermaid
---
title: Interactive
---
sequenceDiagram
    Workers->>Factory: Register
    Factory->>Workers: Dockerfile & runner.sh
    Workers->>Factory: Ready
    Factory->>admin: Ready
    admin->>Data Center: Set Question
    Data Center->>admin: 
    students->>Data Center: Get Question
    Data Center->>students: 
    students->>Factory: Code + any info
    Factory->>Workers: Code
    Workers->>Factory: Outputs & Others info
    Factory->>students: trated info
    Workers->>Data Center: Send submit info
```