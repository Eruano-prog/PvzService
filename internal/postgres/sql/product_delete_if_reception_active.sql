WITH locked_reception AS (
    SELECT id
    FROM receptions
    WHERE pvz_id = :pvz_id AND status = :status
    FOR UPDATE
    LIMIT 1
),
last_production AS (
    SELECT id
    FROM products
    WHERE reception_id IN (SELECT id FROM locked_reception)
    ORDER BY datetime DESC
    LIMIT 1
    FOR UPDATE
)
DELETE FROM products
WHERE id IN (SELECT id FROM last_production)
RETURNING id;