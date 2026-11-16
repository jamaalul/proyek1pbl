-- Migration 000001 Down: Rollback Skema Autentikasi dan Autorisasi (Siakad Mini)

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
