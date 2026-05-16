-- Seed runner account for native runner app smoke tests.
-- Credentials:
--   email: runner.test@kilat.my
--   password: TestRunner123!

\connect kilat_identity

INSERT INTO users (
    id,
    email,
    phone,
    password_hash,
    full_name,
    role,
    is_verified,
    version,
    created_at,
    updated_at
)
VALUES (
    '11111111-1111-4111-8111-111111111111',
    'runner.test@kilat.my',
    '+60123456780',
    '$2a$10$jufJku7AhVYg5808vjEXAO0COaclKkPiOgj3nLdpryDcrNOYlS4Qy',
    'Test Runner',
    'runner',
    true,
    1,
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE SET
    phone = EXCLUDED.phone,
    password_hash = EXCLUDED.password_hash,
    full_name = EXCLUDED.full_name,
    role = EXCLUDED.role,
    is_verified = EXCLUDED.is_verified,
    updated_at = NOW();

\connect kilat_runner

INSERT INTO runners (
    id,
    user_id,
    full_name,
    phone,
    vehicle_type,
    vehicle_plate,
    vehicle_model,
    vehicle_year,
    air_conditioned,
    session_status,
    rating,
    total_trips,
    completion_rate,
    version,
    created_at,
    updated_at
)
VALUES (
    '22222222-2222-4222-8222-222222222222',
    '11111111-1111-4111-8111-111111111111',
    'Test Runner',
    '+60123456780',
    'motorcycle',
    'KLT1234',
    'Yamaha Y15ZR',
    2022,
    false,
    'inactive',
    5.0,
    0,
    1.0,
    1,
    NOW(),
    NOW()
)
ON CONFLICT (user_id) DO UPDATE SET
    full_name = EXCLUDED.full_name,
    phone = EXCLUDED.phone,
    vehicle_type = EXCLUDED.vehicle_type,
    vehicle_plate = EXCLUDED.vehicle_plate,
    vehicle_model = EXCLUDED.vehicle_model,
    vehicle_year = EXCLUDED.vehicle_year,
    air_conditioned = EXCLUDED.air_conditioned,
    updated_at = NOW();

INSERT INTO crate_specs (
    id,
    runner_id,
    size,
    pet_types,
    max_weight_kg,
    width_cm,
    height_cm,
    depth_cm,
    ventilated,
    temperature_controlled,
    created_at
)
VALUES (
    '33333333-3333-4333-8333-333333333333',
    '22222222-2222-4222-8222-222222222222',
    'medium',
    '["cat", "dog"]'::jsonb,
    15.0,
    60.0,
    45.0,
    45.0,
    true,
    false,
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    pet_types = EXCLUDED.pet_types,
    max_weight_kg = EXCLUDED.max_weight_kg,
    width_cm = EXCLUDED.width_cm,
    height_cm = EXCLUDED.height_cm,
    depth_cm = EXCLUDED.depth_cm,
    ventilated = EXCLUDED.ventilated,
    temperature_controlled = EXCLUDED.temperature_controlled;
