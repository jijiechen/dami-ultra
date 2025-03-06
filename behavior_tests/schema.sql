DROP TABLE IF EXISTS `konnect_services`;

CREATE TABLE `konnect_services` (
  `id` INTEGER PRIMARY KEY,
  `name` TEXT NOT NULL,
  `description` TEXT,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS `konnect_service_versions`;

CREATE TABLE `konnect_service_versions` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `konnect_service_id` INTEGER NOT NULL,
  `version` TEXT NOT NULL,
  `host` TEXT NOT NULL,
  `port` INTEGER NOT NULL,
  `path` TEXT NOT NULL,
  `protocol` TEXT NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);
