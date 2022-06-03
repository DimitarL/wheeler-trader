CREATE TYPE vehicle_type AS ENUM ('car', 'motorbike', 'truck');

CREATE TABLE IF NOT EXISTS vehicle (
  id SERIAL,
  type vehicle_type NOT NULL,
  make VARCHAR(32) NOT NULL,
  model VARCHAR(32) NOT NULL,
  horsepower INT NOT NULL,
  CONSTRAINT pk_vehicle PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS sale (
  vehicle_id SERIAL,
  price NUMERIC(10, 2) NOT NULL,
  CONSTRAINT pk_sale PRIMARY KEY(vehicle_id),
  CONSTRAINT fk_vehicle FOREIGN KEY(vehicle_id) REFERENCES vehicle(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rent (
  vehicle_id SERIAL,
  daily_price NUMERIC(10, 2) NOT NULL,
  weekly_price NUMERIC(10, 2) NOT NULL,
  monthly_price NUMERIC(10, 2) NOT NULL,
  CONSTRAINT pk_rent PRIMARY KEY(vehicle_id),
  CONSTRAINT fk_vehicle FOREIGN KEY(vehicle_id) REFERENCES vehicle(id) ON DELETE CASCADE
);
