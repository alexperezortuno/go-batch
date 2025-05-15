# go-batch

Project structure
```
/data-loader/
├── cmd/
│   └── loader/
│       └── main.go          # Punto de entrada principal
├── internal/
│   ├── config/
│   │   └── config.go       # Manejo de configuración
│   ├── domain/
│   │   └── models.go       # Modelos de datos
│   ├── repository/
│   │   ├── database.go     # Conexión DB
│   │   └── loader_repo.go  # Operaciones DB específicas
│   ├── service/
│   │   ├── loader_service.go # Lógica de negocio
│   │   └── validator.go    # Validación de datos
│   ├── handler/
│   │   ├── file_handler.go # Manejo de archivos
│   │   └── metrics.go     # Endpoints de métricas
│   └── utils/
│       ├── logger.go       # Utilidades de logging
│       └── error.go        # Manejo de errores
├── pkg/
│   └── fileprocessor/      # Librería reusable
│       ├── csv/
│       └── excel/
├── migrations/             # Scripts de migración DB
├── scripts/               # Scripts auxiliares
├── configs/
│   ├── config.yaml         # Configuración base
│   └── config.prod.yaml    # Configuración producción
├── Makefile               # Automatización
├── go.mod
└── go.sum
```

### Run application 

in local

```bash
go run cmd/loader/main.go
```

build

```bash
go build -o dist/go_batch cmd/loader/main.go
```
