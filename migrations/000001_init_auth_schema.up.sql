-- Migration 000001: Inisialisasi Skema Autentikasi dan Autorisasi (Siakad Mini)
-- Disesuaikan dengan entitas app/model/user.go dan app/repository/role_repository.go

-- 1. Tabel roles
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Tabel permissions
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Tabel relasi role_permissions (Many-to-Many)
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- 4. Tabel users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL REFERENCES roles(name) ON UPDATE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indeks untuk optimasi query pencarian
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- 5. Seed Data Awal
-- Inisialisasi Roles default
INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator Sistem Pengelola Data'),
    ('mahasiswa', 'Pengguna Mahasiswa')
ON CONFLICT (name) DO NOTHING;

-- Inisialisasi Permissions default
INSERT INTO permissions (name, description) VALUES
    ('*', 'Akses penuh ke seluruh modul sistem'),
    ('krs:read', 'Melihat daftar rencana studi'),
    ('krs:write', 'Mengisi dan mengajukan rencana studi'),
    ('khs:read', 'Melihat kartu hasil studi')
ON CONFLICT (name) DO NOTHING;

-- Relasi Role 'admin' -> Akses Wildcard ('*')
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin' AND p.name = '*'
ON CONFLICT DO NOTHING;

-- Relasi Role 'mahasiswa' -> Izin Operasional KRS dan KHS
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'mahasiswa' AND p.name IN ('krs:read', 'krs:write', 'khs:read')
ON CONFLICT DO NOTHING;
