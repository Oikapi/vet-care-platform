CREATE DATABASE IF NOT EXISTS clinic_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE clinic_db;

CREATE TABLE schedules (
    id INT AUTO_INCREMENT PRIMARY KEY,
    doctor_id INT NOT NULL,
    clinic_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL
);

CREATE TABLE inventories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    medicine_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    threshold INT NOT NULL,
    clinic_id INT NOT NULL
);

CREATE TABLE doctors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    clinic_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO schedules (doctor_id, clinic_id, start_time, end_time) VALUES
(1, 1,'2025-04-30 09:00:00', '2025-04-30 10:00:00'),
(1, 1,'2025-04-30 10:00:00', '2025-04-30 11:00:00');

INSERT INTO inventories (medicine_name, quantity, threshold, clinic_id) VALUES
('Paracetamol', 50, 20, 1),
('Antibiotic', 10, 15, 1);

INSERT INTO doctors (name, email, clinic_id) VALUES
('Dr.Andre','andre@example.com', 1),
('Dr.Williams','williams@example.com', 1);

CREATE DATABASE IF NOT EXISTS pet_profiles;
USE pet_profiles;

CREATE TABLE IF NOT EXISTS pets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    species VARCHAR(50),
    name VARCHAR(50),
    breed VARCHAR(50),
    gender VARCHAR(10),
    age INT,
    user_id INT
);


CREATE DATABASE IF NOT EXISTS vetcare_appointments;
USE vetcare_appointments;

CREATE TABLE IF NOT EXISTS slots (
  id INT AUTO_INCREMENT PRIMARY KEY,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME DEFAULT NULL,
  doctor_id INT NOT NULL,
  slot_time DATETIME NOT NULL,
  is_booked BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS appointments (
  id INT AUTO_INCREMENT PRIMARY KEY,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME DEFAULT NULL,
  client_id INT NOT NULL,
  doctor_id INT NOT NULL,
  clinic_id INT NOT NULL,
  slot_id INT NOT NULL,
  status VARCHAR(50) NOT NULL,
  telegram_id VARCHAR(50),
  FOREIGN KEY (slot_id) REFERENCES slots(id)
);

INSERT INTO slots (id, created_at, updated_at, deleted_at, doctor_id, slot_time, is_booked)
VALUES
  (1, NOW(), NOW(), NULL, 101, '2025-05-04 09:00:00', FALSE),
  (2, NOW(), NOW(), NULL, 101, '2025-05-04 10:00:00', TRUE),
  (3, NOW(), NOW(), NULL, 102, '2025-05-05 11:00:00', FALSE),
  (4, NOW(), NOW(), NULL, 103, '2025-05-05 12:00:00', TRUE);

INSERT INTO appointments (id, created_at, updated_at, deleted_at, client_id, doctor_id, clinic_id, slot_id, status, telegram_id)
VALUES
  (1, NOW(), NOW(), NULL, 201, 101, 301, 2, 'confirmed', '123456789'),
  (2, NOW(), NOW(), NULL, 202, 103, 301, 4, 'confirmed', '987654321');

CREATE DATABASE IF NOT EXISTS auth_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE auth_db;

CREATE TABLE IF NOT EXISTS clinic (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  phone VARCHAR(50) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  photo VARCHAR(255) DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS doctor (
  id INT AUTO_INCREMENT PRIMARY KEY,
  firstName VARCHAR(100) NOT NULL,
  lastName VARCHAR(100) NOT NULL,
  specialization VARCHAR(100) NOT NULL,
  clinicId INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (clinicId) REFERENCES clinic(id) ON DELETE SET NULL
);
INSERT INTO clinic (name, phone, email, password, photo)
VALUES 
  ('ВетКлиника Север', '+7-900-111-2233', 'northvet@example.com', 'hashedpassword1', 'photo1.jpg'),
  ('ВетКлиника Юг', '+7-900-222-3344', 'southvet@example.com', 'hashedpassword2', NULL);

INSERT INTO doctor (firstName, lastName, specialization, clinicId)
VALUES 
  ('Анна', 'Иванова', 'Терапевт', 1),
  ('Борис', 'Смирнов', 'Хирург', 1),
  ('Вера', 'Кузнецова', 'Кардиолог', 2);

CREATE DATABASE IF NOT EXISTS forum_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE forum_db;