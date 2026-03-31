BEGIN;

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    '30 Min Plan',
    '30 minutes unlimited access',
    'time',
    5.00,
    30,
    NULL,
    0,
    0,
    1,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = '30 Min Plan'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    '1 Hour Plan',
    '1 hour unlimited access',
    'time',
    10.00,
    60,
    NULL,
    0,
    0,
    1,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = '1 Hour Plan'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    '3 Hour Plan',
    '3 hours unlimited access',
    'time',
    15.00,
    180,
    NULL,
    0,
    0,
    1,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = '3 Hour Plan'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    '5 Hour Plan',
    '5 hours unlimited access',
    'time',
    20.00,
    300,
    NULL,
    0,
    0,
    1,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = '5 Hour Plan'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    '12 Hour Plan',
    '12 hours unlimited access',
    'time',
    25.00,
    720,
    NULL,
    0,
    0,
    1,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = '12 Hour Plan'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    '24 Hour Plan',
    '24 hours unlimited access (1 day)',
    'time',
    35.00,
    1440,
    NULL,
    0,
    0,
    1,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = '24 Hour Plan'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    '2 Day Plan',
    '2 days unlimited access',
    'time',
    60.00,
    2880,
    NULL,
    0,
    0,
    2,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = '2 Day Plan'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    '1 Week Plan',
    '7 days unlimited access',
    'time',
    160.00,
    10080,
    NULL,
    0,
    0,
    7,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = '1 Week Plan'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    '2 Week Plan',
    '14 days unlimited access',
    'time',
    260.00,
    20160,
    NULL,
    0,
    0,
    14,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = '2 Week Plan'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    '2 Week Plan (2 Devices)',
    '14 days unlimited access for 2 devices',
    'time',
    350.00,
    20160,
    NULL,
    0,
    0,
    14,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = '2 Week Plan (2 Devices)'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    'Monthly Unlimited',
    '30 days unlimited access',
    'time',
    450.00,
    43200,
    NULL,
    0,
    0,
    30,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = 'Monthly Unlimited'
);

INSERT INTO packages (
    name,
    description,
    package_type,
    price,
    duration_minutes,
    data_limit_mb,
    speed_limit_up,
    speed_limit_down,
    validity_days,
    is_active
)
SELECT
    'Monthly (2 Devices)',
    '30 days unlimited access for 2 devices',
    'time',
    720.00,
    43200,
    NULL,
    0,
    0,
    30,
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM packages WHERE name = 'Monthly (2 Devices)'
);

COMMIT;
