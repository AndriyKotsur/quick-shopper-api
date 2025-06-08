-- +goose Up
-- +goose StatementBegin
CREATE TABLE venues (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location TEXT NOT NULL,
    type VARCHAR(20) CHECK (type IN ('BAR', 'CAFFE', 'KARAOKE', 'RESTAURANT')) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS venues;
-- +goose StatementEnd
