-- +goose Up
INSERT INTO 
"subscriptions" (id, service_name, price, user_id, start_date, end_date)
VALUES
(1, 'Service1', 100, gen_random_uuid(), '2026-03-12 14:47:51', '2026-03-15 14:47:51'),
(2, 'Service2', 200, gen_random_uuid(), '2026-03-13 14:47:51', NULL),
(3, 'Service3', 300, gen_random_uuid(), '2026-03-14 14:47:51', '2026-03-18 14:47:51'),
(4, 'Service4', 400, gen_random_uuid(), '2026-03-15 14:47:51', '2026-03-20 14:47:51'),
(5, 'Service5', 500, gen_random_uuid(), '2026-03-16 14:47:51', NULL),
(6, 'Service6', 600, gen_random_uuid(), '2026-03-17 14:47:51', '2026-03-22 14:47:51'),
(7, 'Service7', 700, gen_random_uuid(), '2026-03-18 14:47:51', '2026-03-25 14:47:51'),
(8, 'Service8', 800, gen_random_uuid(), '2026-03-19 14:47:51', NULL),
(9, 'Service9', 900, gen_random_uuid(), '2026-03-20 14:47:51', '2026-03-28 14:47:51'),
(10, 'Service10', 1000, gen_random_uuid(), '2026-03-21 14:47:51', '2026-03-30 14:47:51');

-- +goose Down
DELETE FROM "subscriptions";