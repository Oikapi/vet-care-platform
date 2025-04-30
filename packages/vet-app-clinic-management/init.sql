CREATE DATABASE IF NOT EXISTS clinic_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE clinic_db;

CREATE TABLE schedules (
    id INT AUTO_INCREMENT PRIMARY KEY,
    doctor_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL
);

CREATE TABLE inventories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    medicine_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    threshold INT NOT NULL
);

INSERT INTO schedules (doctor_id, start_time, end_time) VALUES
(1, '2025-04-30 09:00:00', '2025-04-30 10:00:00'),
(1, '2025-04-30 10:00:00', '2025-04-30 11:00:00');

INSERT INTO inventories (medicine_name, quantity, threshold) VALUES
('Paracetamol', 50, 20),
('Antibiotic', 10, 15);