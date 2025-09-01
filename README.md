# 🚀 Plantilla de Proyecto Go con Arquitectura Limpia

Plantilla para proyectos Go que sigue los principios de Clean Architecture, diseñada para ser el punto de partida de aplicaciones escalables y mantenibles.

## 🌟 Características

- 🏗️ **Arquitectura Limpia** con separación clara de capas
- 🔒 **Autenticación JWT** integrada
- 🗄️ **MongoDB** como base de datos principal
- 🧪 **Pruebas unitarias** con ejemplos
- 📦 **Docker** y Docker Compose listos para producción
- 🔄 **GitHub Actions** para CI/CD
- 📝 **Documentación** detallada

## 🚀 Cómo Usar Esta Plantilla

### 1. Crear un Nuevo Proyecto

```bash
# Usar la plantilla con GitHub CLI
gh repo create mi-nuevo-proyecto --template=tu-usuario/go-clean-architecture-template

# O clonar directamente
git clone --depth=1 https://github.com/tu-usuario/go-clean-architecture-template.git mi-nuevo-proyecto
cd mi-nuevo-proyecto
```

### 2. Inicializar el Proyecto

```bash
# Hacer ejecutable el script de inicialización
chmod +x init.sh

# Ejecutar el script de inicialización
./init.sh

# Seguir las instrucciones en pantalla para configurar tu proyecto
```

### 3. Configurar el Entorno

1. Copia el archivo de configuración de ejemplo:
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   ```

2. Actualiza las variables de configuración en `configs/config.yaml` según sea necesario.

### 3.1 Configuración JSON

La aplicación también utiliza un archivo JSON para configuraciones específicas. Este archivo debe estar ubicado en la ruta especificada en la configuración YAML bajo `app.json_config_path` (por defecto `./config/parameters.json`).

#### Estructura del archivo JSON:

```json
{
  "integrationPaths": [
    {
      "name": "ejemplo.ruta",
      "value": "/ruta/ejemplo"
    }
  ],
  "certificates": [
    {
      "name": "nombre.certificado",
      "value": "contenido-del-certificado"
    }
  ],
  "params": [
    {
      "name": "parametro.ejemplo",
      "value": "valor-del-parametro"
    }
  ]
}
```

#### Uso en el código:

```go
import "api-ptf-core-business-orchestrator-go-ms/internal/pkg/utils"

// Obtener un parámetro
if value, found := utils.GetParam("parametro.ejemplo"); found {
    // Usar el valor
}

// Obtener un parámetro con valor por defecto
value := utils.GetParamOrDefault("parametro.inexistente", "valor-por-defecto")
```

### 4. Iniciar el Servidor de Desarrollo

```bash
# Modo desarrollo
go run cmd/server/main.go

# O con variables de entorno
PORT=8080 go run cmd/server/main.go
```

## 🚦 Requisitos Previos

- Go 1.21 o superior
- MongoDB 6.0 o superior
- Git

## ⚡ Instalación Rápida

1. Clonar el repositorio:
   ```bash
   git clone https://github.com/tu-usuario/api-ptf-core-business-orchestrator-go-ms.git
   cd api-ptf-core-business-orchestrator-go-ms
   ```

2. Instalar dependencias:
   ```bash
   go mod download
   ```

3. Configurar variables de entorno:
   ```bash
   cp .env.example .env
   # Editar el archivo .env con tus credenciales
   ```

4. Iniciar el servidor:
   ```bash
   go run cmd/server/main.go
   ```

## ⚙️ Configuración

El archivo `configs/config.yaml` contiene la configuración base. Las variables de entorno tienen prioridad:

```env
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=ptf-core
JWT_SECRET=tu_clave_secreta_aqui
```

## 📦 Constantes del Proyecto

El proyecto utiliza constantes para mantener consistencia en los nombres de rutas y métodos HTTP. Estas constantes se encuentran en el paquete `internal/pkg/constants/`.

### Métodos HTTP

```go
// internal/pkg/constants/Methods.go
package constants

const (
    GET     = "GET"
    POST    = "POST"
    PUT     = "PUT"
    DELETE  = "DELETE"
    PATCH   = "PATCH"
    OPTIONS = "OPTIONS"
    HEAD    = "HEAD"
)
```

### Rutas de la API

```go
// internal/pkg/constants/Paths.go
package constants

const (
    UTILS_GROUP  = ""
    HEALTH_CHECK = "/health"
    RSYNC        = "/rsync"
    USER_GROUP   = "/users"
    PREFIX       = USER_GROUP + "/examples/"
)
```

### Uso en las rutas

Al definir nuevas rutas, se recomienda utilizar estas constantes para mantener la consistencia:

```go
import "api-ptf-core-business-orchestrator-go-ms/internal/pkg/constants"

// Ejemplo de uso en la definición de rutas
router.HandleFunc(constants.HEALTH_CHECK, healthHandler).Methods(constants.GET)
```

## 🏗️ Estructura del Proyecto

```
.
├── cmd/
│   └── server/           # Punto de entrada
├── configs/              # Configuraciones
├── internal/
│   ├── application/      # Lógica de negocio
│   ├── config/           # Configuración
│   ├── domain/           # Modelos y reglas de negocio
│   ├── infrastructure/   # Implementaciones concretas
│   ├── interfaces/       # Controladores y rutas
│   └── pkg/              # Utilidades
└── scripts/              # Scripts de utilidad
```

## 🛠️ Guía de Desarrollo

### 🔄 Implementación de Repositorios

### Repositorio Genérico

El proyecto incluye un `GenericRepository` que proporciona operaciones CRUD básicas para cualquier entidad.

#### Uso Básico

```go
// 1. Definir tu entidad en el dominio
// internal/domain/tu_entidad.go
type TuEntidad struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Nombre      string             `bson:"nombre"`
    Descripcion string             `bson:"descripcion"`
    DateCreated time.Time          `bson:"date_created"`
    DateUpdated time.Time          `bson:"date_updated"`
}

// 2. Crear un repositorio específico
// internal/infrastructure/repository/tu_entidad_repo.go
type TuEntidadRepository struct {
    Repo *GenericRepository[domain.TuEntidad]
}

func NewTuEntidadRepository(db *database.Database) *TuEntidadRepository {
    return &TuEntidadRepository{
        Repo: NewGenericRepository[domain.TuEntidad](db, "tu_coleccion"),
    }
}
```

#### Métodos Disponibles en GenericRepository

Los siguientes métodos están disponibles por defecto en el `GenericRepository`:
- `FindByID(ctx, id)`: Busca por ID
- `FindAll(ctx, filter, opts)`: Lista con filtros
- `Create(ctx, entity)`: Crea nueva entidad
- `Update(ctx, id, entity)`: Actualiza por ID
- `Delete(ctx, id)`: Elimina por ID

### Repositorio Personalizado

Para operaciones más complejas, puedes extender el repositorio genérico con métodos personalizados:

```go
// internal/infrastructure/repository/tu_entidad_repo.go
func (r *TuEntidadRepository) MetodoPersonalizado(ctx context.Context, parametro string) ([]domain.TuEntidad, error) {
    filter := bson.M{"campo": parametro}
    
    cursor, err := r.Repo.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var resultados []domain.TuEntidad
    if err := cursor.All(ctx, &resultados); err != nil {
        return nil, err
    }
    
    return resultados, nil
}
```

## Registro en la Aplicación

Para registrar tu nuevo repositorio en la aplicación:

```go
// En tu archivo de rutas
func RegisterTusRutas(router *mux.Router, app *models.Application) {
    // Inicializar con el DB() getter
    repo := repository.NewTuEntidadRepository(app.DB())
    
    // Resto de la inicialización...
}
```


## 🛣️ Cómo Agregar un Nuevo Manejador (Handler) y sus Rutas

Para agregar un nuevo conjunto de rutas y sus manejadores al proyecto, sigue estos pasos:

### 1. Crear un nuevo paquete de manejadores (Handlers)

Primero, crea un nuevo archivo de manejador en `internal/interfaces/http/handlers/`. Los manejadores son los responsables de procesar las peticiones HTTP y devolver las respuestas:

```go
// internal/interfaces/routes/domain/product.go
package domainRoutes

import (
    "api-ptf-core-business-orchestrator-go-ms/internal/application"
    "api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/repository"
    "api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/handlers"
    "api-ptf-core-business-orchestrator-go-ms/internal/models"
    "github.com/gorilla/mux"
)

// RegisterProductRoutes registra las rutas de productos
func RegisterProductRoutes(router *mux.Router, a *models.Application) {
    // Inicializar repositorios
    productRepo := repository.NewMongoProductRepository(a.DB())

    // Inicializar servicios
    productService := application.NewProductService(productRepo)

    // Inicializar manejadores
    productHandler := handlers.NewProductHandler(productService)

    // Configurar rutas
    productRouter := router.PathPrefix("/products").Subrouter()
    productRouter.HandleFunc("", productHandler.ListProducts).Methods("GET")
    productRouter.HandleFunc("", productHandler.CreateProduct).Methods("POST")
    productRouter.HandleFunc("/{id}", productHandler.GetProduct).Methods("GET")
    productRouter.HandleFunc("/{id}", productHandler.UpdateProduct).Methods("PUT")
    productRouter.HandleFunc("/{id}", productHandler.DeleteProduct).Methods("DELETE")
}
```

### 2. Configurar las rutas con sus manejadores

Crea un nuevo archivo en `internal/interfaces/routes/domain/` para configurar las rutas y asociarlas con sus respectivos manejadores. Este es un ejemplo para un módulo de productos:

```go
// internal/interfaces/routes/domain/product.go
package domainRoutes

import (
    "api-ptf-core-business-orchestrator-go-ms/internal/application"
    "api-ptf-core-business-orchestrator-go-ms/internal/infrastructure/repository"
    "api-ptf-core-business-orchestrator-go-ms/internal/interfaces/http/handlers"
    "api-ptf-core-business-orchestrator-go-ms/internal/models"
    "github.com/gorilla/mux"
)

// RegisterProductRoutes configura las rutas para los productos y sus manejadores correspondientes
func RegisterProductRoutes(router *mux.Router, a *models.Application) {
    // 1. Inicializar repositorio
    productRepo := repository.NewMongoProductRepository(a.DB())
    
    // 2. Inicializar servicio de aplicación
    productService := application.NewProductService(productRepo)
    
    // 3. Inicializar manejador (handler)
    productHandler := handlers.NewProductHandler(productService)
    
    // 4. Configurar rutas y asociarlas con los métodos del manejador
    productRouter := router.PathPrefix("/products").Subrouter()
    
    // Asociar rutas HTTP con los métodos del manejador
    productRouter.HandleFunc("", productHandler.ListProducts).Methods("GET")       // GET /products
    productRouter.HandleFunc("", productHandler.CreateProduct).Methods("POST")    // POST /products
    productRouter.HandleFunc("/{id}", productHandler.GetProduct).Methods("GET")   // GET /products/{id}
    productRouter.HandleFunc("/{id}", productHandler.UpdateProduct).Methods("PUT") // PUT /products/{id}
    productRouter.HandleFunc("/{id}", productHandler.DeleteProduct).Methods("DELETE") // DELETE /products/{id}
}
```

### 3. Registrar las nuevas rutas en el enrutador principal

Abre el archivo `internal/interfaces/routes/routes.go` y agrega la función de registro de tus nuevas rutas al enrutador principal:

```go
// internal/interfaces/routes/routes.go
package routes

import (
    uD "api-ptf-core-business-orchestrator-go-ms/internal/interfaces/routes/domain"
    uR "api-ptf-core-business-orchestrator-go-ms/internal/interfaces/routes/utils"
    "api-ptf-core-business-orchestrator-go-ms/internal/models"
    "github.com/gorilla/mux"
)

// SetupRoutes configura todas las rutas
func SetupRoutes(router *mux.Router, a *models.Application) {
    uR.RegisterInfoRoutes(router)
    uR.RegisterRysncRoutes(router)
    uD.RegisterUserRoutes(router, a)
    uD.RegisterProductRoutes(router, a) // Nueva ruta agregada
}
```

### 3. Crear el manejador (handler)

Crea un nuevo manejador en `internal/interfaces/http/handlers/`:

```go
// internal/interfaces/http/handlers/product_handler.go
package handlers

import (
    "api-ptf-core-business-orchestrator-go-ms/internal/application"
    "net/http"
)

type ProductHandler struct {
    productService application.ProductService
}

func NewProductHandler(ps application.ProductService) *ProductHandler {
    return &ProductHandler{
        productService: ps,
    }
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
    // Implementar lógica para listar productos
}

// Implementar otros métodos (GetProduct, CreateProduct, etc.)
```

### 4. Crear el servicio de aplicación

Crea el servicio de aplicación correspondiente:

```go
// internal/application/product_service.go
package application

type ProductService interface {
    ListProducts() ([]domain.Product, error)
    GetProduct(id string) (*domain.Product, error)
    CreateProduct(product *domain.Product) error
    UpdateProduct(id string, product *domain.Product) error
    DeleteProduct(id string) error
}

type productService struct {
    repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) ProductService {
    return &productService{
        repo: repo,
    }
}

// Implementar métodos de la interfaz...
```

### 5. Crear el repositorio

Implementa el repositorio en `internal/domain/` e `internal/infrastructure/repository/` siguiendo el patrón existente.

## 📚 Documentación de la API

La documentación de la API está disponible en formato OpenAPI (Swagger).

1. Inicia el servidor
2. Navega a: `http://localhost:8080/api/v1/docs`

## 🔧 Endpoints

### Usuarios

- `GET /api/v1/users` - Listar usuarios (paginado)
- `GET /api/v1/users/{id}` - Obtener usuario por ID
- `POST /api/v1/users` - Crear nuevo usuario
- `PUT /api/v1/users/{id}` - Actualizar usuario
- `DELETE /api/v1/users/{id}` - Eliminar usuario
- `GET /api/v1/users/email/{email}` - Buscar usuario por email

## 🧪 Pruebas

Para ejecutar las pruebas:

```bash
go test -v ./...
```

## 🐳 Docker

Construir la imagen:

```bash
docker build -t ptf-core-business-orchestrator .
```

Ejecutar el contenedor:

```bash
docker run -p 8080:8080 --env-file .env ptf-core-business-orchestrator
```

## 🚀 Despliegue

### Construir la imagen Docker:

```bash
docker build -t ptf-core-business-orchestrator .
```

### Ejecutar en producción:

```bash
docker run -p 8080:8080 --env-file .env ptf-core-business-orchestrator
```

## 🧪 Pruebas

Para ejecutar las pruebas:

```bash
go test -v ./...
```

## 📚 Documentación de la API

La documentación interactiva está disponible en:
- `http://localhost:8080/api/v1/docs` (después de iniciar el servidor)

### Endpoints disponibles:

#### Usuarios
- `GET /api/v1/users` - Listar usuarios (paginado)
- `GET /api/v1/users/{id}` - Obtener usuario por ID
- `POST /api/v1/users` - Crear nuevo usuario
- `PUT /api/v1/users/{id}` - Actualizar usuario
- `DELETE /api/v1/users/{id}` - Eliminar usuario
- `GET /api/v1/users/email/{email}` - Buscar usuario por email

## 🏗️ Estructura del Proyecto

```
.
├── cmd/                 # Punto de entrada de la aplicación
├── configs/             # Archivos de configuración YAML
├── internal/
│   ├── application/     # Lógica de negocio y casos de uso
│   ├── client/          # Clientes HTTP/API externas
│   ├── config/          # Configuración de la aplicación
│   ├── domain/          # Entidades del dominio
│   ├── interfaces/      # Controladores HTTP y rutas
│   └── pkg/             # Utilidades compartidas
```

## ⚙️ Configuración

### Archivo de Configuración YAML

El archivo principal de configuración está en `configs/config.yaml` y sigue esta estructura:

```yaml
application_name: "api-core-template-go-ms"
description: "Template for new microservices in Go"
application_version: "1.0.0"
environment: "dev"

http:
  port: "8426"
  base_path: "/api/business-orchestrator/v1"
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

app:
  mongodb:
    uri: ${MONGO_URI}
    database: ${MONGO_DATABASE}
    timeout: ${TIMEOUT}
  
  # Ruta al archivo de configuración JSON
  json_config_path: "${JSON_CONFIG_PATH}"
```

### Estructura del Archivo de Configuración JSON

El archivo de configuración JSON (especificado en `app.json_config_path`) debe seguir esta estructura:

```json
{
  "integrationPaths": [
    {
      "name": "examples.one.domain",
      "value": "api.restful-api.dev"
    },
    {
      "name": "examples.one.port",
      "value": "443"
    },
    {
      "name": "examples.one.path",
      "value": "/objects"
    }
  ]
}
```

## 🌐 Endpoints de la API

### Ejemplos
- `GET /api/business-orchestrator/v1/examples` - Obtener lista de ejemplos de API externa

### Salud
- `GET /api/business-orchestrator/v1/health` - Verificar estado del servicio

## 🚀 Despliegue

### Local con Docker

```bash
docker-compose up -d
```

### Producción

1. Construir la imagen:
   ```bash
   docker build -t myapp .
   ```

2. Ejecutar el contenedor:
   ```bash
   docker run -d -p 8426:8426 --env-file .env myapp
   ```

## 🧪 Pruebas

Ejecutar todas las pruebas:
```bash
go test ./...
```

Con cobertura:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📄 Licencia

Este proyecto está bajo la [Licencia MIT](LICENSE).

## 🤝 Contribución

¡Las contribuciones son bienvenidas! Por favor:
1. Haz un fork del proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Haz commit de tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Haz push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Changelog

Ver [CHANGELOG.md](CHANGELOG.md) para el historial de cambios.

## 📧 Contacto

Para consultas o soporte, contacta al equipo de desarrollo o abre un issue en el repositorio.
