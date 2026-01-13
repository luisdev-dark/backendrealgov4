# RealGo MVP - Backend

Backend ligero escrito en Go para el MVP de transporte con rutas fijas. Proporciona endpoints para listar rutas, consultar detalles y crear viajes siguiendo un flujo básico de reserva.

## Stack tecnológico

- **Go 1.21+**
- **Chi Router** – enrutamiento HTTP minimalista
- **pgxpool** – pool de conexiones PostgreSQL
- **PostgreSQL (compatible con Neon, Supabase u otro proveedor)**
- **Vercel** – despliegues serverless escritos para Go

## Estructura principal

```
backend/
├── api/                # Handlers HTTP expuestos a través de Vercel
├── internal/           # Lógica interna por dominio (db, http helpers, models, rutas, viajes)
├── sql/                # Scripts de esquema y seed para poblar la base
├── go.mod
go.sum
├── vercel.json         # Configuración de despliegue
└── README.md           # Este archivo
```

## Configuración local

1. Instala dependencias:

```bash
cd backend
go mod download
```

2. Crea un archivo `.env` con la siguiente variable:

```env
DATABASE_URL=postgresql://<usuario>:<contraseña>@<host>:<puerto>/<base_de_datos>?sslmode=require
```

3. Aplica el esquema y los datos de prueba desde la terminal o el panel de la base:

```bash
psql "$DATABASE_URL" -f sql/schema.sql
psql "$DATABASE_URL" -f sql/seed.sql
```

## Endpoints disponibles

- `GET /api/health` – confirma que el servicio responde
- `GET /api/routes` – lista rutas activas
- `GET /api/routes/{route_id}` – detalla una ruta con paradas
- `POST /api/trips` – crea un viaje en estado `requested`
- `GET /api/trips/{trip_id}` – obtiene el viaje por su identificador

Todos los endpoints devuelven JSON estándar y tratan errores usando `{ "error": "mensaje" }`.

## Despliegue en Vercel

1. Sube el repositorio a GitHub y conéctalo a Vercel.
2. Vercel detectará el entorno Go.
3. Configura variables de entorno desde el dashboard:

```env
DATABASE_URL=postgresql://<usuario>:<contraseña>@<host>:<puerto>/<base_de_datos>?sslmode=require
```

4. Cada push a la rama principal gatilla un nuevo deploy serverless.

## Notas adicionales

- El `user_id` usado en las pruebas está hardcodeado mientras no exista autenticación.
- `pickup_stop_id` y `dropoff_stop_id` son opcionales, pero si se envían deben ser distintos.
- Los métodos de pago admitidos son `cash`, `yape` y `plin`.

## Próximos pasos

- Implementar autenticación y autorización reales
- Añadir endpoints para drivers y matching
- Registrar historial de viajes y pagos
- Agregar notificaciones y tracking en tiempo real
