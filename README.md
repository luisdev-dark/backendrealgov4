# RealGo MVP - Backend

Esta carpeta contiene el servicio Go que alimenta la API del MVP RealGo. Sirve rutas fijas, expone necesidades de viaje básicas y mantiene la lógica interna separada en dominios claros.

## Stack tecnológico

- **Go 1.21+** – lenguaje principal del servicio
- **Chi Router** – manejo ligero de rutas HTTP
- **pgxpool** – pool seguro de conexiones con PostgreSQL
- **PostgreSQL compatible** – base de datos relacional para rutas, paradas y viajes
- **Vercel** – despliegues serverless con integración directa desde GitHub

## Organización del código

```
backend/
├── api/                # Entry point HTTP para Vercel
├── internal/           # Dominios: db, http helpers, models, routes, trips
├── sql/                # Esquema y datos semilla reutilizables
├── go.mod              # Dependencias Go
├── go.sum              # Sumas de dependencias
└── vercel.json         # Configuración de despliegue para Vercel
```

## Desarrollo local (alto nivel)

- Instala Go 1.21+ y asegúrate de que `go` esté disponible en tu PATH.
- Define la variable de entorno `DATABASE_URL` con la conexión hacia tu PostgreSQL (la forma exacta depende del proveedor que uses).
- Usa los scripts dentro de `sql/` para recrear el esquema y poblar datos de referencia en la base (psql o tu herramienta preferida).
- Ejecuta `go run ./api` o usa un comando equivalente en tu editor para levantar el servidor en modo local.

## Principales responsabilidades del backend

- Validar y responder `GET /api/health` para que el despliegue pueda comprobar el estado del servicio.
- Listar rutas activas y los detalles completos de cada ruta con sus paradas.
- Manejar la creación de viajes en estado `requested`, incluyendo validaciones básicas de pagos y paradas de subida/bajada.
- Permitir la consulta de viajes por identificador, devolviendo objetos JSON consistentes junto con errores estructurados.

## Pruebas y calidad

- El código usa paquetes internos para mantener separación de responsabilidades; puedes agregar tests en carpetas como `internal/routes` o `internal/trips` según lo necesites.
- Sigue añadiendo tests unitarios y de integración conforme se agreguen nuevos flows como autenticación o matching con drivers.

## Notas rápidas

- `DATABASE_URL` debe ofrecer SSL si tu proveedor lo requiere.
- Por ahora el `user_id` dentro de los seeds está estático; cualquier sistema de autenticación debe reemplazar este valor.
- Los métodos de pago permitidos y las tarifas se definen en los modelos, no en configuraciones externas.

## Futuro cercano

- Añadir autenticación real y permisos por rol.
- Exponer endpoints y lógica para drivers, matching y tracking.
- Mejorar el monitoreo de viajes y los registros de pagos para auditoría.
