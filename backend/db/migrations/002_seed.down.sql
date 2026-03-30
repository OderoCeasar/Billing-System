BEGIN;

DELETE FROM packages
WHERE name IN (
    '30 Min Plan',
    '1 Hour Plan',
    '3 Hour Plan',
    '5 Hour Plan',
    '12 Hour Plan',
    '24 Hour Plan',
    '2 Day Plan',
    '1 Week Plan',
    '2 Week Plan',
    '2 Week Plan (2 Devices)',
    'Monthly Unlimited',
    'Monthly (2 Devices)'
)
  AND NOT EXISTS (
      SELECT 1 FROM payments p WHERE p.package_id = packages.id
  )
  AND NOT EXISTS (
      SELECT 1 FROM sessions s WHERE s.package_id = packages.id
  );

COMMIT;
