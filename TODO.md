# Backend Fixes Plan

## Steps to Complete

- [x] A1: Adjust backend/api/index.go for correct Vercel routing (remove /api prefix inside chi router)
- [x] A2: Ensure backend/internal/trips/handlers.go payment_method accepts only cash/yape/plin (already done)
- [x] A3: Ensure backend/sql/schema.sql enum includes cash/yape/plin (already done)
- [x] A4: Adjust backend/sql/seed.sql and README.md with UUIDs and curl examples (update README POST to use plin)
- [x] A5: Create minimal backend/vercel.json (already exists)
