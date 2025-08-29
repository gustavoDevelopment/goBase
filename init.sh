#!/bin/bash

# Colores para la salida
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Función para mostrar mensajes de error
error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Función para mostrar mensajes de éxito
success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

# Función para mostrar advertencias
warning() {
    echo -e "${YELLOW}[ADVERTENCIA]${NC} $1"
}

echo -e "\n${GREEN}🚀 Inicialización del proyecto Go a partir de la plantilla${NC}\n"

# 1. Verificar que se ejecute desde la raíz del proyecto
if [ ! -f "go.mod" ]; then
    error "Por favor, ejecuta este script desde el directorio raíz del proyecto."
fi

# 2. Obtener el nombre del módulo actual
CURRENT_MODULE=$(grep '^module' go.mod | awk '{print $2}')
echo "Módulo actual: $CURRENT_MODULE"

# 3. Solicitar el nuevo nombre del módulo
echo -e "\nIngresa el nombre del módulo para tu proyecto (ej: github.com/tu-usuario/mi-proyecto):"
read -r NEW_MODULE

if [ -z "$NEW_MODULE" ]; then
    error "Debes proporcionar un nombre de módulo válido."
fi

# 4. Solicitar el nombre del proyecto
echo -e "\nIngresa el nombre del proyecto (usado en documentación y configuraciones):"
read -r PROJECT_NAME

if [ -z "$PROJECT_NAME" ]; then
    error "Debes proporcionar un nombre de proyecto."
fi

# 5. Actualizar el nombre del módulo en los archivos
if [ "$CURRENT_MODULE" != "$NEW_MODULE" ]; then
    echo -e "\nActualizando referencias del módulo..."
    
    # Reemplazar en archivos .go
    find . -type f -name "*.go" -exec sed -i '' "s|$CURRENT_MODULE|$NEW_MODULE|g" {} \;
    
    # Reemplazar en otros archivos relevantes
    find . -type f \( -name "*.mod" -o -name "*.yaml" -o -name "*.yml" -o -name "*.json" \) \
        -exec sed -i '' "s|$CURRENT_MODULE|$NEW_MODULE|g" {} \;
    
    # Actualizar el go.mod
    sed -i '' "s|module .*|module $NEW_MODULE|g" go.mod
    
    success "Módulo actualizado a: $NEW_MODULE"
else
    warning "El módulo ya está configurado como $NEW_MODULE"
fi

# 6. Actualizar el nombre del proyecto en los archivos de configuración
if [ -n "$PROJECT_NAME" ]; then
    echo -e "\nActualizando nombre del proyecto..."
    
    # Actualizar en archivos de configuración
    find . -type f \( -name "*.yaml" -o -name "*.yml" -o -name "*.json" -o -name "*.md" \) \
        -exec sed -i '' "s/ProjectName:.*/ProjectName: $PROJECT_NAME/g" {} 2>/dev/null \;
    
    success "Nombre del proyecto actualizado a: $PROJECT_NAME"
fi

# 7. Inicializar un nuevo repositorio Git si no existe
if [ ! -d ".git" ]; then
    echo -e "\nInicializando nuevo repositorio Git..."
    git init
    
    # Crear .gitignore si no existe
    if [ ! -f ".gitignore" ]; then
        echo "Creando archivo .gitignore..."
        cat > .gitignore << 'EOL'
# Binarios para programas y plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Archivos de entorno
.env
.env.local

# Directorios de dependencias
vendor/

# Directorios de IDE
.idea/
.vscode/

# Archivos de sistema operativo
.DS_Store
Thumbs.db
EOL
    fi
    
    # Hacer commit inicial
    git add .
    git commit -m "chore: initial commit from template"
    
    success "Repositorio Git inicializado"
else
    warning "El repositorio Git ya está inicializado"
fi

# 8. Instalar dependencias
echo -e "\nInstalando dependencias..."
go mod download

# 9. Limpiar dependencias no utilizadas
go mod tidy

# 10. Mensaje final
echo -e "\n${GREEN}✅ Proyecto inicializado exitosamente!${NC}"
echo -e "\nPróximos pasos:"
echo "1. Revisa y actualiza la configuración en el directorio 'configs/'"
echo "2. Agrega tus credenciales en un archivo .env (ver .env.example)"
echo "3. Ejecuta 'go run cmd/server/main.go' para iniciar el servidor"
echo -e "\n¡Feliz codificación! 🚀\n"

exit 0
