CREATE TABLE car_groups (
    id          UUID                        NOT NULL DEFAULT gen_random_uuid(),
    name        TEXT                        NOT NULL,
    description TEXT,
    size_category TEXT,
    created_at  TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMP WITH TIME ZONE,
    created_by  TEXT                        NOT NULL,
    updated_by  TEXT                        NOT NULL,
    deleted     BOOL                        NOT NULL DEFAULT FALSE,

    CONSTRAINT car_groups_pkey PRIMARY KEY (id),
    CONSTRAINT car_groups_name_key UNIQUE (name)
);
