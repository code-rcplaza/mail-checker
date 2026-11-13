# FOKKO-0112C-SUBS/TASK-001 · Dos usuarios no se ven
Feature: `cmd/api`

Como usuario quiero que nadie vea mis datos, para mi privacidad.

## Criterios
- Dados dos usuarios con cartolas distintas, cuando cada uno pide su resumen, entonces ve solo lo suyo.
- Dada una clave de línea del usuario A, cuando el usuario B intenta decidir sobre ella, entonces responde 404.
- Dado un `publicId` de una manual de A, cuando B intenta editarla, entonces responde 404.

## Fuera
- Nada.

## Notas
- Ninguna.
