# RealGo MVP - Backend

MVP de transporte con rutas fijas A→B y paradas intermedias (anexos).

## Stack Tecnológico

- **Go 1.21+**
- **Chi Router** - Router HTTP ligero
- **pgxpool** - Connection pool para PostgreSQL
- **PostgreSQL (Neon)** - Base de datos serverless
- **Vercel** - Deployment serverless

## Estructura del Proyecto

```
backend/
├── api/
│   └── index.go          # Handler principal para Vercel
├── internal/
│   ├── db/
│   │   └── db.go         # Pool de conexiones (singleton)
│   ├── httpx/
│   │   └── respond.go    # Helpers para respuestas JSON
│   ├── models/
│   │   └── models.go     # Structs de datos
│   ├── routes/
│   │   └── handlers.go   # Handlers de rutas
│   └── trips/
│       └── handlers.go   # Handlers de viajes
├── sql/
│   ├── schema.sql        # Esquema de base de datos
│   └── seed.sql          # Datos de prueba
├── go.mod
├── go.sum
├── vercel.json           # Configuración de Vercel
└── README.md
```

## Setup Local (Opcional)

### 1. Instalar dependencias

```bash
cd backend
go mod download
```

### 2. Configurar variable de entorno

Crea un archivo `.env` en la carpeta `backend/`:

```bash
DATABASE_URL=postgresql://user:password@host:port/dbname?sslmode=require
```

### 3. Aplicar schema y seed en Neon

Conéctate a tu base de datos Neon y ejecuta:

```bash
# Aplicar schema
psql $DATABASE_URL -f sql/schema.sql

# Aplicar seed data
psql $DATABASE_URL -f sql/seed.sql
```

O desde la consola web de Neon, copia y pega el contenido de cada archivo.

## Configuración de Base de Datos (Neon)

### Variable de Entorno Requerida

```
DATABASE_URL=postgresql://user:password@ep-xxx.region.aws.neon.tech/dbname?sslmode=require
```

**Importante:** La URL debe incluir `sslmode=require` para conexiones seguras.

### UUIDs de Prueba (del seed.sql)

- **Usuario (Passenger):** `11111111-1111-1111-1111-111111111111`
- **Ruta 1 (Centro-Norte):** `22222222-2222-2222-2222-222222222222`
- **Ruta 2 (Sur-Este):** `33333333-3333-3333-3333-333333333333`

## Endpoints Disponibles

### Health Check

```bash
GET /api/health
```

**Respuesta:**
```json
{
  "status": "ok",
  "service": "realgo-mvp"
}
```

### Listar Rutas

```bash
GET /api/routes
```

**Respuesta:**
```json
[
  {
    "id": "22222222-2222-2222-2222-222222222222",
    "name": "Ruta Centro - Norte",
    "origin_name": "Plaza de Armas",
    "destination_name": "Terminal Norte",
    "base_price_cents": 1550,
    "currency": "PEN"
  }
]
```

### Obtener Ruta con Paradas

```bash
GET /api/routes/{id}
```

**Ejemplo:**
```bash
curl https://your-app.vercel.app/api/routes/22222222-2222-2222-2222-222222222222
```

**Respuesta:**
```json
{
  "route": {
    "id": "22222222-2222-2222-2222-222222222222",
    "name": "Ruta Centro - Norte",
    "origin_name": "Plaza de Armas",
    "origin_lat": -12.046374,
    "origin_lon": -77.042793,
    "destination_name": "Terminal Norte",
    "destination_lat": -12.001234,
    "destination_lon": -77.051234,
    "base_price_cents": 1550,
    "currency": "PEN",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "stops": [
    {
      "id": "44444444-4444-4444-4444-444444444444",
      "name": "Plaza de Armas",
      "stop_order": 1
    },
    {
      "id": "55555555-5555-5555-5555-555555555555",
      "name": "Plaza Mayor",
      "stop_order": 2
    },
    {
      "id": "66666666-6666-6666-6666-666666666666",
      "name": "Terminal Norte",
      "stop_order": 3
    }
  ]
}
```

### Crear Viaje

```bash
POST /api/trips
Content-Type: application/json

{
  "route_id": "22222222-2222-2222-2222-222222222222",
  "pickup_stop_id": "44444444-4444-4444-4444-444444444444",
  "dropoff_stop_id": "66666666-6666-6666-6666-666666666666",
  "payment_method": "plin"
}
```

**Notas:**
- `pickup_stop_id` y `dropoff_stop_id` son opcionales (pueden ser `null`)
- Si ambos están presentes, deben ser diferentes
- `payment_method` debe ser: `cash`, `yape`, o `plin`
- El `user_id` está hardcodeado como `11111111-1111-1111-1111-111111111111`

**Respuesta:**
```json
{
  "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "route_id": "22222222-2222-2222-2222-222222222222",
  "passenger_id": "11111111-1111-1111-1111-111111111111",
  "pickup_stop_id": "44444444-4444-4444-4444-444444444444",
  "dropoff_stop_id": "66666666-6666-6666-6666-666666666666",
  "status": "requested",
  "payment_method": "plin",
  "price_cents": 1550,
  "currency": "PEN",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### Obtener Viaje

```bash
GET /api/trips/{id}
```

**Ejemplo:**
```bash
curl https://your-app.vercel.app/api/trips/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

## Deploy en Vercel

### Desde GitHub

1. **Push tu código a GitHub:**

```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/tu-usuario/tu-repo.git
git push -u origin main
```

2. **Conectar con Vercel:**

- Ve a [vercel.com](https://vercel.com)
- Click en "New Project"
- Importa tu repositorio de GitHub
- Vercel detectará automáticamente el proyecto Go

3. **Configurar variables de entorno:**

En el dashboard de Vercel, ve a Settings → Environment Variables y agrega:

```
DATABASE_URL=postgresql://user:password@ep-xxx.region.aws.neon.tech/dbname?sslmode=require
```

4. **Deploy:**

Vercel desplegará automáticamente. Cada push a `main` generará un nuevo deploy.

### Desde CLI

```bash
# Instalar Vercel CLI
npm i -g vercel

# Login
vercel login

# Deploy
vercel

# Agregar variable de entorno
vercel env add DATABASE_URL
```

## Testing con curl

### Health Check

```bash
curl https://your-app.vercel.app/api/health
```

### Listar Rutas

```bash
curl https://your-app.vercel.app/api/routes
```

### Obtener Ruta Específica

```bash
curl https://your-app.vercel.app/api/routes/22222222-2222-2222-2222-222222222222
```

### Crear Viaje

```bash
curl -X POST https://your-app.vercel.app/api/trips \
  -H "Content-Type: application/json" \
  -d '{
    "route_id": "22222222-2222-2222-2222-222222222222",
    "pickup_stop_id": "44444444-4444-4444-4444-444444444444",
    "dropoff_stop_id": "66666666-6666-6666-6666-666666666666",
    "payment_method": "plin"
  }'
```

### Obtener Viaje

```bash
curl https://your-app.vercel.app/api/trips/{trip_id}
```

## Notas Importantes

- **Autenticación:** No implementada en MVP. El `user_id` está hardcodeado.
- **Status de Viajes:** Siempre inicia como `"requested"`.
- **Precios:** Se toman del `base_price_cents` de la ruta.
- **CORS:** Configurado para permitir todos los orígenes (`*`).
- **Errores:** Respuestas JSON con formato `{"error": "mensaje"}`.

## Códigos de Estado HTTP

- `200 OK` - Solicitud exitosa
- `201 Created` - Recurso creado exitosamente
- `400 Bad Request` - Datos de entrada inválidos
- `404 Not Found` - Recurso no encontrado
- `500 Internal Server Error` - Error del servidor

## Próximos Pasos (Post-MVP)

- [ ] Implementar autenticación real (JWT)
- [ ] Agregar endpoints para drivers
- [ ] Sistema de matching driver-passenger
- [ ] Tracking en tiempo real
- [ ] Notificaciones push
- [ ] Historial de viajes
- [ ] Sistema de calificaciones
- [ ] Pagos integrados (Yape/Plin)

## Soporte

Para problemas o preguntas, contacta al equipo de desarrollo.
