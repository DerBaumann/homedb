-- +goose Up
-- +goose StatementBegin
BEGIN;

CREATE TYPE item_unit AS ENUM('g', 'l', 'pcs');

CREATE TABLE items (
  id uuid DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  amount INT NOT NULL CHECK (amount > 0),
  unit item_unit NOT NULL,
  PRIMARY KEY (id)
);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

DROP TABLE items;
DROP TYPE item_unit;

COMMIT;
-- +goose StatementEnd
