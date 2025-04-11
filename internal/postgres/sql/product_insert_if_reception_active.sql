WITH check_reception AS (
    SELECT id 
    FROM receptions 
    WHERE id = :reception_id AND status = :reception_status
    FOR UPDATE
)
INSERT INTO productions (id, reception_id, type, datetime)
SELECT :product_id, :reception_id, :type, :datetime
FROM check_reception
WHERE EXISTS (SELECT 1 FROM check_reception)
RETURNING id;