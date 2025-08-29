# Plantilla de Proyecto Go con Arquitectura Limpia

Esta es una plantilla para proyectos Go que sigue los principios de Clean Architecture, diseñada para ser el punto de partida de aplicaciones escalables y mantenibles.

## 🚀 Características Principales

- ✅ **Arquitectura Limpia** con separación clara de capas
- 🏗️ **Estructura de Proyecto** organizada y escalable
- 🔒 **Autenticación JWT** integrada
- 🗄️ **MongoDB** como base de datos principal
- 🧪 **Pruebas unitarias** con ejemplos
- 📦 **Docker** y Docker Compose listos para producción
- 🔄 **GitHub Actions** para CI/CD
- 📝 **Documentación** detallada

## 📋 Requisitos Previos

- Go {{.go_version}} o superior
- MongoDB 5.0+
- Git
- Docker (opcional, para desarrollo con contenedores)

## 🛠️ Cómo Usar Esta Plantilla

### 1. Crear un Nuevo Proyecto

```bash
# Usar la plantilla con GitHub
gh repo create mi-nuevo-proyecto --template=tu-usuario/go-clean-architecture-template

# O clonar directamente
git clone --depth=1 https://github.com/tu-usuario/go-clean-architecture-template.git mi-nuevo-proyecto
cd mi-nuevo-proyecto
```

### 2. Configurar el Proyecto

1. Copia el archivo de configuración de ejemplo:
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   ```

2. Actualiza las variables de configuración en `configs/config.yaml` según sea necesario.

### 3. Instalar Dependencias

```bash
go mod download
```

### 4. Ejecutar el Proyecto

```bash
# Modo desarrollo
go run cmd/server/main.go

# O con variables de entorno
PORT=8080 go run cmd/server/main.go
```

## 🏗️ Estructura del Proyecto

```
.
├── cmd/                  # Punto de entrada de la aplicación
├── configs/             # Archivos de configuración
├── internal/            # Código fuente interno de la aplicación
│   ├── application/     # Casos de uso y lógica de negocio
│   ├── domain/          # Entidades y reglas de negocio
│   ├── infrastructure/  # Implementaciones concretas (DB, servicios externos)
│   ├── interfaces/      # Controladores, rutas y presentación
│   └── pkg/             # Librerías compartidas
├── scripts/             # Scripts de utilidad
└── test/                # Pruebas de integración
```

## 🔄 Variables de Entorno

Crea un archivo `.env` en la raíz del proyecto con las siguientes variables:

```env
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=mi_proyecto
JWT_SECRET=tu_clave_secreta_aqui
PORT=8080
```

## 🧪 Ejecutar Pruebas

```bash
# Todas las pruebas
go test -v ./...

# Pruebas con cobertura
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## 🐳 Docker

```bash
# Construir la imagen
docker build -t mi-proyecto .

# Ejecutar con Docker Compose
docker-compose up -d
```

## 🤝 Contribuir

1. Haz un fork del proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Haz commit de tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Haz push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Distribuido bajo la licencia MIT. Ver `LICENSE` para más información.

## 📧 Contacto

Tu Nombre - [@tuusuario](https://twitter.com/tuusuario) - email@ejemplo.com

Enlace del Proyecto: [https://github.com/tuusuario/go-clean-architecture-template](https://github.com/tuusuario/go-clean-architecture-template)
